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

var PawnValue = S(94, 209)
var KnightValue = S(479, 755)
var BishopValue = S(474, 750)
var RookValue = S(690, 1181)
var QueenValue = S(1744, 2124)

// Piece Square Values
var pieceScores = [King + 1][8][4]Score{
	{},
	{ // knight
		{S(-168, -58), S(-37, -82), S(-46, -38), S(-36, -10)},
		{S(-19, -130), S(-15, -54), S(-22, -44), S(-3, -23)},
		{S(-29, -73), S(-12, -27), S(-15, -19), S(35, -5)},
		{S(-15, -44), S(19, -49), S(-5, 5), S(13, 2)},
		{S(30, -56), S(-2, -55), S(8, -11), S(27, 13)},
		{S(-12, -77), S(-70, -29), S(4, 12), S(52, -3)},
		{S(-46, -84), S(0, -28), S(114, -65), S(64, -36)},
		{S(-191, -120), S(-41, -127), S(-88, -24), S(13, -51)},
	},
	{ // Bishop
		{S(69, -46), S(95, -25), S(29, -30), S(-23, 3)},
		{S(12, 13), S(55, -48), S(49, -35), S(23, -15)},
		{S(5, 10), S(38, -14), S(39, -5), S(39, -7)},
		{S(29, -52), S(17, -45), S(20, -12), S(21, -6)},
		{S(-42, -18), S(1, -29), S(-96, 20), S(9, -15)},
		{S(-95, 0), S(-153, -1), S(-256, 82), S(-99, 8)},
		{S(-89, 1), S(64, -28), S(-4, 4), S(-37, 9)},
		{S(82, -7), S(80, -19), S(-142, 3), S(-52, 12)},
	},
	{ // Rook
		{S(-68, -28), S(-51, -11), S(-50, -4), S(-36, -23)},
		{S(-116, -10), S(-62, -37), S(-59, -34), S(-55, -38)},
		{S(-44, -21), S(-11, -38), S(-43, -28), S(-28, -48)},
		{S(-60, -5), S(-19, -5), S(-29, 1), S(-23, -4)},
		{S(-55, 16), S(6, 9), S(31, 22), S(-9, 4)},
		{S(43, 2), S(36, 13), S(70, 10), S(29, 12)},
		{S(-23, 39), S(19, 22), S(86, 8), S(80, -6)},
		{S(9, 23), S(47, 17), S(4, 29), S(-13, 24)},
	},
	{ // Queen
		{S(24, -133), S(-4, -103), S(-21, -85), S(27, -135)},
		{S(-2, -75), S(20, -107), S(16, -79), S(14, -82)},
		{S(6, -53), S(-4, -9), S(-6, 0), S(-18, 0)},
		{S(-26, 40), S(-16, 19), S(-18, 42), S(-39, 73)},
		{S(-27, 34), S(-25, 40), S(-29, 23), S(-61, 66)},
		{S(6, 29), S(39, -2), S(-3, 21), S(-5, 68)},
		{S(14, -18), S(-68, 28), S(-25, 48), S(-16, 96)},
		{S(47, -21), S(1, -1), S(21, 27), S(66, 21)},
	},
	{ // King
		{S(244, -37), S(172, 40), S(98, 92), S(112, 74)},
		{S(253, 20), S(218, 45), S(107, 112), S(38, 130)},
		{S(130, 61), S(156, 75), S(53, 117), S(57, 135)},
		{S(54, 80), S(104, 85), S(46, 130), S(46, 144)},
		{S(108, 82), S(214, 88), S(193, 122), S(128, 133)},
		{S(276, 87), S(451, 79), S(483, 96), S(393, 72)},
		{S(64, 61), S(200, 128), S(256, 122), S(306, 108)},
		{S(-27, -6), S(318, 6), S(160, 101), S(40, 65)},
	},
}

// Pawns Square scores
var pawnScores = [7][8]Score{
	{},
	{S(-28, 23), S(3, -3), S(-31, 33), S(-2, -9), S(-30, 32), S(-9, 24), S(-7, 9), S(-47, 12)},
	{S(-35, -1), S(-16, -14), S(-26, 5), S(-20, 6), S(-35, 17), S(0, -12), S(-25, -15), S(-50, 5)},
	{S(-7, 15), S(14, 3), S(-3, -5), S(51, -18), S(37, 1), S(20, -3), S(22, 0), S(-28, 11)},
	{S(29, 18), S(32, 19), S(38, -6), S(68, -20), S(71, -32), S(27, -4), S(67, -8), S(-12, 31)},
	{S(3, 74), S(41, 50), S(141, 0), S(106, -13), S(132, -22), S(67, 7), S(63, 10), S(-37, 53)},
	{S(-1, 105), S(12, 95), S(-27, 80), S(200, -7), S(157, 41), S(84, 100), S(-24, 50), S(-132, 113)},
}

