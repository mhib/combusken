package engine

import (
	"context"

	. "github.com/mhib/combusken/backend"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1
const Mate = 1000000

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

func depthToMate(val int) int {
	if val <= -Mate+500 {
		return val - Mate
	}
	return Mate - val
}

func (e *Engine) EvaluateMoves(pos *Position, moves []EvaledMove, fromTrans Move, height int) {
	var counter Move
	if pos.LastMove != NullMove {
		counter = e.CounterMoves[pos.LastMove.From()][pos.LastMove.To()]
	}
	for i := range moves {
		if moves[i].Move == fromTrans {
			moves[i].Value = 100000
		} else if moves[i].Move.IsPromotion() {
			moves[i].Value = 70000 + SEEValues[moves[i].Move.CapturedPiece()]
		} else if moves[i].Move.IsCapture() {
			moves[i].Value = SEEValues[moves[i].Move.CapturedPiece()]*8 - SEEValues[moves[i].Move.MovedPiece()] + 10000
		} else {
			if moves[i].Move == e.KillerMoves[height][0] {
				moves[i].Value = 9000
			} else if moves[i].Move == e.KillerMoves[height][1] {
				moves[i].Value = 8000
			} else if moves[i].Move == counter {
				moves[i].Value = 7999
			} else {
				moves[i].Value = e.EvalHistory[moves[i].Move.From()][moves[i].Move.To()]
			}
		}
	}
}

func (e *Engine) EvaluateQsMoves(pos *Position, moves []EvaledMove) {
	for i := range moves {
		moves[i].Value = PieceValues[moves[i].Move.CapturedPiece()] - PieceValues[moves[i].Move.MovedPiece()]
	}
}

func (e *Engine) quiescence(pos *Position, alpha, beta, height int, timedOut *bool) int {
	if *timedOut {
		return 0
	}

	if e.isDraw(pos) {
		return contempt(pos)
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
		if val <= -Mate+500 {
			return val
		}
		evaled := pos.GenerateAllMoves(buffer[:])
		for i := range evaled {
			if pos.MakeMove(evaled[i].Move, &child) {
				val = -extensiveEval(&child, Evaluate(&child), height+1)
			}
			if val > alpha {
				alpha = val
				if alpha >= beta {
					return beta
				}
			}
		}
		return val
	}

	evaled := pos.GenerateAllCaptures(buffer[:])
	e.EvaluateQsMoves(pos, evaled)

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
		val = valueFromTrans(ttEntry.value, height)
		if ttEntry.depth >= int32(depth) {
			if ttEntry.flag == TransExact {
				if !hashMove.IsCapture() {
					e.EvalHistory[uint(hashMove.From())][uint(hashMove.To())] += depth
				}
				return val
			}
			if ttEntry.flag == TransAlpha && val <= alpha {
				return alpha
			}
			if ttEntry.flag == TransBeta && val >= beta {
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
		val = -e.alphaBeta(&child, depth-3, -beta, -beta+1, height+1, timedOut)
		if val >= beta {
			return beta
		}
	}

	var buffer [256]EvaledMove

	evaled := pos.GenerateAllMoves(buffer[:])
	var tmpVal int
	e.EvaluateMoves(pos, evaled, hashMove, height)
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
			alpha = tmpVal
			if alpha >= beta {
				if !evaled[i].Move.IsCapture() {
					e.KillerMoves[height][0], e.KillerMoves[height][1] = evaled[i].Move, e.KillerMoves[height][0]
					if pos.LastMove != NullMove {
						e.CounterMoves[pos.LastMove.From()][pos.LastMove.To()] = evaled[i].Move
					}
				}
				e.TransTable.Set(depth, beta, TransBeta, pos.Key, evaled[i].Move, height)
				return beta
			}

			if !evaled[i].Move.IsCapture() {
				e.EvalHistory[uint(evaled[i].Move.From())][uint(evaled[i].Move.To())] += depth
			}
		}
	}

	if moveCount == 0 {
		countPositions++
		if pos.IsInCheck() {
			val = -Mate + height
			return val
		}
		val = contempt(pos)
		return val
	}

	if alpha == alphaOrig {
		e.TransTable.Set(depth, alpha, TransAlpha, pos.Key, bestMove, height)
	} else {
		e.TransTable.Set(depth, alpha, TransExact, pos.Key, bestMove, height)
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

type result struct {
	Move
	int
}

func (e *Engine) depSearch(pos *Position, depth int, lastBestMove Move, resultChan chan result, timedOut *bool) {
	var buffer [256]EvaledMove
	var child Position
	var bestMove = NullMove
	evaled := pos.GenerateAllMoves(buffer[:])
	e.EvaluateMoves(pos, evaled, lastBestMove, 0)
	var alpha = -MaxInt
	countPositions = 0
	if depth == 1 {
		var val int
		for i := range evaled {
			maxMoveToFirst(evaled[i:])
			if !pos.MakeMove(evaled[i].Move, &child) {
				continue
			}
			if e.isDraw(&child) {
				val = contempt(pos)
			} else {
				val = -extensiveEval(&child, Evaluate(&child), 1)
			}
			if val > alpha {
				alpha = val
				bestMove = evaled[i].Move
			}
		}
		e.TransTable.Set(depth, alpha, TransExact, pos.Key, bestMove, 1)
		resultChan <- result{bestMove, alpha}
		return
	}
	for i := range evaled {
		maxMoveToFirst(evaled[i:])
		if !pos.MakeMove(evaled[i].Move, &child) {
			continue
		}
		val := -e.alphaBeta(&child, depth-1, -MaxInt, -alpha, 1, timedOut)
		if val > alpha {
			alpha = val
			bestMove = evaled[i].Move
		}
	}
	e.TransTable.Set(depth, alpha, TransExact, pos.Key, bestMove, 0)
	resultChan <- result{bestMove, alpha}
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
			e.callUpdate(SearchInfo{res.int, i})
			if res.int > Mate-500 && depthToMate(res.int) <= i {
				return res.Move
			}
			if res.Move == 0 {
				return lastBestMove
			}
			lastBestMove = res.Move
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
		e.callUpdate(SearchInfo{res.int, i})
		if res.int > Mate-500 && depthToMate(res.int) <= i {
			return res.Move
		}
		if res.Move == 0 {
			return lastBestMove
		}
		lastBestMove = res.Move
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
		e.callUpdate(SearchInfo{res.int, i})
		if res.int > Mate-500 && depthToMate(res.int) <= i {
			return res.Move
		}
		if res.Move == 0 {
			return lastBestMove
		}
		lastBestMove = res.Move
		if i > 70 || countPositions > count {
			return lastBestMove
		}
	}
}
