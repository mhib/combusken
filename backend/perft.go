package backend

import "fmt"

func Perft(pos *Position, depth int) int {
	result := 0
	var child Position
	var buffer [1000]Move

	for _, move := range pos.GenerateAllMoves(buffer[:]) {
		if pos.MakeMove(move, &child) {
			if depth > 1 {
				result += Perft(&child, depth-1)
			} else {
				result++
			}
		}
	}
	return result
}
