package engine

import . "github.com/mhib/combusken/backend"
import . "github.com/mhib/combusken/evaluation"

const MinSpecialMoveValue = 1499
const MaxBadCapture = -100000 + 2500 // 2500 represents max mvvlva value

const HistoryMax = 400
const HistoryMultiplier = 32
const HistoryDivisor = 512

type MoveEvaluator struct {
	KillerMoves  [STACK_SIZE][2]Move
	CounterMoves [2][64][64]Move
	EvalHistory  [2][64][64]int
}

func (mv *MoveEvaluator) ResetKillers(height int) {
	mv.KillerMoves[height+1][0] = NullMove
	mv.KillerMoves[height+1][1] = NullMove
}

func (mv *MoveEvaluator) Clear() {
	for side := 0; side < 2; side++ {
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				mv.EvalHistory[side][y][x] = 0
			}
		}
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				mv.CounterMoves[side][y][x] = NullMove
			}
		}

	}
	for y := 0; y < MAX_HEIGHT; y++ {
		for x := 0; x < 2; x++ {
			mv.KillerMoves[y][x] = NullMove
		}
	}
}

func (mv *MoveEvaluator) Update(pos *Position, moves []Move, bestMove Move, depth, height int) {
	side := pos.IntSide()
	if pos.LastMove != NullMove {
		if mv.KillerMoves[height][0] != bestMove {
			mv.KillerMoves[height][0], mv.KillerMoves[height][1] = bestMove, mv.KillerMoves[height][0]
		}
		mv.CounterMoves[side][pos.LastMove.From()][pos.LastMove.To()] = bestMove
	}
	bonus := min(depth*depth, HistoryMax)

	for _, move := range moves {
		if move == bestMove {
			entry := mv.EvalHistory[side][move.From()][move.To()]
			entry += HistoryMultiplier*bonus - entry*bonus/HistoryDivisor
			mv.EvalHistory[side][move.From()][move.To()] = entry
			break
		} else {
			entry := mv.EvalHistory[side][move.From()][move.To()]
			entry += HistoryMultiplier*-bonus - entry*bonus/HistoryDivisor
			mv.EvalHistory[side][move.From()][move.To()] = entry
		}
	}
}

func (mv *MoveEvaluator) EvaluateMoves(pos *Position, moves []EvaledMove, fromTrans Move, height, depth int) {
	side := pos.IntSide()
	var counter Move
	if pos.LastMove != NullMove {
		counter = mv.CounterMoves[side][pos.LastMove.From()][pos.LastMove.To()]
	}
	for i := range moves {
		if moves[i].Move == fromTrans {
			moves[i].Value = 100000
		} else if moves[i].Move.IsCaptureOrPromotion() {
			if SeeSign(pos, moves[i].Move) {
				moves[i].Value = mvvlva(moves[i].Move) + 50000
			} else {
				moves[i].Value = mvvlva(moves[i].Move) - 100000
			}
		} else {
			if moves[i].Move == mv.KillerMoves[height][0] {
				moves[i].Value = 20000
			} else if moves[i].Move == mv.KillerMoves[height][1] {
				moves[i].Value = 15000
			} else if moves[i].Move == counter {
				moves[i].Value = 14999
			} else {
				moves[i].Value = mv.EvalHistory[side][moves[i].Move.From()][moves[i].Move.To()]
			}
		}
	}
}

var mvvlvaScores = [...]int{0, 10, 40, 45, 68, 145, 256}

func mvvlva(move Move) int {
	captureScore := mvvlvaScores[move.CapturedPiece()]
	if move.IsPromotion() {
		captureScore += mvvlvaScores[move.PromotedPiece()] - mvvlvaScores[Pawn]
	}
	return captureScore*8 - mvvlvaScores[move.MovedPiece()]
}

// In Quiescent search it is expected that SEE will be check anyway
func (mv *MoveEvaluator) EvaluateQsMoves(pos *Position, moves []EvaledMove, inCheck bool) {
	side := pos.IntSide()
	if inCheck {
		for i := range moves {
			if moves[i].Move.IsCaptureOrPromotion() {
				moves[i].Value = mvvlva(moves[i].Move) + 50000
			} else {
				moves[i].Value = mv.EvalHistory[side][moves[i].Move.From()][moves[i].Move.To()]
			}
		}
	} else {
		for i := range moves {
			moves[i].Value = mvvlva(moves[i].Move)
		}
	}
}
