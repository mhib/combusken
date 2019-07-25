package backend

func Perft(pos *Position, depth int) int {
	result := 0
	var child Position
	var buffer [1000]EvaledMove
	counter := uint8(0)
	pos.GenerateQuiets(buffer[:], &counter)
	for _, move := range buffer[:counter] {
		if pos.MakeMove(move.Move, &child) {
			if depth > 1 {
				result += Perft(&child, depth-1)
			} else {
				result++
			}
		}
	}
	counter = 0
	pos.GenerateAllCaptures(buffer[:], &counter)
	for _, move := range buffer[:counter] {
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
