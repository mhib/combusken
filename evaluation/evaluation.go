package evaluation

import (
	. "github.com/mhib/combusken/backend"
	. "github.com/mhib/combusken/utils"
)

const tuning = false

var T Trace

const PawnPhase = 0
const KnightPhase = 1
const BishopPhase = 1
const RookPhase = 2
const QueenPhase = 4
const TotalPhase = PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2

var PawnValue = S(94, 140)
var KnightValue = S(502, 478)
var BishopValue = S(461, 466)
var RookValue = S(674, 785)
var QueenValue = S(1485, 1493)

// Piece Square Values
var PieceScores = [King + 1][8][4]Score{
	{},
	{ // knight
		{S(-92, -41), S(-19, -51), S(-28, -36), S(-20, -12)},
		{S(-21, -25), S(-24, -21), S(-13, -36), S(-14, -17)},
		{S(-14, -51), S(-1, -26), S(-8, -19), S(6, 2)},
		{S(-6, -12), S(19, -7), S(10, 13), S(9, 19)},
		{S(13, -9), S(10, -7), S(27, 16), S(21, 28)},
		{S(-23, -16), S(-12, -3), S(12, 22), S(26, 18)},
		{S(6, -31), S(-14, -11), S(33, -27), S(44, 6)},
		{S(-238, -36), S(-92, -8), S(-150, 26), S(-10, -11)},
	},
	{ // Bishop
		{S(28, -35), S(33, -3), S(12, -8), S(7, 5)},
		{S(29, -49), S(25, -43), S(32, -14), S(15, -4)},
		{S(22, -9), S(40, -3), S(12, -16), S(29, 9)},
		{S(14, -14), S(22, -2), S(25, 9), S(39, 11)},
		{S(-6, 5), S(25, 5), S(16, 14), S(39, 18)},
		{S(1, 7), S(-1, 21), S(16, -5), S(25, 15)},
		{S(-39, 31), S(-19, 7), S(12, 23), S(0, 29)},
		{S(-40, 7), S(-38, 33), S(-111, 45), S(-94, 48)},
	},
	{ // Rook
		{S(-15, -27), S(-6, -18), S(1, -17), S(6, -28)},
		{S(-58, -12), S(-10, -32), S(-2, -30), S(-11, -29)},
		{S(-32, -20), S(-5, -13), S(-16, -13), S(-13, -16)},
		{S(-26, 0), S(-12, 9), S(-10, 10), S(-7, 3)},
		{S(-10, 12), S(16, 10), S(28, 10), S(37, 6)},
		{S(-13, 23), S(53, 4), S(47, 17), S(64, 0)},
		{S(4, 27), S(7, 31), S(44, 21), S(50, 29)},
		{S(51, 26), S(70, 26), S(45, 33), S(45, 28)},
	},
	{ // Queen
		{S(13, -104), S(14, -93), S(15, -97), S(21, -65)},
		{S(8, -86), S(22, -83), S(24, -94), S(14, -56)},
		{S(4, -63), S(18, -31), S(9, 1), S(5, -15)},
		{S(3, -23), S(10, 6), S(5, 29), S(-9, 55)},
		{S(17, -10), S(10, 36), S(-5, 48), S(-16, 77)},
		{S(10, 4), S(20, -1), S(-6, 53), S(8, 47)},
		{S(-28, 11), S(-78, 57), S(-12, 49), S(-35, 82)},
		{S(-4, 7), S(30, 2), S(28, 31), S(35, 32)},
	},
	{ // King
		{S(175, -31), S(155, 16), S(57, 63), S(62, 83)},
		{S(168, 27), S(127, 42), S(55, 79), S(27, 135)},
		{S(104, 31), S(143, 35), S(111, 67), S(85, 135)},
		{S(85, 22), S(235, 27), S(157, 73), S(98, 145)},
		{S(89, 41), S(269, 37), S(182, 84), S(112, 144)},
		{S(116, 40), S(296, 55), S(250, 80), S(185, 121)},
		{S(114, -2), S(210, 63), S(189, 74), S(148, 114)},
		{S(157, -98), S(290, -16), S(173, 20), S(161, 78)},
	},
}

// Pawns Square scores
var PawnScores = [7][8]Score{
	{},
	{S(-15, 13), S(0, 2), S(-9, 9), S(-4, 11), S(-14, 18), S(11, 13), S(33, 4), S(2, 4)},
	{S(-14, 3), S(-18, 0), S(-6, 0), S(1, -3), S(3, -2), S(-7, 7), S(6, 2), S(-12, 6)},
	{S(-8, 13), S(-5, 11), S(15, -7), S(22, -8), S(29, -13), S(16, 1), S(10, 11), S(-14, 14)},
	{S(5, 26), S(18, 7), S(4, 5), S(23, -16), S(39, -16), S(40, -4), S(34, 10), S(-4, 25)},
	{S(0, 62), S(19, 37), S(32, 21), S(42, -9), S(68, -18), S(83, 5), S(54, 23), S(5, 54)},
	{S(-6, 14), S(-8, 10), S(3, 5), S(26, -12), S(33, -18), S(30, -9), S(-90, 30), S(-129, 44)},
}

var PawnsConnected = [7][4]Score{
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
	{S(1, -15), S(4, 9), S(4, -1), S(15, 19)},
	{S(11, 8), S(26, 6), S(24, 13), S(26, 17)},
	{S(11, 7), S(25, 7), S(11, 11), S(18, 18)},
	{S(6, 16), S(16, 25), S(36, 25), S(30, 18)},
	{S(41, 24), S(32, 60), S(78, 60), S(100, 71)},
	{S(176, 28), S(300, 20), S(294, 37), S(346, 41)},
}

var MobilityBonus = [...][32]Score{
	{S(-54, -134), S(-36, -83), S(-19, -41), S(-10, -18), S(-2, -7), S(3, 7), // Knights
		S(11, 10), S(20, 5), S(30, -9)},
	{S(-4, -141), S(6, -77), S(15, -38), S(24, -13), S(32, 2), S(38, 20), // Bishops
		S(41, 29), S(41, 37), S(42, 41), S(47, 45), S(51, 43), S(63, 38),
		S(91, 40), S(98, 22)},
	{S(-127, -146), S(-12, -38), S(1, 12), S(1, 39), S(6, 49), S(8, 60), // Rooks
		S(11, 68), S(16, 70), S(21, 75), S(26, 79), S(30, 83), S(33, 85),
		S(39, 84), S(54, 74), S(104, 42)},
	{S(-413, -159), S(-122, -138), S(-31, -174), S(-11, -119), S(2, -84), S(2, -13), // Queens
		S(8, 6), S(12, 24), S(16, 38), S(20, 50), S(24, 53), S(28, 54),
		S(29, 55), S(29, 63), S(32, 56), S(27, 61), S(25, 56), S(20, 57),
		S(25, 45), S(32, 27), S(41, 10), S(35, 0), S(38, -21), S(48, -50),
		S(14, -38), S(-88, -10), S(127, -127), S(49, -81)},
}

