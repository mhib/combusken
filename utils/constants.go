package utils

const Mate = 32000
const UnknownValue = int16(32002)

const (
	TransNone  = iota
	TransBeta  // Lower bound
	TransAlpha // Upper bound
	TransExact
)
