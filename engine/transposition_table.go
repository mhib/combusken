package engine

import "github.com/mhib/combusken/backend"

const (
	TransExact = iota + 1
	TransAlpha
	TransBeta
)

const ValUnknown = 1000002

type TransEntry struct {
	key      uint64
	depth    int32
	flag     int32
	bestMove backend.Move
	value    int
}

type TransTable struct {
	Entries []TransEntry
	Mask    uint64
}

func NewTransTable() TransTable {
	return TransTable{make([]TransEntry, 1<<19), (1 << 19) - 1}
}

func (t *TransTable) Get(key uint64) *TransEntry {
	return &t.Entries[key&t.Mask]
}

func (t *TransTable) Set(depth, value, flag int, key uint64, bestMove backend.Move, height int) {
	var element = t.Get(key)
	element.key = key
	if value >= Mate-500 {
		value -= height
	} else if value <= -Mate+500 {
		value += height
	}
	element.value = value
	element.flag = int32(flag)
	element.depth = int32(depth)
	element.bestMove = bestMove
}

func (t *TransTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = TransEntry{}
	}
}
