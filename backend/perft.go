package backend

func Perft(pos *Position, depth int) int {
	result := 0
	var child Position
	var buffer [1000]EvaledMove
	noisySize := pos.GenerateNoisy(buffer[:])
	quietsSize := pos.GenerateQuiet(buffer[noisySize:])
	for _, move := range buffer[:noisySize+quietsSize] {
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
