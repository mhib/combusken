package evaluation

import (
	"unsafe"

	. "github.com/mhib/combusken/utils"
)

var GlobalPawnKingTable PawnKingTable

type PKTableEntry struct {
	key   uint64
	score Score
}

type PawnKingTable struct {
	Entries []PKTableEntry
	Mask    uint64
}

func NewPawnKingTable(megabytes int) PawnKingTable {
	size := NearestPowerOfTwo(1024 * 1024 * megabytes / int(unsafe.Sizeof(PKTableEntry{})))
	return PawnKingTable{make([]PKTableEntry, size), size - 1}
}

func (t *PawnKingTable) Get(key uint64) (ok bool, score Score) {
	var element = &t.Entries[key&t.Mask]
	if element.key != key {
		return
	}
	ok = true
	score = element.score
	return
}

func (t *PawnKingTable) Set(key uint64, score Score) {
	var element = &t.Entries[key&t.Mask]
	element.key = key
	element.score = score
}

func (t *PawnKingTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = PKTableEntry{}
	}
}
