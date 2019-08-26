package evaluation

import "unsafe"
import . "github.com/mhib/combusken/utils"

type PawnKingTable interface {
	Get(key uint64) (ok bool, value Score)
	Set(key uint64, value Score)
	Clear()
}

type PKTableEntry struct {
	key   uint64
	value Score
}

type PKTable struct {
	Entries []PKTableEntry
	Mask    uint64
}

func NewPKTable(megabytes int) *PKTable {
	size := NearestPowerOfTwo(1024 * 1024 * megabytes / int(unsafe.Sizeof(PKTableEntry{})))
	return &PKTable{make([]PKTableEntry, size), size - 1}
}

func (t *PKTable) Get(key uint64) (ok bool, value Score) {
	var element = &t.Entries[key&t.Mask]
	if element.key != key {
		return
	}
	ok = true
	value = element.value
	return
}

func (t *PKTable) Set(key uint64, value Score) {
	var element = &t.Entries[key&t.Mask]
	element.key = key
	element.value = value
}

func (t *PKTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = PKTableEntry{}
	}
}
