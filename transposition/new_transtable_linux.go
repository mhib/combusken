// +build linux

package transposition

import (
	"reflect"
	"unsafe"

	. "github.com/mhib/combusken/utils"
	"golang.org/x/sys/unix"
)

func NewTransTable(megabytes int) TranspositionTable {
	sizeOfEntry := uint64(unsafe.Sizeof(transBucket{}))
	entriesCount := NearestPowerOfTwo(1024 * 1024 * megabytes / int(sizeOfEntry))
	table := TranspositionTable{make([]transBucket, entriesCount), entriesCount - 1, 0}
	unix.Syscall(unix.SYS_MADVISE, uintptr((*reflect.SliceHeader)(unsafe.Pointer(&table)).Data), uintptr(entriesCount*sizeOfEntry), uintptr(unix.MADV_HUGEPAGE))
	return table
}
