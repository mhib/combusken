package backend

func Perft(pos *Position, depth int) int {
	result := 0
	var undo Undo
	var buffer [1000]EvaledMove

	for _, move := range pos.GenerateAllMoves(buffer[:]) {
		if pos.MakeMove(move.Move, &undo) {
			if depth > 1 {
				result += Perft(pos, depth-1)
			} else {
				result++
			}
		}
		pos.Undo(move.Move, &undo)
	}
	return result
}
