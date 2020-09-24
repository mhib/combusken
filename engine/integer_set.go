package engine

import "sort"

type IntegerSet struct {
	mask    uint64
	members []uint64
}

func (is *IntegerSet) Includes(value uint64) bool {
	size := len(is.members)
	for i := 0; i < size; {
		currentVal := is.members[i]
		if currentVal == value {
			return true
		}
		if currentVal < value {
			i = i*2 + 2
		} else {
			i = i*2 + 1
		}
	}
	return false
}

func (is *IntegerSet) MayInclude(value uint64) bool {
	return value&is.mask == value
}

func IntegerSetFromSlice(slice []uint64) IntegerSet {
	sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
	mask := uint64(0)
	for _, x := range slice {
		mask |= x
	}
	return IntegerSet{mask, heapify(slice)}
}

func heapify(src []uint64) []uint64 {
	dst := make([]uint64, len(src))
	id := 0
	size := len(src)
	for jmp := size + 1; jmp >= 2; jmp /= 2 {
		for i := jmp/2 - 1; i < size; i += jmp {
			dst[id] = src[i]
			id++
		}
	}
	return dst
}
