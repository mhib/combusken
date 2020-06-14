package engine

import "sync/atomic"
import "unsafe"

type breadcrumb struct {
	threadPtr uintptr
	zobrist   uint64
}

const breadCrumbsSize = 1024

var breadCrumbs [breadCrumbsSize]breadcrumb

func markPosition(t *thread, zobrist uint64, height int) (marked, owning bool) {
	if height >= 8 {
		return false, false
	}
	entry := &breadCrumbs[zobrist&(breadCrumbsSize-1)]
	ptr := atomic.LoadUintptr(&entry.threadPtr)
	if ptr == 0 {
		atomic.StoreUint64(&entry.zobrist, zobrist)
		atomic.StoreUintptr(&entry.threadPtr, (uintptr)(unsafe.Pointer(t)))
		return false, true
	} else if ptr != (uintptr)(unsafe.Pointer(t)) && atomic.LoadUint64(&entry.zobrist) == zobrist {
		return true, false
	}
	return false, false
}

func unmarkPosition(zobrist uint64) {
	atomic.StoreUintptr(&breadCrumbs[zobrist&(breadCrumbsSize-1)].threadPtr, 0)
}
