package utils

const Mate = 32000

const (
	TransBeta  = iota + 1 // Lower bound
	TransAlpha            // Upper bound
	TransExact
)

var SEEValues = [...]int{100, 450, 450, 675, 1300, Mate / 2, 0}
