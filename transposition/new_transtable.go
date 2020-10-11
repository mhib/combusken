// +build !linux

package transposition

import (
	"unsafe"

	. "github.com/mhib/combusken/utils"
)

func NewTransTable(megabytes int) TranspositionTable {
	size := NearestPowerOfTwo(1024 * 1024 * megabytes / int(unsafe.Sizeof(transBucket{})))
	return TranspositionTable{make([]transBucket, size), size - 1, uint8(0)}
}
