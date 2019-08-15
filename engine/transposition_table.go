package engine

import "github.com/mhib/combusken/backend"
import "unsafe"
import "sync/atomic"
import . "github.com/mhib/combusken/evaluation"
import . "github.com/mhib/combusken/utils"

const maxValue = Mate

const (
	TransBeta  = iota + 1 // Lower bound
	TransAlpha            // Upper bound
	TransExact
)

const NoneDepth = -6

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
	key      uint64
	bestMove backend.Move
	value    int16
	flag     uint8
	depth    uint8
}

type SingleThreadTransTable struct {
	Entries []singleThreadTransEntry
	Mask    uint64
}

func NewSingleThreadTransTable(megabytes int) *SingleThreadTransTable {
	size := NearestPowerOfTwo(1024 * 1024 * megabytes / int(unsafe.Sizeof(singleThreadTransEntry{})))
	return &SingleThreadTransTable{make([]singleThreadTransEntry, size), size - 1}
}

func (t *SingleThreadTransTable) Get(key uint64, height int) (ok bool, value int16, depth int16, move backend.Move, flag uint8) {
	var element = &t.Entries[key&t.Mask]
	if element.key != key {
		return
	}
	ok = true
	value = valueFromTrans(element.value, height)
	depth = int16(element.depth) + NoneDepth
	move = element.bestMove
	flag = element.flag
	return
}

func (t *SingleThreadTransTable) Set(key uint64, value, depth int, bestMove backend.Move, flag int, height int) {
	var element = &t.Entries[key&t.Mask]
	element.key = key
	element.value = valueToTrans(value, height)
	element.flag = uint8(flag)
	element.depth = uint8(depth - NoneDepth)
	element.bestMove = bestMove
}

func (t *SingleThreadTransTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = singleThreadTransEntry{}
	}
}

// As in https://www.chessprogramming.org/Shared_Hash_Table#Lockless

// bits:
// value 16
// depth 8
// flag 2
// move 32

type atomicTransEntry struct {
	key  uint64
	data uint64
}

type AtomicTransTable struct {
	Entries []atomicTransEntry
	Mask    uint64
}

func (t *AtomicTransTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = atomicTransEntry{}
	}
}

func NewAtomicTransTable(megabytes int) *AtomicTransTable {
	size := NearestPowerOfTwo(1024 * 1024 * megabytes / int(unsafe.Sizeof(singleThreadTransEntry{})))
	return &AtomicTransTable{make([]atomicTransEntry, size), size - 1}
}

func (t *AtomicTransTable) Get(key uint64, height int) (ok bool, value int16, depth int16, move backend.Move, flag uint8) {
	idx := key & t.Mask
	data := atomic.LoadUint64(&t.Entries[idx].data)
	if data^atomic.LoadUint64(&t.Entries[idx].key) != key {
		return
	}
	ok = true
	value = valueFromTrans(int16(int(data>>42)-maxValue), height)
	depth = int16((data>>34)&0xFF) + NoneDepth
	flag = uint8((data >> 32) & 3)
	move = backend.Move(data & 0xFFFFFFFF)
	return
}

func (t *AtomicTransTable) Set(key uint64, value, depth int, bestMove backend.Move, flag int, height int) {
	idx := key & t.Mask
	var data uint64
	data |= uint64(valueToTrans(value, height)+maxValue) << 42
	data |= uint64((depth - NoneDepth) << 34)
	data |= uint64(flag << 32)
	data |= uint64(bestMove)
	atomic.StoreUint64(&t.Entries[idx].key, key^data)
	atomic.StoreUint64(&t.Entries[idx].data, data)
}
