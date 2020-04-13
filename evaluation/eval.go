package evaluation

import . "github.com/mhib/combusken/backend"
import . "github.com/mhib/combusken/utils"
import "github.com/mhib/combusken/transposition"

const tuning = false

var T Trace

const pawnPhase = 0
const knightPhase = 1
const bishopPhase = 1
const rookPhase = 2
const queenPhase = 4
const totalPhase = pawnPhase*16 + knightPhase*4 + bishopPhase*4 + rookPhase*4 + queenPhase*2

var PawnValue = S(100, 123)
var KnightValue = S(480, 440)
var BishopValue = S(435, 435)
var RookValue = S(646, 696)
var QueenValue = S(1435, 1322)

// Piece Square Values
var pieceScores = [King + 1][8][4]Score{
	{},
	{ // knight
		{S(-59, -40), S(-10, -54), S(-40, -35), S(-13, -24)},
		{S(-18, -56), S(-25, -37), S(-10, -31), S(-3, -22)},
		{S(-11, -39), S(2, -32), S(0, -19), S(6, -4)},
		{S(-17, -28), S(14, -25), S(11, -2), S(0, 3)},
		{S(-2, -32), S(-1, -19), S(14, 1), S(20, 4)},
		{S(-37, -51), S(22, -43), S(-52, -4), S(28, -6)},
		{S(-76, -49), S(-42, -27), S(54, -53), S(-2, -21)},
		{S(-210, -69), S(-76, -73), S(-127, -30), S(1, -55)},
	},
	{ // Bishop
		{S(-9, -11), S(6, -4), S(26, -11), S(7, -4)},
		{S(3, -25), S(58, -30), S(29, -17), S(26, -6)},
		{S(15, -18), S(33, -12), S(37, 2), S(28, 7)},
		{S(0, -23), S(4, -18), S(20, -6), S(26, -1)},
		{S(-21, -13), S(3, -17), S(-2, -1), S(25, -2)},
		{S(-71, -1), S(-16, -13), S(-50, -1), S(-18, -8)},
		{S(-75, 2), S(20, -9), S(-14, 2), S(2, -14)},
		{S(-1, -28), S(-36, -19), S(-113, -5), S(-92, -4)},
	},
	{ // Rook
		{S(-7, -20), S(-14, -6), S(10, -13), S(14, -18)},
		{S(-42, 0), S(-5, -17), S(-9, -12), S(5, -14)},
		{S(-37, -8), S(-14, -8), S(-3, -17), S(-11, -12)},
		{S(-39, 3), S(-10, -3), S(-16, 1), S(-12, -3)},
		{S(-39, 9), S(-23, 2), S(11, 7), S(-6, 1)},
		{S(-25, 5), S(16, 1), S(20, -2), S(-8, 4)},
		{S(7, 12), S(6, 15), S(48, 2), S(60, -7)},
		{S(6, 11), S(8, 10), S(-21, 18), S(14, 11)},
	},
	{ // Queen
		{S(0, -58), S(16, -72), S(16, -60), S(32, -68)},
		{S(-3, -41), S(7, -50), S(28, -56), S(28, -44)},
		{S(-6, -14), S(18, -30), S(-2, 5), S(4, -6)},
		{S(-4, -18), S(-23, 34), S(-7, 18), S(-18, 44)},
		{S(-20, 12), S(-28, 27), S(-29, 31), S(-33, 51)},
		{S(20, -26), S(3, -8), S(-4, 13), S(-8, 45)},
		{S(-9, -24), S(-54, 17), S(-19, 21), S(-35, 52)},
		{S(8, -20), S(2, 0), S(17, 6), S(19, 14)},
	},
	{ // King
		{S(184, -15), S(158, 29), S(77, 74), S(76, 62)},
		{S(177, 24), S(135, 50), S(71, 85), S(32, 97)},
		{S(82, 50), S(108, 60), S(42, 87), S(39, 97)},
		{S(14, 48), S(72, 58), S(30, 91), S(-3, 102)},
		{S(25, 61), S(102, 77), S(63, 95), S(85, 88)},
		{S(94, 65), S(237, 71), S(221, 86), S(160, 68)},
		{S(43, 69), S(114, 80), S(107, 100), S(167, 75)},
		{S(26, 10), S(129, 37), S(112, 73), S(25, 59)},
	},
}

