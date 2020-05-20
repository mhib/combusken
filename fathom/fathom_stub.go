// +build !cgo

package fathom

import "github.com/mhib/combusken/backend"

var MAX_PIECE_COUNT = 0
var MIN_PROBE_DEPTH = 0

func SetPath(path string) {
}

func Clear() {
}

func ProbeWDL(pos *backend.Position, depth int) int64 {
	return 0
}

func IsWDLProbeable(pos *backend.Position, depth int) bool {
	return false
}

func IsDTZProbeable(pos *backend.Position) bool {
	return false
}

func ProbeDTZ(pos *backend.Position, moves []backend.EvaledMove) (bool, backend.Move, int, int) {
	return false, backend.NullMove, 0, 0
}
