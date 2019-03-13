package engine

import "github.com/mhib/combusken/backend"
import "unsafe"

const (
	TransExact = iota + 1
	TransAlpha
	TransBeta
)

const ValUnknown = 1000002

type TransEntry struct {
	key      uint64
	bestMove backend.Move
	value    int16
	flag     int8
	depth    uint8
}

type TransTable struct {
	Entries []TransEntry
	Mask    uint64
}

func nearestPowerOfTwo(input int) uint64 {
	res := uint64(1)
	for int((res << 1)) <= input {
		res <<= 1
	}
	return res
}

func NewTransTable(megabytes int) TransTable {
	size := nearestPowerOfTwo(1024 * 1024 * megabytes / int(unsafe.Sizeof(TransEntry{})))
	return TransTable{make([]TransEntry, size), size - 1}
}

func (t *TransTable) Get(key uint64) *TransEntry {
	return &t.Entries[key&t.Mask]
}

func (t *TransTable) Set(depth, value, flag int, key uint64, bestMove backend.Move, height int) {
	var element = t.Get(key)
	element.key = key
	if value >= Mate-500 {
		value += height
	} else if value <= -Mate+500 {
		value -= height
	}
	element.value = int16(value)
	element.flag = int8(flag)
	element.depth = uint8(depth)
	element.bestMove = bestMove
}

func (t *TransTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = TransEntry{}
	}
}

func valueFromTrans(value, height int) int {
	if value >= Mate-500 {
		return value - height

	}
	if value <= -Mate+500 {
		return value + height
	}

	return value
}
