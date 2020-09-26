package transposition

import (
	"github.com/mhib/combusken/backend"
	. "github.com/mhib/combusken/utils"
)

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

type flagPv uint8

func newFlagPv(flag int, pvNode bool) flagPv {
	return flagPv(uint8(flag) | (uint8(BoolToInt(pvNode)) << 2))
}

func (f flagPv) getFlag() uint8 {
	return uint8(f) & 3
}

func (f flagPv) isPvNode() bool {
	return uint8(f)>>2 != 0
}

type transEntry struct {
	key      uint32
	bestMove backend.Move
	value    int16
	eval     int16
	flagPv   flagPv
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

func (t *TranspositionTable) Get(key uint64) (ok bool, value int16, eval int16, depth int16, move backend.Move, flag uint8, pvNode bool) {
	var element = &t.Entries[key&t.Mask]
	if element.key != uint32(key>>32) {
		return
	}
	ok = true
	value = element.value
	eval = element.eval
	depth = int16(element.depth) + NoneDepth
	move = element.bestMove
	flag = element.flagPv.getFlag()
	pvNode = element.flagPv.isPvNode()
	return
}

func (t *TranspositionTable) Set(key uint64, value int16, eval int16, depth int, bestMove backend.Move, flag int, pvNode bool) {
	var element = &t.Entries[key&t.Mask]
	element.key = uint32(key >> 32)
	element.value = value
	element.eval = eval
	element.flagPv = newFlagPv(flag, pvNode)
	element.depth = uint8(depth - NoneDepth)
	element.bestMove = bestMove
}

func (t *TranspositionTable) Prefetch(key uint64) {
	prefetch(&t.Entries[key&t.Mask])
}
