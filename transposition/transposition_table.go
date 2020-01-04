package transposition

import "github.com/mhib/combusken/backend"
import "unsafe"
import . "github.com/mhib/combusken/utils"

const NoneDepth = -6

var GlobalTransTable TranspositionTable

type ageFlag uint8

func newAgeFlag(age, flag uint8) ageFlag {
	return ageFlag((age << 2) | uint8(flag))
}

func (a ageFlag) age() uint8 {
	return uint8(a >> 2)
}

func (a ageFlag) flag() uint8 {
	return uint8(a & 3)
}

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
	ageFlag
	depth uint8
}

type transBucket [2]transEntry

type TranspositionTable struct {
	Entries []transBucket
	Mask    uint64
	age     uint8
}

func (t *TranspositionTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = transBucket{}
	}
}

func (t *TranspositionTable) IncrementAge() {
	t.age++
}

func NewTransTable(megabytes int) TranspositionTable {
	size := NearestPowerOfTwo(1024 * 1024 * megabytes / int(unsafe.Sizeof(transBucket{})))
	return TranspositionTable{make([]transBucket, size), size - 1, 0}
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
	flag = bucket[foundIdx].flag()
	return
}

func (entry *transEntry) set(key uint64, value int16, depth int, bestMove backend.Move, flag int, age uint8) {
	entry.key = key
	entry.bestMove = bestMove
	entry.value = value
	entry.ageFlag = newAgeFlag(age, uint8(flag))
	entry.depth = uint8(depth - NoneDepth)
}

// Replacement scheme from laser
func (t *TranspositionTable) Set(key uint64, value int16, depth int, bestMove backend.Move, flag int) {
	var bucket = &t.Entries[key&t.Mask]
	if bucket[0].key == key {
		bucket[0].set(key, value, depth, bestMove, flag, t.age)
	} else if bucket[1].key == key {
		bucket[1].set(key, value, depth, bestMove, flag, t.age)
	} else {
		score0 := 16*int(t.age-bucket[0].age()) + depth - int(bucket[0].depth)
		score1 := 16*int(t.age-bucket[1].age()) + depth - int(bucket[1].depth)
		toReplaceIdx := BoolToInt(score0 <= score1)
		if score0 >= -2 || score1 >= -2 {
			bucket[toReplaceIdx].set(key, value, depth, bestMove, flag, t.age)
		}
	}
}

func (t *TranspositionTable) Prefetch(key uint64) {
	prefetch(&t.Entries[key&t.Mask])
}
