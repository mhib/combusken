// +build linux

package transposition

import "unsafe"
import . "github.com/mhib/combusken/utils"
import "golang.org/x/sys/unix"
import "reflect"

func NewTransTable(megabytes int) TranspositionTable {
	sizeOfEntry := uint64(unsafe.Sizeof(transEntry{}))
	entriesCount := NearestPowerOfTwo(1024 * 1024 * megabytes / int(sizeOfEntry))
	table := TranspositionTable{make([]transEntry, entriesCount), entriesCount - 1}
	unix.Syscall(unix.SYS_MADVISE, uintptr((*reflect.SliceHeader)(unsafe.Pointer(&table)).Data), uintptr(entriesCount*sizeOfEntry), uintptr(unix.MADV_HUGEPAGE))
	return table
}
