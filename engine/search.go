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

func alphaBeta(pos *Position, depth, alpha, beta, evaluation, mate int, timedOut *bool) int {
	var child Position
	var val int
	if *timedOut {
		return 0
	}
	if depth == 0 {
		countPositions++
		return extensiveEval(pos, evaluation, mate)
	}
	evaled := generateAllLegalMoves(pos)
	if len(evaled) == 0 {
		if pos.IsInCheck() {
			return -mate
		}
		return contempt(pos)
	}
	for i := range evaled {
		child = evaled[i].position
		val = -alphaBeta(&child, depth-1, -beta, -alpha, -evaled[i].value, mate-1, timedOut)
		if val >= beta {
			return beta
		}
		if val > alpha {
			alpha = val
		}
	}
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

func depSearch(pos *Position, depth int, lastBestMove Move, mate int, resultChan chan result, timedOut *bool) {
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
		resultChan <- result{bestMove, alpha > Mate-500}
		return
	}
	for i := range evaled {
		child = evaled[i].position
		val := -alphaBeta(&child, depth-1, -beta, -alpha, -evaled[i].value, mate-1, timedOut)
		if val >= beta {
			resultChan <- result{bestMove, false}
			return
		}
		if val > alpha {
			alpha = val
			bestMove = evaled[i].move
		}
	}
	resultChan <- result{bestMove, alpha > Mate-500}
}

func TimeSearch(ctx context.Context, pos *Position) Move {
	var lastBestMove Move
	for i := 1; ; i++ {
		resultChan := make(chan result, 1)
		timedOut := false
		go depSearch(pos, i, lastBestMove, Mate, resultChan, &timedOut)
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

func DepthSearch(pos *Position, depth int) Move {
	timedOut := false
	var lastBestMove Move
	for i := 1; i < depth; i++ {
		resultChan := make(chan result, 1)
		go depSearch(pos, i, lastBestMove, Mate, resultChan, &timedOut)
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
