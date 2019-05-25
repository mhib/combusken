package evaluation

import . "github.com/mhib/combusken/backend"
import "fmt"

const Mate = 32000

const pawnPhase = 0
const knightPhase = 1
const bishopPhase = 1
const rookPhase = 2
const queenPhase = 4
const totalPhase = pawnPhase*16 + knightPhase*4 + bishopPhase*4 + rookPhase*4 + queenPhase*2

type Score struct {
	Middle int16
	End    int16
}

func (s Score) String() string {
	return fmt.Sprintf("Score{%d, %d}", s.Middle, s.End)
}

func addScore(first, second Score) Score {
	return Score{
		Middle: first.Middle + second.Middle,
		End:    first.End + second.End,
	}
}

var PawnValue = Score{174, 246}
var KnightValue = Score{757, 763}
var BishopValue = Score{674, 761}
var RookValue = Score{1009, 1182}
var QueenValue = Score{2312, 2275}

// values from stockfish 10
var pieceScores = [7][8][4]Score{
	{},
	{},
	{ // knight
		{Score{-30, 21}, Score{-17, -108}, Score{0, -55}, Score{-32, -61}},
		{Score{-86, -164}, Score{-59, -9}, Score{-77, -111}, Score{-16, -66}},
		{Score{-54, -66}, Score{-21, -69}, Score{-26, -51}, Score{-11, -15}},
		{Score{-5, -9}, Score{30, -71}, Score{-14, 1}, Score{14, 11}},
		{Score{45, -89}, Score{22, -10}, Score{84, -14}, Score{31, 17}},
		{Score{20, -18}, Score{31, -39}, Score{-8, -11}, Score{161, -9}},
		{Score{-46, -69}, Score{-40, -2}, Score{45, -40}, Score{19, 0}},
		{Score{-293, -172}, Score{-12, 0}, Score{14, 2}, Score{0, -48}},
	},
	{ // Bishop
		{Score{41, -42}, Score{-38, -20}, Score{3, -74}, Score{11, -7}},
		{Score{-56, 8}, Score{23, -84}, Score{-7, -45}, Score{16, -32}},
		{Score{6, -13}, Score{24, -36}, Score{28, -14}, Score{17, -10}},
		{Score{5, -22}, Score{-10, 7}, Score{18, -1}, Score{36, -4}},
		{Score{37, -27}, Score{-6, -17}, Score{69, 4}, Score{50, 9}},
		{Score{0, 24}, Score{29, 7}, Score{66, 13}, Score{77, -15}},
		{Score{-90, 4}, Score{3, 5}, Score{3, -5}, Score{3, 7}},
		{Score{-2, 22}, Score{-7, 88}, Score{-27, -14}, Score{-88, 0}},
	},
	{ // Rook
		{Score{-24, -52}, Score{-17, -50}, Score{-6, -32}, Score{20, -60}},
		{Score{-63, -57}, Score{-50, -36}, Score{-16, -51}, Score{0, -39}},
		{Score{-75, -28}, Score{-61, 15}, Score{-82, -54}, Score{-50, -41}},
		{Score{-79, -2}, Score{-23, 24}, Score{-31, 4}, Score{-19, -22}},
		{Score{-36, 13}, Score{-48, 23}, Score{21, 16}, Score{3, 1}},
		{Score{22, 35}, Score{55, 10}, Score{176, 12}, Score{108, 26}},
		{Score{-1, 24}, Score{-5, 31}, Score{78, 18}, Score{127, 15}},
		{Score{5, 46}, Score{83, 35}, Score{57, 56}, Score{87, 8}},
	},
	{ // Queen
		{Score{-8, -98}, Score{14, -200}, Score{-40, -129}, Score{18, -153}},
		{Score{26, -138}, Score{25, -241}, Score{41, -195}, Score{7, -104}},
		{Score{28, 18}, Score{21, -59}, Score{-16, -1}, Score{-38, 4}},
		{Score{31, -30}, Score{0, 37}, Score{-25, 36}, Score{-38, 38}},
		{Score{51, 33}, Score{3, 52}, Score{37, 42}, Score{-45, 107}},
		{Score{3, 5}, Score{77, 0}, Score{-11, 21}, Score{-1, 75}},
		{Score{34, -13}, Score{-91, 65}, Score{30, 64}, Score{-53, 86}},
		{Score{4, -44}, Score{65, 20}, Score{113, 9}, Score{114, 28}},
	},
	{ // King
		{Score{260, -28}, Score{329, 13}, Score{189, 62}, Score{244, 3}},
		{Score{313, 36}, Score{269, 83}, Score{141, 128}, Score{112, 133}},
		{Score{176, 70}, Score{181, 97}, Score{99, 150}, Score{-15, 176}},
		{Score{-39, 83}, Score{124, 145}, Score{34, 179}, Score{-1, 197}},
		{Score{3, 158}, Score{61, 179}, Score{9, 232}, Score{-130, 219}},
		{Score{118, 172}, Score{43, 231}, Score{-14, 216}, Score{2, 222}},
		{Score{68, 133}, Score{98, 218}, Score{42, 271}, Score{90, 271}},
		{Score{-11, -63}, Score{20, 100}, Score{2, 175}, Score{0, 339}},
	},
}
var pawnScores = [7][8]Score{
	{},
	{Score{-3, -35}, Score{39, -17}, Score{4, -8}, Score{24, -11}, Score{4, -6}, Score{10, -2}, Score{12, 0}, Score{-17, -18}},
	{Score{-7, -34}, Score{-10, -17}, Score{23, -3}, Score{-8, -19}, Score{-4, -4}, Score{-16, 9}, Score{-21, -10}, Score{-8, -32}},
	{Score{-10, -3}, Score{6, 2}, Score{-6, -12}, Score{39, -14}, Score{26, -5}, Score{23, -33}, Score{-30, 17}, Score{-37, -14}},
	{Score{-2, 19}, Score{35, 29}, Score{-23, -12}, Score{40, -31}, Score{22, 0}, Score{-26, 2}, Score{25, 18}, Score{-43, 45}},
	{Score{77, 14}, Score{4, 36}, Score{77, -28}, Score{55, -39}, Score{134, -26}, Score{85, -12}, Score{40, 15}, Score{6, -1}},
	{Score{64, 3}, Score{-1, 106}, Score{1, 53}, Score{32, 1}, Score{16, 77}, Score{102, 96}, Score{-193, 204}, Score{0, 91}},
}

