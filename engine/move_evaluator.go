package engine

import . "github.com/mhib/combusken/backend"

const MinSpecialMoveValue = 7999

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
		} else if moves[i].Move.IsPromotion() {
			moves[i].Value = 70000 + SEEValues[moves[i].Move.CapturedPiece()]
		} else if moves[i].Move.IsCapture() {
			moves[i].Value = SEEValues[moves[i].Move.CapturedPiece()]*8 - SEEValues[moves[i].Move.MovedPiece()] + 10000
		} else {
			if moves[i].Move == mv.KillerMoves[height][0] {
				moves[i].Value = 9000
			} else if moves[i].Move == mv.KillerMoves[height][1] {
				moves[i].Value = 8000
			} else if moves[i].Move == counter {
				moves[i].Value = 7999
			} else {
				moves[i].Value = mv.EvalHistory[moves[i].Move.From()][moves[i].Move.To()]
			}
		}
	}
}

func (mv *MoveEvaluator) EvaluateQsMoves(moves []EvaledMove) {
	for i := range moves {
		moves[i].Value = PieceValues[moves[i].Move.CapturedPiece()] - PieceValues[moves[i].Move.MovedPiece()]
	}
}
