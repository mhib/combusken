package engine

import (
	"fmt"
	"sort"

	. "github.com/mhib/combusken/backend"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

type EvaledPosition struct {
	position Position
	move     Move
	value    int
}

var buffer [256]EvaledPosition

type EvaledPositions []EvaledPosition

var evaledPositions = EvaledPositions(buffer[:])
var generatedMoves [256]Move

var positionCount = 0

func (s EvaledPositions) Len() int {
	return positionCount
}

func (s EvaledPositions) Less(i, j int) bool {
	return s[i].value > s[j].value
}

func (s EvaledPositions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func generateAllLegalMoves(pos *Position) {
	positionCount = 0
	var child Position
	for _, move := range pos.GenerateAllMoves(generatedMoves[:]) {
		if pos.MakeMove(move, &child) {
			evaledPositions[positionCount] = EvaledPosition{child, move, -Evaluate(&child)}
			positionCount++
		}
	}
	sort.Sort(evaledPositions)
}

var countPositions int

func alphaBeta(pos *Position, depth, alpha, beta, evaluation int) int {
	var counter int
	var child Position
	var evaled [256]EvaledPosition
	var val int
	if depth == 0 {
		countPositions++
		return evaluation
	}
	generateAllLegalMoves(pos)
	counter = positionCount
	copy(evaled[:], evaledPositions[0:positionCount])
	for i := 0; i < counter; i++ {
		child = evaled[i].position
		val = -alphaBeta(&child, depth-1, -beta, -alpha, -evaled[i].value)
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

func depSearch(pos *Position, depth int, lastBestMove Move) Move {
	var counter int
	var child Position
	var bestMove Move
	var evaled [256]EvaledPosition
	generateAllLegalMoves(pos)
	counter = positionCount
	copy(evaled[:], evaledPositions[0:positionCount])
	moveToFirst(evaled[:], lastBestMove)
	var alpha = -MaxInt
	var beta = MaxInt
	countPositions = 0
	for i := 0; i < counter; i++ {
		child = evaled[i].position
		val := -alphaBeta(&child, depth-1, -beta, -alpha, -evaled[i].value)
		if val >= beta {
			return bestMove
		}
		if val > alpha {
			alpha = val
			bestMove = evaled[i].move
		}
	}
	return bestMove
}

func Search(pos *Position) Move {
	var lastBestMove, bestMove Move
	for i := 1; ; i++ {
		bestMove = depSearch(pos, i, lastBestMove)
		if bestMove == 0 {
			return lastBestMove
		} else {
			lastBestMove = bestMove
		}
		if countPositions >= 7000000 || i > 70 {
			fmt.Println(i)
			return bestMove
		}
	}
}
