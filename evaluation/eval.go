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

var PawnValue = Score{148, 254}
var KnightValue = Score{677, 827}
var BishopValue = Score{622, 811}
var RookValue = Score{873, 1213}
var QueenValue = Score{2078, 2271}

// values from stockfish 10
var pieceScores = [...][8][4]Score{
	{},
	{},
	{ // knight
		{Score{-168, -48}, Score{-40, -128}, Score{-52, -67}, Score{-52, -90}},
		{Score{-33, -77}, Score{-99, 0}, Score{-79, -104}, Score{-34, -62}},
		{Score{-62, -91}, Score{21, -83}, Score{-43, -58}, Score{-40, -9}},
		{Score{-51, -2}, Score{32, -52}, Score{-22, -5}, Score{-3, -33}},
		{Score{60, -37}, Score{-8, -8}, Score{28, -6}, Score{28, -6}},
		{Score{18, -33}, Score{-24, -42}, Score{-12, -35}, Score{107, -21}},
		{Score{-102, 2}, Score{-123, -79}, Score{-21, -55}, Score{-1, 12}},
		{Score{-263, -110}, Score{-12, 8}, Score{-50, -41}, Score{0, 0}},
	},
	{ // Bishop
		{Score{-29, -37}, Score{-33, 2}, Score{-11, -22}, Score{3, -56}},
		{Score{11, -119}, Score{9, -27}, Score{1, -72}, Score{-11, -42}},
		{Score{-29, -50}, Score{44, -55}, Score{28, -31}, Score{-9, -3}},
		{Score{62, -59}, Score{-7, -22}, Score{7, -1}, Score{41, -1}},
		{Score{3, -57}, Score{3, -24}, Score{24, -30}, Score{88, 2}},
		{Score{17, -13}, Score{135, -62}, Score{28, -19}, Score{-5, -24}},
		{Score{10, -63}, Score{-52, 1}, Score{-1, 1}, Score{-38, -25}},
		{Score{12, -83}, Score{-78, -35}, Score{-33, -24}, Score{-1, -5}},
	},
	{ // Rook
		{Score{-52, -54}, Score{-25, -26}, Score{-9, -27}, Score{-16, -34}},
		{Score{-41, -25}, Score{-88, -39}, Score{-21, -51}, Score{-34, -46}},
		{Score{-46, -31}, Score{-44, -47}, Score{-76, -38}, Score{-48, -51}},
		{Score{-69, -17}, Score{31, 0}, Score{13, -10}, Score{-46, 4}},
		{Score{-33, 0}, Score{24, -1}, Score{13, 26}, Score{21, -11}},
		{Score{40, 15}, Score{46, -2}, Score{36, 14}, Score{77, -21}},
		{Score{-7, 5}, Score{-2, 23}, Score{25, 17}, Score{75, 24}},
		{Score{60, -1}, Score{42, 42}, Score{58, 50}, Score{101, 42}},
	},
	{ // Queen
		{Score{-52, 0}, Score{-33, -93}, Score{-30, -159}, Score{-6, -106}},
		{Score{46, -56}, Score{0, -98}, Score{13, -162}, Score{10, -118}},
		{Score{-46, -31}, Score{9, -18}, Score{-18, -43}, Score{-11, -80}},
		{Score{2, -22}, Score{0, 1}, Score{-31, 18}, Score{-57, 41}},
		{Score{39, 2}, Score{-43, 80}, Score{-3, 60}, Score{-18, 54}},
		{Score{75, -39}, Score{89, 15}, Score{59, 52}, Score{59, 0}},
		{Score{25, -46}, Score{-135, 102}, Score{-24, 8}, Score{4, 105}},
		{Score{31, -75}, Score{0, -15}, Score{0, 7}, Score{57, 28}},
	},
	{ // King
		{Score{248, -19}, Score{305, 17}, Score{229, 43}, Score{271, 41}},
		{Score{216, 69}, Score{245, 80}, Score{145, 136}, Score{132, 125}},
		{Score{168, 49}, Score{189, 109}, Score{103, 149}, Score{46, 173}},
		{Score{196, 119}, Score{185, 143}, Score{215, 162}, Score{-7, 193}},
		{Score{10, 125}, Score{205, 191}, Score{115, 191}, Score{-1, 203}},
		{Score{185, 193}, Score{184, 239}, Score{142, 262}, Score{99, 233}},
		{Score{87, 132}, Score{52, 284}, Score{192, 320}, Score{89, 269}},
		{Score{0, -139}, Score{88, 224}, Score{51, 214}, Score{0, 171}},
	},
}
var pawnScores = [7][8]Score{
	{},
	{Score{4, -10}, Score{35, -5}, Score{21, 0}, Score{3, 32}, Score{20, -7}, Score{1, 10}, Score{26, -33}, Score{-31, -16}},
	{Score{-8, -18}, Score{-5, -13}, Score{15, 5}, Score{-7, -6}, Score{-4, -19}, Score{8, -2}, Score{-3, -27}, Score{-14, -1}},
	{Score{-17, 2}, Score{-46, 34}, Score{12, -13}, Score{18, -13}, Score{34, -18}, Score{11, -10}, Score{-22, 0}, Score{-29, 4}},
	{Score{-40, 30}, Score{15, 18}, Score{20, -38}, Score{24, -31}, Score{22, -49}, Score{43, -12}, Score{-19, -13}, Score{-6, 43}},
	{Score{38, 60}, Score{93, 50}, Score{-15, 12}, Score{27, 51}, Score{151, -30}, Score{60, 38}, Score{31, 133}, Score{1, 53}},
	{Score{-10, 86}, Score{68, 72}, Score{-27, 63}, Score{48, 69}, Score{68, 84}, Score{23, 113}, Score{-50, 61}, Score{-67, 60}},
}

