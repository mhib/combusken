package backend

func Perft(pos *Position, depth int) int {
	result := 0
	var child Position
	var buffer [1000]EvaledMove
	evaled := pos.GenerateAllMoves(buffer[:])
	for _, move := range evaled {
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
