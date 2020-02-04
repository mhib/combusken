package evaluation

import . "github.com/mhib/combusken/backend"
import . "github.com/mhib/combusken/utils"
import "github.com/mhib/combusken/transposition"

const pawnPhase = 0
const knightPhase = 1
const bishopPhase = 1
const rookPhase = 2
const queenPhase = 4
const totalPhase = pawnPhase*16 + knightPhase*4 + bishopPhase*4 + rookPhase*4 + queenPhase*2

var PawnValue = S(100, 121)
var KnightValue = S(465, 434)
var BishopValue = S(422, 428)
var RookValue = S(626, 692)
var QueenValue = S(1426, 1305)

// Piece Square Values
var pieceScores = [King + 1][8][4]Score{
	{},
	{ // knight
		{S(-50, -44), S(-7, -56), S(-28, -39), S(-1, -28)},
		{S(-9, -65), S(-17, -39), S(-9, -30), S(0, -19)},
		{S(-18, -41), S(2, -30), S(0, -15), S(12, -2)},
		{S(-17, -32), S(16, -28), S(13, -1), S(9, 4)},
		{S(-10, -35), S(6, -23), S(10, 4), S(26, 4)},
		{S(-52, -49), S(12, -45), S(-87, 22), S(14, -3)},
		{S(-74, -54), S(-45, -26), S(65, -56), S(-1, -18)},
		{S(-213, -70), S(-81, -74), S(-136, -29), S(13, -60)},
	},
	{ // Bishop
		{S(16, -11), S(31, -4), S(33, -10), S(32, -4)},
		{S(26, -25), S(66, -31), S(47, -15), S(29, -6)},
		{S(34, -19), S(50, -12), S(49, -3), S(34, 5)},
		{S(3, -29), S(10, -27), S(25, -15), S(28, -11)},
		{S(-34, -18), S(15, -30), S(-28, -6), S(24, -12)},
		{S(-128, 9), S(-70, -12), S(-221, 42), S(-67, -6)},
		{S(-61, 7), S(37, -9), S(3, 2), S(18, -12)},
		{S(-2, -18), S(-27, -16), S(-119, 2), S(-98, 3)},
	},
	{ // Rook
		{S(-2, -21), S(-10, -7), S(12, -13), S(14, -15)},
		{S(-42, -2), S(-2, -20), S(-2, -15), S(7, -15)},
		{S(-37, -7), S(-11, -10), S(2, -20), S(-10, -13)},
		{S(-37, 3), S(-10, -3), S(-14, 0), S(-13, -3)},
		{S(-46, 13), S(-29, 5), S(4, 9), S(-10, 3)},
		{S(-31, 7), S(8, 4), S(8, 2), S(-18, 7)},
		{S(-4, 16), S(3, 17), S(56, 0), S(56, -5)},
		{S(-3, 15), S(20, 6), S(-35, 23), S(38, 5)},
	},
	{ // Queen
		{S(4, -64), S(19, -79), S(21, -69), S(36, -69)},
		{S(-4, -42), S(18, -61), S(32, -51), S(30, -51)},
		{S(-5, -15), S(23, -38), S(0, 4), S(6, -11)},
		{S(-1, -18), S(-22, 32), S(-5, 14), S(-16, 34)},
		{S(-23, 17), S(-30, 30), S(-32, 33), S(-32, 47)},
		{S(18, -30), S(-6, -4), S(-5, 15), S(-16, 53)},
		{S(-9, -24), S(-52, 14), S(-21, 25), S(-42, 57)},
		{S(0, -10), S(-14, 12), S(7, 17), S(28, 11)},
	},
	{ // King
		{S(146, -11), S(126, 31), S(46, 75), S(63, 55)},
		{S(145, 29), S(99, 61), S(45, 88), S(7, 95)},
		{S(78, 47), S(99, 64), S(38, 87), S(30, 92)},
		{S(27, 42), S(108, 55), S(36, 91), S(8, 97)},
		{S(55, 50), S(140, 72), S(97, 93), S(97, 86)},
		{S(112, 58), S(257, 70), S(247, 85), S(178, 63)},
		{S(59, 62), S(131, 80), S(118, 101), S(180, 70)},
		{S(24, 5), S(146, 38), S(116, 75), S(31, 55)},
	},
}