var PassedFriendlyDistance = [8]Score{
	S(0, 0), S(-6, 37), S(-11, 22), S(-10, 8),
	S(-4, -7), S(-2, -15), S(14, -26), S(-5, -37),
}

var PassedEnemyDistance = [8]Score{
	S(0, 0), S(-130, -58), S(-11, -8), S(7, 10),
	S(18, 20), S(16, 28), S(9, 37), S(17, 42),
}

var Psqt [2][King + 1][64]Score

var PawnsConnectedSquare [2][64]Score
var pawnsConnectedMask [2][64]uint64

// PassedRank[Rank] contains a bonus according to the rank of a passed pawn
var PassedRank = [7]Score{S(0, 0), S(-16, -26), S(-25, -9), S(-13, 28), S(26, 66), S(55, 140), S(197, 223)}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var PassedFile = [8]Score{S(-2, 29), S(-12, 33), S(-10, 14), S(-9, -2),
	S(-13, 1), S(-14, 8), S(-11, 21), S(13, 9),
}

var PassedStacked = [8]Score{S(0, 0), S(-16, -52), S(-21, -35), S(-32, -58), S(-3, -87), S(28, -207), S(0, 0), S(0, 0)}

var Isolated = S(-8, -14)
var Doubled = S(-11, -30)
var Backward = S(8, -2)
var BackwardOpen = S(-5, -14)

var BishopPair = S(23, 87)
var BishopRammedPawns = S(-6, -20)

var BishopOutpostUndefendedBonus = S(18, 1)
var BishopOutpostDefendedBonus = S(52, 15)

var LongDiagonalBishop = S(21, 23)

var KnightOutpostUndefendedBonus = S(15, -13)
var KnightOutpostDefendedBonus = S(44, 19)

var DistantKnight = [4]Score{S(-12, -2), S(-19, -9), S(-28, -16), S(-53, -20)}

var MinorBehindPawn = S(8, 34)

var Tempo int16 = 15

// Rook on semiopen, open file
var RookOnFile = [2]Score{S(9, 17), S(36, 10)}
var RookOnQueenFile = S(7, 29)

var KingDefenders = [12]Score{
	S(-11, -9), S(-5, -7), S(1, -7), S(6, -6),
	S(10, -3), S(12, 5), S(15, 10), S(17, 10),
	S(16, 8), S(17, -70), S(-1782, -10), S(11, 0),
}

var KingShelter = [2][8][8]Score{
	{{S(-25, 4), S(11, -17), S(21, -8), S(22, -4),
		S(23, -15), S(13, -8), S(6, -37), S(-96, 33)},
		{S(13, -3), S(27, -15), S(2, -2), S(-7, -3),
			S(-3, -6), S(21, -7), S(35, -22), S(-47, 15)},
		{S(14, 9), S(1, 5), S(-31, 7), S(-27, 6),
			S(-12, -7), S(-9, 12), S(10, 0), S(-37, 6)},
		{S(-29, 30), S(10, 7), S(-4, -4), S(2, -5),
			S(10, -24), S(4, -5), S(28, 5), S(-15, 3)},
		{S(-25, 19), S(-13, 9), S(-33, 10), S(-22, 9),
			S(-2, -7), S(-30, 5), S(3, -9), S(-28, 4)},
		{S(43, -15), S(23, -13), S(3, -9), S(3, -18),
			S(11, -27), S(22, -16), S(30, -20), S(-14, -1)},
		{S(23, -9), S(-3, -15), S(-25, -5), S(-17, -2),
			S(-13, -19), S(4, -7), S(4, -19), S(-23, 12)},
		{S(-36, 6), S(-22, -4), S(-10, 18), S(-17, 19),
			S(-25, 15), S(-19, 15), S(-27, -6), S(-71, 43)}},
	{{S(48, 42), S(-42, -10), S(-25, -1), S(-40, -1),
		S(-51, -12), S(-22, 16), S(-70, 4), S(-89, 37)},
		{S(149, 9), S(4, -11), S(-2, 0), S(-26, 14),
			S(-16, -13), S(11, -6), S(3, -6), S(-85, 27)},
		{S(5, 23), S(44, -11), S(16, -12), S(24, -17),
			S(29, -18), S(14, -11), S(45, -10), S(-28, 3)},
		{S(6, 50), S(-12, 31), S(-19, 26), S(-25, 28),
			S(-3, 15), S(-16, 16), S(-12, 21), S(-45, 4)},
		{S(-10, 43), S(-3, 21), S(-6, 14), S(-11, 12),
			S(-3, 12), S(4, -11), S(18, -10), S(-32, 6)},
		{S(35, 6), S(0, -5), S(-7, -2), S(-1, -3),
			S(3, -9), S(-44, -2), S(-7, -13), S(-41, 3)},
		{S(44, 0), S(5, -18), S(1, -16), S(-21, -9),
			S(-10, -19), S(-3, -16), S(-4, -18), S(-70, 18)},
		{S(-5, -4), S(-7, -35), S(0, -22), S(-5, -20),
			S(-14, -29), S(-8, -21), S(-18, -41), S(-69, 28)}},
}