var pawnsConnected = [8][4]Score{
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
	{Score{-4, -14}, Score{5, 7}, Score{11, 37}, Score{36, 11}},
	{Score{17, 3}, Score{53, -5}, Score{19, 17}, Score{44, 30}},
	{Score{2, 23}, Score{6, 21}, Score{26, 21}, Score{51, 10}},
	{Score{0, 0}, Score{-5, 52}, Score{84, 31}, Score{46, 47}},
	{Score{-17, 113}, Score{191, 108}, Score{188, 0}, Score{220, 27}},
	{Score{4, 263}, Score{129, 88}, Score{6, 0}, Score{0, 47}},
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
}

var mobilityBonus = [...][32]Score{
	{Score{-83, -179}, Score{-61, -136}, Score{-45, -94}, Score{-36, -52}, Score{-8, -53}, Score{0, -12}, // Knights
		Score{16, -19}, Score{32, -14}, Score{47, -52}},
	{Score{-11, -107}, Score{-20, -121}, Score{25, -63}, Score{16, -22}, Score{34, -7}, Score{54, 13}, // Bishops
		Score{53, 19}, Score{65, 37}, Score{60, 57}, Score{72, 30}, Score{96, 19}, Score{80, 32},
		Score{92, 39}, Score{75, 4}},
	{Score{-51, -70}, Score{-65, -97}, Score{-24, 34}, Score{-20, 94}, Score{-10, 79}, Score{-11, 103}, // Rooks
		Score{8, 114}, Score{13, 124}, Score{17, 129}, Score{20, 142}, Score{54, 131}, Score{70, 134},
		Score{96, 133}, Score{98, 138}, Score{152, 61}},
	{Score{-39, -36}, Score{-21, -14}, Score{-53, 0}, Score{-36, 0}, Score{-10, -68}, Score{-13, -46}, // Queens
		Score{-8, -52}, Score{4, -4}, Score{16, 34}, Score{49, 20}, Score{55, 26}, Score{50, 76},
		Score{63, 73}, Score{61, 110}, Score{53, 123}, Score{62, 141}, Score{70, 130}, Score{89, 135},
		Score{82, 105}, Score{99, 147}, Score{71, 107}, Score{55, 115}, Score{125, 103}, Score{75, 0},
		Score{38, -7}, Score{3, -6}, Score{0, 0}, Score{0, -1}},
}

