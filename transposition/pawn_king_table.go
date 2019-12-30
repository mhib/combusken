package transposition

import "unsafe"
import . "github.com/mhib/combusken/utils"

type PawnKingTable interface {
	Get(key uint64) (ok bool, score Score)
	Set(key uint64, score Score)
	Clear()
}

var GlobalPawnKingTable PawnKingTable

type PKTableEntry struct {
	key   uint64
	score Score
}

type PKTable struct {
	Entries []PKTableEntry
	Mask    uint64
}

func NewPKTable(megabytes int) *PKTable {
	size := NearestPowerOfTwo(1024 * 1024 * megabytes / int(unsafe.Sizeof(PKTableEntry{})))
	return &PKTable{make([]PKTableEntry, size), size - 1}
}

func (t *PKTable) Get(key uint64) (ok bool, score Score) {
	var element = &t.Entries[key&t.Mask]
	if element.key != key {
		return
	}
	ok = true
	score = element.score
	return
}

func (t *PKTable) Set(key uint64, score Score) {
	var element = &t.Entries[key&t.Mask]
	element.key = key
	element.score = score
}

func (t *PKTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = PKTableEntry{}
	}
}
