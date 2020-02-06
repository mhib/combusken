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

var PawnValue = S(102, 121)
var KnightValue = S(463, 434)
var BishopValue = S(421, 428)
var RookValue = S(623, 693)
var QueenValue = S(1426, 1305)

// Piece Square Values
var pieceScores = [King + 1][8][4]Score{
	{},
	{ // knight
		{S(-55, -41), S(-8, -55), S(-29, -38), S(-2, -27)},
		{S(-11, -63), S(-19, -38), S(-10, -29), S(0, -18)},
		{S(-19, -39), S(1, -29), S(0, -15), S(11, -1)},
		{S(-18, -32), S(16, -28), S(13, 0), S(9, 4)},
		{S(-9, -35), S(5, -23), S(11, 3), S(26, 3)},
		{S(-46, -53), S(14, -47), S(-79, 18), S(17, -5)},
		{S(-76, -53), S(-45, -26), S(63, -56), S(-2, -18)},
		{S(-213, -70), S(-78, -75), S(-133, -30), S(7, -58)},
	},
	{ // Bishop
		{S(13, -8), S(27, -2), S(32, -9), S(31, -3)},
		{S(23, -23), S(65, -30), S(46, -13), S(28, -4)},
		{S(33, -18), S(49, -11), S(48, -2), S(33, 6)},
		{S(2, -29), S(9, -28), S(25, -15), S(27, -11)},
		{S(-33, -19), S(14, -31), S(-25, -8), S(26, -13)},
		{S(-126, 7), S(-63, -16), S(-211, 36), S(-58, -11)},
		{S(-65, 9), S(35, -8), S(0, 3), S(12, -9)},
		{S(-2, -18), S(-32, -15), S(-116, 1), S(-95, 3)},
	},
	{ // Rook
		{S(-2, -21), S(-10, -7), S(11, -13), S(14, -14)},
		{S(-42, -2), S(-3, -19), S(-3, -15), S(7, -15)},
		{S(-38, -7), S(-11, -10), S(1, -19), S(-10, -13)},
		{S(-37, 3), S(-10, -3), S(-14, 0), S(-13, -3)},
		{S(-46, 12), S(-29, 5), S(3, 9), S(-10, 2)},
		{S(-32, 7), S(9, 4), S(10, 1), S(-16, 7)},
		{S(-3, 15), S(3, 17), S(56, 1), S(57, -6)},
		{S(-1, 15), S(14, 9), S(-29, 21), S(28, 9)},
	},
	{ // Queen
		{S(2, -61), S(17, -77), S(20, -67), S(35, -68)},
		{S(-4, -40), S(16, -58), S(31, -50), S(29, -50)},
		{S(-6, -14), S(22, -37), S(-1, 5), S(5, -10)},
		{S(-2, -18), S(-23, 33), S(-6, 16), S(-17, 35)},
		{S(-24, 17), S(-31, 30), S(-32, 34), S(-32, 47)},
		{S(16, -31), S(-5, -7), S(-4, 16), S(-13, 49)},
		{S(-9, -24), S(-53, 15), S(-19, 24), S(-38, 54)},
		{S(4, -14), S(-5, 6), S(13, 12), S(24, 13)},
	},
	{ // King
		{S(151, -3), S(131, 35), S(54, 76), S(67, 60)},
		{S(156, 30), S(118, 52), S(61, 82), S(17, 95)},
		{S(87, 46), S(121, 53), S(54, 81), S(40, 93)},
		{S(25, 45), S(106, 50), S(44, 86), S(9, 98)},
		{S(39, 57), S(124, 70), S(84, 92), S(94, 88)},
		{S(101, 63), S(247, 67), S(234, 84), S(169, 66)},
		{S(49, 67), S(121, 79), S(112, 100), S(173, 73)},
		{S(25, 8), S(135, 36), S(114, 73), S(27, 57)},
	},
}

