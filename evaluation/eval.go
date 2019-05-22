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

var PawnValue = Score{152, 258}
var KnightValue = Score{661, 807}
var BishopValue = Score{562, 810}
var RookValue = Score{785, 1229}
var QueenValue = Score{2104, 2143}

// values from stockfish 10
var pieceScores = [7][8][4]Score{
	{},
	{},
	{ // knight
		{Score{-170, -74}, Score{-50, -113}, Score{-78, -83}, Score{-48, -115}},
		{Score{-12, -133}, Score{-69, -92}, Score{-64, -43}, Score{-41, -98}},
		{Score{-65, -134}, Score{-68, -53}, Score{-54, -89}, Score{-31, -17}},
		{Score{-15, -98}, Score{58, -65}, Score{-14, -50}, Score{-21, 8}},
		{Score{98, -42}, Score{-33, -8}, Score{22, -3}, Score{0, 10}},
		{Score{17, -47}, Score{-26, -51}, Score{31, -12}, Score{61, -32}},
		{Score{-47, -82}, Score{-89, -47}, Score{-80, -10}, Score{16, -19}},
		{Score{-391, -73}, Score{-10, -2}, Score{-130, -54}, Score{-13, -14}},
	},
	{ // Bishop
		{Score{75, -73}, Score{71, -92}, Score{-13, -27}, Score{-46, -24}},
		{Score{12, -105}, Score{12, -69}, Score{1, -41}, Score{-17, -20}},
		{Score{12, -78}, Score{23, -12}, Score{20, -48}, Score{3, -19}},
		{Score{-17, -38}, Score{-8, -23}, Score{1, -7}, Score{64, -15}},
		{Score{-27, -18}, Score{-9, -1}, Score{-1, -33}, Score{65, -1}},
		{Score{19, -121}, Score{1, -12}, Score{49, -19}, Score{-21, -16}},
		{Score{-93, 0}, Score{-19, -1}, Score{-48, 2}, Score{-38, 0}},
		{Score{-39, 46}, Score{-3, 2}, Score{3, 4}, Score{-23, -22}},
	},
	{ // Rook
		{Score{-52, -62}, Score{-25, -46}, Score{-31, -31}, Score{-11, -52}},
		{Score{-100, -42}, Score{-51, -51}, Score{-58, -16}, Score{-37, -61}},
		{Score{-80, -26}, Score{-29, -33}, Score{-66, 20}, Score{-43, -35}},
		{Score{-29, 3}, Score{-5, 0}, Score{-2, -3}, Score{-29, -8}},
		{Score{-26, 2}, Score{-14, 21}, Score{73, -5}, Score{23, -12}},
		{Score{-77, 41}, Score{3, 5}, Score{69, -6}, Score{154, -25}},
		{Score{-1, 11}, Score{-10, 37}, Score{90, 22}, Score{97, -5}},
		{Score{47, 20}, Score{112, 22}, Score{59, 53}, Score{37, 27}},
	},
	{ // Queen
		{Score{-40, -104}, Score{-35, -89}, Score{-35, -113}, Score{3, -160}},
		{Score{-8, -172}, Score{4, -34}, Score{29, -110}, Score{8, -92}},
		{Score{-30, -74}, Score{4, -20}, Score{-12, -11}, Score{-25, -37}},
		{Score{-29, -39}, Score{-8, 34}, Score{-27, 56}, Score{-42, 74}},
		{Score{34, -8}, Score{36, 50}, Score{23, 30}, Score{-1, 125}},
		{Score{67, -35}, Score{97, -17}, Score{67, 10}, Score{79, 52}},
		{Score{-2, -47}, Score{-97, 70}, Score{-1, 0}, Score{-1, 92}},
		{Score{6, -65}, Score{63, -13}, Score{1, 1}, Score{56, 28}},
	},
	{ // King
		{Score{276, -36}, Score{317, 25}, Score{208, 64}, Score{258, 26}},
		{Score{306, 20}, Score{237, 103}, Score{173, 130}, Score{120, 144}},
		{Score{204, 61}, Score{166, 118}, Score{99, 154}, Score{54, 168}},
		{Score{100, 119}, Score{126, 138}, Score{117, 180}, Score{114, 175}},
		{Score{0, 119}, Score{160, 199}, Score{42, 189}, Score{1, 213}},
		{Score{121, 174}, Score{-3, 256}, Score{140, 253}, Score{0, 220}},
		{Score{77, 168}, Score{160, 291}, Score{105, 274}, Score{86, 286}},
		{Score{59, 0}, Score{86, 190}, Score{49, 227}, Score{0, 213}},
	},
}
var pawnScores = [7][8]Score{
	{},
	{Score{-13, 13}, Score{27, 1}, Score{10, 11}, Score{24, 32}, Score{12, 17}, Score{22, -3}, Score{20, 5}, Score{-7, -19}},
	{Score{-3, -2}, Score{2, -14}, Score{-15, -1}, Score{1, -10}, Score{-14, 0}, Score{22, -6}, Score{-4, 2}, Score{-17, -21}},
	{Score{-18, 10}, Score{-14, 1}, Score{10, -40}, Score{27, 1}, Score{26, -20}, Score{19, -15}, Score{-38, 4}, Score{-20, 5}},
	{Score{-6, 7}, Score{-13, 18}, Score{-37, -22}, Score{17, -21}, Score{18, -22}, Score{-22, -10}, Score{-3, 17}, Score{-23, 5}},
	{Score{0, 64}, Score{-10, 56}, Score{29, 48}, Score{34, 29}, Score{132, 0}, Score{-1, 39}, Score{-18, 68}, Score{-15, 67}},
	{Score{18, 31}, Score{0, 103}, Score{56, 49}, Score{63, 55}, Score{23, 167}, Score{-27, 112}, Score{69, 91}, Score{-23, 70}},
}

