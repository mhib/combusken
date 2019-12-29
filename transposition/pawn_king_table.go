package transposition

import "unsafe"
import . "github.com/mhib/combusken/utils"

type PawnKingTable interface {
	Get(key uint64) (ok bool, middleScore int, endScore int)
	Set(key uint64, middleScore int, endScore int)
	Clear()
}

var GlobalPawnKingTable PawnKingTable

type PKTableEntry struct {
	key    uint64
	middle int16
	end    int16
}

type PKTable struct {
	Entries []PKTableEntry
	Mask    uint64
}

func NewPKTable(megabytes int) *PKTable {
	size := NearestPowerOfTwo(1024 * 1024 * megabytes / int(unsafe.Sizeof(PKTableEntry{})))
	return &PKTable{make([]PKTableEntry, size), size - 1}
}

func (t *PKTable) Get(key uint64) (ok bool, middleScore int, endScore int) {
	var element = &t.Entries[key&t.Mask]
	if element.key != key {
		return
	}
	ok = true
	middleScore = int(element.middle)
	endScore = int(element.end)
	return
}

func (t *PKTable) Set(key uint64, middleScore int, endScore int) {
	var element = &t.Entries[key&t.Mask]
	element.key = key
	element.middle = int16(middleScore)
	element.end = int16(endScore)
}

func (t *PKTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = PKTableEntry{}
	}
}
