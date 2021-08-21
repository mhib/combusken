// +build !nocgo

package fathom

// #cgo CFLAGS: -O3 -std=gnu11 -w
// #include "tbprobe.h"
// #include <stdlib.h>
import "C"
import (
	"strings"
	"unsafe"

	"github.com/mhib/combusken/backend"
)

var MaxPieceCount = 0
var MinProbeDepth = 0

func SetPath(path string) {
	cPath := C.CString(strings.TrimSpace(path))
	defer C.free(unsafe.Pointer(cPath))
	C.tb_init(cPath)
	MaxPieceCount = int(C.TB_LARGEST)
}

func Clear() {
	C.tb_free()
}

func ProbeWDL(pos *backend.Position, depth int) int64 {
	return int64(C.tb_probe_wdl(
		C.uint64_t(pos.Colours[backend.White]),
		C.uint64_t(pos.Colours[backend.Black]),
		C.uint64_t(pos.Pieces[backend.King]),
		C.uint64_t(pos.Pieces[backend.Queen]),
		C.uint64_t(pos.Pieces[backend.Rook]),
		C.uint64_t(pos.Pieces[backend.Bishop]),
		C.uint64_t(pos.Pieces[backend.Knight]),
		C.uint64_t(pos.Pieces[backend.Pawn]),
		C.uint(0),
		C.uint(0),
		C.uint(0),
		C.bool(pos.SideToMove == backend.White),
	))
}

func IsWDLProbeable(pos *backend.Position, depth int) bool {
	return MaxPieceCount != 0 &&
		pos.FiftyMove == 0 &&
		pos.EpSquare == 0 &&
		pos.Flags == 0xF &&
		depthCardinalityCheck(pos, depth)
}

func depthCardinalityCheck(pos *backend.Position, depth int) bool {
	cardinality := backend.PopCount(pos.Colours[backend.White] | pos.Colours[backend.Black])
	return cardinality < MaxPieceCount || (cardinality == MaxPieceCount && depth >= MinProbeDepth)
}

func IsDTZProbeable(pos *backend.Position) bool {
	return pos.Flags == 0xF && backend.PopCount(pos.Colours[backend.White]|pos.Colours[backend.Black]) <= MaxPieceCount
}

var promoteTranslation = [...]int{backend.None, backend.Queen, backend.Rook, backend.Bishop, backend.Knight}

func ProbeDTZ(pos *backend.Position, moves []backend.EvaledMove) (bool, backend.Move, int, int) {
	var epSquare int

	if pos.EpSquare == 0 {
		epSquare = 0
	} else if pos.SideToMove == backend.White {
		epSquare = pos.EpSquare + 8
	} else {
		epSquare = pos.EpSquare - 8
	}
	result := uint(C.tb_probe_root(
		C.uint64_t(pos.Colours[backend.White]),
		C.uint64_t(pos.Colours[backend.Black]),
		C.uint64_t(pos.Pieces[backend.King]),
		C.uint64_t(pos.Pieces[backend.Queen]),
		C.uint64_t(pos.Pieces[backend.Rook]),
		C.uint64_t(pos.Pieces[backend.Bishop]),
		C.uint64_t(pos.Pieces[backend.Knight]),
		C.uint64_t(pos.Pieces[backend.Pawn]),
		C.uint(pos.FiftyMove),
		C.uint(0),
		C.uint(epSquare),
		C.bool(pos.SideToMove == backend.White),
		nil,
	))
	if result == uint(C.TB_RESULT_FAILED) || result == uint(C.TB_RESULT_CHECKMATE) || result == uint(C.TB_RESULT_STALEMATE) {
		return false, backend.NullMove, 0, 0
	}

	wdl := int(C.tb_get_wdl_go(C.uint(result)))
	dtz := int(C.tb_get_dtz_go(C.uint(result)))
	from := int(C.tb_get_from_go(C.uint(result)))
	to := int(C.tb_get_to_go(C.uint(result)))
	promotion := promoteTranslation[uint(C.tb_get_promotes_go(C.uint(result)))]

	for _, move := range moves {
		if move.From() == from && move.To() == to {
			if promotion != backend.None {
				if move.IsPromotion() && move.PromotedPiece() == promotion {
					return true, move.Move, wdl, dtz
				}
			} else {
				return true, move.Move, wdl, dtz
			}
		}
	}
	return false, backend.NullMove, 0, 0
}