// Pawns Square scores
var pawnScores = [7][8]Score{
	{},
	{S(-19, 8), S(13, -4), S(-13, 13), S(8, 4), S(0, 15), S(16, 8), S(34, -3), S(-3, -5)},
	{S(-12, -8), S(-16, -3), S(1, -7), S(6, -2), S(2, -3), S(-3, 3), S(12, -8), S(-6, -8)},
	{S(-14, 2), S(-13, 1), S(18, -12), S(26, -12), S(27, -10), S(5, 0), S(3, 1), S(-20, 1)},
	{S(-2, 15), S(27, -4), S(11, -4), S(39, -16), S(40, -14), S(20, -3), S(37, 7), S(-8, 13)},
	{S(10, 41), S(23, 29), S(53, 13), S(47, 0), S(90, -11), S(94, -2), S(56, 19), S(10, 42)},
	{S(-2, 57), S(100, 37), S(20, 37), S(17, 34), S(82, 38), S(1, 36), S(-34, 62), S(-116, 105)},
}

var pawnsConnected = [7][4]Score{
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
	{S(10, -19), S(9, 7), S(7, -9), S(6, 19)},
	{S(10, 3), S(32, 5), S(16, 12), S(17, 20)},
	{S(11, 9), S(24, 8), S(17, 11), S(18, 16)},
	{S(13, 16), S(9, 27), S(26, 24), S(31, 22)},
	{S(15, 56), S(38, 62), S(61, 61), S(74, 45)},
	{S(7, 58), S(150, 0), S(168, 22), S(352, 41)},
}

var mobilityBonus = [...][32]Score{
	{S(-42, -117), S(-35, -67), S(-24, -37), S(-22, -17), S(-7, -17), S(1, -8), // Knights
		S(11, -11), S(21, -12), S(31, -27)},
	{S(-26, -73), S(-11, -59), S(11, -27), S(17, -12), S(30, -2), S(38, 5), // Bishops
		S(41, 9), S(47, 12), S(47, 18), S(47, 18), S(66, 0), S(91, 4),
		S(50, 23), S(66, 8)},
	{S(-28, -35), S(-36, -29), S(-18, 10), S(-7, 32), S(3, 45), S(10, 56), // Rooks
		S(9, 62), S(23, 64), S(22, 65), S(36, 67), S(42, 69), S(43, 73),
		S(50, 75), S(62, 71), S(92, 62)},
	{S(-22, -20), S(-48, -8), S(-4, -141), S(-10, -118), S(-4, -19), S(3, -23), // Queens
		S(2, -21), S(12, -5), S(13, 19), S(21, 22), S(18, 36), S(21, 46),
		S(27, 37), S(27, 61), S(29, 64), S(28, 75), S(31, 69), S(26, 71),
		S(37, 67), S(37, 80), S(65, 55), S(65, 55), S(69, 46), S(61, 33),
		S(68, 20), S(28, 36), S(-7, -6), S(5, 3)},
}

var passedFriendlyDistance = [8]Score{
	S(0, 0), S(2, 31), S(-1, 15), S(-2, -9),
	S(-12, -19), S(-17, -19), S(8, -26), S(-31, -17),
}

