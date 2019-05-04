package engine

import "github.com/mhib/combusken/backend"
import "unsafe"

const (
	TransExact = iota + 1
	TransAlpha
	TransBeta
)

const maxValue = 32000
const Mate = maxValue

func nearestPowerOfTwo(input int) uint64 {
	res := uint64(1)
	for int((res << 1)) <= input {
		res <<= 1
	}
	return res
}

func valueFromTrans(value int16, height int) int16 {
	if value >= Mate-500 {
		return value - int16(height)
	}
	if value <= -Mate+500 {
		return value + int16(height)
	}

	return value
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

func (t *SingleThreadTransTable) Get(key uint64, height int) (ok bool, value int16, eval int16, depth uint8, move backend.Move, flag uint8) {

	var element = &t.Entries[key&t.Mask]
	if element.key != uint32(key>>32) {
		return
	}
	ok = true
	value = valueFromTrans(element.value, height)
	eval = element.eval
	depth = element.depth
	move = element.bestMove
	flag = element.flag
	return
}

func (t *SingleThreadTransTable) Set(key uint64, value, eval, depth int, bestMove backend.Move, flag int, height int) {
	var element = &t.Entries[key&t.Mask]
	element.key = uint32(key >> 32)
	element.value = valueToTrans(value, height)
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
