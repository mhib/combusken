package backend

func Perft(pos *Position, depth int) int {
	result := 0
	var child Position
	var buffer [1000]EvaledMove

	for _, move := range pos.GenerateAllMoves(buffer[:]) {
		if pos.MakeMove(move.Move, &child) {
			if depth > 1 {
				result += Perft(&child, depth-1)
			} else {
				result++
			}
		}
	}
	return result
}
