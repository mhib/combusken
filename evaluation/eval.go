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

var PawnValue = Score{148, 255}
var KnightValue = Score{693, 836}
var BishopValue = Score{646, 806}
var RookValue = Score{897, 1224}
var QueenValue = Score{2115, 2287}

// values from stockfish 10
var pieceScores = [...][8][4]Score{
	{},
	{},
	{ // knight
		{Score{-166, -63}, Score{-37, -129}, Score{-50, -65}, Score{-47, -96}},
		{Score{-34, -71}, Score{-97, 0}, Score{-75, -113}, Score{-31, -65}},
		{Score{-58, -89}, Score{22, -82}, Score{-41, -61}, Score{-37, -13}},
		{Score{-46, -2}, Score{37, -51}, Score{-19, -4}, Score{1, -36}},
		{Score{69, -42}, Score{-5, -12}, Score{31, -8}, Score{31, -9}},
		{Score{15, -19}, Score{-23, -43}, Score{-6, -38}, Score{107, -24}},
		{Score{-90, 0}, Score{-126, -82}, Score{-22, -57}, Score{0, 11}},
		{Score{-265, -99}, Score{-13, 9}, Score{-50, -43}, Score{0, -9}},
	},
	{ // Bishop
		{Score{-28, -45}, Score{-31, 2}, Score{-12, -23}, Score{0, -52}},
		{Score{8, -112}, Score{9, -25}, Score{-1, -70}, Score{-12, -44}},
		{Score{-27, -51}, Score{43, -57}, Score{23, -26}, Score{-11, -5}},
		{Score{61, -59}, Score{-12, -21}, Score{7, -3}, Score{40, 0}},
		{Score{1, -55}, Score{1, -20}, Score{23, -29}, Score{84, 3}},
		{Score{15, -11}, Score{125, -55}, Score{25, -15}, Score{-2, -27}},
		{Score{11, -66}, Score{-54, 5}, Score{-1, 2}, Score{-36, -31}},
		{Score{3, -82}, Score{-73, -39}, Score{-19, -29}, Score{0, -9}},
	},
	{ // Rook
		{Score{-54, -54}, Score{-26, -30}, Score{-8, -35}, Score{-16, -34}},
		{Score{-42, -25}, Score{-89, -39}, Score{-18, -55}, Score{-33, -48}},
		{Score{-49, -30}, Score{-44, -46}, Score{-74, -43}, Score{-48, -51}},
		{Score{-72, -16}, Score{32, 0}, Score{12, -12}, Score{-51, 7}},
		{Score{-35, 1}, Score{23, -3}, Score{5, 28}, Score{22, -9}},
		{Score{36, 15}, Score{45, -2}, Score{33, 11}, Score{79, -25}},
		{Score{-10, 1}, Score{-4, 23}, Score{23, 19}, Score{73, 24}},
		{Score{44, -1}, Score{47, 39}, Score{57, 48}, Score{99, 40}},
	},
	{ // Queen
		{Score{-50, 0}, Score{-31, -89}, Score{-27, -161}, Score{-4, -121}},
		{Score{44, -56}, Score{2, -99}, Score{15, -167}, Score{11, -121}},
		{Score{-51, -20}, Score{11, -18}, Score{-19, -35}, Score{-11, -78}},
		{Score{1, -22}, Score{0, 1}, Score{-27, 11}, Score{-60, 47}},
		{Score{39, 1}, Score{-38, 75}, Score{-1, 58}, Score{-16, 53}},
		{Score{74, -38}, Score{88, 14}, Score{62, 52}, Score{59, 0}},
		{Score{24, -46}, Score{-131, 86}, Score{-27, 9}, Score{13, 89}},
		{Score{31, -74}, Score{0, -16}, Score{3, 5}, Score{57, 28}},
	},
	{ // King
		{Score{241, -11}, Score{301, 25}, Score{225, 45}, Score{269, 43}},
		{Score{220, 71}, Score{249, 81}, Score{145, 137}, Score{133, 126}},
		{Score{170, 50}, Score{190, 109}, Score{103, 150}, Score{48, 176}},
		{Score{214, 111}, Score{190, 142}, Score{201, 165}, Score{-6, 193}},
		{Score{12, 130}, Score{176, 198}, Score{150, 190}, Score{0, 205}},
		{Score{186, 193}, Score{193, 238}, Score{145, 263}, Score{98, 235}},
		{Score{88, 136}, Score{-5, 294}, Score{193, 320}, Score{25, 272}},
		{Score{0, -155}, Score{87, 223}, Score{50, 211}, Score{1, 169}},
	},
}