var passedEnemyDistance = [8]Score{
	S(0, 0), S(-58, -63), S(30, -34), S(14, 5),
	S(8, 27), S(8, 34), S(0, 39), S(-12, 47),
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
var passedRank = [7]Score{S(0, 0), S(-6, -34), S(-9, -15), S(-4, 26), S(21, 74), S(36, 161), S(124, 257)}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{S(3, 22), S(-14, 25), S(-19, 15), S(-24, -2),
	S(-21, -3), S(5, 0), S(-14, 17), S(-10, 17),
}

var passedStacked = [8]Score{S(0, 0), S(17, -37), S(-6, -34), S(-32, -34), S(-68, -35), S(1, -16), S(0, 0), S(0, 0)}

var isolated = S(-9, -11)
var doubled = S(-12, -30)
var backward = S(4, 0)
var backwardOpen = S(-15, -3)

var bishopPair = S(34, 63)
var bishopRammedPawns = S(-9, -13)

var bishopOutpostUndefendedBonus = S(52, 1)
var bishopOutpostDefendedBonus = S(80, 9)

var knightOutpostUndefendedBonus = S(35, -13)
var knightOutpostDefendedBonus = S(58, 13)

var distantKnight = [4]Score{S(-19, 14), S(-22, 3), S(-25, -5), S(-18, -2)}

var minorBehindPawn = S(3, 26)

var tempo = S(40, 35)

// Rook on semiopen, open file
var rookOnFile = [2]Score{S(4, 29), S(51, -4)}

var kingDefenders = [12]Score{
	S(-80, 6), S(-70, 9), S(-29, 1), S(-5, -3),
	S(5, -3), S(19, -4), S(29, 0), S(36, 0),
	S(41, 0), S(41, 0), S(11, 0), S(11, 0),
}

var kingShelter = [2][8][8]Score{
	{{S(-29, 3), S(-8, -9), S(0, 7), S(33, -8),
		S(3, -17), S(2, 1), S(5, -10), S(-29, 9)},
		{S(21, -1), S(34, -13), S(-4, -2), S(-2, 6),
			S(-24, -2), S(8, -8), S(19, -39), S(-29, 6)},
		{S(10, 7), S(1, 2), S(-15, 4), S(-19, 7),
			S(-35, 0), S(-15, 1), S(-8, -6), S(-19, 1)},
		{S(-16, 21), S(6, 6), S(-8, -6), S(2, -2),
			S(7, -22), S(0, -13), S(13, -41), S(-20, 0)},
		{S(0, 7), S(-1, 3), S(-21, 0), S(-23, 9),
			S(-14, -8), S(-23, 2), S(-24, -1), S(-24, 5)},
		{S(46, -8), S(22, -12), S(0, -9), S(3, -10),
			S(11, -22), S(5, -5), S(35, -29), S(-7, 1)},
		{S(25, -2), S(-5, -5), S(-20, -8), S(-8, -4),
			S(-11, -10), S(11, -5), S(6, -17), S(-31, 13)},
		{S(-22, 1), S(-23, 0), S(-8, 5), S(-16, 13),
			S(-6, 7), S(-16, 17), S(-36, 9), S(-60, 33)}},
	{{S(0, -2), S(-42, -17), S(-18, -5), S(-78, -22),
		S(-5, -17), S(-33, -14), S(-77, 1), S(-67, 19)},
		{S(6, 23), S(2, -15), S(-20, -5), S(-6, -4),
			S(-3, -1), S(12, -39), S(2, -16), S(-57, 20)},
		{S(14, 31), S(42, -5), S(12, -5), S(14, -9),
			S(14, 3), S(-21, -9), S(55, -14), S(-28, 11)},
		{S(6, 24), S(-31, 19), S(-19, 9), S(-20, 4),
			S(-23, 19), S(-71, 31), S(-23, -1), S(-44, 3)},
		{S(2, 54), S(1, 4), S(-1, 0), S(-12, -1),
			S(-12, 7), S(1, -10), S(1, -12), S(-35, 9)},
		{S(69, -10), S(22, -10), S(-14, 1), S(-5, -14),
			S(-4, -6), S(-27, -12), S(17, -28), S(-31, 8)},
		{S(1, 11), S(7, -16), S(0, -16), S(-20, -8),
			S(-17, -12), S(-2, -16), S(-1, -19), S(-70, 23)},
		{S(1, 0), S(-5, -26), S(-4, -13), S(-23, -8),
			S(-21, -4), S(3, -4), S(-30, -25), S(-62, 26)}},
}

var kingStorm = [2][4][8]Score{
	{{S(19, 1), S(14, -2), S(17, 2), S(3, 8),
		S(1, 14), S(13, 9), S(3, 15), S(5, -10)},
		{S(15, 3), S(12, 4), S(27, -1), S(4, 8),
			S(13, 6), S(17, -1), S(8, -4), S(0, -12)},
		{S(16, 12), S(12, 6), S(1, 9), S(-9, 16),
			S(-4, 10), S(6, 0), S(24, -16), S(9, -6)},
		{S(24, 10), S(4, 4), S(9, 3), S(-3, 6),
			S(-6, 11), S(9, 8), S(5, 6), S(-1, 0)}},
	{{S(0, 0), S(8, 15), S(-17, 7), S(21, -4),
		S(16, 12), S(-9, 13), S(7, 44), S(11, -21)},
		{S(0, 0), S(7, -28), S(-6, -6), S(58, -11),
			S(48, -15), S(-19, 1), S(2, 19), S(4, -19)},
		{S(0, 0), S(-73, 0), S(-33, -1), S(13, 2),
			S(2, 1), S(-3, -11), S(70, -50), S(3, -3)},
		{S(0, 0), S(0, -21), S(12, -17), S(-8, 0),
			S(-4, 2), S(8, -20), S(-3, 1), S(-7, 16)}},
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
var blackForwardRanksMask [8]uint64

var whiteForwardFileMask [64]uint64
var blackForwardFileMask [64]uint64

// Outpost bitboards
const whiteOutpustRanks = RANK_4_BB | RANK_5_BB | RANK_6_BB
const blackOutpustRanks = RANK_5_BB | RANK_4_BB | RANK_3_BB

var kingSafetyAttacksWeights = [King + 1]int16{0, -1, -3, -4, -2, 0}
var kingSafetyAttackValue int16 = 121
var kingSafetyWeakSquares int16 = 27
var kingSafetyFriendlyPawns int16 = -2
var kingSafetyNoEnemyQueens int16 = -157
var kingSafetySafeQueenCheck int16 = 72
var kingSafetySafeRookCheck int16 = 125
var kingSafetySafeBishopCheck int16 = 94
var kingSafetySafeKnightCheck int16 = 134
var kingSafetyAdjustment int16 = -31

var hanging = S(75, 54)
var threatByKing = S(19, 45)
var threatByMinor = [King + 1]Score{S(0, 0), S(24, 23), S(30, 34), S(64, 22), S(49, 45), S(79, 161)}
var threatByRook = [King + 1]Score{S(0, 0), S(-1, 19), S(-1, 28), S(6, 46), S(91, 31), S(51, 41)}

func loadScoresToPieceSquares() {
	for x := 0; x < 4; x++ {
		for y := 0; y < 8; y++ {
			whiteKnightsPos[y*8+x] = pieceScores[Knight][y][x] + KnightValue
			whiteKnightsPos[y*8+(7-x)] = pieceScores[Knight][y][x] + KnightValue
			blackKnightsPos[(7-y)*8+x] = pieceScores[Knight][y][x] + KnightValue
			blackKnightsPos[(7-y)*8+(7-x)] = pieceScores[Knight][y][x] + KnightValue

			whiteBishopsPos[y*8+x] = pieceScores[Bishop][y][x] + BishopValue
			whiteBishopsPos[y*8+(7-x)] = pieceScores[Bishop][y][x] + BishopValue
			blackBishopsPos[(7-y)*8+x] = pieceScores[Bishop][y][x] + BishopValue
			blackBishopsPos[(7-y)*8+(7-x)] = pieceScores[Bishop][y][x] + BishopValue

			whiteRooksPos[y*8+x] = pieceScores[Rook][y][x] + RookValue
			whiteRooksPos[y*8+(7-x)] = pieceScores[Rook][y][x] + RookValue
			blackRooksPos[(7-y)*8+x] = pieceScores[Rook][y][x] + RookValue
			blackRooksPos[(7-y)*8+(7-x)] = pieceScores[Rook][y][x] + RookValue

			whiteQueensPos[y*8+x] = pieceScores[Queen][y][x] + QueenValue
			whiteQueensPos[y*8+(7-x)] = pieceScores[Queen][y][x] + QueenValue
			blackQueensPos[(7-y)*8+x] = pieceScores[Queen][y][x] + QueenValue
			blackQueensPos[(7-y)*8+(7-x)] = pieceScores[Queen][y][x] + QueenValue

			whiteKingPos[y*8+x] = pieceScores[King][y][x]
			whiteKingPos[y*8+(7-x)] = pieceScores[King][y][x]
			blackKingPos[(7-y)*8+x] = pieceScores[King][y][x]
			blackKingPos[(7-y)*8+(7-x)] = pieceScores[King][y][x]

			if y != 7 {
				whitePawnsConnected[y*8+x] = pawnsConnected[y][x]
				whitePawnsConnected[y*8+(7-x)] = pawnsConnected[y][x]
				blackPawnsConnected[(7-y)*8+x] = pawnsConnected[y][x]
				blackPawnsConnected[(7-y)*8+(7-x)] = pawnsConnected[y][x]
			}
		}
	}

	for y := 1; y < 7; y++ {
		for x := 0; x < 8; x++ {
			whitePawnsPos[y*8+x] = pawnScores[y][x] + PawnValue
			blackPawnsPos[(7-y)*8+x] = pawnScores[y][x] + PawnValue
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
		whitePawnsConnectedMask[i] = PawnAttacks[Black][i] | PawnAttacks[Black][i+8]
		blackPawnsConnectedMask[i] = PawnAttacks[White][i] | PawnAttacks[White][i-8]
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
		blackForwardRanksMask[rank] = (^whiteForwardRanksMask[rank]) | RANKS[rank]
	}

	for y := 0; y < 64; y++ {
		whiteForwardFileMask[y] = whiteForwardRanksMask[Rank(y)] & FILES[File(y)] & ^SquareMask[y]
		blackForwardFileMask[y] = blackForwardRanksMask[Rank(y)] & FILES[File(y)] & ^SquareMask[y]
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

		score += whitePawnsPos[fromId]
		if tuning {
			T.PawnValue++
			T.PawnScores[Rank(fromId)][File(fromId)]++
		}

		// Passed bonus
		if whitePassedMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			// Bonus is calculated based on rank, file, distance from friendly and enemy king
			score +=
				passedRank[Rank(fromId)] +
					passedFile[File(fromId)] +
					passedFriendlyDistance[distanceBetween[whiteKingLocation][fromId]] +
					passedEnemyDistance[distanceBetween[blackKingLocation][fromId]]

			if tuning {
				T.PassedRank[Rank(fromId)]++
				T.PassedFile[File(fromId)]++
				T.PassedFriendlyDistance[distanceBetween[whiteKingLocation][fromId]]++
				T.PassedEnemyDistance[distanceBetween[blackKingLocation][fromId]]++
			}

			if pos.Pieces[Pawn]&pos.Colours[White]&whiteForwardFileMask[fromId] != 0 {
				score += passedStacked[Rank(fromId)]
				if tuning {
					T.PassedStacked[Rank(fromId)]++
				}
			}
		}

		// Isolated pawn penalty
		if adjacentFilesMask[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			score += isolated
			if tuning {
				T.Isolated++
			}
		}

		// Pawn is backward if there are no pawns behind it and cannot increase rank without being attacked by enemy pawn
		if blackPassedMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 &&
			PawnAttacks[White][fromId+8]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
			if FILES[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
				score += backwardOpen
				if tuning {
					T.BackwardOpen++
				}
			} else {
				score += backward
				if tuning {
					T.Backward++
				}
			}
		} else if whitePawnsConnectedMask[fromId]&(pos.Colours[White]&pos.Pieces[Pawn]) != 0 {
			score += whitePawnsConnected[fromId]
			if tuning {
				T.PawnsConnected[Rank(fromId)][FileMirror[File(fromId)]]++
			}
		}
	}

	// white doubled pawns
	score += Score(PopCount(pos.Pieces[Pawn]&pos.Colours[White]&South(pos.Pieces[Pawn]&pos.Colours[White]))) * doubled
	if tuning {
		T.Doubled += PopCount(pos.Pieces[Pawn] & pos.Colours[White] & South(pos.Pieces[Pawn]&pos.Colours[White]))
	}

	// black pawns
	for fromBB = pos.Pieces[Pawn] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		score -= blackPawnsPos[fromId]

		if tuning {
			T.PawnValue--
			T.PawnScores[7-Rank(fromId)][File(fromId)]--
		}
		if blackPassedMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			score -=
				passedRank[7-Rank(fromId)] +
					passedFile[File(fromId)] +
					passedFriendlyDistance[distanceBetween[blackKingLocation][fromId]] +
					passedEnemyDistance[distanceBetween[whiteKingLocation][fromId]]
			if tuning {
				T.PassedRank[7-Rank(fromId)]--
				T.PassedFile[File(fromId)]--
				T.PassedFriendlyDistance[distanceBetween[blackKingLocation][fromId]]--
				T.PassedEnemyDistance[distanceBetween[whiteKingLocation][fromId]]--
			}

			if pos.Pieces[Pawn]&pos.Colours[Black]&blackForwardFileMask[fromId] != 0 {
				score -= passedStacked[7-Rank(fromId)]
				if tuning {
					T.PassedStacked[7-Rank(fromId)]--
				}
			}
		}
		if adjacentFilesMask[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			score -= isolated
			if tuning {
				T.Isolated--
			}
		}
		if whitePassedMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 &&
			PawnAttacks[Black][fromId-8]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
			if FILES[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
				score -= backwardOpen
				if tuning {
					T.BackwardOpen--
				}
			} else {
				score -= backward
				if tuning {
					T.Backward--
				}
			}
		} else if blackPawnsConnectedMask[fromId]&(pos.Colours[Black]&pos.Pieces[Pawn]) != 0 {
			score -= blackPawnsConnected[fromId]
			if tuning {
				T.PawnsConnected[7-Rank(fromId)][FileMirror[File(fromId)]]--
			}
		}
	}

	// black doubled pawns
	score -= Score(PopCount(pos.Pieces[Pawn]&pos.Colours[Black]&North(pos.Pieces[Pawn]&pos.Colours[Black]))) * doubled
	if tuning {
		T.Doubled -= PopCount(pos.Pieces[Pawn] & pos.Colours[Black] & North(pos.Pieces[Pawn]&pos.Colours[Black]))
	}

	// White king storm shelter
	for file := Max(File(whiteKingLocation)-1, FILE_A); file <= Min(File(whiteKingLocation)+1, FILE_H); file++ {
		ours := pos.Pieces[Pawn] & FILES[file] & pos.Colours[White] & whiteForwardRanksMask[Rank(whiteKingLocation)]
		var ourDist int
		if ours == 0 {
			ourDist = 7
		} else {
			ourDist = Abs(Rank(whiteKingLocation) - Rank(BitScan(ours)))
		}
		theirs := pos.Pieces[Pawn] & FILES[file] & pos.Colours[Black] & whiteForwardRanksMask[Rank(whiteKingLocation)]
		var theirDist int
		if theirs == 0 {
			theirDist = 7
		} else {
			theirDist = Abs(Rank(whiteKingLocation) - Rank(BitScan(theirs)))
		}
		sameFile := BoolToInt(file == File(whiteKingLocation))
		score += kingShelter[sameFile][file][ourDist]
		if tuning {
			T.KingShelter[sameFile][file][ourDist]++
		}

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		score += kingStorm[blocked][FileMirror[file]][theirDist]

		if tuning {
			T.KingStorm[blocked][FileMirror[file]][theirDist]++
		}
	}

	// Black king storm / shelter
	for file := Max(File(blackKingLocation)-1, FILE_A); file <= Min(File(blackKingLocation)+1, FILE_H); file++ {
		ours := pos.Pieces[Pawn] & FILES[file] & pos.Colours[Black] & blackForwardRanksMask[Rank(blackKingLocation)]
		var ourDist int
		if ours == 0 {
			ourDist = 7
		} else {
			ourDist = Abs(Rank(blackKingLocation) - Rank(MostSignificantBit(ours)))
		}
		theirs := pos.Pieces[Pawn] & FILES[file] & pos.Colours[White] & blackForwardRanksMask[Rank(blackKingLocation)]
		var theirDist int
		if theirs == 0 {
			theirDist = 7
		} else {
			theirDist = Abs(Rank(blackKingLocation) - Rank(MostSignificantBit(theirs)))
		}
		sameFile := BoolToInt(file == File(blackKingLocation))
		score -= kingShelter[sameFile][file][ourDist]
		if tuning {
			T.KingShelter[sameFile][file][ourDist]--
		}

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		score -= kingStorm[blocked][FileMirror[file]][theirDist]
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

	phase := totalPhase
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
		phase -= knightPhase
		fromId = BitScan(fromBB)

		attacks = KnightAttacks[fromId]
		mobility := PopCount(whiteMobilityArea & attacks)
		score += whiteKnightsPos[fromId]
		score += mobilityBonus[0][mobility]
		if tuning {
			T.KnightValue++
			T.PieceScores[Knight][Rank(fromId)][FileMirror[File(fromId)]]++
			T.MobilityBonus[0][mobility]++
		}

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Knight] |= attacks

		if (pos.Pieces[Pawn]>>8)&SquareMask[fromId] != 0 {
			score += minorBehindPawn
			if tuning {
				T.MinorBehindPawn++
			}
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && whiteOutpostMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			if PawnAttacks[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
				score += knightOutpostDefendedBonus
				if tuning {
					T.KnightOutpostDefendedBonus++
				}
			} else {
				score += knightOutpostUndefendedBonus
				if tuning {
					T.KnightOutpostUndefendedBonus++
				}
			}
		}

		kingDistance := Min(int(distanceBetween[fromId][whiteKingLocation]), int(distanceBetween[fromId][blackKingLocation]))
		if kingDistance >= 4 {
			score += distantKnight[kingDistance-4]
			if tuning {
				T.DistantKnight[kingDistance-4]++
			}
		}
		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += kingSafetyAttacksWeights[Knight]
		}
	}

	// black knights
	for fromBB = pos.Pieces[Knight] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= knightPhase
		fromId = BitScan(fromBB)

		attacks = KnightAttacks[fromId]
		mobility := PopCount(blackMobilityArea & attacks)
		score -= blackKnightsPos[fromId]
		score -= mobilityBonus[0][mobility]
		if tuning {
			T.KnightValue--
			T.PieceScores[Knight][7-Rank(fromId)][FileMirror[File(fromId)]]--
			T.MobilityBonus[0][mobility]--
		}

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Knight] |= attacks

		if (pos.Pieces[Pawn]<<8)&SquareMask[fromId] != 0 {
			score -= minorBehindPawn
			if tuning {
				T.MinorBehindPawn--
			}
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && blackOutpostMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			if PawnAttacks[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
				score -= knightOutpostDefendedBonus
				if tuning {
					T.KnightOutpostDefendedBonus--
				}
			} else {
				score -= knightOutpostUndefendedBonus
				if tuning {
					T.KnightOutpostUndefendedBonus--
				}
			}
		}
		kingDistance := Min(int(distanceBetween[fromId][whiteKingLocation]), int(distanceBetween[fromId][blackKingLocation]))
		if kingDistance >= 4 {
			score -= distantKnight[kingDistance-4]
			if tuning {
				T.DistantKnight[kingDistance-4]--
			}
		}
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += kingSafetyAttacksWeights[Knight]
		}
	}

	// white bishops
	whiteRammedPawns := South(pos.Pieces[Pawn]&pos.Colours[Black]) & (pos.Pieces[Pawn] & pos.Colours[White])
	for fromBB = pos.Pieces[Bishop] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= bishopPhase
		fromId = BitScan(fromBB)

		attacks = BishopAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		score += mobilityBonus[1][mobility]
		score += whiteBishopsPos[fromId]
		if tuning {
			T.BishopValue++
			T.PieceScores[Bishop][Rank(fromId)][FileMirror[File(fromId)]]++
			T.MobilityBonus[1][mobility]++
		}

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Bishop] |= attacks

		if (pos.Pieces[Pawn]>>8)&SquareMask[fromId] != 0 {
			score += minorBehindPawn
			if tuning {
				T.MinorBehindPawn++
			}
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && whiteOutpostMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			if PawnAttacks[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
				score += bishopOutpostDefendedBonus
				if tuning {
					T.BishopOutpostDefendedBonus++
				}
			} else {
				score += bishopOutpostUndefendedBonus
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
		score += bishopRammedPawns * rammedCount
		if tuning {
			T.BishopRammedPawns += int(rammedCount)
		}
		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += kingSafetyAttacksWeights[Bishop]
		}
	}

	// Bishop pair bonus
	// It is not checked if bishops have opposite colors, but that is almost always the case
	if MoreThanOne(pos.Pieces[Bishop] & pos.Colours[White]) {
		score += bishopPair
		if tuning {
			T.BishopPair++
		}
	}

	// black bishops
	blackRammedPawns := North(pos.Pieces[Pawn]&pos.Colours[White]) & (pos.Pieces[Pawn] & pos.Colours[Black])
	for fromBB = pos.Pieces[Bishop] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= bishopPhase
		fromId = BitScan(fromBB)

		attacks = BishopAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		score -= mobilityBonus[1][mobility]
		score -= blackBishopsPos[fromId]
		if tuning {
			T.BishopValue--
			T.PieceScores[Bishop][7-Rank(fromId)][FileMirror[File(fromId)]]--
			T.MobilityBonus[1][mobility]--
		}

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Bishop] |= attacks

		if (pos.Pieces[Pawn]<<8)&SquareMask[fromId] != 0 {
			score -= minorBehindPawn
			if tuning {
				T.MinorBehindPawn--
			}
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && blackOutpostMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			if PawnAttacks[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
				score -= bishopOutpostDefendedBonus
				if tuning {
					T.BishopOutpostDefendedBonus--
				}
			} else {
				score -= bishopOutpostUndefendedBonus
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
		score -= bishopRammedPawns * rammedCount
		if tuning {
			T.BishopRammedPawns -= int(rammedCount)
		}
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += kingSafetyAttacksWeights[Bishop]
		}
	}

	if MoreThanOne(pos.Pieces[Bishop] & pos.Colours[Black]) {
		score -= bishopPair

		if tuning {
			T.BishopPair--
		}
	}

	// white rooks
	for fromBB = pos.Pieces[Rook] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= rookPhase
		fromId = BitScan(fromBB)

		attacks = RookAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		score += mobilityBonus[2][mobility]
		score += whiteRooksPos[fromId]

		if tuning {
			T.RookValue++
			T.PieceScores[Rook][Rank(fromId)][FileMirror[File(fromId)]]++
			T.MobilityBonus[2][mobility]++
		}

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Rook] |= attacks

		if pos.Pieces[Pawn]&FILES[File(fromId)] == 0 {
			score += rookOnFile[1]
			if tuning {
				T.RookOnFile[1]++
			}
		} else if (pos.Pieces[Pawn]&pos.Colours[White])&FILES[File(fromId)] == 0 {
			score += rookOnFile[0]
			if tuning {
				T.RookOnFile[0]++
			}
		}

		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += kingSafetyAttacksWeights[Rook]
		}
	}

	// black rooks
	for fromBB = pos.Pieces[Rook] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= rookPhase
		fromId = BitScan(fromBB)

		attacks = RookAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		score -= mobilityBonus[2][mobility]
		score -= blackRooksPos[fromId]

		if tuning {
			T.RookValue--
			T.PieceScores[Rook][7-Rank(fromId)][FileMirror[File(fromId)]]--
			T.MobilityBonus[2][mobility]--
		}

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Rook] |= attacks

		if pos.Pieces[Pawn]&FILES[File(fromId)] == 0 {
			score -= rookOnFile[1]
			if tuning {
				T.RookOnFile[1]--
			}
		} else if (pos.Pieces[Pawn]&pos.Colours[Black])&FILES[File(fromId)] == 0 {
			score -= rookOnFile[0]
			if tuning {
				T.RookOnFile[0]--
			}
		}

		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += kingSafetyAttacksWeights[Rook]
		}
	}

	//white queens
	for fromBB = pos.Pieces[Queen] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= queenPhase
		fromId = BitScan(fromBB)

		attacks = QueenAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		score += mobilityBonus[3][mobility]
		score += whiteQueensPos[fromId]

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
			whiteKingAttackersWeight += kingSafetyAttacksWeights[Queen]
		}
	}

	// black queens
	for fromBB = pos.Pieces[Queen] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= queenPhase
		fromId = BitScan(fromBB)

		attacks = QueenAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		score -= mobilityBonus[3][mobility]
		score -= blackQueensPos[fromId]

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
			blackKingAttackersWeight += kingSafetyAttacksWeights[Queen]
		}
	}

	// tempo bonus
	if pos.SideToMove == White {
		score += tempo
		if tuning {
			T.Tempo++
		}
	} else {
		score -= tempo
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
	score += whiteKingPos[whiteKingLocation]
	score += kingDefenders[whiteKingDefenders]
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
		count += int(kingSafetyAttackValue) * 9 * int(blackKingAttackersCount) / PopCount(whiteKingArea)
		count += int(kingSafetyWeakSquares) * PopCount(whiteKingArea&weakForWhite)
		count += int(kingSafetyFriendlyPawns) * PopCount(pos.Colours[White]&pos.Pieces[Pawn]&whiteKingArea & ^weakForWhite)
		count += int(kingSafetyNoEnemyQueens) * BoolToInt(pos.Colours[Black]&pos.Pieces[Queen] == 0)
		count += int(kingSafetySafeQueenCheck) * PopCount(queenChecks)
		count += int(kingSafetySafeRookCheck) * PopCount(rookChecks)
		count += int(kingSafetySafeBishopCheck) * PopCount(bishopChecks)
		count += int(kingSafetySafeKnightCheck) * PopCount(knightChecks)
		count += int(kingSafetyAdjustment)
		if count > 0 {
			score -= S(int16(count*count/720), int16(count/20))
		}
	}

	// black king
	blackKingDefenders := PopCount(
		(pos.Pieces[Pawn] | pos.Pieces[Bishop] | pos.Pieces[Knight]) & pos.Colours[Black] & blackKingAreaMask[blackKingLocation],
	)
	score -= blackKingPos[blackKingLocation]
	score -= kingDefenders[blackKingDefenders]
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
		count += int(kingSafetyAttackValue) * int(whiteKingAttackersCount) * 9 / PopCount(blackKingArea) // Scale value to king area size
		count += int(kingSafetyWeakSquares) * PopCount(blackKingArea&weakForBlack)
		count += int(kingSafetyFriendlyPawns) * PopCount(pos.Colours[Black]&pos.Pieces[Pawn]&blackKingArea & ^weakForBlack)
		count += int(kingSafetyNoEnemyQueens) * BoolToInt(pos.Colours[White]&pos.Pieces[Queen] == 0)
		count += int(kingSafetySafeQueenCheck) * PopCount(queenChecks)
		count += int(kingSafetySafeRookCheck) * PopCount(rookChecks)
		count += int(kingSafetySafeBishopCheck) * PopCount(bishopChecks)
		count += int(kingSafetySafeKnightCheck) * PopCount(knightChecks)
		count += int(kingSafetyAdjustment)
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
			score += threatByMinor[threatenedPiece]
			if tuning {
				T.ThreatByMinor[threatenedPiece]++
			}
		}

		for fromBB = pos.Colours[Black] & (blackDefended | weakForBlack) & whiteAttackedBy[Rook] & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) & ^pos.Pieces[Pawn] {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score += threatByRook[threatenedPiece]
			if tuning {
				T.ThreatByRook[threatenedPiece]++
			}
		}

		if weakForBlack&pos.Colours[Black]&whiteAttackedBy[King] != 0 {
			score += threatByKing
			if tuning {
				T.ThreatByKing++
			}
		}

		// Bonus if enemy has a hanging piece
		score += hanging *
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
			score -= threatByMinor[threatenedPiece]
			if tuning {
				T.ThreatByMinor[threatenedPiece]--
			}
		}

		for fromBB = pos.Colours[White] & (whiteDefended | weakForWhite) & blackAttackedBy[Rook] & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score -= threatByRook[threatenedPiece]
			if tuning {
				T.ThreatByRook[threatenedPiece]--
			}
		}

		if weakForWhite&pos.Colours[White]&blackAttackedBy[King] != 0 {
			score -= threatByKing
			if tuning {
				T.ThreatByKing--
			}
		}

		// Bonus if enemy has a hanging piece
		score -= hanging *
			Score(PopCount(pos.Colours[White] & ^pos.Pieces[Pawn] & blackAttackedByTwo & weakForWhite))

		if tuning {
			T.Hanging -= PopCount(pos.Colours[White] & ^pos.Pieces[Pawn] & blackAttackedByTwo & weakForWhite)
		}
	}

	scale := scaleFactor(pos, score.End())

	// tapering eval
	phase = (phase*256 + (totalPhase / 2)) / totalPhase
	result := (int(score.Middle())*(256-phase) + (int(score.End()) * phase * scale / SCALE_NORMAL)) / 256

	if pos.SideToMove == White {
		return result
	}
	return -result
}

const SCALE_NORMAL = 1
const SCALE_DRAW = 0

func scaleFactor(pos *Position, endResult int16) int {
	if (endResult > 0 && PopCount(pos.Colours[White]) == 2 && (pos.Colours[White]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) ||
		(endResult < 0 && PopCount(pos.Colours[Black]) == 2 && (pos.Colours[Black]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) {
		return SCALE_DRAW
	}
	return SCALE_NORMAL
}
