package engine

import . "github.com/mhib/combusken/backend"

const MinSpecialMoveValue = 1499

type MoveEvaluator struct {
	KillerMoves  [STACK_SIZE][2]Move
	CounterMoves [64][64]Move
	EvalHistory  [64][64]int
}

func (mv *MoveEvaluator) Clear() {
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			mv.EvalHistory[y][x] = 0
		}
	}
	for y := 0; y < MAX_HEIGHT; y++ {
		for x := 0; x < 2; x++ {
			mv.KillerMoves[y][x] = NullMove
		}
	}
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			mv.CounterMoves[y][x] = NullMove
		}
	}
}

func (mv *MoveEvaluator) EvaluateMoves(pos *Position, moves []EvaledMove, fromTrans Move, height int) {
	var counter Move
	if pos.LastMove != NullMove {
		counter = mv.CounterMoves[pos.LastMove.From()][pos.LastMove.To()]
	}
	for i := range moves {
		if moves[i].Move == fromTrans {
			moves[i].Value = 100000
		} else if moves[i].Move.IsCaptureOrPromotion() {
			if seeSign(pos, moves[i].Move) {
				moves[i].Value = mvvlva(moves[i].Move) + 50000
			} else {
				moves[i].Value = mv.EvalHistory[moves[i].Move.From()][moves[i].Move.To()]
			}
		} else {
			if moves[i].Move == mv.KillerMoves[height][0] {
				moves[i].Value = 20000
			} else if moves[i].Move == mv.KillerMoves[height][1] {
				moves[i].Value = 15000
			} else if moves[i].Move == counter {
				moves[i].Value = 14999
			} else {
				moves[i].Value = mv.EvalHistory[moves[i].Move.From()][moves[i].Move.To()]
			}
		}
	}
}

var mvvlvaScores = [...]int{0, 1, 4, 4, 6, 12, 20}

func mvvlva(move Move) int {
	captureScore := mvvlvaScores[move.CapturedPiece()]
	if move.IsPromotion() {
		captureScore += mvvlvaScores[move.PromotedPiece()] - mvvlvaScores[Pawn]
	}
	return captureScore*8 - mvvlvaScores[move.MovedPiece()]
}

func (mv *MoveEvaluator) EvaluateQsMoves(pos *Position, moves []EvaledMove, inCheck bool) {
	if inCheck {
		for i := range moves {
			if moves[i].Move.IsCaptureOrPromotion() {
				moves[i].Value = mvvlva(moves[i].Move) + 50000
			} else {
				moves[i].Value = mv.EvalHistory[moves[i].Move.From()][moves[i].Move.To()]
			}
		}
	} else {
		for i := range moves {
			moves[i].Value = mvvlva(moves[i].Move)
		}
	}
}
