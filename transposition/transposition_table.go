package transposition

import "github.com/mhib/combusken/backend"
import "unsafe"
import . "github.com/mhib/combusken/utils"

const maxValue = Mate

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

type transBucket [2]transEntry

type TranspositionTable struct {
	Entries []transBucket
	Mask    uint64
}

func (t *TranspositionTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = transBucket{}
	}
}

func NewTransTable(megabytes int) TranspositionTable {
	size := NearestPowerOfTwo(1024 * 1024 * megabytes / int(unsafe.Sizeof(transBucket{})))
	return TranspositionTable{make([]transBucket, size), size - 1}
}

func (t *TranspositionTable) Get(key uint64) (ok bool, value int16, depth int16, move backend.Move, flag uint8) {
	var bucket = &t.Entries[key&t.Mask]
	var foundIdx int
	if bucket[0].key == key {
		// foundIdx = 0
	} else if bucket[1].key == key {
		foundIdx = 1
	} else {
		return
	}
	ok = true
	value = bucket[foundIdx].value
	depth = int16(bucket[foundIdx].depth) + NoneDepth
	move = bucket[foundIdx].bestMove
	flag = bucket[foundIdx].flag
	return
}

func (entry *transEntry) set(key uint64, value int16, depth int, bestMove backend.Move, flag int) {
	entry.key = key
	entry.bestMove = bestMove
	entry.value = value
	entry.flag = uint8(flag)
	entry.depth = uint8(depth - NoneDepth)
}

const msBit = uint64(1) << 63

// Replacement scheme from laser
func (t *TranspositionTable) Set(key uint64, value int16, depth int, bestMove backend.Move, flag int) {
	var bucket = &t.Entries[key&t.Mask]
	bucket[BoolToInt((key&msBit) == 0)].set(key, value, depth, bestMove, flag)
}
