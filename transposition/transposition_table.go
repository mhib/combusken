package transposition

import "github.com/mhib/combusken/backend"
import . "github.com/mhib/combusken/utils"

const NoneDepth = -6

var GlobalTransTable TranspositionTable

func ValueFromTrans(value int16, height int) int16 {
	if value >= Mate-500 {
		return value - int16(height)
	}
	if value <= -Mate+500 {
		return value + int16(height)
	}

	return value
}

func ValueToTrans(value int, height int) int16 {
	if value >= Mate-500 {
		return int16(value + height)
	}
	if value <= -Mate+500 {
		return int16(value - height)
	}

	return int16(value)

}

type transEntry struct {
	key      uint64
	bestMove backend.Move
	value    int16
	flag     uint8
	depth    uint8
}

type TranspositionTable struct {
	Entries []transEntry
	Mask    uint64
}

func (t *TranspositionTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = transEntry{}
	}
}

func (t *TranspositionTable) Get(key uint64) (ok bool, value int16, depth int16, move backend.Move, flag uint8) {
	var element = &t.Entries[key&t.Mask]
	if element.key != key {
		return
	}
	ok = true
	value = element.value
	depth = int16(element.depth) + NoneDepth
	move = element.bestMove
	flag = element.flag
	return
}

func (t *TranspositionTable) Set(key uint64, value int16, depth int, bestMove backend.Move, flag int) {
	var element = &t.Entries[key&t.Mask]
	element.key = key
	element.value = value
	element.flag = uint8(flag)
	element.depth = uint8(depth - NoneDepth)
	element.bestMove = bestMove
}

func (t *TranspositionTable) Prefetch(key uint64) {
	prefetch(&t.Entries[key&t.Mask])
}