// Pawns Square scores
var pawnScores = [7][8]Score{
	{},
	{S(-28, 11), S(6, 0), S(-17, 13), S(4, 5), S(-4, 17), S(20, 6), S(33, -2), S(-6, -2)},
	{S(-20, -3), S(-20, -1), S(-3, -5), S(3, -2), S(-1, -1), S(-6, 5), S(14, -8), S(-10, -4)},
	{S(-21, 6), S(-19, 4), S(10, -7), S(20, -12), S(22, -8), S(0, 0), S(5, 0), S(-23, 4)},
	{S(-9, 17), S(22, -2), S(3, -4), S(33, -18), S(38, -15), S(13, -1), S(43, 3), S(-9, 15)},
	{S(7, 42), S(8, 34), S(44, 17), S(43, -6), S(88, -14), S(100, -4), S(56, 18), S(7, 41)},
	{S(-13, 60), S(48, 50), S(-3, 45), S(12, 30), S(80, 39), S(-9, 39), S(-15, 57), S(-89, 94)},
}

var pawnsConnected = [8][4]Score{
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
	{S(12, -20), S(8, 7), S(6, -9), S(6, 15)},
	{S(9, 1), S(28, 4), S(14, 9), S(16, 19)},
	{S(8, 8), S(23, 7), S(16, 11), S(18, 14)},
	{S(12, 17), S(6, 25), S(27, 24), S(30, 23)},
	{S(11, 57), S(41, 63), S(68, 59), S(67, 49)},
	{S(9, 60), S(147, -2), S(162, 21), S(220, 49)},
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
}

var mobilityBonus = [...][32]Score{
	{S(-46, -120), S(-38, -61), S(-28, -33), S(-26, -13), S(-9, -15), S(0, -6), // Knights
		S(10, -11), S(22, -14), S(34, -29)},
	{S(-36, -74), S(-19, -60), S(2, -31), S(8, -11), S(20, 1), S(31, 7), // Bishops
		S(37, 10), S(42, 12), S(45, 16), S(46, 17), S(67, 3), S(92, 4),
		S(49, 31), S(66, 16)},
	{S(-28, -35), S(-37, -32), S(-20, 10), S(-10, 36), S(-1, 48), S(4, 59), // Rooks
		S(7, 67), S(19, 66), S(20, 66), S(35, 68), S(40, 71), S(42, 73),
		S(52, 73), S(58, 71), S(69, 67)},
	{S(-23, -20), S(-50, -9), S(0, -145), S(-10, -126), S(-4, -25), S(2, -33), // Queens
		S(0, -33), S(10, -11), S(12, 24), S(15, 28), S(17, 48), S(17, 51),
		S(23, 45), S(23, 68), S(27, 66), S(28, 74), S(30, 73), S(28, 71),
		S(36, 67), S(45, 76), S(67, 53), S(64, 54), S(78, 41), S(66, 31),
		S(71, 17), S(32, 33), S(-7, -6), S(6, 3)},
}

var passedFriendlyDistance = [8]Score{
	S(0, 0), S(-3, 30), S(-5, 12), S(-9, -11),
	S(-14, -22), S(-17, -23), S(6, -29), S(-29, -18),
}

