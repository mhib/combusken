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

type transEntry struct {
	bestMove   backend.Move
	key        uint32
	value      int16
	evaluation int16
	flag       uint8
	depth      uint8
}

type simpleTransTable struct {
	Entries []transEntry
	Mask    uint64
}

func NewTransTable(megabytes int) *simpleTransTable {
	size := nearestPowerOfTwo(1024 * 1024 * megabytes / int(unsafe.Sizeof(transEntry{})))
	return &simpleTransTable{make([]transEntry, size), size - 1}
}

func (t *simpleTransTable) Get(key uint64, height int) (ok bool, value int16, depth uint8, move backend.Move, evaluation int16, flag uint8) {
	var element = &t.Entries[key&t.Mask]
	key32 := uint32(key >> 32)
	if element.key != key32 {
		return
	}
	ok = true
	value = valueFromTrans(element.value, height)
	depth = element.depth
	move = element.bestMove
	evaluation = element.evaluation
	flag = element.flag
	return
}

func (t *simpleTransTable) Set(key uint64, value, depth int, bestMove backend.Move, evaluation int, flag int, height int) {
	var element = &t.Entries[key&t.Mask]
	element.key = uint32(key >> 32)
	element.value = valueToTrans(value, height)
	element.evaluation = int16(evaluation)
	element.flag = uint8(flag)
	element.depth = uint8(depth)
	element.bestMove = bestMove
}

func (t *simpleTransTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = transEntry{}
	}
}
