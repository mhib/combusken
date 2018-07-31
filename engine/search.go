package engine

import (
	"fmt"
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
	return evaledValue
}

var countPositions int

func alphaBeta(pos *Position, depth, alpha, beta, evaluation, mate int) int {
	var child Position
	var val int
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
		val = -alphaBeta(&child, depth-1, -beta, -alpha, -evaled[i].value, mate-1)
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

func depSearch(pos *Position, depth int, lastBestMove Move, mate int) (Move, bool) {
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
		return bestMove, alpha > Mate-500
	}
	for i := range evaled {
		child = evaled[i].position
		val := -alphaBeta(&child, depth-1, -beta, -alpha, -evaled[i].value, mate-1)
		if val >= beta {
			return bestMove, false
		}
		if val > alpha {
			alpha = val
			bestMove = evaled[i].move
		}
	}
	return bestMove, alpha > Mate-500
}

func Search(pos *Position) Move {
	var lastBestMove, bestMove Move
	var mate bool
	for i := 1; ; i++ {
		bestMove, mate = depSearch(pos, i, lastBestMove, Mate)
		if mate {
			fmt.Println(i, countPositions)
			return bestMove
		}
		if bestMove == 0 {
			return lastBestMove
		} else {
			lastBestMove = bestMove
		}
		if countPositions >= 700000 || i > 70 {
			fmt.Println(i, countPositions)
			return bestMove
		}
	}
}
