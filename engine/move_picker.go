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

const (
	kindNormal        kind = iota // 00
	kindNoBadCaptures             // 01
	kindNoQuiets                  // 10
)

const kindQs kind = kindNoBadCaptures | kindNoQuiets // 11

const noneMove = backend.Move(1)

type movePicker struct {
	buffer      [256]backend.EvaledMove
	hashMove    backend.Move
	killerMove1 backend.Move
	killerMove2 backend.Move
	counterMove backend.Move
	stage
	kind
	noisySize    uint8
	quietsSize   uint8
	split        uint8
	badNoisySize uint8
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

func (mp *movePicker) bestMoveToFirst(start, end uint8) {
	best := start
	for i := start + 1; i < end; i++ {
		if mp.buffer[i].Value > mp.buffer[best].Value {
			best = i
		}
	}
	mp.buffer[start], mp.buffer[best] = mp.buffer[best], mp.buffer[start]
}

func (mp *movePicker) popBadNoisy() backend.Move {
	ret := mp.buffer[mp.noisySize].Move
	mp.noisySize++
	mp.badNoisySize--
	return ret
}

func (mp *movePicker) popMove(size *uint8, index int) backend.Move {
	ret := mp.buffer[index].Move
	*size--
	mp.buffer[index] = mp.buffer[*size]
	return ret
}

func (mp *movePicker) noteBadNoisyMove() {
	mp.noisySize--
	mp.badNoisySize++
}

func (mp *movePicker) popGoodNoisyMove() backend.Move {
	ret := mp.buffer[mp.badNoisySize].Move
	mp.noisySize--
	mp.buffer[mp.badNoisySize], mp.buffer[mp.noisySize] = mp.buffer[mp.noisySize], mp.buffer[mp.badNoisySize]
	return ret
}

func (mp *movePicker) popQuietMove() backend.Move {
	ret := mp.buffer[mp.split].Move
	mp.quietsSize--
	mp.buffer[mp.split], mp.buffer[mp.split+mp.quietsSize] = mp.buffer[mp.split+mp.quietsSize], mp.buffer[mp.split]
	return ret
}

func (mp *movePicker) nextMove(pos *backend.Position, mv *MoveEvaluator, height int) backend.Move {
	var bestMove backend.Move
Top:
	switch mp.stage {
	case stageTT:
		mp.stage = stageGenerateNoisy
		if pos.IsMovePseudoLegal(mp.hashMove) {
			return mp.hashMove
		}
		fallthrough
	case stageGenerateNoisy:
		mp.noisySize = 0
		mp.badNoisySize = 0
		pos.GenerateAllCaptures(mp.buffer[:], &mp.noisySize)
		EvaluateNoisy(mp.buffer[:mp.noisySize])
		mp.split = mp.noisySize
		mp.stage = stageGoodNoisy
		fallthrough
	case stageGoodNoisy:
	GoodNoisy:
		if mp.noisySize > 0 {
			mp.bestMoveToFirst(mp.badNoisySize, mp.badNoisySize+mp.noisySize)
			if !evaluation.SeeSign(pos, mp.buffer[mp.badNoisySize].Move) {
				mp.noteBadNoisyMove()
				goto GoodNoisy
			} else {
				bestMove = mp.popGoodNoisyMove()
				if bestMove == mp.hashMove {
					goto GoodNoisy
				}
				return bestMove
			}
		}
		if mp.kind&kindNoQuiets != 0 {
			if mp.kind&kindNoBadCaptures != 0 {
				mp.stage = stageDone
				return noneMove
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
		mp.stage = stageQuiets
		mp.quietsSize = 0
		pos.GenerateQuiets(mp.buffer[mp.split:], &mp.quietsSize)
		mv.EvaluateQuiets(pos, mp.buffer[mp.split:mp.split+mp.quietsSize], height)
		fallthrough
	case stageQuiets:
	Quiets:
		if mp.quietsSize > 0 {
			mp.bestMoveToFirst(mp.split, mp.split+mp.quietsSize)
			bestMove = mp.popQuietMove()
			if bestMove == mp.hashMove || bestMove == mp.killerMove1 || bestMove == mp.killerMove2 || bestMove == mp.counterMove {
				goto Quiets
			}
			return bestMove
		}
		mp.stage = stageBadNoisy
		fallthrough
	case stageBadNoisy:
	badNoisy:
		if mp.badNoisySize > 0 && mp.kind&kindNoBadCaptures == 0 {
			bestMove = mp.popBadNoisy()
			if bestMove == mp.hashMove {
				goto badNoisy
			}
			return bestMove
		}
		mp.stage = stageDone
		fallthrough
	case stageDone:
		return noneMove
	}
	return noneMove
}
