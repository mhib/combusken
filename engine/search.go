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
	var buffer [256]EvaledMove
	var evaledPositions = EvaledPositions(make([]EvaledPosition, 0, 40))
	var child Position
	for _, move := range pos.GenerateAllMoves(buffer[:]) {
		if pos.MakeMove(move.Move, &child) {
			evaledPositions = append(evaledPositions, EvaledPosition{child, move.Move, -Evaluate(&child)})
		}
	}
	return evaledPositions
}

func areAnyLegalMoves(pos *Position) bool {
	var buffer [256]EvaledMove
	var child Position
	for _, move := range pos.GenerateAllMoves(buffer[:]) {
		if pos.MakeMove(move.Move, &child) {
			return true
		}
	}
	return false
}

func (e *Engine) EvaluateMoves(pos *Position, moves []EvaledMove, fromTrans Move) {
	for i := range moves {
		if moves[i].Move == fromTrans {
			moves[i].Value = 10000
		} else if moves[i].Move.IsPromotion() {
			moves[i].Value = 8000
		} else if moves[i].Move.IsCapture() {
			moves[i].Value = PieceValues[moves[i].Move.CapturedPiece()] - PieceValues[moves[i].Move.MovedPiece()] + 10000
		} else {
			moves[i].Value = e.EvalHistory[moves[i].Move.From()][moves[i].Move.To()]
		}
	}
}

func (e *Engine) quiescence(pos *Position, alpha, beta, height int, timedOut *bool) int {
	if *timedOut {
		return 0
	}

	var buffer [200]EvaledMove
	val := extensiveEval(pos, Evaluate(pos), height)

	if val >= beta {
		return beta
	}
	if alpha < val {
		alpha = val
	}

	var child Position
	var moveCount = 0

	if pos.IsInCheck() {
		evaled := pos.GenerateAllMoves(buffer[:])
		for i := range evaled {
			maxMoveToFirst(evaled[i:])
			if pos.MakeMove(evaled[i].Move, &child) {
				val = -extensiveEval(&child, Evaluate(&child), height+1)
				moveCount++
			}
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

	evaled := pos.GenerateAllCaptures(buffer[:])

	for i := range evaled {
		maxMoveToFirst(evaled[i:])
		if !pos.MakeMove(evaled[i].Move, &child) {
			continue
		}
		moveCount++
		val = -e.quiescence(&child, -beta, -alpha, height+1, timedOut)
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
func extensiveEval(pos *Position, evaledValue, height int) int {
	if !areAnyLegalMoves(pos) {
		if pos.IsInCheck() {
			return -Mate + height
		}
		return contempt(pos)
	}
	return evaledValue
}

var countPositions int

func (e *Engine) alphaBeta(pos *Position, depth, alpha, beta, height int, timedOut *bool) int {
	var child Position
	var val int
	if *timedOut {
		return 0
	}

	if e.isDraw(pos) {
		return contempt(pos)
	}

	val = MinInt

	var alphaOrig = alpha
	hashMove := NullMove
	ttEntry := e.TransTable.Get(pos.Key)
	if ttEntry.key == pos.Key {
		hashMove = ttEntry.bestMove
		if ttEntry.depth >= int32(depth) {
			val = ttEntry.value
			if val >= Mate-500 {
				val -= height + 1
			} else if val <= -Mate+500 {
				val += height - 1
			}
			if ttEntry.flag == TransExact {
				e.EvalHistory[uint(hashMove.From())][uint(hashMove.To())] += depth
				return val
			}
			if ttEntry.flag == TransAlpha && ttEntry.value <= alpha {
				return alpha
			}
			if ttEntry.flag == TransBeta && ttEntry.value >= beta {
				return beta
			}
		}
	}

	if depth == 0 {
		countPositions++
		val = e.quiescence(pos, alpha, beta, height, timedOut)
		return val
	}

	if pos.LastMove != NullMove && depth >= 4 && !pos.IsInCheck() && !isLateEndGame(pos) {
		pos.MakeNullMove(&child)
		val = -e.alphaBeta(&child, depth-3, -beta, -beta+1, height, timedOut)
		if val >= beta {
			return beta
		}
	}

	var buffer [256]EvaledMove

	evaled := pos.GenerateAllMoves(buffer[:])
	var tmpVal int
	e.EvaluateMoves(pos, evaled, hashMove)
	bestMove := NullMove
	moveCount := 0
	for i := range evaled {
		maxMoveToFirst(evaled[i:])
		if !pos.MakeMove(evaled[i].Move, &child) {
			continue
		}
		tmpVal = -e.alphaBeta(&child, depth-1, -beta, -alpha, height+1, timedOut)
		moveCount++
		if tmpVal > val {
			val = tmpVal
			bestMove = evaled[i].Move
		}
		if tmpVal > alpha {
			if !evaled[i].Move.IsCapture() {
				e.EvalHistory[uint(evaled[i].Move.From())][uint(evaled[i].Move.To())] += depth
			}
			alpha = tmpVal

			if alpha >= beta {
				e.TransTable.Set(depth, beta, TransBeta, pos.Key, evaled[i].Move)
				return beta
			}
		}
	}

	if moveCount == 0 {
		countPositions++
		if pos.IsInCheck() {
			val = -Mate + height
			e.TransTable.Set(depth, val, TransExact, pos.Key, NullMove)
			return val
		}
		val = contempt(pos)
		e.TransTable.Set(depth, val, TransExact, pos.Key, NullMove)
		return val
	}

	if alpha == alphaOrig {
		e.TransTable.Set(depth, alpha, TransAlpha, pos.Key, bestMove)
	} else {
		e.TransTable.Set(depth, val, TransExact, pos.Key, bestMove)
	}
	return alpha
}

func (e *Engine) isDraw(pos *Position) bool {
	if pos.FiftyMove > 100 {
		return true
	}

	if (pos.Pawns|pos.Rooks|pos.Queens) == 0 && !MoreThanOne(pos.Knights|pos.Bishops) {
		return true
	}

	if e.MoveHistory[pos.Key] >= 2 {
		return true
	}

	return false
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

func (e *Engine) depSearch(pos *Position, depth int, lastBestMove Move, resultChan chan result, timedOut *bool) {
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
			var val int
			if e.isDraw(&child) {
				val = contempt(pos)
			} else {
				val = -extensiveEval(&child, -evaled[i].value, 1)
			}
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
		val := -e.alphaBeta(&child, depth-1, -beta, -alpha, 1, timedOut)
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
		go e.depSearch(pos, i, lastBestMove, resultChan, &timedOut)
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
		go e.depSearch(pos, i, lastBestMove, resultChan, &timedOut)
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
		go e.depSearch(pos, i, lastBestMove, resultChan, &timedOut)
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
