package engine

import (
	. "github.com/mhib/combusken/backend"
	. "github.com/mhib/combusken/evaluation"

	. "github.com/mhib/combusken/utils"
)

const MinSpecialMoveValue = 53000
const MaxBadCapture = -100000 + 4096 // 4096 represents max mvvlva value

const HistoryMax = 416
const HistoryMultiplier = 35
const HistoryDivisor = 706

type MoveHistory struct {
	KillerMoves      [STACK_SIZE + 1][2]Move
	CounterMoves     [2][64][64]Move
	ButterflyHistory [2][64][64]int32
	FollowUpHistory  [King + 1][64][King + 1][64]int32
	CounterHistory   [King + 1][64][King + 1][64]int32
	CurrentMove      [STACK_SIZE + 1]Move
}

func (mv *MoveHistory) ResetKillers(height int) {
	mv.KillerMoves[height][0] = NullMove
	mv.KillerMoves[height][1] = NullMove
}

func (mv *MoveHistory) CounterHistoryValue(lastMove Move, move Move) int32 {
	return mv.CounterHistory[lastMove.MovedPiece()][lastMove.To()][move.MovedPiece()][move.To()]
}

func (mv *MoveHistory) SetCurrentMove(height int, move Move) {
	mv.CurrentMove[height] = move
}

func (mv *MoveHistory) GetPreviousMove(height int) Move {
	return mv.CurrentMove[height-1]
}

func (mv *MoveHistory) GetPreviousMoveFromCurrentSide(height int) Move {
	return mv.CurrentMove[height-2]
}

func (mv *MoveHistory) Clear() {
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

func (mv *MoveHistory) Update(pos *Position, moves []Move, bestMove Move, depth, height int) {
	if pos.LastMove != NullMove {
		if mv.KillerMoves[height][0] != bestMove {
			mv.KillerMoves[height][0], mv.KillerMoves[height][1] = bestMove, mv.KillerMoves[height][0]
		}
		mv.CounterMoves[pos.SideToMove][pos.LastMove.From()][pos.LastMove.To()] = bestMove
	}
	unsignedBonus := int32(Min(depth*depth, HistoryMax))

	followUp := NullMove
	if height > 1 {
		followUp = mv.CurrentMove[height-2]
	}

	for _, move := range moves {
		var signedBonus int32
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

const MinGoodCapture = int32(55001)

func (mv *MoveHistory) EvaluateMoves(pos *Position, moves []EvaledMove, fromTrans Move, height, depth int) {
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

func (mv *MoveHistory) EvaluateQuiets(pos *Position, moves []EvaledMove, height int) {
	followUp := NullMove
	if height > 1 {
		followUp = mv.CurrentMove[height-2]
	}

	for i := range moves {
		moves[i].Value = mv.ButterflyHistory[pos.SideToMove][moves[i].From()][moves[i].To()]
		if pos.LastMove != NullMove {
			moves[i].Value += mv.CounterHistory[pos.LastMove.MovedPiece()][pos.LastMove.To()][moves[i].MovedPiece()][moves[i].To()]
		}
		if followUp != NullMove {
			moves[i].Value += mv.FollowUpHistory[followUp.MovedPiece()][followUp.To()][moves[i].MovedPiece()][moves[i].To()]
		}
	}
}