var pawnScores = [7][8]Score{
	{},
	{Score{0, -6}, Score{31, -5}, Score{23, -1}, Score{-1, 33}, Score{17, -9}, Score{-1, 10}, Score{28, -37}, Score{-33, -14}},
	{Score{-12, -10}, Score{-8, -10}, Score{15, 5}, Score{-6, -6}, Score{-6, -15}, Score{6, 0}, Score{-4, -30}, Score{-16, -1}},
	{Score{-21, 5}, Score{-46, 33}, Score{14, -16}, Score{14, -13}, Score{29, -14}, Score{13, -12}, Score{-16, -1}, Score{-29, 4}},
	{Score{-42, 34}, Score{15, 14}, Score{15, -38}, Score{21, -33}, Score{23, -53}, Score{46, -12}, Score{-17, -17}, Score{-7, 45}},
	{Score{40, 57}, Score{91, 48}, Score{-21, 11}, Score{23, 51}, Score{148, -37}, Score{54, 36}, Score{30, 132}, Score{1, 54}},
	{Score{-11, 82}, Score{70, 65}, Score{-1, 44}, Score{0, 77}, Score{62, 71}, Score{-8, 116}, Score{-54, 56}, Score{-67, 53}},
}

var pawnsConnected = [8][4]Score{
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
	{Score{4, -15}, Score{-3, 5}, Score{1, 21}, Score{15, 44}},
	{Score{15, -18}, Score{37, 3}, Score{8, 5}, Score{27, 43}},
	{Score{-4, 11}, Score{16, -12}, Score{8, 4}, Score{32, 15}},
	{Score{39, -25}, Score{56, 32}, Score{31, 29}, Score{42, 49}},
	{Score{123, 61}, Score{-8, -11}, Score{79, 71}, Score{179, -18}},
	{Score{40, 0}, Score{200, 7}, Score{227, -4}, Score{216, -12}},
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
}

var mobilityBonus = [...][32]Score{
	{Score{-51, -82}, Score{-51, -194}, Score{-44, -123}, Score{-45, -76}, Score{-27, -34}, Score{-10, -32}, // Knights
		Score{2, 3}, Score{9, -7}, Score{14, -47}},
	{Score{-22, -61}, Score{-24, -124}, Score{4, -75}, Score{14, -58}, Score{28, 3}, Score{34, 32}, // Bishops
		Score{35, 43}, Score{45, 45}, Score{50, 54}, Score{51, 43}, Score{75, 36}, Score{84, 44},
		Score{96, 5}, Score{95, 0}},
	{Score{-58, -76}, Score{-49, -19}, Score{-15, -1}, Score{-6, 54}, Score{-5, 57}, Score{-2, 90}, // Rooks
		Score{7, 118}, Score{10, 134}, Score{6, 139}, Score{18, 150}, Score{24, 148}, Score{70, 137},
		Score{77, 150}, Score{65, 137}, Score{122, 83}},
	{Score{-39, -36}, Score{11, 0}, Score{-9, 0}, Score{-5, 17}, Score{-10, -1}, Score{-5, 0}, // Queens
		Score{5, 1}, Score{25, -40}, Score{29, 45}, Score{45, 60}, Score{65, 68}, Score{59, 85},
		Score{58, 89}, Score{69, 120}, Score{67, 123}, Score{62, 162}, Score{64, 167}, Score{70, 137},
		Score{98, 141}, Score{66, 143}, Score{88, 148}, Score{62, 165}, Score{69, 120}, Score{79, 172},
		Score{98, 118}, Score{41, 188}, Score{104, 109}, Score{42, 208}},
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
var passedRank = [7]Score{Score{0, 0}, Score{38, 11}, Score{-4, 36}, Score{11, 84}, Score{74, 142}, Score{117, 214}, Score{239, 350}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{Score{-67, 45}, Score{-81, 61}, Score{-72, 8}, Score{-46, -15},
	Score{-45, -16}, Score{21, 5}, Score{65, -10}, Score{14, 11},
}

var isolated = Score{-35, -18}
var doubled = Score{-11, -68}
var backward = Score{0, 8}
var backwardOpen = Score{-11, -10}

var bishopPair = Score{63, 121}

var minorBehindPawn = Score{7, 57}

var tempo = Score{19, 14}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{18, 18}, Score{100, -16}}

// this bonus only improves midScore
var pawnShieldBonus = [...]Score{Score{2, 0}, Score{11, 0}} // score for every pawn

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
