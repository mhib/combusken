//go:build !cgo
// +build !cgo

package fathom

import "github.com/mhib/combusken/chess"

var MaxPieceCount = 0
var MinProbeDepth = 0

func SetPath(path string) {
}

func Clear() {
}

func ProbeWDL(pos *chess.Position, depth int) int64 {
	return 0
}

func IsWDLProbeable(pos *chess.Position, depth int) bool {
	return false
}

func IsDTZProbeable(pos *chess.Position) bool {
	return false
}

func ProbeDTZ(pos *chess.Position, moves []chess.EvaledMove) (bool, chess.Move, int, int) {
	return false, chess.NullMove, 0, 0
}