var pawnsConnected = [8][4]Score{
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
	{S(25, -33), S(5, 5), S(-4, -10), S(11, 24)},
	{S(37, 1), S(25, 11), S(21, 10), S(42, 22)},
	{S(-10, 7), S(2, 10), S(24, 19), S(14, 15)},
	{S(35, 0), S(13, 36), S(44, 15), S(26, 39)},
	{S(-20, 53), S(186, 39), S(260, 70), S(167, 75)},
	{S(-185, 140), S(179, -67), S(-649, 192), S(182, 26)},
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
}

var mobilityBonus = [...][32]Score{
	{S(-54, -165), S(-41, -114), S(-30, -40), S(-21, -19), S(-1, -18), S(3, 4), // Knights
		S(12, -3), S(27, -19), S(26, -35)},
	{S(-41, -171), S(-13, -109), S(9, -73), S(13, -28), S(30, -7), S(41, 14), // Bishops
		S(42, 23), S(46, 23), S(45, 45), S(47, 43), S(87, 12), S(122, 18),
		S(36, 79), S(127, 46)},
	{S(-40, -18), S(-54, -38), S(-47, 30), S(-27, 51), S(-17, 64), S(-17, 88), // Rooks
		S(-9, 101), S(7, 105), S(12, 107), S(28, 114), S(19, 124), S(50, 113),
		S(52, 113), S(90, 104), S(126, 88)},
	{S(-31, -28), S(1, 4), S(-19, -204), S(-24, -164), S(-1, 13), S(0, -63), // Queens
		S(7, 5), S(18, -6), S(26, 31), S(21, 32), S(32, 55), S(28, 66),
		S(41, 52), S(32, 86), S(24, 92), S(20, 107), S(44, 100), S(30, 113),
		S(54, 99), S(53, 115), S(92, 79), S(46, 78), S(105, 59), S(73, 54),
		S(100, 6), S(77, 24), S(58, 0), S(-1, -10)},
}

var passedFriendlyDistance = [8]Score{
	S(0, 0), S(-1, 44), S(-23, 18), S(-2, -28),
	S(-17, -36), S(-35, -33), S(-9, -34), S(26, -44),
}

