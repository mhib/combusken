package engine

import (
	"context"

	. "github.com/mhib/combusken/backend"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1
const Mate = 1000000
const ValueWin = Mate - 150

func lossIn(height int) int {
	return -Mate + height
}

func areAnyLegalMoves(pos *Position, moves []EvaledMove) bool {
	var child Position
	for _, move := range moves {
		if pos.MakeMove(move.Move, &child) {
			return true
		}
	}
	return false
}

func depthToMate(val int) int {
	if val >= ValueWin {
		return Mate - val
	}
	return val - Mate
}

func (e *Engine) EvaluateMoves(moves []EvaledMove, fromTrans Move, height int) {
	var pos *Position = &e.Stack[height].position
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

func (e *Engine) EvaluateQsMoves(moves []EvaledMove) {
	for i := range moves {
		moves[i].Value = PieceValues[moves[i].Move.CapturedPiece()] - PieceValues[moves[i].Move.MovedPiece()]
	}
}

func (e *Engine) quiescence(alpha, beta, height int) int {
	e.incNodes()
	e.Stack[height].PV.clear()
	var pos = &e.Stack[height].position

	if height >= MAX_HEIGHT {
		return contempt(pos)
	}
	var val int

	var child = &e.Stack[height+1].position
	var moveCount = 0

	val = Evaluate(pos)

	if val >= beta {
		return beta
	}
	if alpha < val {
		alpha = val
	}

	evaled := pos.GenerateAllCaptures(e.Stack[height].moves[:])
	e.EvaluateQsMoves(evaled)

	for i := range evaled {
		maxMoveToFirst(evaled[i:])
		if !pos.MakeMove(evaled[i].Move, child) {
			continue
		}
		moveCount++
		val = -e.quiescence(-beta, -alpha, height+1)
		if val > alpha {
			alpha = val
			if val >= beta {
				return beta
			}
			e.Stack[height].PV.assign(evaled[i].Move, &e.Stack[height+1].PV)
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

func (e *Engine) alphaBeta(depth, alpha, beta, height int) int {
	e.incNodes()
	e.Stack[height].PV.clear()

	var pos = &e.Stack[height].position

	if e.isDraw(height) {
		return contempt(pos)
	}

	var tmpVal int

	var alphaOrig = alpha
	hashMove := NullMove
	ttEntry := e.TransTable.Get(pos.Key)
	if ttEntry.key == pos.Key {
		hashMove = ttEntry.bestMove
		tmpVal = valueFromTrans(ttEntry.value, height)
		if ttEntry.depth >= int32(depth) {
			if ttEntry.flag == TransExact {
				if !hashMove.IsCapture() {
					e.EvalHistory[uint(hashMove.From())][uint(hashMove.To())] += depth
				}
				return tmpVal
			}
			if ttEntry.flag == TransAlpha && tmpVal <= alpha {
				return alpha
			}
			if ttEntry.flag == TransBeta && tmpVal >= beta {
				return beta
			}
		}
	}

	if depth == 0 {
		return e.quiescence(alpha, beta, height)
	}

	var child = &e.Stack[height+1].position

	if pos.LastMove != NullMove && depth >= 4 && !pos.IsInCheck() && !isLateEndGame(pos) {
		pos.MakeNullMove(child)
		tmpVal = -e.alphaBeta(depth-3, -beta, -beta+1, height+1)
		if tmpVal >= beta {
			return beta
		}
	}

	var val = MinInt

	evaled := pos.GenerateAllMoves(e.Stack[height].moves[:])
	e.EvaluateMoves(evaled, hashMove, height)
	bestMove := NullMove
	moveCount := 0
	for i := range evaled {
		maxMoveToFirst(evaled[i:])
		if !pos.MakeMove(evaled[i].Move, child) {
			continue
		}
		tmpVal = -e.alphaBeta(depth-1, -beta, -alpha, height+1)
		moveCount++

		if tmpVal > val {
			val = tmpVal
			bestMove = evaled[i].Move
			if val > alpha {
				alpha = val

				// Maybe move this out of loop?
				if !evaled[i].Move.IsCapture() {
					e.EvalHistory[uint(evaled[i].Move.From())][uint(evaled[i].Move.To())] += depth
				}

				if alpha >= beta {
					if !evaled[i].Move.IsCapture() && pos.LastMove != NullMove {
						e.KillerMoves[height][0], e.KillerMoves[height][1] = evaled[i].Move, e.KillerMoves[height][0]
						e.CounterMoves[pos.LastMove.From()][pos.LastMove.To()] = evaled[i].Move
					}
					e.TransTable.Set(depth, beta, TransBeta, pos.Key, evaled[i].Move, height)
					return beta
				}
			}
			e.Stack[height].PV.assign(evaled[i].Move, &e.Stack[height+1].PV)
		}
	}

	if moveCount == 0 {
		if pos.IsInCheck() {
			return lossIn(height)
		}
		return contempt(pos)
	}

	if alpha == alphaOrig {
		e.TransTable.Set(depth, alpha, TransAlpha, pos.Key, bestMove, height)
	} else {
		e.TransTable.Set(depth, alpha, TransExact, pos.Key, bestMove, height)
	}
	return alpha
}

func (e *Engine) isDraw(height int) bool {
	var pos *Position = &e.Stack[height].position
	if pos.FiftyMove > 100 {
		return true
	}

	if (pos.Pawns|pos.Rooks|pos.Queens) == 0 && !MoreThanOne(pos.Knights|pos.Bishops) {
		return true
	}

	for i := height - 1; i >= 0; i-- {
		desc := &e.Stack[i].position
		if desc.Key == pos.Key {
			return true
		}
		if desc.FiftyMove == 0 || desc.LastMove == NullMove {
			return false
		}
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

func (e *Engine) depSearch(depth int, lastBestMove Move, resultChan chan result) {
	var pos = &e.Stack[0].position
	defer recoverFromTimeout()
	e.Nodes = 0
	var child = &e.Stack[1].position
	var bestMove = NullMove
	evaled := pos.GenerateAllMoves(e.Stack[0].moves[:])
	e.EvaluateMoves(evaled, lastBestMove, 0)
	var alpha = -MaxInt
	if depth == 1 {
		e.Stack[0].PV.clear()
		var val int
		for i := range evaled {
			maxMoveToFirst(evaled[i:])
			if !pos.MakeMove(evaled[i].Move, child) {
				continue
			}
			if e.isDraw(1) {
				val = contempt(pos)
			} else {
				val = -Evaluate(child)
			}
			if val > alpha {
				alpha = val
				bestMove = evaled[i].Move
			}
		}
		e.Stack[0].PV.assign(bestMove, &e.Stack[1].PV)
		e.TransTable.Set(depth, alpha, TransExact, pos.Key, bestMove, 1)
		resultChan <- result{bestMove, alpha}
		return
	}
	for i := range evaled {
		maxMoveToFirst(evaled[i:])
		if !pos.MakeMove(evaled[i].Move, child) {
			continue
		}
		val := -e.alphaBeta(depth-1, -MaxInt, -alpha, 1)
		if val > alpha {
			alpha = val
			bestMove = evaled[i].Move
			e.Stack[0].PV.assign(evaled[i].Move, &e.Stack[1].PV)
		}
	}
	e.TransTable.Set(depth, alpha, TransExact, pos.Key, bestMove, 0)
	resultChan <- result{bestMove, alpha}
}

func (e *Engine) TimeSearch(ctx context.Context, pos *Position) Move {
	var lastBestMove Move
	for i := 1; ; i++ {
		resultChan := make(chan result, 1)
		e.Stack[0].position = *pos
		go e.depSearch(i, lastBestMove, resultChan)
		select {
		case <-ctx.Done():
			e.timedOut <- true
			return lastBestMove
		case res := <-resultChan:
			if i >= 3 {
				e.callUpdate(SearchInfo{res.int, i, e.Nodes, e.Stack[0].PV})
			}
			if res.int >= ValueWin && depthToMate(res.int) <= i {
				return res.Move
			}
			if res.Move == 0 {
				return lastBestMove
			}
			lastBestMove = res.Move
			if i >= MAX_HEIGHT {
				return res.Move
			}
		}
	}
}

func (e *Engine) DepthSearch(pos *Position, depth int) Move {
	resultChan := make(chan result, 1)
	e.Stack[0].position = *pos
	go e.depSearch(depth, NullMove, resultChan)
	res := <-resultChan
	e.callUpdate(SearchInfo{res.int, depth, e.Nodes, e.Stack[0].PV})
	return res.Move
}

func (e *Engine) CountSearch(ctx context.Context, pos *Position, count int) Move {
	var lastBestMove Move
	for i := 1; ; i++ {
		resultChan := make(chan result, 1)
		e.Stack[0].position = *pos
		go e.depSearch(i, lastBestMove, resultChan)
		res := <-resultChan
		e.callUpdate(SearchInfo{res.int, i, e.Nodes, e.Stack[0].PV})
		if res.int >= ValueWin && depthToMate(res.int) <= i {
			return res.Move
		}
		if res.Move == 0 {
			return lastBestMove
		}
		lastBestMove = res.Move
		if i > 70 || e.Nodes >= count {
			return lastBestMove
		}
	}
}

func recoverFromTimeout() {
	var err = recover()
	if err != nil && err != errTimeout {
		panic(err)
	}
}