var blackPawnsPos [64]Score
var whitePawnsPos [64]Score

var blackPawnsConnected [64]Score
var blackPawnsConnectedMask [64]uint64
var whitePawnsConnected [64]Score
var whitePawnsConnectedMask [64]uint64

var blackKnightsPos [64]Score
var whiteKnightsPos [64]Score

var blackBishopsPos [64]Score
var whiteBishopsPos [64]Score

var blackRooksPos [64]Score
var whiteRooksPos [64]Score

var blackQueensPos [64]Score
var whiteQueensPos [64]Score

var blackKingPos [64]Score
var whiteKingPos [64]Score

var blackPassedMask [64]uint64
var whitePassedMask [64]uint64

var adjacentFilesMask [8]uint64

// PassedRank[Rank] contains a bonus according to the rank of a passed pawn
var passedRank = [7]Score{Score{0, 0}, Score{17, -3}, Score{24, -2}, Score{9, 67}, Score{16, 130}, Score{71, 262}, Score{204, 367}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{Score{-12, 39}, Score{-92, 51}, Score{-59, 27}, Score{-36, -18},
	Score{-51, -3}, Score{14, 18}, Score{-35, 29}, Score{54, 18},
}

var isolated = Score{-16, -22}
var doubled = Score{-47, -90}
var backward = Score{15, 3}
var backwardOpen = Score{-18, -15}

var bishopPair = Score{82, 106}

var minorBehindPawn = Score{5, 69}

var tempo = Score{21, 14}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{10, 41}, Score{74, 1}}

// this bonus only improves midScore
var pawnShieldBonus = [...]Score{Score{7, 0}, Score{4, 0}} // score for every pawn

const whiteKingKingSide = F1 | G1 | H1
const whiteKingKingSideShield1 = (whiteKingKingSide << 8)  // one rank up
const whiteKingKingSideShield2 = (whiteKingKingSide << 16) // two ranks up
const whiteKingQueenSide = A1 | B1 | C1
const whiteKingQueenSideShield1 = (whiteKingQueenSide << 8)  // one rank up
const whiteKingQueenSideShield2 = (whiteKingQueenSide << 16) // two ranks up

const blackKingKingSide = F8 | G8 | H8
const blackKingKingSideShield1 = (blackKingKingSide >> 8)  // one rank down
const blackKingKingSideShield2 = (blackKingKingSide >> 16) // two ranks down
const blackKingQueenSide = A8 | B8 | C8
const blackKingQueenSideShield1 = (blackKingQueenSide >> 8)  // one rank down
const blackKingQueenSideShield2 = (blackKingQueenSide >> 16) // two ranks down

