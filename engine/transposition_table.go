package engine

import "github.com/mhib/combusken/backend"
import "unsafe"

const (
	TransNone = iota
	TransExact
	TransAlpha
	TransBeta
)

const maxValue = 32000
const TRANS_VALUE_NONE = maxValue + 1
const Mate = maxValue

func nearestPowerOfTwo(input int) uint64 {
	res := uint64(1)
	for int((res << 1)) <= input {
		res <<= 1
	}
	return res
}

func valueFromTrans(value int16, height int) int {
	if value >= Mate-500 {
		return int(value) - height
	}
	if value <= -Mate+500 {
		return int(value) + height
	}

	return int(value)
}

func valueToTrans(value int, height int) int16 {
	if value >= Mate-500 {
		return int16(value + height)
	}
	if value <= -Mate+500 {
		return int16(value - height)
	}

	return int16(value)

}

type singleThreadTransEntry struct {
	key      uint32
	bestMove backend.Move
	value    int16
	eval     int16
	flag     uint8
	depth    uint8
}

type SingleThreadTransTable struct {
	Entries []singleThreadTransEntry
	Mask    uint64
}

func NewSingleThreadTransTable(megabytes int) *SingleThreadTransTable {
	size := nearestPowerOfTwo(1024 * 1024 * megabytes / int(unsafe.Sizeof(singleThreadTransEntry{})))
	return &SingleThreadTransTable{make([]singleThreadTransEntry, size), size - 1}
}

func (t *SingleThreadTransTable) Get(key uint64) (ok bool, value int16, eval int16, depth uint8, move backend.Move, flag uint8) {
	var element = &t.Entries[key&t.Mask]
	if element.key != uint32(key>>32) {
		return
	}
	ok = true
	value = element.value
	eval = element.eval
	depth = element.depth
	move = element.bestMove
	flag = element.flag
	return
}

func (t *SingleThreadTransTable) Set(key uint64, value int16, eval, depth int, bestMove backend.Move, flag int) {
	var element = &t.Entries[key&t.Mask]
	element.key = uint32(key >> 32)
	element.value = value
	element.eval = int16(eval)
	element.flag = uint8(flag)
	element.depth = uint8(depth)
	element.bestMove = bestMove
}

func (t *SingleThreadTransTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = singleThreadTransEntry{}
	}
}