// Pawns Square scores
var pawnScores = [7][8]Score{
	{},
	{S(-27, 11), S(9, -1), S(-15, 12), S(9, 2), S(-3, 18), S(13, 6), S(31, -3), S(-8, -3)},
	{S(-22, -3), S(-18, 1), S(-2, -4), S(6, -4), S(-1, -3), S(-5, 2), S(12, -11), S(-8, -6)},
	{S(-21, 6), S(-14, 4), S(14, -10), S(26, -15), S(22, -11), S(1, -1), S(3, -2), S(-22, 3)},
	{S(-3, 17), S(24, -2), S(9, -7), S(36, -17), S(36, -18), S(10, 2), S(43, 3), S(-13, 15)},
	{S(11, 42), S(20, 32), S(38, 17), S(38, 2), S(84, -14), S(92, 7), S(39, 26), S(-4, 46)},
	{S(14, 54), S(94, 49), S(-2, 44), S(-2, 28), S(88, 41), S(-22, 45), S(-26, 53), S(-134, 101)},
}

var pawnsConnected = [8][4]Score{
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
	{S(12, -21), S(8, 7), S(6, -9), S(5, 16)},
	{S(9, 1), S(28, 3), S(14, 9), S(16, 19)},
	{S(8, 8), S(24, 7), S(16, 11), S(18, 14)},
	{S(12, 16), S(5, 25), S(28, 24), S(30, 23)},
	{S(7, 58), S(30, 70), S(56, 67), S(58, 56)},
	{S(14, 56), S(150, -4), S(167, 30), S(228, 67)},
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
}

var mobilityBonus = [...][32]Score{
	{S(-47, -121), S(-38, -60), S(-27, -32), S(-26, -12), S(-9, -15), S(1, -6), // Knights
		S(11, -11), S(23, -15), S(34, -30)},
	{S(-36, -74), S(-18, -59), S(3, -31), S(8, -10), S(20, 1), S(31, 7), // Bishops
		S(37, 10), S(42, 12), S(46, 16), S(47, 17), S(69, 2), S(94, 3),
		S(48, 31), S(60, 17)},
	{S(-28, -35), S(-36, -33), S(-20, 9), S(-9, 35), S(-1, 48), S(4, 59), // Rooks
		S(8, 67), S(19, 66), S(20, 66), S(35, 68), S(40, 70), S(43, 73),
		S(53, 73), S(58, 71), S(66, 68)},
	{S(-24, -20), S(-52, -10), S(-2, -153), S(-11, -137), S(-6, -31), S(0, -38), // Queens
		S(-2, -35), S(7, -9), S(9, 28), S(13, 31), S(14, 52), S(14, 54),
		S(21, 48), S(20, 71), S(25, 70), S(26, 77), S(29, 75), S(26, 74),
		S(34, 70), S(46, 76), S(67, 54), S(64, 56), S(91, 35), S(72, 29),
		S(78, 16), S(42, 31), S(-7, -4), S(8, 5)},
}

var passedFriendlyDistance = [8]Score{
	S(0, 0), S(-14, 31), S(-11, 13), S(-10, -10),
	S(-10, -23), S(-9, -26), S(10, -29), S(-23, -19),
}