func loadScoresToPieceSquares() {
	for x := 0; x < 4; x++ {
		for y := 0; y < 8; y++ {
			whiteKnightsPos[y*8+x] = addScore(pieceScores[2][y][x], KnightValue)
			whiteKnightsPos[y*8+(7-x)] = addScore(pieceScores[2][y][x], KnightValue)
			blackKnightsPos[(7-y)*8+x] = addScore(pieceScores[2][y][x], KnightValue)
			blackKnightsPos[(7-y)*8+(7-x)] = addScore(pieceScores[2][y][x], KnightValue)

			whiteBishopsPos[y*8+x] = addScore(pieceScores[3][y][x], BishopValue)
			whiteBishopsPos[y*8+(7-x)] = addScore(pieceScores[3][y][x], BishopValue)
			blackBishopsPos[(7-y)*8+x] = addScore(pieceScores[3][y][x], BishopValue)
			blackBishopsPos[(7-y)*8+(7-x)] = addScore(pieceScores[3][y][x], BishopValue)

			whiteRooksPos[y*8+x] = addScore(pieceScores[4][y][x], RookValue)
			whiteRooksPos[y*8+(7-x)] = addScore(pieceScores[4][y][x], RookValue)
			blackRooksPos[(7-y)*8+x] = addScore(pieceScores[4][y][x], RookValue)
			blackRooksPos[(7-y)*8+(7-x)] = addScore(pieceScores[4][y][x], RookValue)

			whiteQueensPos[y*8+x] = addScore(pieceScores[5][y][x], QueenValue)
			whiteQueensPos[y*8+(7-x)] = addScore(pieceScores[5][y][x], QueenValue)
			blackQueensPos[(7-y)*8+x] = addScore(pieceScores[5][y][x], QueenValue)
			blackQueensPos[(7-y)*8+(7-x)] = addScore(pieceScores[5][y][x], QueenValue)

			whiteKingPos[y*8+x] = pieceScores[6][y][x]
			whiteKingPos[y*8+(7-x)] = pieceScores[6][y][x]
			blackKingPos[(7-y)*8+x] = pieceScores[6][y][x]
			blackKingPos[(7-y)*8+(7-x)] = pieceScores[6][y][x]
		}
	}

	for y := 1; y < 7; y++ {
		for x := 0; x < 8; x++ {
			whitePawnsPos[y*8+x] = addScore(pawnScores[y][x], PawnValue)
			blackPawnsPos[(7-y)*8+(7-x)] = addScore(pawnScores[y][x], PawnValue)
		}
	}
	for x := 0; x < 4; x++ {
		for y := 0; y < 8; y++ {
			whitePawnsConnected[y*8+x] = pawnsConnected[y][x]
			whitePawnsConnected[y*8+(7-x)] = pawnsConnected[y][x]
			blackPawnsConnected[(7-y)*8+x] = pawnsConnected[y][x]
			blackPawnsConnected[(7-y)*8+(7-x)] = pawnsConnected[y][x]
		}
	}
}

func init() {
	loadScoresToPieceSquares()
	for i := 8; i <= 55; i++ {
		whitePassedMask[i] = 0
		for file := File(i) - 1; file <= File(i)+1; file++ {
			if file < FILE_A || file > FILE_H {
				continue
			}
			for rank := Rank(i) + 1; rank < RANK_8; rank++ {
				whitePassedMask[i] |= 1 << uint(rank*8+file)
			}
		}
	}
	for i := 55; i >= 8; i-- {
		blackPassedMask[i] = 0
		for file := File(i) - 1; file <= File(i)+1; file++ {
			if file < FILE_A || file > FILE_H {
				continue
			}
			for rank := Rank(i) - 1; rank > RANK_1; rank-- {
				blackPassedMask[i] |= 1 << uint(rank*8+file)
			}
		}
	}
	for i := 8; i <= 55; i++ {
		whitePawnsConnectedMask[i] = BlackPawnAttacks[i] | BlackPawnAttacks[i+8]
		blackPawnsConnectedMask[i] = WhitePawnAttacks[i] | WhitePawnAttacks[i-8]
	}
	for i := range FILES {
		adjacentFilesMask[i] = 0
		if i != 0 {
			adjacentFilesMask[i] |= FILES[i-1]
		}
		if i != 7 {
			adjacentFilesMask[i] |= FILES[i+1]
		}
	}
}

