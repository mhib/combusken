package engine

const (
	TransExact = iota + 1
	TransAlpha
	TransBeta
)

const ValUnknown = 1000002

type TransEntry struct {
	key   uint64
	depth int32
	flag  int32
	value int
}

type TransTable struct {
	Entries []TransEntry
	Length  uint64
}

func NewTransTable() TransTable {
	return TransTable{make([]TransEntry, 1<<21), 1 << 21}
}

func (t *TransTable) Get(key uint64) *TransEntry {
	return &t.Entries[key%t.Length]
}

func (t *TransTable) Set(depth, value, flag int, key uint64) {
	var element = t.Get(key)
	element.key = key
	element.value = value
	element.flag = int32(flag)
	element.depth = int32(depth)
}

func (t *TransTable) Clear() {
	for i := range t.Entries {
		t.Entries[i] = TransEntry{}
	}
}
