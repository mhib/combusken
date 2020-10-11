package transposition

import (
	"github.com/mhib/combusken/backend"
	. "github.com/mhib/combusken/utils"
)

const NoneDepth = -6

var GlobalTransTable TranspositionTable

type AgeFlag uint8

func (a AgeFlag) Flag() uint8 {
	return uint8(a) & 3
}

func (a AgeFlag) Age() uint8 {
	return uint8(a) >> 2
}

func NewAgeFlag(age, flag uint8) AgeFlag {
	return AgeFlag((age << 2) | flag)
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
	key      uint32
	bestMove backend.Move
	value    int16
	eval     int16
	ageFlag  AgeFlag
	depth    uint8
}

type transBucket [2]transEntry

type TranspositionTable struct {
	Entries []transBucket
	Mask    uint64
	Age     uint8
}

func (t *TranspositionTable) Clear() {
	for i := range t.Entries {
		t.Entries[i][0] = transEntry{}
		t.Entries[i][1] = transEntry{}
	}
}

func (t *TranspositionTable) Get(key uint64) (ok bool, value int16, eval int16, depth int16, move backend.Move, flag uint8) {
	var bucket = &t.Entries[key&t.Mask]
	var element *transEntry
	matchKey := uint32(key >> 32)
	if bucket[0].key == matchKey {
		element = &bucket[0]
	} else if bucket[1].key == matchKey {
		element = &bucket[1]
	} else {
		return
	}
	ok = true
	value = element.value
	eval = element.eval
	depth = int16(element.depth) + NoneDepth
	move = element.bestMove
	flag = element.ageFlag.Flag()
	return
}

func (t *TranspositionTable) Set(key uint64, value int16, eval int16, depth int, bestMove backend.Move, flag int) {
	var bucket = &t.Entries[key&t.Mask]
	matchKey := uint32(key >> 32)
	var element *transEntry
	if bucket[0].key == matchKey {
		element = &bucket[0]
	} else if bucket[1].key == matchKey {
		element = &bucket[1]
	} else {
		element = &bucket[0]
		firstScore := (int(element.ageFlag.Age()) - int(t.Age)) + (int(element.depth)-depth)*4
		secondScore := (int(bucket[1].ageFlag.Age()) - int(t.Age)) + (int(bucket[1].depth)-depth)*4
		if firstScore > 0 && secondScore > 0 {
			return
		}
		if secondScore < firstScore {
			element = &bucket[1]
		}
	}
	element.key = uint32(key >> 32)
	element.value = value
	element.eval = eval
	element.ageFlag = NewAgeFlag(t.Age, uint8(flag))
	element.depth = uint8(depth - NoneDepth)
	element.bestMove = bestMove
}

func (t *TranspositionTable) Prefetch(key uint64) {
	prefetch(&t.Entries[key&t.Mask])
}

func (t *TranspositionTable) IncrementAge() {
	t.Age = (t.Age + 1) & uint8(0x3f)
}
