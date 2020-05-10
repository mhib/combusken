package engine

import (
	. "github.com/mhib/combusken/backend"
	"github.com/mhib/combusken/evaluation"
)

const (
	TT_MOVE uint8 = iota
	GENERATE_NOISY
	GOOD_NOISY
	KILLER_1
	KILLER_2
	COUNTER
	GENERATE_QUIET
	QUIET
	BAD_NOISY
	DONE
)

type MoveProvider struct {
	Moves      [MAX_MOVES]EvaledMove
	ttMove     Move
	counter    Move
	killer1    Move
	killer2    Move
	kind       uint8
	stage      uint8
	split      uint8
	noisySize  uint8
	quietsSize uint8
}

const (
	NORMAL uint8 = iota
	NOISY
)

var mvvlvaScores = [None + 1]int32{10, 40, 45, 68, 145, 256, 0}

const badNoisyValue = -4096

func mvvlva(move Move) int32 {
	captureScore := mvvlvaScores[move.CapturedPiece()]
	if move.IsPromotion() && move.PromotedPiece() == Queen {
		captureScore += mvvlvaScores[Queen] - mvvlvaScores[Pawn]
	}
	return captureScore*8 - mvvlvaScores[move.MovedPiece()]
}

func (mp *MoveProvider) InitQs(pos *Position) {
	mp.kind = NOISY
	mp.ttMove = NullMove
	mp.stage = GENERATE_NOISY
}

func (mp *MoveProvider) InitNormal(pos *Position, mh *MoveHistory, height int, ttMove Move) {
	mp.kind = NORMAL
	mp.stage = TT_MOVE
	mp.ttMove = ttMove
	if pos.LastMove != NullMove {
		mp.counter = mh.CounterMoves[pos.SideToMove][pos.LastMove.From()][pos.LastMove.To()]
	}
	mp.killer1 = mh.KillerMoves[height][0]
	mp.killer2 = mh.KillerMoves[height][1]
}

func (mp *MoveProvider) InitSingular() {
	mp.stage = GENERATE_NOISY
}

func (mp *MoveProvider) RestoreFromSingular() {
	mp.stage = GOOD_NOISY
	mp.noisySize = mp.split
}

// In Quiescent search it is expected that SEE will be checked anyway
func evaluateNoisy(Moves []EvaledMove) {
	for i := range Moves {
		Moves[i].Value = mvvlva(Moves[i].Move)
	}
}

func (mp *MoveProvider) GetMoveStage() uint8 {
	return mp.stage
}

func (mp *MoveProvider) dropNoisy(bestIdx int) {
	mp.noisySize--
	mp.Moves[mp.noisySize], mp.Moves[bestIdx] = mp.Moves[bestIdx], mp.Moves[mp.noisySize]
}

func (mp *MoveProvider) GetNextMove(pos *Position, mh *MoveHistory, height int) Move {
	var move EvaledMove
	var bestIdx int
	switch mp.stage {
	case TT_MOVE:
		mp.stage++
		if mp.ttMove != NullMove && pos.IsMovePseudoLegal(mp.ttMove) {
			return mp.ttMove
		}
		fallthrough
	case GENERATE_NOISY:
		mp.stage++
		mp.noisySize = pos.GenerateNoisy(mp.Moves[:])
		mp.split = mp.noisySize
		evaluateNoisy(mp.Moves[:mp.noisySize])
		fallthrough
	case GOOD_NOISY:
		for mp.noisySize > 0 {
			bestIdx = 0
			for i := 1; i < int(mp.noisySize); i++ {
				if mp.Moves[i].Value > mp.Moves[bestIdx].Value {
					bestIdx = i
				}
			}
			move = mp.Moves[bestIdx]
			if move.Value == badNoisyValue {
				break
			}
			if move.Move == mp.ttMove {
				mp.dropNoisy(bestIdx)
				continue
			}
			if !evaluation.SeeSign(pos, move.Move) {
				mp.Moves[bestIdx].Value = badNoisyValue
				continue
			}
			mp.dropNoisy(bestIdx)
			return move.Move
		}
		if mp.kind == NOISY {
			mp.stage = DONE
			return NullMove
		}
		mp.stage++
		fallthrough
	case KILLER_1:
		mp.stage++
		if mp.killer1 != NullMove && mp.killer1 != mp.ttMove && pos.IsMovePseudoLegal(mp.killer1) {
			return mp.killer1
		}
		fallthrough
	case KILLER_2:
		mp.stage++
		if mp.killer2 != NullMove && mp.killer2 != mp.ttMove && pos.IsMovePseudoLegal(mp.killer2) {
			return mp.killer2
		}
		fallthrough
	case COUNTER:
		mp.stage++
		if mp.counter != NullMove && mp.counter != mp.ttMove && mp.counter != mp.killer1 && mp.counter != mp.killer2 && pos.IsMovePseudoLegal(mp.counter) {
			return mp.counter
		}
		fallthrough
	case GENERATE_QUIET:
		mp.stage++
		mp.quietsSize = pos.GenerateQuiet(mp.Moves[mp.split:])
		mh.EvaluateQuiets(pos, mp.Moves[mp.split:mp.split+mp.quietsSize], height)
		// Insertion sort
		for i := mp.split + 1; i < mp.split+mp.quietsSize; i++ {
			j, t := i, mp.Moves[i]
			for ; j >= mp.split+1 && mp.Moves[j-1].Value > t.Value; j -= 1 {
				mp.Moves[j] = mp.Moves[j-1]
			}
			mp.Moves[j] = t
		}
		fallthrough
	case QUIET:
		for mp.quietsSize > 0 {
			mp.quietsSize--
			move = mp.Moves[mp.split+mp.quietsSize]
			if move.Move != mp.ttMove && move.Move != mp.killer1 && move.Move != mp.killer2 && move.Move != mp.counter {
				return move.Move
			}
		}
		mp.stage++
		fallthrough
	case BAD_NOISY:
		for mp.noisySize > 0 {
			mp.noisySize--
			move = mp.Moves[mp.noisySize]
			if move.Move != mp.ttMove {
				return move.Move
			}
		}
		mp.stage++
		fallthrough
	default:
		return NullMove
	}
}