var passedEnemyDistance = [8]Score{
	S(0, 0), S(-98, -130), S(74, -55), S(31, 7),
	S(-15, 44), S(-38, 61), S(-24, 61), S(7, 59),
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
var passedRank = [7]Score{S(0, 0), S(5, -49), S(10, -17), S(-13, 55), S(10, 123), S(54, 247), S(169, 357)}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{S(-15, 27), S(-20, 34), S(-18, 7), S(-33, -19),
	S(-26, -13), S(-8, -6), S(31, 32), S(-18, 22),
}

var isolated = S(2, -23)
var doubled = S(-21, -59)
var backward = S(14, -2)
var backwardOpen = S(-13, -15)

var bishopPair = S(1, 88)
var bishopRammedPawns = S(-35, -16)

var bishopOutpostUndefendedBonus = S(107, -3)
var bishopOutpostDefendedBonus = S(94, 28)

var knightOutpostUndefendedBonus = S(58, -5)
var knightOutpostDefendedBonus = S(83, 51)

var minorBehindPawn = S(11, 60)

var tempo = S(30, 33)

// Rook on semiopen, open file
var rookOnFile = [2]Score{S(-25, 58), S(10, 29)}

var kingDefenders = [12]Score{
	S(-97, 12), S(-68, 14), S(-37, -3), S(-28, -8),
	S(-19, -1), S(-2, -1), S(19, 2), S(35, -16),
	S(56, -100), S(49, 1), S(16, 0), S(16, 0),
}

var kingShelter = [2][8][8]Score{
	{{S(-99, 11), S(22, -9), S(38, 12), S(85, -10),
		S(-27, -1), S(143, -51), S(28, 3), S(-139, 42)},
		{S(129, -28), S(105, -40), S(28, -11), S(-4, 37),
			S(-17, -36), S(28, -19), S(21, -48), S(-51, 11)},
		{S(-50, 35), S(11, 5), S(-37, 19), S(-28, 4),
			S(-79, 32), S(-65, 18), S(8, 2), S(-9, -6)},
		{S(-24, 67), S(-3, 27), S(-36, -12), S(-5, -1),
			S(1, -32), S(6, -44), S(27, -61), S(-20, -8)},
		{S(33, 12), S(-18, 14), S(-31, -1), S(-22, -8),
			S(-33, -28), S(1, 3), S(-56, 34), S(-36, 5)},
		{S(80, -29), S(27, -21), S(-31, -5), S(9, -11),
			S(5, -30), S(-31, 16), S(55, -27), S(-28, 6)},
		{S(35, -17), S(-29, -3), S(-38, -12), S(19, 0),
			S(-17, -30), S(-15, -4), S(-16, -32), S(-11, 3)},
		{S(-41, -14), S(-27, -20), S(-4, -4), S(9, 19),
			S(1, 13), S(-22, 25), S(-38, -12), S(-101, 55)}},
	{{S(-2, 123), S(75, -64), S(27, -31), S(35, -81),
		S(91, 4), S(-43, -24), S(-148, 19), S(-100, 23)},
		{S(68, 59), S(13, -18), S(-31, 34), S(-107, 26),
			S(8, -19), S(-16, -40), S(-39, -23), S(-49, 21)},
		{S(25, 47), S(4, 3), S(25, -24), S(-29, -28),
			S(20, -4), S(-102, 13), S(85, 5), S(-24, 13)},
		{S(7, 40), S(-54, 32), S(-72, 34), S(-61, 31),
			S(-98, 9), S(-67, 34), S(-50, 34), S(-42, 6)},
		{S(-3, 63), S(20, 10), S(-2, 6), S(2, 8),
			S(-24, 14), S(-13, -33), S(1, -3), S(-34, 11)},
		{S(-1, 26), S(23, -24), S(20, -6), S(24, -5),
			S(33, 0), S(-48, -11), S(26, -40), S(-45, 4)},
		{S(-131, 17), S(1, -10), S(2, -30), S(-13, -11),
			S(-40, -26), S(-43, -21), S(-18, -18), S(-108, 30)},
		{S(-290, 85), S(-3, -51), S(-7, -19), S(-41, -19),
			S(4, -6), S(13, -3), S(-47, -37), S(-118, 36)}},
}

var kingStorm = [2][4][8]Score{
	{{S(-19, 16), S(27, 4), S(34, 0), S(-7, 27),
		S(6, 29), S(15, 4), S(8, 24), S(11, -22)},
		{S(12, 0), S(35, -5), S(36, -1), S(-15, 30),
			S(6, 9), S(59, -21), S(37, -14), S(3, -9)},
		{S(7, 28), S(19, 10), S(25, 14), S(-8, 31),
			S(-39, 32), S(-3, 21), S(35, -21), S(6, -1)},
		{S(-1, 32), S(-4, 13), S(7, 1), S(-12, 3),
			S(-10, 26), S(11, 4), S(14, 27), S(0, 2)}},
	{{S(0, 0), S(-180, 87), S(87, -22), S(58, -32),
		S(-3, 22), S(-2, 62), S(11, 64), S(-4, -21)},
		{S(0, 0), S(162, 6), S(-84, 2), S(86, -15),
			S(-17, -4), S(86, -60), S(37, 91), S(16, -28)},
		{S(0, 0), S(-59, -24), S(17, -8), S(-6, -1),
			S(3, -8), S(-24, 49), S(257, -102), S(-12, -6)},
		{S(0, 0), S(0, -28), S(26, -29), S(-27, 26),
			S(-21, 7), S(5, -31), S(1, -42), S(0, 23)}},
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

var kingSafetyAttacksWeights = [King + 1]int16{0, 3, 7, -1, 8, 0}
var kingSafetyAttackValue int16 = 121
var kingSafetyWeakSquares int16 = 23
var kingSafetyFriendlyPawns int16 = -30
var kingSafetyNoEnemyQueens int16 = 54
var kingSafetySafeQueenCheck int16 = 89
var kingSafetySafeRookCheck int16 = 109
var kingSafetySafeBishopCheck int16 = 122
var kingSafetySafeKnightCheck int16 = 147
var kingSafetyAdjustment int16 = -1

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

// CounterGO's version
func IsLateEndGame(pos *Position) bool {
	return ((pos.Pieces[Rook]|pos.Pieces[Queen])&pos.Colours[pos.SideToMove]) == 0 &&
		!MoreThanOne((pos.Pieces[Knight]|pos.Pieces[Bishop])&pos.Colours[pos.SideToMove])
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
