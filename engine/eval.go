package engine

import . "github.com/mhib/combusken/backend"

// For now as in https://chessprogramming.wikispaces.com/Simplified+evaluation+function

const (
	pawnValue   = 100
	knightValue = 320
	bishopValue = 330
	rookValue   = 500
	queenValue  = 900
	kingValue   = 2000
)

var blackPawnsPos = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	50, 50, 50, 50, 50, 50, 50, 50,
	10, 10, 20, 30, 30, 20, 10, 10,
	5, 5, 10, 25, 25, 10, 5, 5,
	0, 0, 0, 20, 20, 0, 0, 0,
	5, -5, -10, 0, 0, -10, -5, 5,
	5, 10, 10, -20, -20, 10, 10, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
}
var whitePawnsPos [64]int

var blackKnightsPos = [...]int{
	-50, -40, -30, -30, -30, -30, -40, -50,
	-40, -20, 0, 0, 0, 0, -20, -40,
	-30, 0, 10, 15, 15, 10, 0, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 0, 15, 20, 20, 15, 0, -30,
	-30, 5, 10, 15, 15, 10, 5, -30,
	-40, -20, 0, 5, 5, 0, -20, -40,
	-50, -40, -30, -30, -30, -30, -40, -50,
}
var whiteKnightsPos [64]int

var blackBishopsPos = [64]int{
	-20, -10, -10, -10, -10, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 5, 5, 10, 10, 5, 5, -10,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-10, 10, 10, 10, 10, 10, 10, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-20, -10, -10, -10, -10, -10, -10, -20,
}
var whiteBishopsPos [64]int

var blackRooksPos = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 10, 10, 10, 10, 10, 10, 5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	0, 0, 5, 7, 7, 5, 0, 0,
}

var whiteRooksPos [64]int

var blackQueensPos = [64]int{
	-20, -10, -10, -5, -5, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-5, 0, 5, 5, 5, 5, 0, -5,
	0, 0, 5, 5, 5, 5, 0, -5,
	-10, 5, 5, 5, 5, 5, 0, -10,
	-10, 0, 5, 0, 0, 0, 0, -10,
	-20, -10, -10, -5, -5, -10, -10, -20,
}
var whiteQueensPos [64]int

var blackKingMiddleGamePos = [64]int{
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-20, -30, -30, -40, -40, -30, -30, -20,
	-10, -20, -20, -20, -20, -20, -20, -10,
	20, 20, 0, 0, 0, 0, 20, 20,
	20, 30, 10, 0, 0, 10, 30, 20,
}
var whiteKingMiddleGamePos [64]int

var blackKingEndGamePos = [64]int{
	-50, -40, -30, -20, -20, -30, -40, -50,
	-30, -20, -10, 0, 0, -10, -20, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -30, 0, 0, 0, 0, -30, -30,
	-50, -30, -30, -30, -30, -30, -30, -50,
}
var whiteKingEndGamePos [64]int

func rotateArray(input *[64]int, res *[64]int) {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			res[y*8+x] = input[(7-y)*8+x]
		}
	}
}

func init() {
	rotateArray(&blackPawnsPos, &whitePawnsPos)
	rotateArray(&blackBishopsPos, &whiteBishopsPos)
	rotateArray(&blackKnightsPos, &whiteKnightsPos)
	rotateArray(&blackRooksPos, &whiteRooksPos)
	rotateArray(&blackQueensPos, &whiteQueensPos)
	rotateArray(&blackKingMiddleGamePos, &whiteKingMiddleGamePos)
	rotateArray(&blackKingEndGamePos, &whiteKingEndGamePos)
}

func isEndGame(pos *Position) bool {
	return PopCount(pos.White|pos.Black) < 16
}

// Only pawns and kings left
func isLateEndGame(pos *Position) bool {
	return PopCount(pos.White|pos.Black)-2 == PopCount(pos.Pawns)
}

func Evaluate(pos *Position) int {
	var result = 0
	var fromId int
	var fromBB uint64

	for fromBB = pos.Pawns & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result += pawnValue + whitePawnsPos[fromId]
	}
	result -= PopCount(pos.Pawns&pos.White&South(pos.Pawns&pos.White)) * 12
	for fromBB = pos.Knights & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result += knightValue + whiteKnightsPos[fromId]
	}
	for fromBB = pos.Bishops & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result += bishopValue + whiteBishopsPos[fromId]
	}
	if PopCount(pos.Bishops&pos.White) > 1 {
		result += 50
	}
	for fromBB = pos.Rooks & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result += rookValue + whiteRooksPos[fromId]
	}
	for fromBB = pos.Queens & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result += queenValue + whiteQueensPos[fromId]
	}

	for fromBB = pos.Pawns & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result -= pawnValue + blackPawnsPos[fromId]
	}
	result += PopCount((pos.Pawns&pos.Black)&North(pos.Pawns&pos.Black)) * 12
	for fromBB = pos.Knights & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result -= knightValue + blackKnightsPos[fromId]
	}
	for fromBB = pos.Bishops & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result -= bishopValue + blackBishopsPos[fromId]
	}
	if PopCount(pos.Bishops&pos.Black) > 1 {
		result -= 50
	}
	for fromBB = pos.Rooks & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result -= rookValue + blackRooksPos[fromId]
	}
	for fromBB = pos.Queens & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result -= queenValue + blackQueensPos[fromId]
	}

	if isEndGame(pos) {
		result += kingValue + whiteKingEndGamePos[BitScan(pos.Kings&pos.White)]
		result -= kingValue + blackKingEndGamePos[BitScan(pos.Kings&pos.Black)]
	} else {
		result += kingValue + whiteKingMiddleGamePos[BitScan(pos.Kings&pos.White)]
		result -= kingValue + blackKingMiddleGamePos[BitScan(pos.Kings&pos.Black)]
	}

	if pos.WhiteMove {
		return result
	}
	return -result
}