var pawnsConnected = [8][4]Score{
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
	{Score{0, 1}, Score{10, 0}, Score{1, 4}, Score{17, 1}},
	{Score{11, 3}, Score{37, 5}, Score{0, 3}, Score{34, 58}},
	{Score{15, -6}, Score{33, 10}, Score{14, 9}, Score{29, 10}},
	{Score{68, 4}, Score{56, 15}, Score{51, 30}, Score{65, 35}},
	{Score{61, 78}, Score{118, 128}, Score{70, 69}, Score{148, 32}},
	{Score{110, 267}, Score{197, 138}, Score{227, 157}, Score{-16, 48}},
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
}

var mobilityBonus = [...][32]Score{
	{Score{-90, -192}, Score{-67, -284}, Score{-57, -102}, Score{-46, -44}, Score{-35, -39}, Score{-23, -11}, // Knights
		Score{3, -19}, Score{13, -25}, Score{44, -37}},
	{Score{-14, -289}, Score{-21, -171}, Score{8, -75}, Score{8, -34}, Score{28, -24}, Score{39, 14}, // Bishops
		Score{51, 33}, Score{51, 43}, Score{53, 60}, Score{54, 47}, Score{57, 43}, Score{34, 40},
		Score{90, 13}, Score{0, -49}},
	{Score{-55, -76}, Score{-47, -84}, Score{-22, 0}, Score{-12, 48}, Score{-1, 75}, Score{-10, 98}, // Rooks
		Score{3, 118}, Score{15, 126}, Score{18, 124}, Score{25, 139}, Score{28, 148}, Score{53, 145},
		Score{90, 138}, Score{34, 158}, Score{170, 75}},
	{Score{-39, -36}, Score{-21, -15}, Score{-17, 0}, Score{-5, 0}, Score{2, -3}, Score{4, -110}, // Queens
		Score{11, 29}, Score{20, 68}, Score{33, 61}, Score{51, 23}, Score{56, 44}, Score{63, 86},
		Score{66, 113}, Score{73, 119}, Score{74, 140}, Score{75, 127}, Score{81, 138}, Score{49, 136},
		Score{87, 98}, Score{50, 140}, Score{82, 140}, Score{117, 120}, Score{106, 103}, Score{103, 76},
		Score{42, 37}, Score{0, 6}, Score{-12, -2}, Score{39, 63}},
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
var passedRank = [7]Score{Score{0, 0}, Score{29, 3}, Score{4, 1}, Score{2, 83}, Score{68, 110}, Score{121, 194}, Score{236, 347}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{Score{-45, 55}, Score{-72, 82}, Score{-117, 13}, Score{-84, -3},
	Score{-14, -10}, Score{39, 4}, Score{96, 21}, Score{39, 0},
}

var isolated = Score{-22, -23}
var doubled = Score{-31, -56}
var backward = Score{6, -6}
var backwardOpen = Score{-10, -21}

var bishopPair = Score{93, 88}

var minorBehindPawn = Score{10, 39}

var tempo = Score{21, 14}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{28, 17}, Score{72, 8}}

// this bonus only improves midScore
var pawnShieldBonus = [...]Score{Score{14, 0}, Score{2, 0}} // score for every pawn

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
