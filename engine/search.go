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
	var buffer [256]Move
	var evaledPositions = EvaledPositions(make([]EvaledPosition, 0, 40))
	var child Position
	for _, move := range pos.GenerateAllMoves(buffer[:]) {
		if pos.MakeMove(move, &child) {
			evaledPositions = append(evaledPositions, EvaledPosition{child, move, -Evaluate(&child)})
		}
	}
	return evaledPositions
}

func (e *Engine) quiescence(pos *Position, alpha, beta, mate, evaluation int, timedOut *bool) int {
	if *timedOut {
		return evaluation
	}

	var buffer [40]EvaledMove
	val := extensiveEval(pos, evaluation, mate)

	if val >= beta {
		return beta
	}
	if alpha < val {
		alpha = val
	}

	if pos.IsInCheck() {
		evaled := generateAllLegalMoves(pos)
		sort.Sort(evaled)
		moveCount := 0
		for i := range evaled {
			child := evaled[i].position
			val := -extensiveEval(&child, -evaled[i].value, mate-1)
			moveCount++
			if val > alpha {
				alpha = val
			}
			if val >= beta {
				return beta
			}
		}
		if moveCount > 0 {
			return alpha
		}
		return val
	}

	moves := pos.GenerateAllCaptures(buffer[:])

	moveCount := 0
	var child Position

	for i := range moves {
		maxMoveToFirst(moves[i:])
		if pos.MakeMove(moves[i].Move, &child) {
			val = -e.quiescence(&child, -beta, -alpha, mate-1, Evaluate(&child), timedOut)
			moveCount++
		}
		if val > alpha {
			alpha = val
			if val >= beta {
				return beta
			}
		}
	}

	if moveCount == 0 {
		return val
	}

	return alpha
}

func contempt(pos *Position) int {
	return 0
}

func maxMoveToFirst(moves []EvaledMove) {
	maxIdx := 0
	for i := 1; i < len(moves); i++ {
		if moves[i].Value > moves[maxIdx].Value {
			maxIdx = i
		}
	}
	moves[0], moves[maxIdx] = moves[maxIdx], moves[0]
}

func maxToFirst(positions EvaledPositions) {
	maxIdx := 0
	for i := 1; i < len(positions); i++ {
		if positions[i].value > positions[maxIdx].value {
			maxIdx = i
		}
	}
	positions[0], positions[maxIdx] = positions[maxIdx], positions[0]
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

func (e *Engine) alphaBeta(pos *Position, depth, alpha, beta, evaluation, mate int, timedOut *bool) int {
	var child Position
	var val int
	if *timedOut {
		return 0
	}
	var alphaOrig = alpha
	hashMove := NullMove
	ttEntry := e.TransTable.Get(pos.Key)
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
		hashMove = ttEntry.bestMove
	}
	if depth == 0 {
		countPositions++
		val = e.quiescence(pos, alpha, beta, mate, evaluation, timedOut)
		//e.TransTable.Set(depth, val, TransExact, pos.Key, NullMove)
		return val
	}

	if pos.LastMove != NullMove && depth >= 4 && !pos.IsInCheck() && !isLateEndGame(pos) {
		pos.MakeNullMove(&child)
		val = -e.alphaBeta(&child, depth-3, -beta, -beta+1, -evaluation, mate, timedOut)
		if val >= beta {
			return beta
		}
	}

	evaled := generateAllLegalMoves(pos)
	if len(evaled) == 0 {
		countPositions++
		if pos.IsInCheck() {
			val = -mate
			e.TransTable.Set(depth, val, TransExact, pos.Key, NullMove)
			return val
		}
		val = contempt(pos)
		e.TransTable.Set(depth, val, TransExact, pos.Key, NullMove)
		return val
	}
	val = MinInt
	var tmpVal int
	if hashMove != 0 {
		moveToFirst(evaled, hashMove)
		sort.Sort(evaled[1:])
	} else {
		sort.Sort(evaled)
	}
	bestMove := NullMove
	for i := range evaled {
		child = evaled[i].position
		tmpVal = -e.alphaBeta(&child, depth-1, -beta, -alpha, -evaled[i].value, mate-1, timedOut)
		if tmpVal > val {
			val = tmpVal
			bestMove = evaled[i].move
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
	ttEntry.bestMove = bestMove
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

func (e *Engine) depSearch(pos *Position, depth int, lastBestMove Move, mate int, resultChan chan result, timedOut *bool) {
	var child Position
	var bestMove Move
	evaled := generateAllLegalMoves(pos)
	if lastBestMove != 0 {
		moveToFirst(evaled, lastBestMove)
		sort.Sort(evaled[1:])
	} else {
		sort.Sort(evaled)
	}
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
		e.TransTable.Set(1, alpha, TransExact, pos.Key, bestMove)
		resultChan <- result{bestMove, alpha > Mate-500}
		return
	}
	for i := range evaled {
		child = evaled[i].position
		val := -e.alphaBeta(&child, depth-1, -beta, -alpha, -evaled[i].value, mate-1, timedOut)
		if val > alpha {
			alpha = val
			bestMove = evaled[i].move
		}
	}
	e.TransTable.Set(depth, alpha, TransExact, pos.Key, bestMove)
	resultChan <- result{bestMove, alpha > Mate-500}
}

func (e *Engine) TimeSearch(ctx context.Context, pos *Position) Move {
	var lastBestMove Move
	for i := 1; ; i++ {
		resultChan := make(chan result, 1)
		timedOut := false
		go e.depSearch(pos, i, lastBestMove, Mate, resultChan, &timedOut)
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

func (e *Engine) DepthSearch(pos *Position, depth int) Move {
	timedOut := false
	var lastBestMove Move
	for i := 1; i < depth; i++ {
		resultChan := make(chan result, 1)
		go e.depSearch(pos, i, lastBestMove, Mate, resultChan, &timedOut)
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

func (e *Engine) CountSearch(ctx context.Context, pos *Position, count int) Move {
	timedOut := false
	var lastBestMove Move
	for i := 1; ; i++ {
		countPositions = 0
		resultChan := make(chan result, 1)
		go e.depSearch(pos, i, lastBestMove, Mate, resultChan, &timedOut)
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
