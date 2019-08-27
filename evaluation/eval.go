package evaluation

import . "github.com/mhib/combusken/backend"
import . "github.com/mhib/combusken/utils"
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

func (first *Score) Add(second Score) {
	first.Middle += second.Middle
	first.End += second.End
}

func (first *Score) Subtract(second Score) {
	first.Middle -= second.Middle
	first.End -= second.End
}

func addScore(first, second Score) Score {
	return Score{
		Middle: first.Middle + second.Middle,
		End:    first.End + second.End,
	}
}

var PawnValue = Score{166, 207}
var KnightValue = Score{841, 749}
var BishopValue = Score{758, 740}
var RookValue = Score{1090, 1192}
var QueenValue = Score{2363, 2339}

// Piece Square Values
var pieceScores = [7][8][4]Score{
	{},
	{},
	{ // knight
		{Score{-125, -44}, Score{-27, -89}, Score{-77, -60}, Score{-25, -42}},
		{Score{-43, -99}, Score{-59, -62}, Score{-30, -55}, Score{-11, -40}},
		{Score{-44, -71}, Score{-9, -57}, Score{-10, -34}, Score{-5, -4}},
		{Score{-29, -48}, Score{30, -43}, Score{21, 1}, Score{23, 7}},
		{Score{1, -59}, Score{17, -31}, Score{32, 4}, Score{65, 5}},
		{Score{-45, -91}, Score{44, -75}, Score{3, -2}, Score{60, -8}},
		{Score{-109, -89}, Score{-67, -40}, Score{94, -88}, Score{1, -34}},
		{Score{-299, -119}, Score{-74, -132}, Score{-149, -67}, Score{0, -90}},
	},
	{ // Bishop
		{Score{-21, -25}, Score{7, -9}, Score{8, -16}, Score{6, -8}},
		{Score{-1, -47}, Score{63, -54}, Score{44, -33}, Score{6, -12}},
		{Score{16, -35}, Score{43, -22}, Score{43, -6}, Score{20, 8}},
		{Score{0, -37}, Score{9, -28}, Score{14, -2}, Score{50, 2}},
		{Score{-29, -24}, Score{1, -25}, Score{16, -6}, Score{42, 5}},
		{Score{-82, -5}, Score{-4, -28}, Score{0, -7}, Score{-17, -16}},
		{Score{-94, -11}, Score{31, -20}, Score{-17, 0}, Score{3, -22}},
		{Score{-9, -43}, Score{-39, -37}, Score{-96, -31}, Score{-91, -21}},
	},
	{ // Rook
		{Score{-16, -35}, Score{-32, -14}, Score{2, -21}, Score{18, -32}},
		{Score{-81, 0}, Score{-12, -31}, Score{-23, -18}, Score{0, -23}},
		{Score{-68, -12}, Score{-26, -15}, Score{-21, -23}, Score{-17, -24}},
		{Score{-66, 3}, Score{-19, -3}, Score{-26, 3}, Score{-16, -5}},
		{Score{-50, 11}, Score{-34, 3}, Score{17, 11}, Score{4, -2}},
		{Score{-26, 4}, Score{38, 1}, Score{39, -5}, Score{20, -1}},
		{Score{32, 11}, Score{12, 25}, Score{82, 3}, Score{99, -10}},
		{Score{6, 18}, Score{12, 13}, Score{-36, 29}, Score{24, 18}},
	},
	{ // Queen
		{Score{-10, -102}, Score{2, -107}, Score{10, -102}, Score{45, -129}},
		{Score{-8, -100}, Score{-2, -79}, Score{49, -108}, Score{37, -75}},
		{Score{0, -34}, Score{30, -55}, Score{0, 7}, Score{0, -4}},
		{Score{0, -28}, Score{-19, 38}, Score{-6, 36}, Score{-27, 83}},
		{Score{-8, 0}, Score{-42, 50}, Score{-20, 42}, Score{-51, 106}},
		{Score{45, -45}, Score{10, -11}, Score{16, 21}, Score{0, 81}},
		{Score{1, -43}, Score{-81, 31}, Score{0, 18}, Score{-30, 86}},
		{Score{4, -40}, Score{0, -3}, Score{28, 9}, Score{26, 23}},
	},
	{ // King
		{Score{356, -34}, Score{347, 23}, Score{208, 98}, Score{223, 76}},
		{Score{313, 41}, Score{254, 78}, Score{123, 136}, Score{70, 157}},
		{Score{140, 83}, Score{137, 112}, Score{68, 145}, Score{2, 171}},
		{Score{7, 94}, Score{86, 111}, Score{4, 166}, Score{-11, 176}},
		{Score{5, 121}, Score{123, 151}, Score{103, 171}, Score{9, 180}},
		{Score{119, 135}, Score{206, 166}, Score{178, 189}, Score{13, 166}},
		{Score{64, 131}, Score{34, 181}, Score{39, 205}, Score{29, 179}},
		{Score{46, 1}, Score{23, 99}, Score{5, 138}, Score{0, 100}},
	},
}

