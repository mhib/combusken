package engine

import (
	. "github.com/mhib/combusken/backend"
	"github.com/mhib/combusken/evaluation"
)

const (
	StageTtMove uint8 = iota
	StageGenerateNoisy
	StageGoodNoisy
	StageKiller1
	StageKiller2
	StageCounter
	StageGenerateQuiet
	StageQuiet
	StageNoisy
	StageDone
)

type MoveProvider struct {
	Moves      [MaxMoves]EvaledMove
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

func (mp *MoveProvider) InitQs() {
	mp.kind = NOISY
	mp.ttMove = NullMove
	mp.stage = StageGenerateNoisy
}

func (mp *MoveProvider) InitNormal(pos *Position, mh *MoveHistory, height int, ttMove Move) {
	mp.kind = NORMAL
	mp.stage = StageTtMove
	mp.ttMove = ttMove
	if pos.LastMove != NullMove {
		mp.counter = mh.CounterMoves[pos.SideToMove][pos.LastMove.From()][pos.LastMove.To()]
	}
	mp.killer1 = mh.KillerMoves[height][0]
	mp.killer2 = mh.KillerMoves[height][1]
}

func (mp *MoveProvider) InitSingular() {
	mp.stage = StageGenerateNoisy
}

// We can save generating noisy after quiescence
func (mp *MoveProvider) RestoreFromSingular() {
	mp.stage = StageGoodNoisy
	mp.noisySize = mp.split
}

// In Quiescent search it is expected that SEE will be checked anyway
func evaluateNoisy(Moves []EvaledMove) {
	for i := range Moves {
		Moves[i].Value = mvvlva(Moves[i].Move)
	}
}

func (mp *MoveProvider) GetStage() uint8 {
	return mp.stage
}

func (mp *MoveProvider) dropNoisy(bestIdx int) {
	mp.noisySize--
	mp.Moves[mp.noisySize], mp.Moves[bestIdx] = mp.Moves[bestIdx], mp.Moves[mp.noisySize]
}

func (mp *MoveProvider) skipQuiets() {
	mp.stage = StageNoisy
}

func (mp *MoveProvider) GetNextMove(pos *Position, mh *MoveHistory, depth, height int) Move {
	var move EvaledMove
	var bestIdx int
	switch mp.stage {
	case StageTtMove:
		mp.stage++
		if mp.ttMove != NullMove && pos.IsMovePseudoLegal(mp.ttMove) {
			return mp.ttMove
		}
		fallthrough
	case StageGenerateNoisy:
		mp.stage++
		mp.noisySize = GenerateNoisy(pos, mp.Moves[:])
		mp.split = mp.noisySize
		evaluateNoisy(mp.Moves[:mp.noisySize])
		fallthrough
	case StageGoodNoisy:
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
			mp.stage = StageDone
			return NullMove
		}
		mp.stage++
		fallthrough
	case StageKiller1:
		mp.stage++
		if mp.killer1 != NullMove && mp.killer1 != mp.ttMove && pos.IsMovePseudoLegal(mp.killer1) {
			return mp.killer1
		}
		fallthrough
	case StageKiller2:
		mp.stage++
		if mp.killer2 != NullMove && mp.killer2 != mp.ttMove && pos.IsMovePseudoLegal(mp.killer2) {
			return mp.killer2
		}
		fallthrough
	case StageCounter:
		mp.stage++
		if mp.counter != NullMove && mp.counter != mp.ttMove && mp.counter != mp.killer1 && mp.counter != mp.killer2 && pos.IsMovePseudoLegal(mp.counter) {
			return mp.counter
		}
		fallthrough
	case StageGenerateQuiet:
		mp.stage++
		mp.quietsSize = GenerateQuiet(pos, mp.Moves[mp.split:])
		quietMoves := mp.Moves[mp.split : mp.split+mp.quietsSize]
		mh.EvaluateQuiets(pos, quietMoves, height)
		sortTreshold := -2000 * int32(depth)
		// Partial Insertion sort
		for i := len(quietMoves) - 2; i >= 0; i-- {
			if quietMoves[i].Value > sortTreshold {
				j, t := i, quietMoves[i]
				for ; j <= len(quietMoves)-2 && quietMoves[j+1].Value < t.Value; j++ {
					quietMoves[j] = quietMoves[j+1]
				}
				quietMoves[j] = t
			}
		}
		fallthrough
	case StageQuiet:
		for mp.quietsSize > 0 {
			mp.quietsSize--
			move = mp.Moves[mp.split+mp.quietsSize]
			if move.Move != mp.ttMove && move.Move != mp.killer1 && move.Move != mp.killer2 && move.Move != mp.counter {
				return move.Move
			}
		}

		mp.stage++
		fallthrough
	case StageNoisy:
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