var passedEnemyDistance = [8]Score{
	S(0, 0), S(-51, -74), S(39, -31), S(17, 7),
	S(9, 28), S(1, 38), S(-5, 42), S(-21, 49),
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
var passedRank = [7]Score{S(0, 0), S(7, -36), S(-2, -15), S(-4, 27), S(22, 75), S(21, 168), S(110, 257)}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{S(3, 19), S(-8, 19), S(-15, 7), S(-23, -2),
	S(-27, -4), S(-3, -2), S(-31, 20), S(-13, 12),
}

var isolated = S(-9, -11)
var doubled = S(-12, -34)
var backward = S(5, -2)
var backwardOpen = S(-15, -5)

var bishopPair = S(38, 63)
var bishopRammedPawns = S(-8, -13)

var bishopOutpostUndefendedBonus = S(88, 2)
var bishopOutpostDefendedBonus = S(107, 14)

var knightOutpostUndefendedBonus = S(46, -15)
var knightOutpostDefendedBonus = S(67, 14)

var minorBehindPawn = S(5, 28)

var tempo = S(32, 29)

// Rook on semiopen, open file
var rookOnFile = [2]Score{S(7, 26), S(55, -3)}

var kingDefenders = [12]Score{
	S(-80, 11), S(-76, 14), S(-33, 4), S(-8, -3),
	S(2, -3), S(17, -7), S(27, -2), S(34, -7),
	S(39, -8), S(34, 2), S(12, 0), S(11, 0),
}

var kingShelter = [2][8][8]Score{
	{{S(-31, 6), S(-8, -9), S(-10, 11), S(44, -11),
		S(14, -20), S(16, -4), S(9, -14), S(-37, 17)},
		{S(16, 4), S(37, -14), S(-4, -2), S(-13, 11),
			S(-26, -5), S(18, -12), S(36, -44), S(-43, 10)},
		{S(14, 10), S(3, 3), S(-26, 8), S(-21, 10),
			S(-34, 1), S(-6, -2), S(10, -13), S(-25, 2)},
		{S(-24, 26), S(4, 7), S(-10, -4), S(3, -3),
			S(8, -22), S(0, -18), S(12, -43), S(-21, 1)},
		{S(-8, 13), S(-8, 8), S(-27, 3), S(-26, 11),
			S(-16, -7), S(-16, -2), S(-18, -10), S(-29, 8)},
		{S(29, -4), S(15, -12), S(-2, -11), S(6, -10),
			S(19, -26), S(11, -6), S(39, -31), S(-12, 3)},
		{S(19, 0), S(-9, -4), S(-27, -7), S(-14, 0),
			S(-12, -7), S(21, -8), S(17, -25), S(-32, 14)},
		{S(-35, 8), S(-32, 3), S(-15, 9), S(-14, 11),
			S(6, 1), S(-6, 11), S(-27, 0), S(-64, 37)}},
	{{S(-1, -4), S(-50, -13), S(-23, -4), S(-83, -21),
		S(-3, -16), S(-35, -18), S(-75, -2), S(-72, 21)},
		{S(7, 28), S(9, -20), S(-18, -5), S(-3, -1),
			S(-2, -2), S(22, -45), S(19, -21), S(-76, 24)},
		{S(17, 36), S(41, -4), S(11, -1), S(14, -10),
			S(20, -1), S(-21, -11), S(74, -22), S(-35, 11)},
		{S(6, 25), S(-34, 19), S(-19, 13), S(-16, 3),
			S(-26, 25), S(-76, 33), S(-24, -5), S(-46, 3)},
		{S(-1, 54), S(1, 9), S(-2, 2), S(-10, 1),
			S(-11, 11), S(1, -16), S(2, -21), S(-36, 10)},
		{S(71, -19), S(14, -7), S(-16, 3), S(-1, -11),
			S(-1, -7), S(-19, -15), S(17, -29), S(-42, 9)},
		{S(-1, 9), S(3, -12), S(-3, -14), S(-21, -8),
			S(-18, -10), S(4, -16), S(10, -21), S(-71, 25)},
		{S(2, -1), S(-14, -20), S(-5, -12), S(-22, -7),
			S(-23, -4), S(16, -12), S(-24, -27), S(-71, 33)}},
}

var kingStorm = [2][4][8]Score{
	{{S(9, 4), S(7, 2), S(9, 3), S(-3, 10),
		S(3, 10), S(16, 5), S(6, 7), S(4, -12)},
		{S(7, 5), S(5, 5), S(20, 0), S(0, 10),
			S(14, 4), S(22, -2), S(15, -6), S(-1, -11)},
		{S(10, 17), S(6, 9), S(-5, 12), S(-13, 16),
			S(-1, 10), S(13, 1), S(32, -19), S(4, -4)},
		{S(7, 15), S(2, 5), S(5, 3), S(-2, 5),
			S(-6, 12), S(7, 5), S(4, 6), S(-5, 1)}},
	{{S(0, 0), S(8, 15), S(-19, 8), S(22, -4),
		S(15, 12), S(-12, 20), S(15, 45), S(10, -21)},
		{S(0, 0), S(9, -32), S(-23, 2), S(75, -13),
			S(59, -23), S(-15, 2), S(-1, 22), S(-8, -16)},
		{S(0, 0), S(-79, 2), S(-21, 0), S(13, 4),
			S(6, -4), S(-2, -13), S(72, -49), S(-7, 2)},
		{S(0, 0), S(-2, -21), S(14, -19), S(-8, 1),
			S(-6, 1), S(5, -19), S(-4, 4), S(-8, 20)}},
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
		blackForwardRanksMasks[rank] = (^whiteForwardRanksMask[rank]) | RANKS[rank]
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
	if tuning {
		T.PieceScores[King][7-Rank(blackKingLocation)][FileMirror[File(blackKingLocation)]]--
		T.KingDefenders[blackKingDefenders]--
	}
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
