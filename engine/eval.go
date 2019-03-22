package engine

import . "github.com/mhib/combusken/backend"
import "math"

// For now as in https://chessprogramming.wikispaces.com/Simplified+evaluation+function

var mobilityBonus [64]int

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

	for i := range mobilityBonus {
		mobilityBonus[i] = int(math.Round(5 * math.Sqrt(float64(i))))
	}
}

func IsEndGame(pos *Position) bool {
	return PopCount(pos.White|pos.Black) < 16
}

// CounterGO's version
func isLateEndGame(pos *Position) bool {
	if pos.WhiteMove {
		return ((pos.Rooks|pos.Queens)&pos.White) == 0 && !MoreThanOne((pos.Knights|pos.Bishops)&pos.White)

	} else {
		return ((pos.Rooks|pos.Queens)&pos.Black) == 0 && !MoreThanOne((pos.Knights|pos.Bishops)&pos.Black)
	}
}

func Evaluate(pos *Position) int {
	var result = 0
	var fromId int
	var fromBB uint64

	whiteMobilityArea := ^((pos.Pawns & pos.White) | (BlackPawnsAttacks(pos.Pawns & pos.Black)))
	blackMobilityArea := ^((pos.Pawns & pos.Black) | (WhitePawnsAttacks(pos.Pawns & pos.White)))
	allOccupation := pos.White | pos.Black

	for fromBB = pos.Pawns & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result += PawnValue + whitePawnsPos[fromId]
	}

	// white doubled pawns
	result -= PopCount(pos.Pawns&pos.White&South(pos.Pawns&pos.White)) * 12

	// white knights
	result += (mobilityBonus[PopCount(whiteMobilityArea&KnightsAttacks(pos.Knights&pos.White))])
	for fromBB = pos.Knights & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result += KnightValue + whiteKnightsPos[fromId]
	}

	// white bishops
	result += (mobilityBonus[PopCount(whiteMobilityArea&BishopsAttacks(pos.Bishops&pos.White, allOccupation))])
	for fromBB = pos.Bishops & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result += BishopValue + whiteBishopsPos[fromId]
	}
	// bishop pair bonus
	if MoreThanOne(pos.Bishops & pos.White) {
		result += 50
	}

	// white rooks
	result += (mobilityBonus[PopCount(whiteMobilityArea&RooksAttacks(pos.Rooks&pos.White, allOccupation))])
	for fromBB = pos.Rooks & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result += RookValue + whiteRooksPos[fromId]
	}

	//white queens
	result += (mobilityBonus[PopCount(whiteMobilityArea&QueensAttacks(pos.Queens&pos.White, allOccupation))])
	for fromBB = pos.Queens & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result += QueenValue + whiteQueensPos[fromId]
	}

	// black pawns
	for fromBB = pos.Pawns & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result -= PawnValue + blackPawnsPos[fromId]
	}

	// black doubled pawn
	result += PopCount(((pos.Pawns & pos.Black) & North(pos.Pawns&pos.Black)) * 12)

	// black knights
	result -= mobilityBonus[PopCount(blackMobilityArea&KnightsAttacks(pos.Knights&pos.Black))]
	for fromBB = pos.Knights & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result -= KnightValue + blackKnightsPos[fromId]
	}

	// black bishops
	result -= mobilityBonus[PopCount(blackMobilityArea&BishopsAttacks(pos.Bishops&pos.Black, allOccupation))]
	for fromBB = pos.Bishops & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result -= BishopValue + blackBishopsPos[fromId]
	}
	if MoreThanOne(pos.Bishops & pos.Black) {
		result -= 50
	}

	// black rooks
	result -= mobilityBonus[PopCount(blackMobilityArea&RooksAttacks(pos.Rooks&pos.Black, allOccupation))]
	for fromBB = pos.Rooks & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result -= RookValue + blackRooksPos[fromId]
	}

	// black queens
	result -= mobilityBonus[PopCount(blackMobilityArea&QueensAttacks(pos.Queens&pos.Black, allOccupation))]
	for fromBB = pos.Queens & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		result -= QueenValue + blackQueensPos[fromId]
	}

	if IsEndGame(pos) {
		result += KingValue + whiteKingEndGamePos[BitScan(pos.Kings&pos.White)]
		result -= KingValue + blackKingEndGamePos[BitScan(pos.Kings&pos.Black)]
	} else {
		result += KingValue + whiteKingMiddleGamePos[BitScan(pos.Kings&pos.White)]
		result -= KingValue + blackKingMiddleGamePos[BitScan(pos.Kings&pos.Black)]
	}

	if pos.WhiteMove {
		return result
	}
	return -result
}