var KingStorm = [2][4][8]Score{
	{{S(20, 2), S(12, 8), S(16, 2), S(7, 4),
		S(-2, 6), S(2, 7), S(-2, 15), S(-7, -8)},
		{S(27, -5), S(23, -2), S(25, -4), S(6, 6),
			S(14, 0), S(17, -9), S(8, -6), S(-12, -10)},
		{S(24, 10), S(8, 5), S(7, 11), S(-1, 14),
			S(1, 9), S(12, 0), S(13, -17), S(2, 0)},
		{S(11, 12), S(7, 8), S(12, 2), S(-1, -2),
			S(-3, 2), S(15, -1), S(10, -4), S(-9, 4)}},
	{{S(0, 0), S(19, -6), S(-20, 7), S(8, -8),
		S(-3, 18), S(-19, 35), S(4, 53), S(1, -9)},
		{S(0, 0), S(-40, -20), S(0, -8), S(41, -8),
			S(-1, -4), S(-13, -5), S(13, 56), S(-6, -7)},
		{S(0, 0), S(-17, -9), S(-9, -2), S(15, 0),
			S(5, -5), S(-7, -19), S(49, -46), S(-4, 7)},
		{S(0, 0), S(-13, -18), S(-18, -15), S(-10, -1),
			S(1, -19), S(1, -52), S(5, 3), S(-13, 12)}},
}
var KingOnPawnlessFlank = S(19, -70)

var passedMask [2][64]uint64

var outpustMask [2][64]uint64

var distanceBetween [64][64]int16

var adjacentFilesMask [8]uint64

var whiteKingAreaMask [64]uint64
var blackKingAreaMask [64]uint64

var forwardRanksMask [2][8]uint64

var forwardFileMask [2][64]uint64

// Outpost bitboards
const whiteOutpustRanks = RANK_4_BB | RANK_5_BB | RANK_6_BB
const blackOutpustRanks = RANK_5_BB | RANK_4_BB | RANK_3_BB

var KingSafetyAttacksWeights = [King + 1]int16{0, -3, -7, -4, 7, 0}
var KingSafetyAttackValue int16 = 122
var KingSafetyWeakSquares int16 = 43
var KingSafetyFriendlyPawns int16 = -33
var KingSafetyNoEnemyQueens int16 = -160
var KingSafetySafeQueenCheck int16 = 92
var KingSafetySafeRookCheck int16 = 76
var KingSafetySafeBishopCheck int16 = 60
var KingSafetySafeKnightCheck int16 = 113
var KingSafetyAdjustment int16 = -14

var Hanging = S(30, 14)
var ThreatByKing = S(-6, 31)
var ThreatByMinor = [King + 1]Score{S(0, 0), S(22, 40), S(19, 33), S(34, 27), S(31, 30), S(13, 23)}
var ThreatByRook = [King + 1]Score{S(0, 0), S(-2, 11), S(-1, 14), S(-6, -7), S(35, 8), S(23, 1)}

func LoadScoresToPieceSquares() {
	for x := 0; x < 4; x++ {
		for y := 0; y < 8; y++ {
			Psqt[White][Knight][y*8+x] = PieceScores[Knight][y][x] + KnightValue
			Psqt[White][Knight][y*8+(7-x)] = PieceScores[Knight][y][x] + KnightValue
			Psqt[Black][Knight][(7-y)*8+x] = PieceScores[Knight][y][x] + KnightValue
			Psqt[Black][Knight][(7-y)*8+(7-x)] = PieceScores[Knight][y][x] + KnightValue

			Psqt[White][Bishop][y*8+x] = PieceScores[Bishop][y][x] + BishopValue
			Psqt[White][Bishop][y*8+(7-x)] = PieceScores[Bishop][y][x] + BishopValue
			Psqt[Black][Bishop][(7-y)*8+x] = PieceScores[Bishop][y][x] + BishopValue
			Psqt[Black][Bishop][(7-y)*8+(7-x)] = PieceScores[Bishop][y][x] + BishopValue

			Psqt[White][Rook][y*8+x] = PieceScores[Rook][y][x] + RookValue
			Psqt[White][Rook][y*8+(7-x)] = PieceScores[Rook][y][x] + RookValue
			Psqt[Black][Rook][(7-y)*8+x] = PieceScores[Rook][y][x] + RookValue
			Psqt[Black][Rook][(7-y)*8+(7-x)] = PieceScores[Rook][y][x] + RookValue

			Psqt[White][Queen][y*8+x] = PieceScores[Queen][y][x] + QueenValue
			Psqt[White][Queen][y*8+(7-x)] = PieceScores[Queen][y][x] + QueenValue
			Psqt[Black][Queen][(7-y)*8+x] = PieceScores[Queen][y][x] + QueenValue
			Psqt[Black][Queen][(7-y)*8+(7-x)] = PieceScores[Queen][y][x] + QueenValue

			Psqt[White][King][y*8+x] = PieceScores[King][y][x]
			Psqt[White][King][y*8+(7-x)] = PieceScores[King][y][x]
			Psqt[Black][King][(7-y)*8+x] = PieceScores[King][y][x]
			Psqt[Black][King][(7-y)*8+(7-x)] = PieceScores[King][y][x]

			if y != 7 {
				PawnsConnectedSquare[White][y*8+x] = PawnsConnected[y][x]
				PawnsConnectedSquare[White][y*8+(7-x)] = PawnsConnected[y][x]
				PawnsConnectedSquare[Black][(7-y)*8+x] = PawnsConnected[y][x]
				PawnsConnectedSquare[Black][(7-y)*8+(7-x)] = PawnsConnected[y][x]
			}
		}
	}

	for y := 1; y < 7; y++ {
		for x := 0; x < 8; x++ {
			Psqt[White][Pawn][y*8+x] = PawnScores[y][x] + PawnValue
			Psqt[Black][Pawn][(7-y)*8+x] = PawnScores[y][x] + PawnValue
		}
	}
}

func init() {
	LoadScoresToPieceSquares()

	// Pawn is passed if no pawn of opposite color can stop it from promoting
	for i := 8; i <= 55; i++ {
		passedMask[White][i] = 0
		for file := File(i) - 1; file <= File(i)+1; file++ {
			if file < FILE_A || file > FILE_H {
				continue
			}
			for rank := Rank(i) + 1; rank < RANK_8; rank++ {
				passedMask[White][i] |= 1 << uint(rank*8+file)
			}
		}
	}
	// Outpust is similar to passed bitboard bot we do not care about pawns in same file
	for i := 8; i <= 55; i++ {
		outpustMask[White][i] = passedMask[White][i] & ^FILES[File(i)]
	}

	for i := 55; i >= 8; i-- {
		passedMask[Black][i] = 0
		for file := File(i) - 1; file <= File(i)+1; file++ {
			if file < FILE_A || file > FILE_H {
				continue
			}
			for rank := Rank(i) - 1; rank > RANK_1; rank-- {
				passedMask[Black][i] |= 1 << uint(rank*8+file)
			}
		}
	}
	for i := 55; i >= 8; i-- {
		outpustMask[Black][i] = passedMask[Black][i] & ^FILES[File(i)]
	}

	for i := 8; i <= 55; i++ {
		pawnsConnectedMask[White][i] = PawnAttacks[Black][i] | PawnAttacks[Black][i+8]
		pawnsConnectedMask[Black][i] = PawnAttacks[White][i] | PawnAttacks[White][i-8]
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
			forwardRanksMask[White][rank] |= RANKS[y]
		}
		forwardRanksMask[Black][rank] = (^forwardRanksMask[White][rank]) | RANKS[rank]
	}

	for y := 0; y < 64; y++ {
		forwardFileMask[White][y] = forwardRanksMask[White][Rank(y)] & FILES[File(y)] & ^SquareMask[y]
		forwardFileMask[Black][y] = forwardRanksMask[Black][Rank(y)] & FILES[File(y)] & ^SquareMask[y]
	}
}