// Pawns Square scores
var pawnScores = [7][8]Score{
	{},
	{Score{-13, -6}, Score{32, -11}, Score{2, 11}, Score{26, 0}, Score{20, 8}, Score{2, 14}, Score{36, -5}, Score{-21, -5}},
	{Score{-7, -24}, Score{-14, -11}, Score{6, -5}, Score{9, -10}, Score{-3, 4}, Score{7, -7}, Score{-18, -14}, Score{-7, -20}},
	{Score{-30, -5}, Score{-12, -8}, Score{14, -15}, Score{34, -18}, Score{22, -10}, Score{18, -13}, Score{-9, -8}, Score{-24, -3}},
	{Score{-3, 19}, Score{40, -5}, Score{27, -15}, Score{57, -22}, Score{54, -33}, Score{13, -2}, Score{40, -3}, Score{-8, 20}},
	{Score{9, 67}, Score{30, 55}, Score{71, 14}, Score{69, 6}, Score{74, -1}, Score{107, 22}, Score{12, 57}, Score{13, 74}},
	{Score{-1, 117}, Score{6, 113}, Score{0, 59}, Score{0, 72}, Score{41, 79}, Score{-19, 78}, Score{0, 82}, Score{-120, 136}},
}

var pawnsConnected = [8][4]Score{
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
	{Score{19, -32}, Score{8, 2}, Score{11, -9}, Score{2, 19}},
	{Score{9, 4}, Score{58, 2}, Score{22, 11}, Score{27, 28}},
	{Score{23, 11}, Score{40, 12}, Score{40, 12}, Score{45, 16}},
	{Score{21, 24}, Score{18, 41}, Score{47, 39}, Score{63, 32}},
	{Score{3, 93}, Score{55, 86}, Score{97, 93}, Score{128, 72}},
	{Score{2, 101}, Score{121, 0}, Score{152, 0}, Score{0, 48}},
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
}

var mobilityBonus = [...][32]Score{
	{Score{-76, -193}, Score{-63, -121}, Score{-42, -75}, Score{-44, -32}, Score{-16, -34}, Score{1, -19}, // Knights
		Score{18, -24}, Score{38, -28}, Score{60, -53}},
	{Score{-55, -134}, Score{-23, -103}, Score{11, -48}, Score{20, -17}, Score{38, 4}, Score{58, 12}, // Bishops
		Score{71, 16}, Score{74, 22}, Score{84, 26}, Score{94, 22}, Score{117, 0}, Score{144, 7},
		Score{77, 37}, Score{79, 20}},
	{Score{-42, -64}, Score{-67, -55}, Score{-43, 11}, Score{-24, 55}, Score{-13, 82}, Score{-5, 103}, // Rooks
		Score{1, 118}, Score{17, 123}, Score{20, 121}, Score{47, 124}, Score{62, 125}, Score{74, 127},
		Score{84, 131}, Score{109, 122}, Score{159, 105}},
	{Score{-39, -36}, Score{-21, -15}, Score{-20, 0}, Score{-26, -7}, Score{-6, 0}, Score{6, -43}, // Queens
		Score{5, -26}, Score{24, -2}, Score{31, 24}, Score{38, 21}, Score{41, 59}, Score{41, 74},
		Score{52, 62}, Score{52, 105}, Score{55, 111}, Score{52, 121}, Score{51, 128}, Score{45, 131},
		Score{63, 117}, Score{62, 146}, Score{102, 105}, Score{112, 98}, Score{114, 86}, Score{100, 67},
		Score{113, 40}, Score{32, 63}, Score{0, 0}, Score{1, 6}},
}

var passedFriendlyDistance = [8]Score{
	Score{0, 0}, Score{14, 31}, Score{-11, 10}, Score{-17, -23},
	Score{-31, -38}, Score{-32, -40}, Score{1, -49}, Score{-57, -28},
}

