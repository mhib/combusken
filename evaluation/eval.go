package evaluation

import (
	. "github.com/mhib/combusken/backend"
	"github.com/mhib/combusken/transposition"
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

var PawnValue = S(99, 120)
var KnightValue = S(494, 437)
var BishopValue = S(449, 429)
var RookValue = S(643, 693)
var QueenValue = S(1435, 1321)

// Piece Square Values
var PieceScores = [King + 1][8][4]Score{
	{},
	{ // knight
		{S(-51, -40), S(-9, -54), S(-39, -34), S(-11, -25)},
		{S(-10, -61), S(-21, -38), S(-10, -30), S(-3, -19)},
		{S(-11, -38), S(2, -32), S(0, -18), S(6, -3)},
		{S(-11, -31), S(13, -25), S(8, -1), S(4, 5)},
		{S(-10, -31), S(1, -18), S(0, 6), S(19, 7)},
		{S(-35, -52), S(13, -42), S(-20, 1), S(4, -1)},
		{S(-75, -52), S(-42, -26), S(55, -53), S(14, -21)},
		{S(-202, -70), S(-74, -73), S(-92, -34), S(5, -53)},
	},
	{ // Bishop
		{S(0, -10), S(18, -3), S(21, -9), S(23, -5)},
		{S(15, -25), S(50, -29), S(35, -14), S(19, -3)},
		{S(17, -16), S(38, -11), S(36, 1), S(21, 9)},
		{S(1, -23), S(12, -19), S(15, -2), S(28, 0)},
		{S(-20, -15), S(5, -15), S(-2, -1), S(21, 3)},
		{S(-33, -10), S(-20, -12), S(-21, 0), S(-17, -8)},
		{S(-57, 2), S(-11, 7), S(2, 2), S(2, -6)},
		{S(-2, -26), S(-3, -19), S(-86, -5), S(-91, 2)},
	},
	{ // Rook
		{S(-4, -18), S(-10, -6), S(11, -14), S(13, -14)},
		{S(-41, -1), S(-6, -18), S(-5, -14), S(4, -13)},
		{S(-33, -8), S(-11, -9), S(-2, -17), S(-9, -12)},
		{S(-40, 3), S(-10, -3), S(-14, 1), S(-8, -2)},
		{S(-37, 9), S(-22, 4), S(12, 7), S(-5, 4)},
		{S(-25, 4), S(15, 1), S(21, -3), S(-8, 4)},
		{S(5, 12), S(6, 15), S(50, 3), S(59, -8)},
		{S(7, 11), S(9, 10), S(-18, 18), S(22, 11)},
	},
	{ // Queen
		{S(1, -58), S(16, -72), S(16, -58), S(30, -62)},
		{S(-2, -37), S(8, -50), S(29, -56), S(27, -44)},
		{S(-6, -14), S(18, -29), S(-1, 7), S(5, -5)},
		{S(-2, -17), S(-22, 34), S(-6, 18), S(-11, 36)},
		{S(-22, 16), S(-28, 36), S(-24, 30), S(-24, 52)},
		{S(16, -26), S(3, -8), S(24, -5), S(-3, 46)},
		{S(-10, -24), S(-53, 17), S(-11, 23), S(-39, 52)},
		{S(11, -19), S(3, 2), S(19, 6), S(28, 11)},
	},
	{ // King
		{S(188, -14), S(165, 20), S(80, 70), S(84, 57)},
		{S(178, 26), S(139, 45), S(70, 81), S(34, 92)},
		{S(98, 48), S(116, 55), S(58, 81), S(40, 93)},
		{S(17, 49), S(73, 55), S(32, 90), S(-2, 99)},
		{S(27, 59), S(119, 67), S(64, 90), S(80, 85)},
		{S(97, 61), S(237, 63), S(235, 77), S(177, 60)},
		{S(42, 69), S(116, 75), S(108, 98), S(170, 67)},
		{S(27, 15), S(133, 29), S(113, 69), S(26, 53)},
	},
}

// Pawns Square scores
var PawnScores = [7][8]Score{
	{},
	{S(-18, 6), S(11, -2), S(-10, 12), S(11, 4), S(4, 14), S(22, 7), S(32, -1), S(-4, -1)},
	{S(-10, -8), S(-19, -1), S(1, -5), S(6, -1), S(4, -4), S(-4, 5), S(9, -5), S(-8, -4)},
	{S(-14, 4), S(-13, 3), S(15, -8), S(23, -12), S(26, -11), S(4, 3), S(1, 3), S(-22, 5)},
	{S(-1, 14), S(27, -4), S(12, -3), S(37, -16), S(39, -17), S(20, 0), S(39, 5), S(-7, 16)},
	{S(11, 40), S(14, 31), S(49, 13), S(41, 0), S(82, -15), S(95, -3), S(57, 20), S(3, 41)},
	{S(-1, 57), S(84, 37), S(18, 36), S(52, 22), S(50, 38), S(98, 20), S(-48, 70), S(-117, 105)},
}

var PawnsConnected = [7][4]Score{
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
	{S(10, -19), S(7, 7), S(7, -8), S(3, 19)},
	{S(7, 4), S(32, 4), S(15, 9), S(18, 18)},
	{S(10, 8), S(24, 8), S(17, 11), S(18, 15)},
	{S(13, 16), S(6, 26), S(26, 24), S(30, 25)},
	{S(14, 57), S(33, 59), S(59, 61), S(58, 46)},
	{S(74, 41), S(152, 0), S(170, 20), S(355, 9)},
}

var MobilityBonus = [...][32]Score{
	{S(-42, -122), S(-33, -68), S(-21, -38), S(-20, -20), S(-6, -20), S(3, -8), // Knights
		S(13, -14), S(25, -15), S(41, -25)},
	{S(-28, -72), S(-9, -58), S(10, -27), S(16, -7), S(27, 1), S(38, 8), // Bishops
		S(43, 10), S(47, 13), S(50, 16), S(52, 16), S(66, 4), S(90, 6),
		S(49, 23), S(68, 17)},
	{S(-29, -35), S(-35, -26), S(-18, 10), S(-7, 35), S(3, 45), S(8, 55), // Rooks
		S(10, 67), S(21, 67), S(24, 64), S(33, 68), S(41, 69), S(44, 71),
		S(50, 74), S(62, 72), S(90, 59)},
	{S(-22, -20), S(-52, -8), S(4, -206), S(-6, -119), S(-1, -21), S(5, -23), // Queens
		S(4, -21), S(14, -5), S(16, 19), S(20, 22), S(20, 44), S(21, 46),
		S(27, 38), S(28, 61), S(29, 64), S(29, 70), S(31, 69), S(26, 71),
		S(37, 67), S(46, 72), S(62, 49), S(49, 55), S(75, 34), S(64, 21),
		S(68, 4), S(29, 16), S(-4, -16), S(-27, -15)},
}

var PassedFriendlyDistance = [8]Score{
	S(0, 0), S(2, 29), S(-3, 11), S(-6, -11),
	S(-10, -22), S(-16, -23), S(8, -28), S(-22, -19),
}

var PassedEnemyDistance = [8]Score{
	S(0, 0), S(5, -66), S(30, -28), S(14, 9),
	S(11, 30), S(5, 38), S(3, 42), S(-11, 48),
}

var Psqt [2][King + 1][64]Score

var PawnsConnectedSquare [2][64]Score
var pawnsConnectedMask [2][64]uint64

// PassedRank[Rank] contains a bonus according to the rank of a passed pawn
var PassedRank = [7]Score{S(0, 0), S(-4, -33), S(-6, -15), S(-6, 26), S(24, 73), S(31, 164), S(106, 255)}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var PassedFile = [8]Score{S(0, 24), S(-15, 23), S(-18, 12), S(-24, -4),
	S(-19, -3), S(5, -1), S(-20, 18), S(-11, 13),
}

var PassedStacked = [8]Score{S(0, 0), S(-14, -41), S(4, -45), S(-31, -41), S(-60, -43), S(-1, -48), S(0, 0), S(0, 0)}

var Isolated = S(-11, -10)
var Doubled = S(-10, -29)
var Backward = S(6, -2)
var BackwardOpen = S(-16, -2)

var BishopPair = S(38, 63)
var BishopRammedPawns = S(-7, -13)

var BishopOutpostUndefendedBonus = S(40, 1)
var BishopOutpostDefendedBonus = S(80, 6)

var KnightOutpostUndefendedBonus = S(35, -14)
var KnightOutpostDefendedBonus = S(56, 15)

var DistantKnight = [4]Score{S(-17, 14), S(-21, 3), S(-31, 14), S(-81, -2)}

var MinorBehindPawn = S(5, 29)

var Tempo = S(37, 32)

// Rook on semiopen, open file
var RookOnFile = [2]Score{S(4, 28), S(56, -3)}

var KingDefenders = [12]Score{
	S(-80, 7), S(-69, 7), S(-31, -1), S(-5, -8),
	S(5, -7), S(17, -6), S(27, -1), S(34, 0),
	S(41, -16), S(36, 3), S(11, 0), S(11, 0),
}

var KingShelter = [2][8][8]Score{
	{{S(-30, 7), S(-8, -8), S(-1, 7), S(33, -7),
		S(7, -17), S(10, -2), S(6, -10), S(-28, 15)},
		{S(22, -1), S(35, -13), S(-4, -2), S(-6, 6),
			S(-25, -4), S(9, -12), S(19, -39), S(-31, 4)},
		{S(11, 7), S(2, 3), S(-15, 4), S(-18, 9),
			S(-34, 0), S(-14, 2), S(0, -5), S(-21, 2)},
		{S(-15, 22), S(6, 5), S(-9, -2), S(3, 0),
			S(8, -22), S(-2, -14), S(10, -41), S(-21, 2)},
		{S(0, 8), S(-2, 3), S(-21, 1), S(-24, 9),
			S(-18, -7), S(-20, 0), S(-24, -9), S(-26, 5)},
		{S(43, -10), S(19, -12), S(0, -12), S(4, -10),
			S(10, -22), S(3, -4), S(31, -28), S(-8, 2)},
		{S(26, -3), S(-4, -5), S(-21, -8), S(-8, -2),
			S(-11, -8), S(11, -4), S(6, -21), S(-33, 14)},
		{S(-23, 1), S(-26, -1), S(-10, 6), S(-15, 11),
			S(-5, 8), S(-18, 19), S(-38, 8), S(-61, 36)}},
	{{S(-1, -6), S(-41, -17), S(-18, -6), S(-79, -21),
		S(-5, -12), S(-28, -14), S(-62, 1), S(-67, 20)},
		{S(9, 40), S(4, -16), S(-18, -4), S(-7, -2),
			S(-3, -1), S(19, -47), S(4, -16), S(-65, 20)},
		{S(23, 31), S(42, -5), S(11, -3), S(15, -9),
			S(16, 3), S(-11, -8), S(64, -15), S(-27, 10)},
		{S(7, 23), S(-31, 19), S(-14, 12), S(-19, 4),
			S(-16, 21), S(-79, 30), S(-15, -1), S(-44, 4)},
		{S(4, 49), S(3, 4), S(-4, 0), S(-11, -1),
			S(-13, 7), S(-4, -10), S(-3, -20), S(-38, 8)},
		{S(76, -18), S(21, -11), S(-9, 2), S(1, -13),
			S(-1, -6), S(-12, -13), S(10, -26), S(-31, 7)},
		{S(0, 11), S(6, -15), S(2, -15), S(-20, -8),
			S(-18, -11), S(-3, -17), S(0, -18), S(-70, 24)},
		{S(-1, -1), S(-1, -27), S(2, -13), S(-22, -8),
			S(-24, -8), S(13, -10), S(-22, -28), S(-62, 25)}},
}

var KingStorm = [2][4][8]Score{
	{{S(20, -2), S(14, 1), S(15, 4), S(2, 9),
		S(-1, 14), S(10, 12), S(0, 17), S(5, -12)},
		{S(16, 2), S(12, 4), S(26, 0), S(3, 9),
			S(10, 6), S(15, 0), S(6, -1), S(4, -13)},
		{S(16, 12), S(16, 8), S(2, 13), S(-9, 17),
			S(-3, 13), S(7, 5), S(27, -16), S(10, -6)},
		{S(14, 11), S(5, 4), S(7, 2), S(-3, 6),
			S(-7, 12), S(7, 8), S(3, 5), S(-2, 1)}},
	{{S(0, 0), S(10, 15), S(-12, 10), S(21, 0),
		S(17, 11), S(-11, 17), S(27, 45), S(17, -23)},
		{S(0, 0), S(41, -35), S(-1, -1), S(66, -11),
			S(49, -18), S(-20, 1), S(-30, 23), S(5, -20)},
		{S(0, 0), S(-104, 4), S(-2, -2), S(23, 3),
			S(4, 2), S(1, -9), S(72, -50), S(10, -2)},
		{S(0, 0), S(3, -19), S(13, -15), S(-7, 1),
			S(-6, 1), S(7, -19), S(-11, 2), S(-5, 17)}},
}

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

var KingSafetyAttacksWeights = [King + 1]int16{-2, -3, -4, -4, 0, 0}
var KingSafetyAttackValue int16 = 127
var KingSafetyWeakSquares int16 = 27
var KingSafetyFriendlyPawns int16 = -1
var KingSafetyNoEnemyQueens int16 = -157
var KingSafetySafeQueenCheck int16 = 71
var KingSafetySafeRookCheck int16 = 126
var KingSafetySafeBishopCheck int16 = 94
var KingSafetySafeKnightCheck int16 = 132
var KingSafetyAdjustment int16 = -32

var Hanging = S(54, 18)
var ThreatByKing = S(18, 39)
var ThreatByMinor = [King + 1]Score{S(0, 0), S(24, 32), S(30, 36), S(72, 24), S(60, 15), S(133, 55)}
var ThreatByRook = [King + 1]Score{S(0, 0), S(0, 19), S(-1, 28), S(-17, 23), S(96, 15), S(150, 38)}

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
	if ok, score := transposition.GlobalPawnKingTable.Get(pos.PawnKey); ok {
		return score
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
	transposition.GlobalPawnKingTable.Set(pos.PawnKey, score)
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

	// Tempo bonus
	if pos.SideToMove == White {
		score += Tempo
		if tuning {
			T.Tempo++
		}
	} else {
		score -= Tempo
		if tuning {
			T.Tempo--
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

	scale := ScaleFactor(pos, score.End())

	// tapering eval
	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase
	result := (int(score.Middle())*(256-phase) + (int(score.End()) * phase * scale / SCALE_NORMAL)) / 256

	if pos.SideToMove == White {
		return result
	}
	return -result
}

const SCALE_NORMAL = 1
const SCALE_DRAW = 0

func ScaleFactor(pos *Position, endResult int16) int {
	if (endResult > 0 && PopCount(pos.Colours[White]) == 2 && (pos.Colours[White]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) ||
		(endResult < 0 && PopCount(pos.Colours[Black]) == 2 && (pos.Colours[Black]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) {
		return SCALE_DRAW
	}
	return SCALE_NORMAL
}