func IsLateEndGame(pos *Position) bool {
	return ((pos.Pieces[Rook] | pos.Pieces[Queen] | pos.Pieces[Bishop] | pos.Pieces[Knight]) & pos.Colours[pos.SideToMove]) == 0
}

func evaluateKingPawns(pos *Position) Score {
	if !tuning {
		if ok, score := GlobalPawnKingTable.Get(pos.PawnKey); ok {
			return score
		}
	}
	var fromBB uint64
	var fromId int
	whiteKingLocation := BitScan(pos.Pieces[King] & pos.Colours[White])
	blackKingLocation := BitScan(pos.Pieces[King] & pos.Colours[Black])
	score := SCORE_ZERO

	// white pawns
	for fromBB = pos.Pieces[Pawn] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		score += Psqt[White][Pawn][fromId]
		if tuning {
			T.PawnValue++
			T.PawnScores[Rank(fromId)][File(fromId)]++
		}

		// Passed bonus
		if passedMask[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			// Bonus is calculated based on rank, file, distance from friendly and enemy king
			score +=
				PassedRank[Rank(fromId)] +
					PassedFile[File(fromId)] +
					PassedFriendlyDistance[distanceBetween[whiteKingLocation][fromId]] +
					PassedEnemyDistance[distanceBetween[blackKingLocation][fromId]]

			if tuning {
				T.PassedRank[Rank(fromId)]++
				T.PassedFile[File(fromId)]++
				T.PassedFriendlyDistance[distanceBetween[whiteKingLocation][fromId]]++
				T.PassedEnemyDistance[distanceBetween[blackKingLocation][fromId]]++
			}

			if pos.Pieces[Pawn]&pos.Colours[White]&forwardFileMask[White][fromId] != 0 {
				score += PassedStacked[Rank(fromId)]
				if tuning {
					T.PassedStacked[Rank(fromId)]++
				}
			}
		}

		// Isolated pawn penalty
		if adjacentFilesMask[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			score += Isolated
			if tuning {
				T.Isolated++
			}
		}

		// Pawn is backward if there are no pawns behind it and cannot increase rank without being attacked by enemy pawn
		if passedMask[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 &&
			PawnAttacks[White][fromId+8]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
			if FILES[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
				score += BackwardOpen
				if tuning {
					T.BackwardOpen++
				}
			} else {
				score += Backward
				if tuning {
					T.Backward++
				}
			}
		} else if pawnsConnectedMask[White][fromId]&(pos.Colours[White]&pos.Pieces[Pawn]) != 0 {
			score += PawnsConnectedSquare[White][fromId]
			if tuning {
				T.PawnsConnected[Rank(fromId)][FileMirror[File(fromId)]]++
			}
		}
	}

	// white doubled pawns
	score += Score(PopCount(pos.Pieces[Pawn]&pos.Colours[White]&South(pos.Pieces[Pawn]&pos.Colours[White]))) * Doubled
	if tuning {
		T.Doubled += PopCount(pos.Pieces[Pawn] & pos.Colours[White] & South(pos.Pieces[Pawn]&pos.Colours[White]))
	}

	// black pawns
	for fromBB = pos.Pieces[Pawn] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		score -= Psqt[Black][Pawn][fromId]

		if tuning {
			T.PawnValue--
			T.PawnScores[7-Rank(fromId)][File(fromId)]--
		}
		if passedMask[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			score -=
				PassedRank[7-Rank(fromId)] +
					PassedFile[File(fromId)] +
					PassedFriendlyDistance[distanceBetween[blackKingLocation][fromId]] +
					PassedEnemyDistance[distanceBetween[whiteKingLocation][fromId]]
			if tuning {
				T.PassedRank[7-Rank(fromId)]--
				T.PassedFile[File(fromId)]--
				T.PassedFriendlyDistance[distanceBetween[blackKingLocation][fromId]]--
				T.PassedEnemyDistance[distanceBetween[whiteKingLocation][fromId]]--
			}

			if pos.Pieces[Pawn]&pos.Colours[Black]&forwardFileMask[Black][fromId] != 0 {
				score -= PassedStacked[7-Rank(fromId)]
				if tuning {
					T.PassedStacked[7-Rank(fromId)]--
				}
			}
		}
		if adjacentFilesMask[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			score -= Isolated
			if tuning {
				T.Isolated--
			}
		}
		if passedMask[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 &&
			PawnAttacks[Black][fromId-8]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
			if FILES[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
				score -= BackwardOpen
				if tuning {
					T.BackwardOpen--
				}
			} else {
				score -= Backward
				if tuning {
					T.Backward--
				}
			}
		} else if pawnsConnectedMask[Black][fromId]&(pos.Colours[Black]&pos.Pieces[Pawn]) != 0 {
			score -= PawnsConnectedSquare[Black][fromId]
			if tuning {
				T.PawnsConnected[7-Rank(fromId)][FileMirror[File(fromId)]]--
			}
		}
	}

	// black doubled pawns
	score -= Score(PopCount(pos.Pieces[Pawn]&pos.Colours[Black]&North(pos.Pieces[Pawn]&pos.Colours[Black]))) * Doubled
	if tuning {
		T.Doubled -= PopCount(pos.Pieces[Pawn] & pos.Colours[Black] & North(pos.Pieces[Pawn]&pos.Colours[Black]))
	}

	// White king storm shelter
	for file := Max(File(whiteKingLocation)-1, FILE_A); file <= Min(File(whiteKingLocation)+1, FILE_H); file++ {
		ours := pos.Pieces[Pawn] & FILES[file] & pos.Colours[White] & forwardRanksMask[White][Rank(whiteKingLocation)]
		var ourDist int
		if ours == 0 {
			ourDist = 7
		} else {
			ourDist = Abs(Rank(whiteKingLocation) - Rank(BitScan(ours)))
		}
		theirs := pos.Pieces[Pawn] & FILES[file] & pos.Colours[Black] & forwardRanksMask[White][Rank(whiteKingLocation)]
		var theirDist int
		if theirs == 0 {
			theirDist = 7
		} else {
			theirDist = Abs(Rank(whiteKingLocation) - Rank(BitScan(theirs)))
		}
		sameFile := BoolToInt(file == File(whiteKingLocation))
		score += KingShelter[sameFile][file][ourDist]
		if tuning {
			T.KingShelter[sameFile][file][ourDist]++
		}

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		score += KingStorm[blocked][FileMirror[file]][theirDist]

		if tuning {
			T.KingStorm[blocked][FileMirror[file]][theirDist]++
		}
	}
	if KING_FLANK_BB[File(whiteKingLocation)]&pos.Pieces[Pawn] == 0 {
		score += KingOnPawnlessFlank
		if tuning {
			T.KingOnPawnlessFlank++
		}
	}

	// Black king storm / shelter
	for file := Max(File(blackKingLocation)-1, FILE_A); file <= Min(File(blackKingLocation)+1, FILE_H); file++ {
		ours := pos.Pieces[Pawn] & FILES[file] & pos.Colours[Black] & forwardRanksMask[Black][Rank(blackKingLocation)]
		var ourDist int
		if ours == 0 {
			ourDist = 7
		} else {
			ourDist = Abs(Rank(blackKingLocation) - Rank(MostSignificantBit(ours)))
		}
		theirs := pos.Pieces[Pawn] & FILES[file] & pos.Colours[White] & forwardRanksMask[Black][Rank(blackKingLocation)]
		var theirDist int
		if theirs == 0 {
			theirDist = 7
		} else {
			theirDist = Abs(Rank(blackKingLocation) - Rank(MostSignificantBit(theirs)))
		}
		sameFile := BoolToInt(file == File(blackKingLocation))
		score -= KingShelter[sameFile][file][ourDist]
		if tuning {
			T.KingShelter[sameFile][file][ourDist]--
		}

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		score -= KingStorm[blocked][FileMirror[file]][theirDist]
		if tuning {
			T.KingStorm[blocked][FileMirror[file]][theirDist]--
		}
	}
	if KING_FLANK_BB[File(blackKingLocation)]&pos.Pieces[Pawn] == 0 {
		score -= KingOnPawnlessFlank
		if tuning {
			T.KingOnPawnlessFlank--
		}
	}
	if !tuning {
		GlobalPawnKingTable.Set(pos.PawnKey, score)
	}
	return score
}