var passedEnemyDistance = [8]Score{
	Score{0, 0}, Score{-97, -124}, Score{38, -48}, Score{24, 11},
	Score{32, 43}, Score{18, 62}, Score{18, 66}, Score{-23, 87},
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

// PassedRank[Rank] contains a bonus according to the rank of a passed pawn
var passedRank = [7]Score{Score{0, 0}, Score{-10, -46}, Score{-13, -13}, Score{-7, 56}, Score{47, 130}, Score{70, 275}, Score{188, 431}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{Score{-17, 40}, Score{-48, 41}, Score{-55, 20}, Score{-52, -6},
	Score{-30, -6}, Score{17, -5}, Score{-15, 30}, Score{-9, 18},
}

var isolated = Score{-18, -19}
var doubled = Score{-23, -57}
var backward = Score{8, -5}
var backwardOpen = Score{-29, -10}

var bishopPair = Score{97, 103}
var bishopRammedPawns = Score{-11, -24}

var bishopOutpostUndefendedBonus = Score{65, -6}
var bishopOutpostDefendedBonus = Score{133, 1}

var knightOutpostUndefendedBonus = Score{50, -30}
var knightOutpostDefendedBonus = Score{90, 16}

var minorBehindPawn = Score{7, 50}

var tempo = Score{45, 50}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{18, 39}, Score{93, -8}}

var kingDefenders = [12]Score{
	Score{-147, -1}, Score{-111, -8}, Score{-56, -3}, Score{-16, 3},
	Score{8, 6}, Score{34, 6}, Score{58, 3}, Score{67, 9},
	Score{88, 3}, Score{71, 10}, Score{12, 6}, Score{12, 6},
}

var kingShelter = [2][8][8]Score{
	{{Score{-11, 0}, Score{-6, -17}, Score{0, 2}, Score{18, -1},
		Score{0, -24}, Score{0, -1}, Score{1, -30}, Score{-48, 19}},
		{Score{18, 2}, Score{51, -21}, Score{-4, -3}, Score{-2, 4},
			Score{-32, -2}, Score{8, -13}, Score{23, -64}, Score{-51, 6}},
		{Score{14, 12}, Score{0, 1}, Score{-29, 6}, Score{-12, 4},
			Score{-48, -3}, Score{-23, 0}, Score{-23, -1}, Score{-27, 0}},
		{Score{0, 26}, Score{4, 11}, Score{-30, -2}, Score{-9, -2},
			Score{11, -36}, Score{-17, -7}, Score{0, -22}, Score{-50, 7}},
		{Score{0, 9}, Score{-2, 3}, Score{-37, -1}, Score{-34, 13},
			Score{-22, -14}, Score{-57, 5}, Score{-55, 0}, Score{-45, 11}},
		{Score{78, -19}, Score{56, -33}, Score{-14, -18}, Score{-1, -14},
			Score{2, -38}, Score{-5, -12}, Score{48, -40}, Score{-21, -1}},
		{Score{42, -6}, Score{0, -11}, Score{-41, -10}, Score{-13, -8},
			Score{-20, -17}, Score{8, -4}, Score{3, -29}, Score{-58, 19}},
		{Score{-10, -3}, Score{-35, -5}, Score{-25, 12}, Score{-35, 17},
			Score{-20, 10}, Score{-56, 35}, Score{-91, 17}, Score{-91, 53}}},
	{{Score{0, 0}, Score{-55, -24}, Score{0, -16}, Score{-42, -51},
		Score{-22, -10}, Score{-2, -22}, Score{-138, 0}, Score{-74, 16}},
		{Score{0, 0}, Score{-1, -26}, Score{-28, -12}, Score{-13, 0},
			Score{0, -9}, Score{-19, -30}, Score{-31, -8}, Score{-97, 22}},
		{Score{0, 5}, Score{67, -11}, Score{14, -3}, Score{10, -8},
			Score{15, 1}, Score{-45, -14}, Score{64, -9}, Score{-46, 9}},
		{Score{0, 1}, Score{-21, 18}, Score{-32, 14}, Score{-48, 0},
			Score{-26, 11}, Score{-98, 29}, Score{-39, -3}, Score{-72, 2}},
		{Score{0, 65}, Score{1, 9}, Score{-8, 4}, Score{-11, 1},
			Score{-15, 9}, Score{-18, -10}, Score{-20, -6}, Score{-62, 14}},
		{Score{0, 0}, Score{39, -17}, Score{-18, 1}, Score{-12, -13},
			Score{0, -13}, Score{-33, -23}, Score{9, -35}, Score{-50, 6}},
		{Score{0, 0}, Score{14, -22}, Score{2, -22}, Score{-35, -9},
			Score{-33, -21}, Score{-9, -32}, Score{3, -31}, Score{-119, 33}},
		{Score{0, 0}, Score{-3, -40}, Score{-10, -22}, Score{-42, -13},
			Score{-49, -8}, Score{-8, -10}, Score{-73, -34}, Score{-105, 42}}},
}

var kingStorm = [2][4][8]Score{
	{{Score{34, 6}, Score{27, 1}, Score{27, 8}, Score{15, 8},
		Score{1, 13}, Score{1, 12}, Score{-25, 31}, Score{8, -14}},
		{Score{30, 6}, Score{25, 3}, Score{38, 1}, Score{9, 14},
			Score{17, 9}, Score{10, 3}, Score{-2, -2}, Score{-3, -10}},
		{Score{27, 26}, Score{11, 13}, Score{1, 18}, Score{-6, 19},
			Score{-6, 15}, Score{8, 1}, Score{5, -14}, Score{3, -5}},
		{Score{59, 21}, Score{15, 7}, Score{9, 6}, Score{-5, 11},
			Score{-8, 17}, Score{3, 15}, Score{0, 8}, Score{-3, 2}}},
	{{Score{0, 0}, Score{0, 1}, Score{-30, 5}, Score{35, -14},
		Score{30, 12}, Score{-3, 2}, Score{2, 1}, Score{30, -44}},
		{Score{0, 0}, Score{0, -33}, Score{-5, -12}, Score{76, -10},
			Score{39, -1}, Score{-7, -2}, Score{0, 0}, Score{12, -35}},
		{Score{0, 0}, Score{-25, -14}, Score{-56, -1}, Score{13, 1},
			Score{1, -2}, Score{0, -8}, Score{3, -11}, Score{4, -9}},
		{Score{0, 0}, Score{0, -17}, Score{5, -21}, Score{-11, 2},
			Score{-8, 0}, Score{-1, -5}, Score{0, 0}, Score{-7, 1}}},
}

var blackPassedMask [64]uint64
var whitePassedMask [64]uint64

var whiteOutpostMask [64]uint64
var blackOutpostMask [64]uint64

var distanceBetween [64][64]int16

var adjacentFilesMask [8]uint64

var whiteKingAreaMask [64]uint64
var blackKingAreaMask [64]uint64

var whiteForwardRanksMask [8]uint64
var blackForwardRanksMasks [8]uint64

// King shield bitboards
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

// Outpost bitboards
const whiteOutpustRanks = RANK_4_BB | RANK_5_BB | RANK_6_BB
const blackOutpustRanks = RANK_5_BB | RANK_4_BB | RANK_3_BB

var kingSafetyAttacksWeights = [King + 1]int16{0, 0, -13, -10, 8, 75, 0}
var kingSafetyAttackValue int16 = 92
var kingSafetyWeakSquares int16 = 22
var kingSafetyFriendlyPawns int16 = 1
var kingSafetyNoEnemyQueens int16 = -127
var kingSafetySafeQueenCheck int16 = 105
var kingSafetySafeRookCheck int16 = 134
var kingSafetySafeBishopCheck int16 = 136
var kingSafetySafeKnightCheck int16 = 193
var kingSafetyAdjustment int16 = -18

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

	// Pawn is passed if no pawn of opposite color can stop it from promoting
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
	// Outpust is similar to passed bitboard bot we do not care about pawns in same file
	for i := 8; i <= 55; i++ {
		whiteOutpostMask[i] = whitePassedMask[i] & ^FILES[File(i)]
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
	for i := 55; i >= 8; i-- {
		blackOutpostMask[i] = blackPassedMask[i] & ^FILES[File(i)]
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

	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			distanceBetween[y][x] = int16(Max(Abs(Rank(y)-Rank(x)), Abs(File(y)-File(x))))
		}
	}

	for y := 0; y < 64; y++ {
		whiteKingAreaMask[y] = KingAttacks[y] | SquareMask[y] | North(KingAttacks[y])
		blackKingAreaMask[y] = KingAttacks[y] | SquareMask[y] | South(KingAttacks[y])
		if File(y) > FILE_A {
			whiteKingAreaMask[y] |= West(whiteKingAreaMask[y])
			blackKingAreaMask[y] |= West(blackKingAreaMask[y])
		}
		if File(y) < FILE_H {
			whiteKingAreaMask[y] |= East(whiteKingAreaMask[y])
			blackKingAreaMask[y] |= East(blackKingAreaMask[y])
		}
	}

	for rank := RANK_1; rank <= RANK_8; rank++ {
		for y := rank; y <= RANK_8; y++ {
			whiteForwardRanksMask[rank] |= RANKS[y]
		}
		blackForwardRanksMasks[rank] = (^whiteForwardRanksMask[rank]) | RANKS[rank]
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

func evaluateKingPawns(pos *Position, pkTable PawnKingTable) (int, int) {
	if ok, midScore, endScore := pkTable.Get(pos.PawnKey); ok {
		return midScore, endScore
	}
	var fromBB uint64
	var fromId int
	whiteKingLocation := BitScan(pos.Kings & pos.White)
	blackKingLocation := BitScan(pos.Kings & pos.Black)
	midResult := 0
	endResult := 0

	// white pawns
	for fromBB = pos.Pawns & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		midResult += int(whitePawnsPos[fromId].Middle)
		endResult += int(whitePawnsPos[fromId].End)

		// Passed bonus
		if whitePassedMask[fromId]&(pos.Pawns&pos.Black) == 0 {
			// Bonus is calculated based on rank, file, distance from friendly and enemy king
			midResult += int(
				passedRank[Rank(fromId)].Middle +
					passedFile[File(fromId)].Middle +
					passedFriendlyDistance[distanceBetween[whiteKingLocation][fromId]].Middle +
					passedEnemyDistance[distanceBetween[blackKingLocation][fromId]].Middle,
			)
			endResult += int(
				passedRank[Rank(fromId)].End +
					passedFile[File(fromId)].End +
					passedFriendlyDistance[distanceBetween[whiteKingLocation][fromId]].End +
					passedEnemyDistance[distanceBetween[blackKingLocation][fromId]].End,
			)
		}
		// Isolated pawn penalty
		if adjacentFilesMask[File(fromId)]&(pos.Pawns&pos.White) == 0 {
			midResult += int(isolated.Middle)
			endResult += int(isolated.End)
		}

		// Pawn is backward if there are no pawns behind it and cannot increase rank without being attacked by enemy pawn
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
	}

	// white doubled pawns
	doubledCount := PopCount(pos.Pawns & pos.White & South(pos.Pawns&pos.White))
	midResult += doubledCount * int(doubled.Middle)
	endResult += doubledCount * int(doubled.End)

	// black pawns
	for fromBB = pos.Pawns & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		midResult -= int(blackPawnsPos[fromId].Middle)
		endResult -= int(blackPawnsPos[fromId].End)
		if blackPassedMask[fromId]&(pos.Pawns&pos.White) == 0 {
			midResult -= int(
				passedRank[7-Rank(fromId)].Middle +
					passedFile[File(fromId)].Middle +
					passedFriendlyDistance[distanceBetween[blackKingLocation][fromId]].Middle +
					passedEnemyDistance[distanceBetween[whiteKingLocation][fromId]].Middle,
			)
			endResult -= int(
				passedRank[7-Rank(fromId)].End +
					passedFile[File(fromId)].End +
					passedFriendlyDistance[distanceBetween[blackKingLocation][fromId]].End +
					passedEnemyDistance[distanceBetween[whiteKingLocation][fromId]].End,
			)
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
	}

	// black doubled pawns
	doubledCount = PopCount(pos.Pawns & pos.Black & North(pos.Pawns&pos.Black))
	midResult -= doubledCount * int(doubled.Middle)
	endResult -= doubledCount * int(doubled.End)

	// White king storm shelter
	for file := Max(File(whiteKingLocation)-1, FILE_A); file <= Min(File(whiteKingLocation)+1, FILE_H); file++ {
		ours := pos.Pawns & FILES[file] & pos.White & whiteForwardRanksMask[Rank(whiteKingLocation)]
		var ourDist int
		if ours == 0 {
			ourDist = 7
		} else {
			ourDist = Abs(Rank(whiteKingLocation) - Rank(BitScan(ours)))
		}
		theirs := pos.Pawns & FILES[file] & pos.Black & whiteForwardRanksMask[Rank(whiteKingLocation)]
		var theirDist int
		if theirs == 0 {
			theirDist = 7
		} else {
			theirDist = Abs(Rank(whiteKingLocation) - Rank(BitScan(theirs)))
		}
		sameFile := BoolToInt(file == File(whiteKingLocation))
		midResult += int(kingShelter[sameFile][file][ourDist].Middle)
		endResult += int(kingShelter[sameFile][file][ourDist].End)

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		midResult += int(kingStorm[blocked][FileMirror[file]][theirDist].Middle)
		endResult += int(kingStorm[blocked][FileMirror[file]][theirDist].End)
	}

	// Black king storm / shelter
	for file := Max(File(blackKingLocation)-1, FILE_A); file <= Min(File(blackKingLocation)+1, FILE_H); file++ {
		ours := pos.Pawns & FILES[file] & pos.Black & blackForwardRanksMasks[Rank(blackKingLocation)]
		var ourDist int
		if ours == 0 {
			ourDist = 7
		} else {
			ourDist = Abs(Rank(blackKingLocation) - Rank(MostSignificantBit(ours)))
		}
		theirs := pos.Pawns & FILES[file] & pos.White & blackForwardRanksMasks[Rank(blackKingLocation)]
		var theirDist int
		if theirs == 0 {
			theirDist = 7
		} else {
			theirDist = Abs(Rank(blackKingLocation) - Rank(MostSignificantBit(theirs)))
		}
		sameFile := BoolToInt(file == File(blackKingLocation))
		midResult -= int(kingShelter[sameFile][file][ourDist].Middle)
		endResult -= int(kingShelter[sameFile][file][ourDist].End)

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		midResult -= int(kingStorm[blocked][FileMirror[file]][theirDist].Middle)
		endResult -= int(kingStorm[blocked][FileMirror[file]][theirDist].End)
	}
	pkTable.Set(pos.PawnKey, midResult, endResult)
	return midResult, endResult
}

func Evaluate(pos *Position, pkTable PawnKingTable) int {
	var fromId int
	var fromBB uint64
	var attacks uint64

	var whiteAttacked uint64
	var whiteAttackedBy [King + 1]uint64
	var whiteAttackedByTwo uint64
	var blackAttacked uint64
	var whiteKingAttacksCount int16
	var whiteKingAttackersCount int16
	var whiteKingAttackersWeight int16
	var blackAttackedBy [King + 1]uint64
	var blackAttackedByTwo uint64
	var blackKingAttacksCount int16
	var blackKingAttackersCount int16
	var blackKingAttackersWeight int16

	phase := totalPhase
	whiteMobilityArea := ^((pos.Pawns & pos.White) | (BlackPawnsAttacks(pos.Pawns & pos.Black)))
	blackMobilityArea := ^((pos.Pawns & pos.Black) | (WhitePawnsAttacks(pos.Pawns & pos.White)))
	allOccupation := pos.White | pos.Black

	whiteKingLocation := BitScan(pos.Kings & pos.White)
	attacks = KingAttacks[whiteKingLocation]
	whiteAttacked |= attacks
	//whiteAttackedBy[King] |= attacks
	whiteKingArea := whiteKingAreaMask[whiteKingLocation]

	blackKingLocation := BitScan(pos.Kings & pos.Black)
	attacks = KingAttacks[blackKingLocation]
	blackAttacked |= attacks
	//blackAttackedBy[King] |= attacks
	blackKingArea := blackKingAreaMask[blackKingLocation]

	// white pawns
	attacks = WhitePawnsAttacks(pos.Pawns & pos.White)
	whiteAttackedByTwo |= whiteAttacked & attacks
	whiteAttacked |= attacks
	//whiteAttackedBy[Pawn] |= attacks
	whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))

	// black pawns
	attacks = BlackPawnsAttacks(pos.Pawns & pos.Black)
	blackAttackedByTwo |= blackAttacked & attacks
	blackAttacked |= attacks
	//blackAttackedBy[Pawn] |= attacks
	blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))

	midResult, endResult := evaluateKingPawns(pos, pkTable)

	// white knights
	for fromBB = pos.Knights & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= knightPhase
		fromId = BitScan(fromBB)

		attacks = KnightAttacks[fromId]
		mobility := PopCount(whiteMobilityArea & attacks)
		midResult += int(whiteKnightsPos[fromId].Middle)
		endResult += int(whiteKnightsPos[fromId].End)
		midResult += int(mobilityBonus[0][mobility].Middle)
		endResult += int(mobilityBonus[0][mobility].End)

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Knight] |= attacks

		if (pos.Pawns>>8)&SquareMask[fromId] != 0 {
			midResult += int(minorBehindPawn.Middle)
			endResult += int(minorBehindPawn.End)
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && whiteOutpostMask[fromId]&(pos.Pawns&pos.Black) == 0 {
			if BlackPawnAttacks[fromId]&(pos.Pawns&pos.White) != 0 {
				midResult += int(knightOutpostDefendedBonus.Middle)
				endResult += int(knightOutpostDefendedBonus.End)
			} else {
				midResult += int(knightOutpostUndefendedBonus.Middle)
				endResult += int(knightOutpostUndefendedBonus.End)
			}
		}
		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += kingSafetyAttacksWeights[Knight]
		}
	}

	// black knights
	for fromBB = pos.Knights & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= knightPhase
		fromId = BitScan(fromBB)

		attacks = KnightAttacks[fromId]
		mobility := PopCount(blackMobilityArea & attacks)
		midResult -= int(blackKnightsPos[fromId].Middle)
		endResult -= int(blackKnightsPos[fromId].End)
		midResult -= int(mobilityBonus[0][mobility].Middle)
		endResult -= int(mobilityBonus[0][mobility].End)

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Knight] |= attacks

		if (pos.Pawns<<8)&SquareMask[fromId] != 0 {
			midResult -= int(minorBehindPawn.Middle)
			endResult -= int(minorBehindPawn.End)
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && blackOutpostMask[fromId]&(pos.Pawns&pos.White) == 0 {
			if WhitePawnAttacks[fromId]&(pos.Pawns&pos.Black) != 0 {
				midResult -= int(knightOutpostDefendedBonus.Middle)
				endResult -= int(knightOutpostDefendedBonus.End)
			} else {
				midResult -= int(knightOutpostUndefendedBonus.Middle)
				endResult -= int(knightOutpostUndefendedBonus.End)
			}
		}
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += kingSafetyAttacksWeights[Knight]
		}
	}

	// white bishops
	whiteRammedPawns := South(pos.Pawns&pos.Black) & (pos.Pawns & pos.White)
	for fromBB = pos.Bishops & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= bishopPhase
		fromId = BitScan(fromBB)

		attacks = BishopAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		midResult += int(mobilityBonus[1][mobility].Middle)
		endResult += int(mobilityBonus[1][mobility].End)
		midResult += int(whiteBishopsPos[fromId].Middle)
		endResult += int(whiteBishopsPos[fromId].End)

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Bishop] |= attacks

		if (pos.Pawns>>8)&SquareMask[fromId] != 0 {
			midResult += int(minorBehindPawn.Middle)
			endResult += int(minorBehindPawn.End)
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && whiteOutpostMask[fromId]&(pos.Pawns&pos.Black) == 0 {
			if BlackPawnAttacks[fromId]&(pos.Pawns&pos.White) != 0 {
				midResult += int(bishopOutpostDefendedBonus.Middle)
				endResult += int(bishopOutpostDefendedBonus.End)
			} else {
				midResult += int(bishopOutpostUndefendedBonus.Middle)
				endResult += int(bishopOutpostUndefendedBonus.End)
			}
		}

		// Bishop is worth less if there are friendly rammed pawns of its color
		var rammedCount int16
		if SquareMask[fromId]&WHITE_SQUARES != 0 {
			rammedCount = int16(PopCount(whiteRammedPawns & WHITE_SQUARES))
		} else {
			rammedCount = int16(PopCount(whiteRammedPawns & BLACK_SQUARES))
		}
		midResult += int(bishopRammedPawns.Middle * rammedCount)
		endResult += int(bishopRammedPawns.End * rammedCount)
		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += kingSafetyAttacksWeights[Bishop]
		}
	}

	// Bishop pair bonus
	// It is not checked if bishops have opposite colors, but that is almost always the case
	if MoreThanOne(pos.Bishops & pos.White) {
		midResult += int(bishopPair.Middle)
		endResult += int(bishopPair.End)
	}

	// black bishops
	blackRammedPawns := North(pos.Pawns&pos.White) & (pos.Pawns & pos.Black)
	for fromBB = pos.Bishops & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= bishopPhase
		fromId = BitScan(fromBB)

		attacks = BishopAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		midResult -= int(mobilityBonus[1][mobility].Middle)
		endResult -= int(mobilityBonus[1][mobility].End)
		midResult -= int(blackBishopsPos[fromId].Middle)
		endResult -= int(blackBishopsPos[fromId].End)

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Bishop] |= attacks

		if (pos.Pawns<<8)&SquareMask[fromId] != 0 {
			midResult -= int(minorBehindPawn.Middle)
			endResult -= int(minorBehindPawn.End)
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && blackOutpostMask[fromId]&(pos.Pawns&pos.White) == 0 {
			if WhitePawnAttacks[fromId]&(pos.Pawns&pos.Black) != 0 {
				midResult -= int(bishopOutpostDefendedBonus.Middle)
				endResult -= int(bishopOutpostDefendedBonus.End)
			} else {
				midResult -= int(bishopOutpostUndefendedBonus.Middle)
				endResult -= int(bishopOutpostUndefendedBonus.End)
			}
		}
		var rammedCount int16
		if SquareMask[fromId]&WHITE_SQUARES != 0 {
			rammedCount = int16(PopCount(blackRammedPawns & WHITE_SQUARES))
		} else {
			rammedCount = int16(PopCount(blackRammedPawns & BLACK_SQUARES))
		}
		midResult -= int(bishopRammedPawns.Middle * rammedCount)
		endResult -= int(bishopRammedPawns.End * rammedCount)
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += kingSafetyAttacksWeights[Bishop]
		}
	}

	if MoreThanOne(pos.Bishops & pos.Black) {
		midResult -= int(bishopPair.Middle)
		endResult -= int(bishopPair.End)
	}

	// white rooks
	for fromBB = pos.Rooks & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= rookPhase
		fromId = BitScan(fromBB)

		attacks = RookAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		midResult += int(mobilityBonus[2][mobility].Middle)
		endResult += int(mobilityBonus[2][mobility].End)
		midResult += int(whiteRooksPos[fromId].Middle)
		endResult += int(whiteRooksPos[fromId].End)

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Rook] |= attacks

		if pos.Pawns&FILES[File(fromId)] == 0 {
			midResult += int(rookOnFile[1].Middle)
			endResult += int(rookOnFile[1].End)
		} else if (pos.Pawns&pos.White)&FILES[File(fromId)] == 0 {
			midResult += int(rookOnFile[0].Middle)
			endResult += int(rookOnFile[0].End)
		}

		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += kingSafetyAttacksWeights[Rook]
		}
	}

	// black rooks
	for fromBB = pos.Rooks & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= rookPhase
		fromId = BitScan(fromBB)

		attacks = RookAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		midResult -= int(mobilityBonus[2][mobility].Middle)
		endResult -= int(mobilityBonus[2][mobility].End)
		midResult -= int(blackRooksPos[fromId].Middle)
		endResult -= int(blackRooksPos[fromId].End)

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Rook] |= attacks

		if pos.Pawns&FILES[File(fromId)] == 0 {
			midResult -= int(rookOnFile[1].Middle)
			endResult -= int(rookOnFile[1].End)
		} else if (pos.Pawns&pos.Black)&FILES[File(fromId)] == 0 {
			midResult -= int(rookOnFile[0].Middle)
			endResult -= int(rookOnFile[0].End)
		}

		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += kingSafetyAttacksWeights[Rook]
		}
	}

	//white queens
	for fromBB = pos.Queens & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= queenPhase
		fromId = BitScan(fromBB)

		attacks = QueenAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		midResult += int(mobilityBonus[3][mobility].Middle)
		endResult += int(mobilityBonus[3][mobility].End)
		midResult += int(whiteQueensPos[fromId].Middle)
		endResult += int(whiteQueensPos[fromId].End)

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Queen] |= attacks

		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += kingSafetyAttacksWeights[Queen]
		}
	}

	// black queens
	for fromBB = pos.Queens & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= queenPhase
		fromId = BitScan(fromBB)

		attacks = QueenAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		midResult -= int(mobilityBonus[3][mobility].Middle)
		endResult -= int(mobilityBonus[3][mobility].End)
		midResult -= int(blackQueensPos[fromId].Middle)
		endResult -= int(blackQueensPos[fromId].End)

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Queen] |= attacks
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += kingSafetyAttacksWeights[Queen]
		}
	}

	// tempo bonus
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

	// white king
	whiteKingDefenders := PopCount(
		(pos.Pawns | pos.Bishops | pos.Knights) & pos.White & whiteKingAreaMask[whiteKingLocation],
	)
	midResult += int(whiteKingPos[whiteKingLocation].Middle)
	endResult += int(whiteKingPos[whiteKingLocation].End)
	midResult += int(kingDefenders[whiteKingDefenders].Middle)
	midResult += int(kingDefenders[whiteKingDefenders].End)
	if int(blackKingAttackersCount) > 1-PopCount(pos.Black&pos.Queens) {

		// Weak squares are attacked by the enemy, defended no more
		// than once and only defended by our Queens or our King
		weak := blackAttacked & ^whiteAttackedByTwo & (^whiteAttacked | whiteAttackedBy[Queen] | whiteAttackedBy[King])

		safe := ^pos.Black & (^whiteAttacked | (weak & blackAttackedByTwo))

		knightThreats := KnightAttacks[whiteKingLocation]
		bishopThreats := BishopAttacks(whiteKingLocation, allOccupation)
		rookThreats := RookAttacks(whiteKingLocation, allOccupation)
		queenThreats := bishopThreats | rookThreats

		knightChecks := knightThreats & safe & blackAttackedBy[Knight]
		bishopChecks := bishopThreats & safe & blackAttackedBy[Bishop]
		rookChecks := rookThreats & safe & blackAttackedBy[Rook]
		queenChecks := queenThreats & safe & blackAttackedBy[Queen]

		count := int(blackKingAttackersCount) * int(blackKingAttackersWeight)
		count += int(kingSafetyAttackValue) * 9 * int(blackKingAttackersCount) / PopCount(whiteKingArea)
		count += int(kingSafetyWeakSquares) * PopCount(whiteKingArea&weak)
		count += int(kingSafetyFriendlyPawns) * PopCount(pos.White&pos.Pawns&whiteKingArea & ^weak)
		count += int(kingSafetyNoEnemyQueens) * BoolToInt(pos.Black&pos.Queens != 0)
		count += int(kingSafetySafeQueenCheck) * PopCount(queenChecks)
		count += int(kingSafetySafeRookCheck) * PopCount(rookChecks)
		count += int(kingSafetySafeBishopCheck) * PopCount(bishopChecks)
		count += int(kingSafetySafeKnightCheck) * PopCount(knightChecks)
		count += int(kingSafetyAdjustment)
		if count > 0 {
			midResult -= count * count / 720
			endResult -= count / 20
		}
	}

	// black king
	blackKingDefenders := PopCount(
		(pos.Pawns | pos.Bishops | pos.Knights) & pos.Black & blackKingAreaMask[blackKingLocation],
	)
	midResult -= int(blackKingPos[blackKingLocation].Middle)
	endResult -= int(blackKingPos[blackKingLocation].End)
	midResult -= int(kingDefenders[blackKingDefenders].Middle)
	midResult -= int(kingDefenders[blackKingDefenders].End)
	if int(whiteKingAttackersCount) > 1-PopCount(pos.White&pos.Queens) {
		// Weak squares are attacked by the enemy, defended no more
		// than once and only defended by our Queens or our King
		weak := whiteAttacked & ^blackAttackedByTwo & (^blackAttacked | blackAttackedBy[Queen] | blackAttackedBy[King])

		safe := ^pos.White & (^blackAttacked | (weak & whiteAttackedByTwo))

		knightThreats := KnightAttacks[blackKingLocation]
		bishopThreats := BishopAttacks(blackKingLocation, allOccupation)
		rookThreats := RookAttacks(blackKingLocation, allOccupation)
		queenThreats := bishopThreats | rookThreats

		knightChecks := knightThreats & safe & whiteAttackedBy[Knight]
		bishopChecks := bishopThreats & safe & whiteAttackedBy[Bishop]
		rookChecks := rookThreats & safe & whiteAttackedBy[Rook]
		queenChecks := queenThreats & safe & whiteAttackedBy[Queen]

		count := int(whiteKingAttackersCount) * int(whiteKingAttackersWeight)
		count += int(kingSafetyAttackValue) * int(whiteKingAttackersCount) * 9 / PopCount(blackKingArea) // Scale value to king area size
		count += int(kingSafetyWeakSquares) * PopCount(blackKingArea&weak)
		count += int(kingSafetyFriendlyPawns) * PopCount(pos.Black&pos.Pawns&blackKingArea & ^weak)
		count += int(kingSafetyNoEnemyQueens) * BoolToInt(pos.White&pos.Queens != 0)
		count += int(kingSafetySafeQueenCheck) * PopCount(queenChecks)
		count += int(kingSafetySafeRookCheck) * PopCount(rookChecks)
		count += int(kingSafetySafeBishopCheck) * PopCount(bishopChecks)
		count += int(kingSafetySafeKnightCheck) * PopCount(knightChecks)
		count += int(kingSafetyAdjustment)
		if count > 0 {
			midResult += count * count / 720
			endResult += count / 20
		}
	}

	// tapering eval
	phase = (phase*256 + (totalPhase / 2)) / totalPhase
	result := (midResult*(256-phase) + (endResult * phase)) / 256

	if pos.WhiteMove {
		return result
	}
	return -result
}