// CounterGO's version
func IsLateEndGame(pos *Position) bool {
	if pos.WhiteMove {
		return ((pos.Rooks|pos.Queens)&pos.White) == 0 && !MoreThanOne((pos.Knights|pos.Bishops)&pos.White)

	} else {
		return ((pos.Rooks|pos.Queens)&pos.Black) == 0 && !MoreThanOne((pos.Knights|pos.Bishops)&pos.Black)
	}
}

func Evaluate(pos *Position) int {
	var fromId int
	var fromBB uint64

	phase := totalPhase
	midResult := 0
	endResult := 0
	whiteMobilityArea := ^((pos.Pawns & pos.White) | (BlackPawnsAttacks(pos.Pawns & pos.Black)))
	blackMobilityArea := ^((pos.Pawns & pos.Black) | (WhitePawnsAttacks(pos.Pawns & pos.White)))
	allOccupation := pos.White | pos.Black

	for fromBB = pos.Pawns & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		if whitePassedMask[fromId]&(pos.Pawns&pos.Black) == 0 {
			midResult += int(passedRank[Rank(fromId)].Middle + passedFile[File(fromId)].Middle)
			endResult += int(passedRank[Rank(fromId)].End + passedFile[File(fromId)].End)
		}
		if adjacentFilesMask[File(fromId)]&(pos.Pawns&pos.White) == 0 {
			midResult += int(isolated.Middle)
			endResult += int(isolated.End)
		}
		if blackPassedMask[fromId]&(pos.Pawns&pos.White) == 0 &&
			WhitePawnAttacks[fromId+8]&(pos.Pawns&pos.Black) != 0 {
			if FILES[File(fromId)]&(pos.Pawns&pos.Black) == 0 {
				midResult += int(backwardOpen.Middle)
				endResult += int(backwardOpen.End)
			} else {
				midResult += int(backward.Middle)
				endResult += int(backward.End)
			}
		} else if whitePawnsConnectedMask[fromId]&(pos.White&pos.Pawns) != 0 {
			midResult += int(whitePawnsConnected[fromId].Middle)
			endResult += int(whitePawnsConnected[fromId].End)
		}
		midResult += int(whitePawnsPos[fromId].Middle)
		endResult += int(whitePawnsPos[fromId].End)
		phase -= pawnPhase
	}

	// white doubled pawns
	doubledCount := PopCount(pos.Pawns & pos.White & South(pos.Pawns&pos.White))
	midResult += doubledCount * int(doubled.Middle)
	endResult += doubledCount * int(doubled.End)

	// white knights

	for fromBB = pos.Knights & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(whiteMobilityArea & KnightAttacks[fromId])
		midResult += int(whiteKnightsPos[fromId].Middle)
		endResult += int(whiteKnightsPos[fromId].End)
		midResult += int(mobilityBonus[0][mobility].Middle)
		endResult += int(mobilityBonus[0][mobility].End)
		if (pos.Pawns>>8)&SquareMask[fromId] != 0 {
			midResult += int(minorBehindPawn.Middle)
			endResult += int(minorBehindPawn.End)
		}
		phase -= knightPhase
	}

	// white bishops
	for fromBB = pos.Bishops & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(whiteMobilityArea & BishopAttacks(fromId, allOccupation))
		midResult += int(mobilityBonus[1][mobility].Middle)
		endResult += int(mobilityBonus[1][mobility].End)
		midResult += int(whiteBishopsPos[fromId].Middle)
		endResult += int(whiteBishopsPos[fromId].End)
		if (pos.Pawns>>8)&SquareMask[fromId] != 0 {
			midResult += int(minorBehindPawn.Middle)
			endResult += int(minorBehindPawn.End)
		}
		phase -= bishopPhase
	}
	// bishop pair bonus
	if MoreThanOne(pos.Bishops & pos.White) {
		midResult += int(bishopPair.Middle)
		endResult += int(bishopPair.End)
	}

	// white rooks
	for fromBB = pos.Rooks & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(whiteMobilityArea & RookAttacks(fromId, allOccupation))
		midResult += int(mobilityBonus[2][mobility].Middle)
		endResult += int(mobilityBonus[2][mobility].End)
		midResult += int(whiteRooksPos[fromId].Middle)
		endResult += int(whiteRooksPos[fromId].End)
		if pos.Pawns&FILES[File(fromId)] == 0 {
			midResult += int(rookOnFile[1].Middle)
			endResult += int(rookOnFile[1].End)
		} else if (pos.Pawns&pos.White)&FILES[File(fromId)] == 0 {
			midResult += int(rookOnFile[0].Middle)
			endResult += int(rookOnFile[0].End)
		}
		phase -= rookPhase
	}

	//white queens
	for fromBB = pos.Queens & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(whiteMobilityArea & QueenAttacks(fromId, allOccupation))
		midResult += int(mobilityBonus[3][mobility].Middle)
		endResult += int(mobilityBonus[3][mobility].End)
		midResult += int(whiteQueensPos[fromId].Middle)
		endResult += int(whiteQueensPos[fromId].End)
		phase -= queenPhase
	}

	// king
	fromId = BitScan(pos.Kings & pos.White)
	midResult += int(whiteKingPos[fromId].Middle)
	endResult += int(whiteKingPos[fromId].End)

	// shield
	if (pos.Kings&pos.White)&whiteKingKingSide != 0 {
		midResult += PopCount(pos.White&pos.Pawns&whiteKingKingSideShield1) * int(pawnShieldBonus[0].Middle)
		midResult += PopCount(pos.White&pos.Pawns&whiteKingKingSideShield2) * int(pawnShieldBonus[1].Middle)
	}
	if (pos.Kings&pos.White)&whiteKingQueenSide != 0 {
		midResult += PopCount(pos.White&pos.Pawns&whiteKingQueenSideShield1) * int(pawnShieldBonus[0].Middle)
		midResult += PopCount(pos.White&pos.Pawns&whiteKingQueenSideShield2) * int(pawnShieldBonus[1].Middle)
	}

	// black pawns
	for fromBB = pos.Pawns & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		if blackPassedMask[fromId]&(pos.Pawns&pos.White) == 0 {
			midResult -= int(passedRank[7-Rank(fromId)].Middle + passedFile[File(fromId)].Middle)
			endResult -= int(passedRank[7-Rank(fromId)].End + passedFile[File(fromId)].End)
		}
		if adjacentFilesMask[File(fromId)]&(pos.Pawns&pos.Black) == 0 {
			midResult -= int(isolated.Middle)
			endResult -= int(isolated.End)
		}
		if whitePassedMask[fromId]&(pos.Pawns&pos.Black) == 0 &&
			BlackPawnAttacks[fromId-8]&(pos.Pawns&pos.White) != 0 {
			if FILES[File(fromId)]&(pos.Pawns&pos.White) == 0 {
				midResult -= int(backwardOpen.Middle)
				endResult -= int(backwardOpen.End)
			} else {
				midResult -= int(backward.Middle)
				endResult -= int(backward.End)
			}
		} else if blackPawnsConnectedMask[fromId]&(pos.Black&pos.Pawns) != 0 {
			midResult -= int(blackPawnsConnected[fromId].Middle)
			endResult -= int(blackPawnsConnected[fromId].End)
		}
		midResult -= int(blackPawnsPos[fromId].Middle)
		endResult -= int(blackPawnsPos[fromId].End)
		phase -= pawnPhase
	}

	// black doubled pawns
	doubledCount = PopCount(pos.Pawns & pos.Black & North(pos.Pawns&pos.Black))
	midResult -= doubledCount * int(doubled.Middle)
	endResult -= doubledCount * int(doubled.End)

	// black knights
	for fromBB = pos.Knights & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(blackMobilityArea & KnightAttacks[fromId])
		midResult -= int(blackKnightsPos[fromId].Middle)
		endResult -= int(blackKnightsPos[fromId].End)
		midResult -= int(mobilityBonus[0][mobility].Middle)
		endResult -= int(mobilityBonus[0][mobility].End)
		if (pos.Pawns<<8)&SquareMask[fromId] != 0 {
			midResult -= int(minorBehindPawn.Middle)
			endResult -= int(minorBehindPawn.End)
		}
		phase -= knightPhase
	}

	// black bishops
	for fromBB = pos.Bishops & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(blackMobilityArea & BishopAttacks(fromId, allOccupation))
		midResult -= int(mobilityBonus[1][mobility].Middle)
		endResult -= int(mobilityBonus[1][mobility].End)
		midResult -= int(blackBishopsPos[fromId].Middle)
		endResult -= int(blackBishopsPos[fromId].End)
		if (pos.Pawns<<8)&SquareMask[fromId] != 0 {
			midResult -= int(minorBehindPawn.Middle)
			endResult -= int(minorBehindPawn.End)
		}
		phase -= bishopPhase
	}

	if MoreThanOne(pos.Bishops & pos.Black) {
		midResult -= int(bishopPair.Middle)
		endResult -= int(bishopPair.End)
	}

	// black rooks
	for fromBB = pos.Rooks & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(blackMobilityArea & RookAttacks(fromId, allOccupation))
		midResult -= int(mobilityBonus[2][mobility].Middle)
		endResult -= int(mobilityBonus[2][mobility].End)
		midResult -= int(blackRooksPos[fromId].Middle)
		endResult -= int(blackRooksPos[fromId].End)
		if pos.Pawns&FILES[File(fromId)] == 0 {
			midResult -= int(rookOnFile[1].Middle)
			endResult -= int(rookOnFile[1].End)
		} else if (pos.Pawns&pos.Black)&FILES[File(fromId)] == 0 {
			midResult -= int(rookOnFile[0].Middle)
			endResult -= int(rookOnFile[0].End)
		}
		phase -= rookPhase
	}

	// black queens
	for fromBB = pos.Queens & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(blackMobilityArea & QueenAttacks(fromId, allOccupation))
		midResult -= int(mobilityBonus[3][mobility].Middle)
		endResult -= int(mobilityBonus[3][mobility].End)
		midResult -= int(blackQueensPos[fromId].Middle)
		endResult -= int(blackQueensPos[fromId].End)
		phase -= queenPhase
	}

	fromId = BitScan(pos.Kings & pos.Black)
	midResult -= int(blackKingPos[fromId].Middle)
	endResult -= int(blackKingPos[fromId].End)
	// shield
	if (pos.Kings&pos.Black)&blackKingKingSide != 0 {
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingKingSideShield1) * int(pawnShieldBonus[0].Middle)
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingKingSideShield2) * int(pawnShieldBonus[1].Middle)
	}
	if (pos.Kings&pos.Black)&blackKingQueenSide != 0 {
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingQueenSideShield1) * int(pawnShieldBonus[0].Middle)
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingQueenSideShield2) * int(pawnShieldBonus[1].Middle)
	}

	if pos.WhiteMove {
		midResult += int(tempo.Middle)
		endResult += int(tempo.End)
	} else {
		midResult -= int(tempo.Middle)
		endResult -= int(tempo.End)
	}

	if phase < 0 {
		phase = 0
	}

	phase = (phase*256 + (totalPhase / 2)) / totalPhase
	result := ((midResult * (256 - phase)) + (endResult * phase)) / 256

	if pos.WhiteMove {
		return result
	}
	return -result
}
