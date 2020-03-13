package engine

import . "github.com/mhib/combusken/backend"
import . "github.com/mhib/combusken/evaluation"
import . "github.com/mhib/combusken/utils"

const MinSpecialMoveValue = 53000
const MaxBadCapture = -100000 + 4096 // 4096 represents max mvvlva value

const HistoryMax = 400
const HistoryMultiplier = 32
const HistoryDivisor = 512

type MoveEvaluator struct {
	KillerMoves      [STACK_SIZE + 1][2]Move
	CounterMoves     [2][64][64]Move
	ButterflyHistory [2][64][64]int
	FollowUpHistory  [King + 1][64][King + 1][64]int
	CounterHistory   [King + 1][64][King + 1][64]int
	CurrentMove      [STACK_SIZE + 1]Move
}

func (mv *MoveEvaluator) ResetKillers(height int) {
	mv.KillerMoves[height][0] = NullMove
	mv.KillerMoves[height][1] = NullMove
}

func (mv *MoveEvaluator) SetCurrentMove(height int, move Move) {
	mv.CurrentMove[height] = move
}

func (mv *MoveEvaluator) GetPreviousMove(height int) Move {
	return mv.CurrentMove[height-1]
}

func (mv *MoveEvaluator) Clear() {
	for side := 0; side < 2; side++ {
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				mv.ButterflyHistory[side][y][x] = 0
			}
		}
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				mv.CounterMoves[side][y][x] = NullMove
			}
		}

	}
	for y := 0; y <= STACK_SIZE; y++ {
		for x := 0; x < 2; x++ {
			mv.KillerMoves[y][x] = NullMove
		}
		mv.CurrentMove[y] = NullMove
	}

	for a := Pawn; a <= King; a++ {
		for b := 0; b < 64; b++ {
			for c := Pawn; c <= King; c++ {
				for d := 0; d <= 64; d++ {
					mv.FollowUpHistory[a][b][c][d] = 0
					mv.CounterHistory[a][b][c][d] = 0
				}
			}
		}
	}
}

func (mv *MoveEvaluator) Update(pos *Position, moves []Move, bestMove Move, depth, height int, failHigh bool) {
	unsignedBonus := Min(depth*depth, HistoryMax)

	if !failHigh {
		if depth < 4 {
			return
		}
		unsignedBonus /= 4
	} else if pos.LastMove != NullMove {
		if mv.KillerMoves[height][0] != bestMove {
			mv.KillerMoves[height][0], mv.KillerMoves[height][1] = bestMove, mv.KillerMoves[height][0]
		}
		mv.CounterMoves[pos.SideToMove][pos.LastMove.From()][pos.LastMove.To()] = bestMove
	}

	followUp := NullMove
	if height > 1 {
		followUp = mv.CurrentMove[height-2]
	}

	for _, move := range moves {
		var signedBonus int
		if move == bestMove {
			signedBonus = unsignedBonus
		} else {
			signedBonus = -unsignedBonus
		}
		entry := mv.ButterflyHistory[pos.SideToMove][move.From()][move.To()]
		entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
		mv.ButterflyHistory[pos.SideToMove][move.From()][move.To()] = entry

		if pos.LastMove != NullMove {
			entry = mv.CounterHistory[pos.LastMove.MovedPiece()][pos.LastMove.To()][move.MovedPiece()][move.To()]
			entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
			mv.CounterHistory[pos.LastMove.MovedPiece()][pos.LastMove.To()][move.MovedPiece()][move.To()] = entry
		}
		if followUp != NullMove {
			entry = mv.FollowUpHistory[followUp.MovedPiece()][followUp.To()][move.MovedPiece()][move.To()]
			entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
			mv.FollowUpHistory[followUp.MovedPiece()][followUp.To()][move.MovedPiece()][move.To()] = entry
		}
	}
}

const MinGoodCapture = 55001

func (mv *MoveEvaluator) EvaluateMoves(pos *Position, moves []EvaledMove, fromTrans Move, height, depth int) {
	var counter Move
	if pos.LastMove != NullMove {
		counter = mv.CounterMoves[pos.SideToMove][pos.LastMove.From()][pos.LastMove.To()]
	}

	followUp := NullMove
	if height > 1 {
		followUp = mv.CurrentMove[height-2]
	}

	for i := range moves {
		if moves[i].Move == fromTrans {
			moves[i].Value = 120000
		} else if moves[i].Move.IsCaptureOrPromotion() {
			if SeeSign(pos, moves[i].Move) {
				moves[i].Value = mvvlva(moves[i].Move) + 100000
			} else {
				moves[i].Value = mvvlva(moves[i].Move) - 100000
			}
		} else {
			if moves[i].Move == mv.KillerMoves[height][0] {
				moves[i].Value = 55000
			} else if moves[i].Move == mv.KillerMoves[height][1] {
				moves[i].Value = 54000
			} else if moves[i].Move == counter {
				moves[i].Value = 53000
			} else {
				moves[i].Value = mv.ButterflyHistory[pos.SideToMove][moves[i].From()][moves[i].To()]
				if pos.LastMove != NullMove {
					moves[i].Value += mv.CounterHistory[pos.LastMove.MovedPiece()][pos.LastMove.To()][moves[i].MovedPiece()][moves[i].To()]
				}
				if followUp != NullMove {
					moves[i].Value += mv.FollowUpHistory[followUp.MovedPiece()][followUp.To()][moves[i].MovedPiece()][moves[i].To()]
				}
			}
		}
	}
}

var mvvlvaScores = [None + 1]int{10, 40, 45, 68, 145, 256, 0}

func mvvlva(move Move) int {
	captureScore := mvvlvaScores[move.CapturedPiece()]
	if move.IsPromotion() {
		captureScore += mvvlvaScores[move.PromotedPiece()] - mvvlvaScores[Pawn]
	}
	return captureScore*8 - mvvlvaScores[move.MovedPiece()]
}

// In Quiescent search it is expected that SEE will be checked anyway
func (mv *MoveEvaluator) EvaluateQsMoves(pos *Position, moves []EvaledMove, bestMove Move, inCheck bool) {
	if inCheck {
		for i := range moves {
			if moves[i].Move == bestMove {
				moves[i].Value = 100000
			} else if moves[i].Move.IsCaptureOrPromotion() {
				moves[i].Value = mvvlva(moves[i].Move) + 50000
			} else {
				moves[i].Value = mv.ButterflyHistory[pos.SideToMove][moves[i].Move.From()][moves[i].Move.To()]
			}
		}
	} else {
		for i := range moves {
			if moves[i].Move == bestMove {
				moves[i].Value = 100000
			} else {
				moves[i].Value = mvvlva(moves[i].Move)
			}
		}
	}
}
