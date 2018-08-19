package engine

import (
	"context"
	"sort"

	. "github.com/mhib/combusken/backend"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1
const Mate = 1000000

type EvaledPosition struct {
	position Position
	move     Move
	value    int
}

type EvaledPositions []EvaledPosition

var generatedMoves [256]Move

func (s EvaledPositions) Len() int {
	return len(s)
}

func (s EvaledPositions) Less(i, j int) bool {
	return s[i].value > s[j].value
}

func (s EvaledPositions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func generateAllLegalMoves(pos *Position) EvaledPositions {
	var evaledPositions = EvaledPositions(make([]EvaledPosition, 0, 40))
	var child Position
	for _, move := range pos.GenerateAllMoves(generatedMoves[:]) {
		if pos.MakeMove(move, &child) {
			evaledPositions = append(evaledPositions, EvaledPosition{child, move, -Evaluate(&child)})
		}
	}
	sort.Sort(evaledPositions)
	return evaledPositions
}

func contempt(pos *Position) int {
	return 0
}

// Evals that checks for mate
func extensiveEval(pos *Position, evaledValue, mate int) int {
	if len(generateAllLegalMoves(pos)) == 0 {
		if pos.IsInCheck() {
			return -mate
		}
		return contempt(pos)
	}
	if pos.FiftyMove > 100 {
		return contempt(pos)
	}
	return evaledValue
}

var countPositions int

func alphaBeta(pos *Position, transTable *TransTable, depth, alpha, beta, evaluation, mate int, timedOut *bool) int {
	var child Position
	var val int
	if *timedOut {
		return 0
	}
	var alphaOrig = alpha
	ttEntry := transTable.Get(pos.Key)
	if ttEntry.key == pos.Key && ttEntry.depth >= int32(depth) {
		if ttEntry.flag == TransExact {
			return ttEntry.value
		}
		if ttEntry.flag == TransAlpha && ttEntry.value > alpha {
			alpha = ttEntry.value
		}
		if ttEntry.flag == TransBeta && ttEntry.value < beta {
			beta = ttEntry.value
		}
		if alpha >= beta {
			return ttEntry.value
		}
	}
	if depth == 0 {
		countPositions++
		val = extensiveEval(pos, evaluation, mate)
		transTable.Set(depth, val, TransExact, pos.Key)
		return val
	}

	if pos.LastMove != NullMove && depth >= 4 && !pos.IsInCheck() && !isLateEndGame(pos) {
		pos.MakeNullMove(&child)
		val = -alphaBeta(&child, transTable, depth-3, -beta, -beta+1, -evaluation, mate, timedOut)
		if val >= beta {
			return beta
		}
	}

	evaled := generateAllLegalMoves(pos)
	if len(evaled) == 0 {
		countPositions++
		if pos.IsInCheck() {
			val = -mate
			transTable.Set(depth, val, TransExact, pos.Key)
			return val
		}
		val = contempt(pos)
		transTable.Set(depth, val, TransExact, pos.Key)
		return val
	}
	val = MinInt
	var tmpVal int
	for i := range evaled {
		child = evaled[i].position
		tmpVal = -alphaBeta(&child, transTable, depth-1, -beta, -alpha, -evaled[i].value, mate-1, timedOut)
		if tmpVal > val {
			val = tmpVal
		}
		if tmpVal > alpha {
			alpha = tmpVal
		}
		if alpha >= beta {
			break
		}
	}
	if val <= alphaOrig {
		ttEntry.flag = TransBeta
	} else if val >= beta {
		ttEntry.flag = TransAlpha
	} else {
		ttEntry.flag = TransExact
	}
	ttEntry.depth = int32(depth)
	ttEntry.value = val
	return alpha
}

func moveToFirst(list []EvaledPosition, m Move) {
	if m == 0 {
		return
	}
	for i := range list {
		if list[i].move == m {
			list[i], list[0] = list[0], list[i]
			return
		}
	}
}

type result struct {
	Move
	bool
}

func depSearch(pos *Position, transTable *TransTable, depth int, lastBestMove Move, mate int, resultChan chan result, timedOut *bool) {
	var child Position
	var bestMove Move
	evaled := generateAllLegalMoves(pos)
	moveToFirst(evaled, lastBestMove)
	var alpha = -MaxInt
	var beta = MaxInt
	countPositions = 0
	if depth == 1 {
		for i := range evaled {
			child = evaled[i].position
			val := -extensiveEval(&child, -evaled[i].value, mate-1)
			if val > alpha {
				alpha = val
				bestMove = evaled[i].move
			}
		}
		transTable.Set(1, alpha, TransExact, pos.Key)
		resultChan <- result{bestMove, alpha > Mate-500}
		return
	}
	for i := range evaled {
		child = evaled[i].position
		val := -alphaBeta(&child, transTable, depth-1, -beta, -alpha, -evaled[i].value, mate-1, timedOut)
		if val > alpha {
			alpha = val
			bestMove = evaled[i].move
		}
	}
	transTable.Set(depth, alpha, TransExact, pos.Key)
	resultChan <- result{bestMove, alpha > Mate-500}
}

func TimeSearch(ctx context.Context, pos *Position, transTable *TransTable) Move {
	var lastBestMove Move
	for i := 1; ; i++ {
		resultChan := make(chan result, 1)
		timedOut := false
		go depSearch(pos, transTable, i, lastBestMove, Mate, resultChan, &timedOut)
		select {
		case <-ctx.Done():
			timedOut = true
			return lastBestMove
		case res := <-resultChan:
			if res.bool {
				return res.Move
			}
			if res.Move == 0 {
				return lastBestMove
			} else {
				lastBestMove = res.Move
			}
			if i > 70 {
				return res.Move
			}
		}
	}
}

func DepthSearch(pos *Position, transTable *TransTable, depth int) Move {
	timedOut := false
	var lastBestMove Move
	for i := 1; i < depth; i++ {
		resultChan := make(chan result, 1)
		go depSearch(pos, transTable, i, lastBestMove, Mate, resultChan, &timedOut)
		res := <-resultChan
		if res.bool {
			return res.Move
		}
		if res.Move == 0 {
			return lastBestMove
		} else {
			lastBestMove = res.Move
		}
		if i > 70 {
			return res.Move
		}
	}
	return lastBestMove
}

func CountSearch(ctx context.Context, pos *Position, transTable *TransTable, count int) Move {
	timedOut := false
	var lastBestMove Move
	for i := 1; ; i++ {
		countPositions = 0
		resultChan := make(chan result, 1)
		go depSearch(pos, transTable, i, lastBestMove, Mate, resultChan, &timedOut)
		res := <-resultChan
		if res.bool {
			return res.Move
		}
		if res.Move == 0 {
			return lastBestMove
		} else {
			lastBestMove = res.Move
		}
		if i > 70 || countPositions > count {
			return lastBestMove
		}
	}
}