func Evaluate(pos *Position) int {
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

	phase := TotalPhase
	whiteMobilityArea := ^((pos.Pieces[Pawn] & pos.Colours[White]) | (BlackPawnsAttacks(pos.Pieces[Pawn] & pos.Colours[Black])))
	blackMobilityArea := ^((pos.Pieces[Pawn] & pos.Colours[Black]) | (WhitePawnsAttacks(pos.Pieces[Pawn] & pos.Colours[White])))
	allOccupation := pos.Colours[White] | pos.Colours[Black]

	whiteKingLocation := BitScan(pos.Pieces[King] & pos.Colours[White])
	attacks = KingAttacks[whiteKingLocation]
	whiteAttacked |= attacks
	whiteAttackedBy[King] |= attacks
	whiteKingArea := whiteKingAreaMask[whiteKingLocation]

	blackKingLocation := BitScan(pos.Pieces[King] & pos.Colours[Black])
	attacks = KingAttacks[blackKingLocation]
	blackAttacked |= attacks
	blackAttackedBy[King] |= attacks
	blackKingArea := blackKingAreaMask[blackKingLocation]

	// white pawns
	attacks = WhitePawnsAttacks(pos.Pieces[Pawn] & pos.Colours[White])
	whiteAttackedByTwo |= whiteAttacked & attacks
	whiteAttacked |= attacks
	whiteAttackedBy[Pawn] |= attacks
	whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))

	// black pawns
	attacks = BlackPawnsAttacks(pos.Pieces[Pawn] & pos.Colours[Black])
	blackAttackedByTwo |= blackAttacked & attacks
	blackAttacked |= attacks
	blackAttackedBy[Pawn] |= attacks
	blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))

	score := evaluateKingPawns(pos)

	// white knights
	for fromBB = pos.Pieces[Knight] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= KnightPhase
		fromId = BitScan(fromBB)

		attacks = KnightAttacks[fromId]
		mobility := PopCount(whiteMobilityArea & attacks)
		score += Psqt[White][Knight][fromId]
		score += MobilityBonus[0][mobility]
		if tuning {
			T.KnightValue++
			T.PieceScores[Knight][Rank(fromId)][FileMirror[File(fromId)]]++
			T.MobilityBonus[0][mobility]++
		}

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Knight] |= attacks

		if (pos.Pieces[Pawn]>>8)&SquareMask[fromId] != 0 {
			score += MinorBehindPawn
			if tuning {
				T.MinorBehindPawn++
			}
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && outpustMask[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			if PawnAttacks[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
				score += KnightOutpostDefendedBonus
				if tuning {
					T.KnightOutpostDefendedBonus++
				}
			} else {
				score += KnightOutpostUndefendedBonus
				if tuning {
					T.KnightOutpostUndefendedBonus++
				}
			}
		}

		kingDistance := Min(int(distanceBetween[fromId][whiteKingLocation]), int(distanceBetween[fromId][blackKingLocation]))
		if kingDistance >= 4 {
			score += DistantKnight[kingDistance-4]
			if tuning {
				T.DistantKnight[kingDistance-4]++
			}
		}
		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += KingSafetyAttacksWeights[Knight]
		}
	}

	// black knights
	for fromBB = pos.Pieces[Knight] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= KnightPhase
		fromId = BitScan(fromBB)

		attacks = KnightAttacks[fromId]
		mobility := PopCount(blackMobilityArea & attacks)
		score -= Psqt[Black][Knight][fromId]
		score -= MobilityBonus[0][mobility]
		if tuning {
			T.KnightValue--
			T.PieceScores[Knight][7-Rank(fromId)][FileMirror[File(fromId)]]--
			T.MobilityBonus[0][mobility]--
		}

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Knight] |= attacks

		if (pos.Pieces[Pawn]<<8)&SquareMask[fromId] != 0 {
			score -= MinorBehindPawn
			if tuning {
				T.MinorBehindPawn--
			}
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && outpustMask[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			if PawnAttacks[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
				score -= KnightOutpostDefendedBonus
				if tuning {
					T.KnightOutpostDefendedBonus--
				}
			} else {
				score -= KnightOutpostUndefendedBonus
				if tuning {
					T.KnightOutpostUndefendedBonus--
				}
			}
		}
		kingDistance := Min(int(distanceBetween[fromId][whiteKingLocation]), int(distanceBetween[fromId][blackKingLocation]))
		if kingDistance >= 4 {
			score -= DistantKnight[kingDistance-4]
			if tuning {
				T.DistantKnight[kingDistance-4]--
			}
		}
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += KingSafetyAttacksWeights[Knight]
		}
	}

	// white bishops
	whiteRammedPawns := South(pos.Pieces[Pawn]&pos.Colours[Black]) & (pos.Pieces[Pawn] & pos.Colours[White])
	for fromBB = pos.Pieces[Bishop] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= BishopPhase
		fromId = BitScan(fromBB)

		attacks = BishopAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		score += MobilityBonus[1][mobility]
		score += Psqt[White][Bishop][fromId]
		if tuning {
			T.BishopValue++
			T.PieceScores[Bishop][Rank(fromId)][FileMirror[File(fromId)]]++
			T.MobilityBonus[1][mobility]++
		}

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Bishop] |= attacks

		if (pos.Pieces[Pawn]>>8)&SquareMask[fromId] != 0 {
			score += MinorBehindPawn
			if tuning {
				T.MinorBehindPawn++
			}
		}
		if (LONG_DIAGONALS&SquareMask[fromId]) != 0 && (MoreThanOne(BishopAttacks(fromId, pos.Pieces[Pawn]) & CENTER)) {
			score += LongDiagonalBishop
			if tuning {
				T.LongDiagonalBishop++
			}
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && outpustMask[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			if PawnAttacks[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
				score += BishopOutpostDefendedBonus
				if tuning {
					T.BishopOutpostDefendedBonus++
				}
			} else {
				score += BishopOutpostUndefendedBonus
				if tuning {
					T.BishopOutpostUndefendedBonus++
				}
			}
		}

		// Bishop is worth less if there are friendly rammed pawns of its color
		var rammedCount Score
		if SquareMask[fromId]&WHITE_SQUARES != 0 {
			rammedCount = Score(PopCount(whiteRammedPawns & WHITE_SQUARES))
		} else {
			rammedCount = Score(PopCount(whiteRammedPawns & BLACK_SQUARES))
		}
		score += BishopRammedPawns * rammedCount
		if tuning {
			T.BishopRammedPawns += int(rammedCount)
		}
		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += KingSafetyAttacksWeights[Bishop]
		}
	}

	// Bishop pair bonus
	// It is not checked if bishops have opposite colors, but that is almost always the case
	if MoreThanOne(pos.Pieces[Bishop] & pos.Colours[White]) {
		score += BishopPair
		if tuning {
			T.BishopPair++
		}
	}

	// black bishops
	blackRammedPawns := North(pos.Pieces[Pawn]&pos.Colours[White]) & (pos.Pieces[Pawn] & pos.Colours[Black])
	for fromBB = pos.Pieces[Bishop] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= BishopPhase
		fromId = BitScan(fromBB)

		attacks = BishopAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		score -= MobilityBonus[1][mobility]
		score -= Psqt[Black][Bishop][fromId]
		if tuning {
			T.BishopValue--
			T.PieceScores[Bishop][7-Rank(fromId)][FileMirror[File(fromId)]]--
			T.MobilityBonus[1][mobility]--
		}

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Bishop] |= attacks

		if (pos.Pieces[Pawn]<<8)&SquareMask[fromId] != 0 {
			score -= MinorBehindPawn
			if tuning {
				T.MinorBehindPawn--
			}
		}
		if (LONG_DIAGONALS&SquareMask[fromId]) != 0 && (MoreThanOne(BishopAttacks(fromId, pos.Pieces[Pawn]) & CENTER)) {
			score -= LongDiagonalBishop
			if tuning {
				T.LongDiagonalBishop--
			}
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && outpustMask[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			if PawnAttacks[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
				score -= BishopOutpostDefendedBonus
				if tuning {
					T.BishopOutpostDefendedBonus--
				}
			} else {
				score -= BishopOutpostUndefendedBonus
				if tuning {
					T.BishopOutpostUndefendedBonus--
				}
			}
		}
		var rammedCount Score
		if SquareMask[fromId]&WHITE_SQUARES != 0 {
			rammedCount = Score(PopCount(blackRammedPawns & WHITE_SQUARES))
		} else {
			rammedCount = Score(PopCount(blackRammedPawns & BLACK_SQUARES))
		}
		score -= BishopRammedPawns * rammedCount
		if tuning {
			T.BishopRammedPawns -= int(rammedCount)
		}
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += KingSafetyAttacksWeights[Bishop]
		}
	}

	if MoreThanOne(pos.Pieces[Bishop] & pos.Colours[Black]) {
		score -= BishopPair

		if tuning {
			T.BishopPair--
		}
	}

	// white rooks
	for fromBB = pos.Pieces[Rook] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= RookPhase
		fromId = BitScan(fromBB)

		attacks = RookAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		score += MobilityBonus[2][mobility]
		score += Psqt[White][Rook][fromId]

		if tuning {
			T.RookValue++
			T.PieceScores[Rook][Rank(fromId)][FileMirror[File(fromId)]]++
			T.MobilityBonus[2][mobility]++
		}

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Rook] |= attacks

		if pos.Pieces[Pawn]&FILES[File(fromId)] == 0 {
			score += RookOnFile[1]
			if tuning {
				T.RookOnFile[1]++
			}
		} else if (pos.Pieces[Pawn]&pos.Colours[White])&FILES[File(fromId)] == 0 {
			score += RookOnFile[0]
			if tuning {
				T.RookOnFile[0]++
			}
		}

		if FileBB(fromId)&pos.Pieces[Queen] != 0 {
			score += RookOnQueenFile
			if tuning {
				T.RookOnQueenFile++
			}
		}

		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += KingSafetyAttacksWeights[Rook]
		}
	}

	// black rooks
	for fromBB = pos.Pieces[Rook] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= RookPhase
		fromId = BitScan(fromBB)

		attacks = RookAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		score -= MobilityBonus[2][mobility]
		score -= Psqt[Black][Rook][fromId]

		if tuning {
			T.RookValue--
			T.PieceScores[Rook][7-Rank(fromId)][FileMirror[File(fromId)]]--
			T.MobilityBonus[2][mobility]--
		}

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Rook] |= attacks

		if pos.Pieces[Pawn]&FILES[File(fromId)] == 0 {
			score -= RookOnFile[1]
			if tuning {
				T.RookOnFile[1]--
			}
		} else if (pos.Pieces[Pawn]&pos.Colours[Black])&FILES[File(fromId)] == 0 {
			score -= RookOnFile[0]
			if tuning {
				T.RookOnFile[0]--
			}
		}

		if FileBB(fromId)&pos.Pieces[Queen] != 0 {
			score -= RookOnQueenFile
			if tuning {
				T.RookOnQueenFile--
			}
		}

		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += KingSafetyAttacksWeights[Rook]
		}
	}

	//white queens
	for fromBB = pos.Pieces[Queen] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= QueenPhase
		fromId = BitScan(fromBB)

		attacks = QueenAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		score += MobilityBonus[3][mobility]
		score += Psqt[White][Queen][fromId]

		if tuning {
			T.QueenValue++
			T.PieceScores[Queen][Rank(fromId)][FileMirror[File(fromId)]]++
			T.MobilityBonus[3][mobility]++
		}

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Queen] |= attacks

		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += KingSafetyAttacksWeights[Queen]
		}
	}

	// black queens
	for fromBB = pos.Pieces[Queen] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= QueenPhase
		fromId = BitScan(fromBB)

		attacks = QueenAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		score -= MobilityBonus[3][mobility]
		score -= Psqt[Black][Queen][fromId]

		if tuning {
			T.QueenValue--
			T.PieceScores[Queen][7-Rank(fromId)][FileMirror[File(fromId)]]--
			T.MobilityBonus[3][mobility]--
		}

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Queen] |= attacks
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += KingSafetyAttacksWeights[Queen]
		}
	}

	if phase < 0 {
		phase = 0
	}

	// white king
	whiteKingDefenders := PopCount(
		(pos.Pieces[Pawn] | pos.Pieces[Bishop] | pos.Pieces[Knight]) & pos.Colours[White] & whiteKingAreaMask[whiteKingLocation],
	)
	score += Psqt[White][King][whiteKingLocation]
	score += KingDefenders[whiteKingDefenders]
	if tuning {
		T.PieceScores[King][Rank(whiteKingLocation)][FileMirror[File(whiteKingLocation)]]++
		T.KingDefenders[whiteKingDefenders]++
	}

	// Weak squares are attacked by the enemy, defended no more
	// than once and only defended by our Queens or our King
	weakForWhite := blackAttacked & ^whiteAttackedByTwo & (^whiteAttacked | whiteAttackedBy[Queen] | whiteAttackedBy[King])
	if int(blackKingAttackersCount) > 1-PopCount(pos.Colours[Black]&pos.Pieces[Queen]) {
		safe := ^pos.Colours[Black] & (^whiteAttacked | (weakForWhite & blackAttackedByTwo))

		knightThreats := KnightAttacks[whiteKingLocation]
		bishopThreats := BishopAttacks(whiteKingLocation, allOccupation)
		rookThreats := RookAttacks(whiteKingLocation, allOccupation)
		queenThreats := bishopThreats | rookThreats

		knightChecks := knightThreats & safe & blackAttackedBy[Knight]
		bishopChecks := bishopThreats & safe & blackAttackedBy[Bishop]
		rookChecks := rookThreats & safe & blackAttackedBy[Rook]
		queenChecks := queenThreats & safe & blackAttackedBy[Queen]

		count := int(blackKingAttackersCount) * int(blackKingAttackersWeight)
		count += int(KingSafetyAttackValue) * 9 * int(blackKingAttackersCount) / PopCount(whiteKingArea)
		count += int(KingSafetyWeakSquares) * PopCount(whiteKingArea&weakForWhite)
		count += int(KingSafetyFriendlyPawns) * PopCount(pos.Colours[White]&pos.Pieces[Pawn]&whiteKingArea & ^weakForWhite)
		count += int(KingSafetyNoEnemyQueens) * BoolToInt(pos.Colours[Black]&pos.Pieces[Queen] == 0)
		count += int(KingSafetySafeQueenCheck) * PopCount(queenChecks)
		count += int(KingSafetySafeRookCheck) * PopCount(rookChecks)
		count += int(KingSafetySafeBishopCheck) * PopCount(bishopChecks)
		count += int(KingSafetySafeKnightCheck) * PopCount(knightChecks)
		count += int(KingSafetyAdjustment)
		if count > 0 {
			score -= S(int16(count*count/720), int16(count/20))
		}
	}

	// black king
	blackKingDefenders := PopCount(
		(pos.Pieces[Pawn] | pos.Pieces[Bishop] | pos.Pieces[Knight]) & pos.Colours[Black] & blackKingAreaMask[blackKingLocation],
	)
	score -= Psqt[Black][King][blackKingLocation]
	score -= KingDefenders[blackKingDefenders]
	if tuning {
		T.PieceScores[King][7-Rank(blackKingLocation)][FileMirror[File(blackKingLocation)]]--
		T.KingDefenders[blackKingDefenders]--
	}

	// Weak squares are attacked by the enemy, defended no more
	// than once and only defended by our Queens or our King
	weakForBlack := whiteAttacked & ^blackAttackedByTwo & (^blackAttacked | blackAttackedBy[Queen] | blackAttackedBy[King])
	if int(whiteKingAttackersCount) > 1-PopCount(pos.Colours[White]&pos.Pieces[Queen]) {
		safe := ^pos.Colours[White] & (^blackAttacked | (weakForBlack & whiteAttackedByTwo))

		knightThreats := KnightAttacks[blackKingLocation]
		bishopThreats := BishopAttacks(blackKingLocation, allOccupation)
		rookThreats := RookAttacks(blackKingLocation, allOccupation)
		queenThreats := bishopThreats | rookThreats

		knightChecks := knightThreats & safe & whiteAttackedBy[Knight]
		bishopChecks := bishopThreats & safe & whiteAttackedBy[Bishop]
		rookChecks := rookThreats & safe & whiteAttackedBy[Rook]
		queenChecks := queenThreats & safe & whiteAttackedBy[Queen]

		count := int(whiteKingAttackersCount) * int(whiteKingAttackersWeight)
		count += int(KingSafetyAttackValue) * int(whiteKingAttackersCount) * 9 / PopCount(blackKingArea) // Scale value to king area size
		count += int(KingSafetyWeakSquares) * PopCount(blackKingArea&weakForBlack)
		count += int(KingSafetyFriendlyPawns) * PopCount(pos.Colours[Black]&pos.Pieces[Pawn]&blackKingArea & ^weakForBlack)
		count += int(KingSafetyNoEnemyQueens) * BoolToInt(pos.Colours[White]&pos.Pieces[Queen] == 0)
		count += int(KingSafetySafeQueenCheck) * PopCount(queenChecks)
		count += int(KingSafetySafeRookCheck) * PopCount(rookChecks)
		count += int(KingSafetySafeBishopCheck) * PopCount(bishopChecks)
		count += int(KingSafetySafeKnightCheck) * PopCount(knightChecks)
		count += int(KingSafetyAdjustment)
		if count > 0 {
			score += S(int16(count*count/720), int16(count/20))
		}
	}

	// White threats
	blackStronglyProtected := blackAttackedBy[Pawn] | (blackAttackedByTwo & ^whiteAttackedByTwo)
	blackDefended := pos.Colours[Black] & ^pos.Pieces[Pawn] & blackStronglyProtected
	if ((pos.Colours[Black] & weakForBlack) | blackDefended) != 0 {
		for fromBB = pos.Colours[Black] & (blackDefended | weakForBlack) & (whiteAttackedBy[Knight] | whiteAttackedBy[Bishop]) & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score += ThreatByMinor[threatenedPiece]
			if tuning {
				T.ThreatByMinor[threatenedPiece]++
			}
		}

		for fromBB = pos.Colours[Black] & (blackDefended | weakForBlack) & whiteAttackedBy[Rook] & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) & ^pos.Pieces[Pawn] {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score += ThreatByRook[threatenedPiece]
			if tuning {
				T.ThreatByRook[threatenedPiece]++
			}
		}

		if weakForBlack&pos.Colours[Black]&whiteAttackedBy[King] != 0 {
			score += ThreatByKing
			if tuning {
				T.ThreatByKing++
			}
		}

		// Bonus if enemy has a hanging piece
		score += Hanging *
			Score(PopCount((pos.Colours[Black] & ^pos.Pieces[Pawn] & whiteAttackedByTwo)&weakForBlack))

		if tuning {
			T.Hanging += PopCount((pos.Colours[Black] & ^pos.Pieces[Pawn] & whiteAttackedByTwo) & weakForBlack)
		}

	}

	// Black threats
	whiteStronglyProtected := whiteAttackedBy[Pawn] | (whiteAttackedByTwo & ^blackAttackedByTwo)
	whiteDefended := pos.Colours[White] & ^pos.Pieces[Pawn] & whiteStronglyProtected
	if ((pos.Colours[White] & weakForWhite) | whiteDefended) != 0 {
		for fromBB = pos.Colours[White] & (whiteDefended | weakForWhite) & (blackAttackedBy[Knight] | blackAttackedBy[Bishop]) & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score -= ThreatByMinor[threatenedPiece]
			if tuning {
				T.ThreatByMinor[threatenedPiece]--
			}
		}

		for fromBB = pos.Colours[White] & (whiteDefended | weakForWhite) & blackAttackedBy[Rook] & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score -= ThreatByRook[threatenedPiece]
			if tuning {
				T.ThreatByRook[threatenedPiece]--
			}
		}

		if weakForWhite&pos.Colours[White]&blackAttackedBy[King] != 0 {
			score -= ThreatByKing
			if tuning {
				T.ThreatByKing--
			}
		}

		// Bonus if enemy has a hanging piece
		score -= Hanging *
			Score(PopCount(pos.Colours[White] & ^pos.Pieces[Pawn] & blackAttackedByTwo & weakForWhite))

		if tuning {
			T.Hanging -= PopCount(pos.Colours[White] & ^pos.Pieces[Pawn] & blackAttackedByTwo & weakForWhite)
		}
	}

	// Scale Factor inlined
	scale := SCALE_NORMAL
	if OnlyOne(pos.Colours[Black]&pos.Pieces[Bishop]) &&
		OnlyOne(pos.Colours[White]&pos.Pieces[Bishop]) &&
		OnlyOne(pos.Pieces[Bishop]&WHITE_SQUARES) &&
		(pos.Pieces[Knight]|pos.Pieces[Rook]|pos.Pieces[Queen]) == 0 {
		scale = SCALE_HARD
	} else if (score.End() > 0 && PopCount(pos.Colours[White]) == 2 && (pos.Colours[White]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) ||
		(score.End() < 0 && PopCount(pos.Colours[Black]) == 2 && (pos.Colours[Black]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) {
		return SCALE_DRAW
	}

	// tapering eval
	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase
	result := (int(score.Middle())*(256-phase) + (int(score.End()) * phase * scale / SCALE_NORMAL)) / 256

	if pos.SideToMove == White {
		return result + int(Tempo)
	}
	return -result + int(Tempo)
}

const SCALE_NORMAL = 2
const SCALE_HARD = 1
const SCALE_DRAW = 0

func ScaleFactor(pos *Position, endResult int16) int {
	// OCB without other pieces endgame
	if OnlyOne(pos.Colours[Black]&pos.Pieces[Bishop]) &&
		OnlyOne(pos.Colours[White]&pos.Pieces[Bishop]) &&
		OnlyOne(pos.Pieces[Bishop]&WHITE_SQUARES) &&
		(pos.Pieces[Knight]|pos.Pieces[Rook]|pos.Pieces[Queen]) == 0 {
		return SCALE_HARD
	}
	if (endResult > 0 && PopCount(pos.Colours[White]) == 2 && (pos.Colours[White]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) ||
		(endResult < 0 && PopCount(pos.Colours[Black]) == 2 && (pos.Colours[Black]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) {
		return SCALE_DRAW
	}
	return SCALE_NORMAL
}