var pawnsConnected = [8][4]Score{
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
	{Score{2, -15}, Score{-1, 2}, Score{-2, 31}, Score{17, 22}},
	{Score{15, -16}, Score{33, 8}, Score{6, 13}, Score{28, 43}},
	{Score{-4, 13}, Score{16, -7}, Score{9, 4}, Score{28, 18}},
	{Score{46, -24}, Score{59, 30}, Score{27, 36}, Score{41, 49}},
	{Score{124, 60}, Score{0, -18}, Score{92, 61}, Score{181, -22}},
	{Score{94, 0}, Score{200, 3}, Score{227, -2}, Score{231, 0}},
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
}

var mobilityBonus = [...][32]Score{
	{Score{-56, -81}, Score{-55, -199}, Score{-46, -115}, Score{-48, -62}, Score{-29, -28}, Score{-13, -27}, // Knights
		Score{0, 5}, Score{4, -3}, Score{12, -43}},
	{Score{-25, -79}, Score{-23, -135}, Score{4, -83}, Score{16, -59}, Score{28, 1}, Score{33, 30}, // Bishops
		Score{35, 40}, Score{46, 37}, Score{50, 48}, Score{49, 43}, Score{73, 34}, Score{90, 38},
		Score{101, 0}, Score{67, 0}},
	{Score{-58, -76}, Score{-49, -18}, Score{-15, -2}, Score{-6, 55}, Score{-6, 61}, Score{-3, 100}, // Rooks
		Score{7, 123}, Score{10, 137}, Score{6, 145}, Score{20, 150}, Score{26, 151}, Score{72, 141},
		Score{77, 154}, Score{60, 145}, Score{122, 85}},
	{Score{-39, -36}, Score{11, 2}, Score{-9, 0}, Score{-4, 18}, Score{-11, -1}, Score{-5, 2}, // Queens
		Score{5, 6}, Score{26, -54}, Score{30, 45}, Score{47, 43}, Score{65, 61}, Score{62, 72},
		Score{60, 69}, Score{70, 120}, Score{64, 123}, Score{68, 142}, Score{61, 167}, Score{64, 149},
		Score{104, 140}, Score{57, 143}, Score{88, 148}, Score{59, 163}, Score{70, 112}, Score{90, 155},
		Score{103, 102}, Score{36, 187}, Score{13, 139}, Score{41, 211}},
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
var passedRank = [7]Score{Score{0, 0}, Score{37, 14}, Score{-4, 36}, Score{2, 89}, Score{73, 143}, Score{122, 211}, Score{207, 346}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{Score{-66, 45}, Score{-81, 63}, Score{-77, 6}, Score{-46, -15},
	Score{-45, -15}, Score{27, 2}, Score{66, -10}, Score{15, 12},
}

var isolated = Score{-35, -16}
var doubled = Score{-14, -67}
var backward = Score{2, 8}
var backwardOpen = Score{-10, -10}

var bishopPair = Score{65, 116}

var minorBehindPawn = Score{7, 59}

var tempo = Score{19, 14}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{20, 15}, Score{100, -16}}

// this bonus only improves midScore
var pawnShieldBonus = [...]Score{Score{1, 0}, Score{12, 0}} // score for every pawn

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