var passedEnemyDistance = [8]Score{
	S(0, 0), S(-41, -75), S(37, -30), S(13, 8),
	S(9, 28), S(0, 38), S(-3, 42), S(-22, 49),
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
var passedRank = [7]Score{S(0, 0), S(8, -36), S(1, -15), S(-2, 27), S(23, 74), S(20, 166), S(109, 258)}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{S(-5, 18), S(-19, 20), S(-24, 8), S(-26, -4),
	S(-21, -6), S(19, -6), S(-29, 22), S(-8, 16),
}

var isolated = S(-8, -11)
var doubled = S(-10, -34)
var backward = S(5, -2)
var backwardOpen = S(-15, -6)

var bishopPair = S(38, 63)
var bishopRammedPawns = S(-7, -13)

var bishopOutpostUndefendedBonus = S(91, 0)
var bishopOutpostDefendedBonus = S(111, 11)

var knightOutpostUndefendedBonus = S(47, -16)
var knightOutpostDefendedBonus = S(69, 13)

var minorBehindPawn = S(5, 28)

var tempo = S(32, 29)

// Rook on semiopen, open file
var rookOnFile = [2]Score{S(7, 26), S(55, -3)}

var kingDefenders = [12]Score{
	S(-80, 7), S(-76, 12), S(-33, 4), S(-8, -2),
	S(2, -1), S(17, -4), S(26, 0), S(33, -5),
	S(39, -16), S(33, 4), S(14, 0), S(11, 0),
}

var kingShelter = [2][8][8]Score{
	{{S(-29, 12), S(-6, -5), S(-5, 13), S(52, -14),
		S(17, -25), S(12, -5), S(4, -10), S(-41, 19)},
		{S(17, 4), S(42, -13), S(-4, -1), S(-19, 13),
			S(-25, -8), S(17, -16), S(28, -45), S(-47, 10)},
		{S(20, 8), S(2, 6), S(-23, 8), S(-18, 10),
			S(-37, 4), S(-7, -3), S(4, -16), S(-26, 1)},
		{S(-19, 24), S(1, 11), S(-10, -2), S(0, -1),
			S(7, -19), S(-4, -21), S(9, -46), S(-22, 1)},
		{S(-5, 8), S(-11, 11), S(-25, 5), S(-28, 14),
			S(-19, -4), S(-18, -2), S(-19, -18), S(-31, 9)},
		{S(39, -7), S(20, -12), S(5, -11), S(12, -10),
			S(15, -23), S(3, -2), S(25, -23), S(-11, 2)},
		{S(27, -4), S(-3, -4), S(-24, -8), S(-12, 0),
			S(-12, -9), S(17, -7), S(10, -22), S(-31, 13)},
		{S(-35, 12), S(-26, 2), S(-9, 8), S(-9, 11),
			S(5, 3), S(-10, 15), S(-35, 7), S(-70, 39)}},
	{{S(-2, -8), S(-51, -8), S(-23, -2), S(-87, -19),
		S(-1, -19), S(-37, -30), S(-68, -6), S(-73, 20)},
		{S(9, 32), S(14, -16), S(-14, -3), S(0, 0),
			S(-3, 1), S(26, -53), S(10, -18), S(-80, 26)},
		{S(19, 32), S(38, 1), S(14, 1), S(17, -9),
			S(18, 2), S(-27, -8), S(67, -21), S(-38, 12)},
		{S(6, 24), S(-35, 21), S(-16, 11), S(-15, 3),
			S(-29, 27), S(-81, 29), S(-27, -11), S(-44, 1)},
		{S(-6, 52), S(1, 10), S(1, 0), S(-6, 0),
			S(-12, 12), S(0, -16), S(2, -23), S(-34, 8)},
		{S(77, -21), S(17, -6), S(-14, 5), S(1, -10),
			S(-6, -3), S(-18, -16), S(6, -23), S(-42, 9)},
		{S(-6, 6), S(10, -10), S(3, -12), S(-16, -7),
			S(-17, -8), S(-3, -12), S(0, -12), S(-70, 27)},
		{S(2, -6), S(-9, -22), S(1, -15), S(-16, -9),
			S(-22, -3), S(17, -14), S(-28, -27), S(-75, 32)}},
}

var kingStorm = [2][4][8]Score{
	{{S(19, -26), S(14, -24), S(19, -21), S(7, -7),
		S(4, 15), S(-31, 57), S(41, 129), S(6, -31)},
		{S(23, -25), S(21, -24), S(38, -28), S(15, -12),
			S(11, 12), S(-94, 79), S(44, 138), S(18, -35)},
		{S(20, -1), S(17, -11), S(11, -6), S(-3, 0),
			S(-5, 9), S(-45, 57), S(12, 104), S(17, -19)},
		{S(13, -1), S(7, -18), S(10, -18), S(8, -16),
			S(-8, 2), S(-72, 66), S(19, 72), S(-1, -14)}},
	{{S(0, 0), S(4, -6), S(-13, -13), S(18, -20),
		S(11, 11), S(-6, 20), S(-6, 61), S(17, -45)},
		{S(0, 0), S(-4, -46), S(21, -34), S(86, -43),
			S(49, -18), S(-42, 15), S(-17, 24), S(9, -40)},
		{S(0, 0), S(-63, -17), S(-31, -11), S(16, -15),
			S(6, -7), S(-21, 0), S(79, -54), S(7, -13)},
		{S(0, 0), S(3, -42), S(16, -33), S(-2, -21),
			S(-8, -10), S(4, 2), S(0, 25), S(-3, 8)}},
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

// Outpost bitboards
const whiteOutpustRanks = RANK_4_BB | RANK_5_BB | RANK_6_BB
const blackOutpustRanks = RANK_5_BB | RANK_4_BB | RANK_3_BB

var kingSafetyAttacksWeights = [King + 1]int16{0, -3, 0, -1, 15, 0}
var kingSafetyAttackValue int16 = 105
var kingSafetyWeakSquares int16 = 17
var kingSafetyFriendlyPawns int16 = -7
var kingSafetyNoEnemyQueens int16 = 40
var kingSafetySafeQueenCheck int16 = 77
var kingSafetySafeRookCheck int16 = 113
var kingSafetySafeBishopCheck int16 = 97
var kingSafetySafeKnightCheck int16 = 147
var kingSafetyAdjustment int16 = -72

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

			whitePawnsConnected[y*8+x] = pawnsConnected[y][x]
			whitePawnsConnected[y*8+(7-x)] = pawnsConnected[y][x]
			blackPawnsConnected[(7-y)*8+x] = pawnsConnected[y][x]
			blackPawnsConnected[(7-y)*8+(7-x)] = pawnsConnected[y][x]
		}
	}

	for y := 1; y < 7; y++ {
		for x := 0; x < 8; x++ {
			whitePawnsPos[y*8+x] = pawnScores[y][x] + PawnValue
			blackPawnsPos[(7-y)*8+(7-x)] = pawnScores[y][x] + PawnValue
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
		blackForwardRanksMasks[rank] = (^whiteForwardRanksMask[rank]) | RANKS[rank]
	}
}

