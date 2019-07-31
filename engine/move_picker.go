package engine

import "github.com/mhib/combusken/backend"
import "github.com/mhib/combusken/evaluation"

type stage uint8
type kind uint8

const (
	stageTT stage = iota
	stageGenerateNoisy
	stageGoodNoisy
	stageKiller1
	stageKiller2
	stageCounter
	stageGenerateQuiets
	stageQuiets
	stageBadNoisy
	stageDone
)

const kindNormal = 0
const (
	kindNoBadCaptures = 1 << iota
	kindNoQuiets
)

const kindQs kind = kindNoBadCaptures | kindNoQuiets

type movePicker struct {
	stage
	kind
	noisySize   uint8
	quietsSize  uint8
	split       uint8
	buffer      [256]backend.EvaledMove
	hashMove    backend.Move
	killerMove1 backend.Move
	killerMove2 backend.Move
	counterMove backend.Move
}

func (mp *movePicker) loadSpecialMoves(t *thread, hashMove backend.Move, height int) {
	mp.hashMove = hashMove
	mp.killerMove1 = t.KillerMoves[height][0]
	mp.killerMove2 = t.KillerMoves[height][1]
	pos := &t.stack[height].position
	if pos.LastMove != backend.NullMove {
		mp.counterMove = t.CounterMoves[pos.IntSide()][pos.LastMove.From()][pos.LastMove.To()]
	} else {
		mp.counterMove = backend.NullMove
	}
}

func (mp *movePicker) initNormal(t *thread, hashMove backend.Move, height int) {
	mp.stage = stageTT
	mp.kind = kindNormal
	mp.loadSpecialMoves(t, hashMove, height)
}

func (mp *movePicker) initSingular(t *thread, hashMove backend.Move, height int) {
	mp.stage = stageGenerateNoisy
	mp.kind = kindNoBadCaptures
	mp.loadSpecialMoves(t, hashMove, height)
}

func (mp *movePicker) initQs(t *thread, hashMove backend.Move, height int) {
	mp.kind = kindQs
	if hashMove.IsCaptureOrPromotion() {
		mp.hashMove = hashMove
		mp.stage = stageTT
	} else {
		mp.hashMove = backend.NullMove
		mp.stage = stageGenerateNoisy
	}
	mp.killerMove1 = backend.NullMove
	mp.killerMove2 = backend.NullMove
	mp.counterMove = backend.NullMove
}

func (mp *movePicker) bestMoveIdx(start, end uint8) (best uint8) {
	best = start
	for i := start + 1; i < end; i++ {
		if mp.buffer[i].Value > mp.buffer[best].Value {
			best = i
		}
	}
	return
}

func (mp *movePicker) popMove(index, sizeOffset uint8, size *uint8) (move backend.Move) {
	move = mp.buffer[index].Move
	*size--
	mp.buffer[index] = mp.buffer[sizeOffset+*size]
	return
}

func (mp *movePicker) nextMove(pos *backend.Position, mv *MoveEvaluator, height int) backend.Move {
	var bestMove backend.Move
	var idx uint8
Top:
	switch mp.stage {
	case stageTT:
		mp.stage = stageGenerateNoisy
		if pos.IsMovePseudoLegal(mp.hashMove) {
			return mp.hashMove
		}
		fallthrough
	case stageGenerateNoisy:
		moves := pos.GenerateAllCaptures(mp.buffer[:])
		mp.noisySize = uint8(len(moves))
		mp.split = mp.noisySize
		EvaluateNoisy(mp.buffer[:mp.noisySize])
		mp.stage = stageGoodNoisy
		fallthrough
	case stageGoodNoisy:
	GoodNoisy:
		if mp.noisySize > 0 {
			idx = mp.bestMoveIdx(0, mp.split)
			if mp.buffer[idx].Value > 0 {
				if !evaluation.SeeSign(pos, mp.buffer[idx].Move) {
					mp.buffer[idx].Value = -1
					goto GoodNoisy
				} else {
					bestMove = mp.popMove(idx, 0, &mp.noisySize)
					if bestMove == mp.hashMove {
						goto GoodNoisy
					}
					return bestMove
				}
			}
		}
		if mp.kind&kindNoQuiets != 0 {
			if mp.kind&kindNoBadCaptures != 0 {
				mp.stage = stageDone
				return backend.NullMove
			}
			mp.stage = stageBadNoisy
			goto Top
		}

		mp.stage = stageKiller1
		fallthrough
	case stageKiller1:
		mp.stage = stageKiller2
		if mp.killerMove1 != mp.hashMove && pos.IsMovePseudoLegal(mp.killerMove1) {
			return mp.killerMove1
		}
		fallthrough
	case stageKiller2:
		mp.stage = stageCounter
		if mp.killerMove2 != mp.hashMove && mp.killerMove2 != mp.killerMove1 && pos.IsMovePseudoLegal(mp.killerMove2) {
			return mp.killerMove2
		}
		fallthrough
	case stageCounter:
		mp.stage = stageGenerateQuiets
		if mp.counterMove != mp.hashMove && mp.counterMove != mp.killerMove1 && mp.counterMove != mp.killerMove2 && pos.IsMovePseudoLegal(mp.counterMove) {
			return mp.counterMove
		}
		fallthrough
	case stageGenerateQuiets:
		moves := pos.GenerateQuiets(mp.buffer[mp.split:])
		mp.quietsSize = uint8(len(moves))
		mp.stage = stageQuiets
		mv.EvaluateQuiets(pos, mp.buffer[mp.split:mp.split+mp.quietsSize])
		fallthrough
	case stageQuiets:
	Quiets:
		if mp.quietsSize > 0 {
			idx = mp.bestMoveIdx(mp.split, mp.quietsSize)
			bestMove = mp.popMove(idx, mp.split, &mp.quietsSize)
			if bestMove == mp.hashMove || bestMove == mp.killerMove1 || bestMove == mp.killerMove2 || bestMove == mp.counterMove {
				goto Quiets
			}
			return bestMove
		}
		if mp.kind&kindNoBadCaptures != 0 {
			mp.stage = stageDone
			return backend.NullMove
		}
		mp.stage = stageBadNoisy
		fallthrough
	case stageBadNoisy:
	badNoisy:
		if mp.noisySize > 0 {
			bestMove = mp.popMove(0, 0, &mp.noisySize)
			if bestMove == mp.hashMove {
				goto badNoisy
			}
			return bestMove
		}
		mp.stage = stageDone
		fallthrough
	case stageDone:
		fallthrough
	default:
		return backend.NullMove
	}
}