func IsLateEndGame(pos *Position) bool {
	return ((pos.Pieces[Rook]|pos.Pieces[Queen]|pos.Pieces[Bishop]|pos.Pieces[Knight])&pos.Colours[pos.SideToMove]) == 0
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

		// Passed bonus
		if whitePassedMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			// Bonus is calculated based on rank, file, distance from friendly and enemy king
			score +=
				passedRank[Rank(fromId)] +
					passedFile[File(fromId)] +
					passedFriendlyDistance[distanceBetween[whiteKingLocation][fromId]] +
					passedEnemyDistance[distanceBetween[blackKingLocation][fromId]]

		}
		// Isolated pawn penalty
		if adjacentFilesMask[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			score += isolated
		}

		// Pawn is backward if there are no pawns behind it and cannot increase rank without being attacked by enemy pawn
		if blackPassedMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 &&
			PawnAttacks[White][fromId+8]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
			if FILES[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
				score += backwardOpen
			} else {
				score += backward
			}
		} else if whitePawnsConnectedMask[fromId]&(pos.Colours[White]&pos.Pieces[Pawn]) != 0 {
			score += whitePawnsConnected[fromId]
		}
	}

	// white doubled pawns
	score += Score(PopCount(pos.Pieces[Pawn]&pos.Colours[White]&South(pos.Pieces[Pawn]&pos.Colours[White]))) * doubled

	// black pawns
	for fromBB = pos.Pieces[Pawn] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		score -= blackPawnsPos[fromId]
		if blackPassedMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			score -=
				passedRank[7-Rank(fromId)] +
					passedFile[File(fromId)] +
					passedFriendlyDistance[distanceBetween[blackKingLocation][fromId]] +
					passedEnemyDistance[distanceBetween[whiteKingLocation][fromId]]

		}
		if adjacentFilesMask[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			score -= isolated
		}
		if whitePassedMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 &&
			PawnAttacks[Black][fromId-8]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
			if FILES[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
				score -= backwardOpen
			} else {
				score -= backward
			}
		} else if blackPawnsConnectedMask[fromId]&(pos.Colours[Black]&pos.Pieces[Pawn]) != 0 {
			score -= blackPawnsConnected[fromId]
		}
	}

	// black doubled pawns
	score -= Score(PopCount(pos.Pieces[Pawn]&pos.Colours[Black]&North(pos.Pieces[Pawn]&pos.Colours[Black]))) * doubled

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

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		score += kingStorm[blocked][FileMirror[file]][theirDist]
	}

	// Black king storm / shelter
	for file := Max(File(blackKingLocation)-1, FILE_A); file <= Min(File(blackKingLocation)+1, FILE_H); file++ {
		ours := pos.Pieces[Pawn] & FILES[file] & pos.Colours[Black] & blackForwardRanksMasks[Rank(blackKingLocation)]
		var ourDist int
		if ours == 0 {
			ourDist = 7
		} else {
			ourDist = Abs(Rank(blackKingLocation) - Rank(MostSignificantBit(ours)))
		}
		theirs := pos.Pieces[Pawn] & FILES[file] & pos.Colours[White] & blackForwardRanksMasks[Rank(blackKingLocation)]
		var theirDist int
		if theirs == 0 {
			theirDist = 7
		} else {
			theirDist = Abs(Rank(blackKingLocation) - Rank(MostSignificantBit(theirs)))
		}
		sameFile := BoolToInt(file == File(blackKingLocation))
		score -= kingShelter[sameFile][file][ourDist]

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		score -= kingStorm[blocked][FileMirror[file]][theirDist]
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
	//whiteAttackedBy[King] |= attacks
	whiteKingArea := whiteKingAreaMask[whiteKingLocation]

	blackKingLocation := BitScan(pos.Pieces[King] & pos.Colours[Black])
	attacks = KingAttacks[blackKingLocation]
	blackAttacked |= attacks
	//blackAttackedBy[King] |= attacks
	blackKingArea := blackKingAreaMask[blackKingLocation]

	// white pawns
	attacks = WhitePawnsAttacks(pos.Pieces[Pawn] & pos.Colours[White])
	whiteAttackedByTwo |= whiteAttacked & attacks
	whiteAttacked |= attacks
	//whiteAttackedBy[Pawn] |= attacks
	whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))

	// black pawns
	attacks = BlackPawnsAttacks(pos.Pieces[Pawn] & pos.Colours[Black])
	blackAttackedByTwo |= blackAttacked & attacks
	blackAttacked |= attacks
	//blackAttackedBy[Pawn] |= attacks
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

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Knight] |= attacks

		if (pos.Pieces[Pawn]>>8)&SquareMask[fromId] != 0 {
			score += minorBehindPawn
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && whiteOutpostMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			if PawnAttacks[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
				score += knightOutpostDefendedBonus
			} else {
				score += knightOutpostUndefendedBonus
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

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Knight] |= attacks

		if (pos.Pieces[Pawn]<<8)&SquareMask[fromId] != 0 {
			score -= minorBehindPawn
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && blackOutpostMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			if PawnAttacks[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
				score -= knightOutpostDefendedBonus
			} else {
				score -= knightOutpostUndefendedBonus
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

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Bishop] |= attacks

		if (pos.Pieces[Pawn]>>8)&SquareMask[fromId] != 0 {
			score += minorBehindPawn
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && whiteOutpostMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			if PawnAttacks[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
				score += bishopOutpostDefendedBonus
			} else {
				score += bishopOutpostUndefendedBonus
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

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Bishop] |= attacks

		if (pos.Pieces[Pawn]<<8)&SquareMask[fromId] != 0 {
			score -= minorBehindPawn
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && blackOutpostMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			if PawnAttacks[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
				score -= bishopOutpostDefendedBonus
			} else {
				score -= bishopOutpostUndefendedBonus
			}
		}
		var rammedCount Score
		if SquareMask[fromId]&WHITE_SQUARES != 0 {
			rammedCount = Score(PopCount(blackRammedPawns & WHITE_SQUARES))
		} else {
			rammedCount = Score(PopCount(blackRammedPawns & BLACK_SQUARES))
		}
		score -= bishopRammedPawns * rammedCount
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += kingSafetyAttacksWeights[Bishop]
		}
	}

	if MoreThanOne(pos.Pieces[Bishop] & pos.Colours[Black]) {
		score -= bishopPair
	}

	// white rooks
	for fromBB = pos.Pieces[Rook] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= rookPhase
		fromId = BitScan(fromBB)

		attacks = RookAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		score += mobilityBonus[2][mobility]
		score += whiteRooksPos[fromId]

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Rook] |= attacks

		if pos.Pieces[Pawn]&FILES[File(fromId)] == 0 {
			score += rookOnFile[1]
		} else if (pos.Pieces[Pawn]&pos.Colours[White])&FILES[File(fromId)] == 0 {
			score += rookOnFile[0]
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

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Rook] |= attacks

		if pos.Pieces[Pawn]&FILES[File(fromId)] == 0 {
			score -= rookOnFile[1]
		} else if (pos.Pieces[Pawn]&pos.Colours[Black])&FILES[File(fromId)] == 0 {
			score -= rookOnFile[0]
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
	} else {
		score -= tempo
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
	if int(blackKingAttackersCount) > 1-PopCount(pos.Colours[Black]&pos.Pieces[Queen]) {

		// Weak squares are attacked by the enemy, defended no more
		// than once and only defended by our Queens or our King
		weak := blackAttacked & ^whiteAttackedByTwo & (^whiteAttacked | whiteAttackedBy[Queen] | whiteAttackedBy[King])

		safe := ^pos.Colours[Black] & (^whiteAttacked | (weak & blackAttackedByTwo))

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
		count += int(kingSafetyFriendlyPawns) * PopCount(pos.Colours[White]&pos.Pieces[Pawn]&whiteKingArea & ^weak)
		count += int(kingSafetyNoEnemyQueens) * BoolToInt(pos.Colours[Black]&pos.Pieces[Queen] != 0)
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
	if int(whiteKingAttackersCount) > 1-PopCount(pos.Colours[White]&pos.Pieces[Queen]) {
		// Weak squares are attacked by the enemy, defended no more
		// than once and only defended by our Queens or our King
		weak := whiteAttacked & ^blackAttackedByTwo & (^blackAttacked | blackAttackedBy[Queen] | blackAttackedBy[King])

		safe := ^pos.Colours[White] & (^blackAttacked | (weak & whiteAttackedByTwo))

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
		count += int(kingSafetyFriendlyPawns) * PopCount(pos.Colours[Black]&pos.Pieces[Pawn]&blackKingArea & ^weak)
		count += int(kingSafetyNoEnemyQueens) * BoolToInt(pos.Colours[White]&pos.Pieces[Queen] != 0)
		count += int(kingSafetySafeQueenCheck) * PopCount(queenChecks)
		count += int(kingSafetySafeRookCheck) * PopCount(rookChecks)
		count += int(kingSafetySafeBishopCheck) * PopCount(bishopChecks)
		count += int(kingSafetySafeKnightCheck) * PopCount(knightChecks)
		count += int(kingSafetyAdjustment)
		if count > 0 {
			score += S(int16(count*count/720), int16(count/20))
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
