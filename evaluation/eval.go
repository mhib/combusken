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

var PawnValue = S(100, 209)
var KnightValue = S(485, 753)
var BishopValue = S(479, 750)
var RookValue = S(705, 1172)
var QueenValue = S(2080, 1916)

// Piece Square Values
var pieceScores = [King + 1][8][4]Score{
	{},
	{ // knight
		{S(-85, -57), S(-39, -78), S(-55, -47), S(-18, -34)},
		{S(-30, -81), S(-39, -52), S(-21, -44), S(-4, -26)},
		{S(-31, -57), S(-10, -44), S(-18, -17), S(32, -4)},
		{S(-15, -42), S(19, -37), S(-2, 0), S(18, 2)},
		{S(-2, -49), S(-8, -40), S(11, -2), S(32, 6)},
		{S(-53, -75), S(32, -62), S(-15, 0), S(41, -5)},
		{S(-113, -73), S(-63, -40), S(82, -82), S(-4, -33)},
		{S(-305, -101), S(-110, -113), S(-188, -49), S(2, -82)},
	},
	{ // Bishop
		{S(2, -15), S(22, -5), S(28, -30), S(-7, -7)},
		{S(5, -36), S(57, -47), S(47, -27), S(19, -13)},
		{S(19, -26), S(40, -14), S(40, -1), S(40, -2)},
		{S(5, -33), S(6, -28), S(22, -4), S(31, 2)},
		{S(-30, -18), S(3, -20), S(-21, 0), S(27, 2)},
		{S(-81, -1), S(-34, -20), S(-118, 9), S(-26, -11)},
		{S(-108, 3), S(29, -13), S(-20, 3), S(3, -20)},
		{S(-1, -40), S(-52, -27), S(-163, -7), S(-133, -5)},
	},
	{ // Rook
		{S(-63, -24), S(-43, -11), S(-46, -8), S(-33, -23)},
		{S(-110, -2), S(-56, -37), S(-53, -34), S(-55, -39)},
		{S(-39, -10), S(-24, -13), S(-27, -30), S(-21, -43)},
		{S(-57, 2), S(-13, -4), S(-23, 2), S(-20, -4)},
		{S(-55, 12), S(-30, 9), S(16, 14), S(-8, 0)},
		{S(-23, 12), S(24, 3), S(29, 11), S(-11, 11)},
		{S(11, 18), S(9, 24), S(70, 6), S(87, -10)},
		{S(9, 18), S(12, 15), S(-30, 27), S(21, 16)},
	},
	{ // Queen
		{S(0, -84), S(19, -104), S(-14, -86), S(25, -101)},
		{S(-5, -59), S(9, -72), S(15, -81), S(11, -66)},
		{S(-8, -20), S(2, -43), S(-7, 8), S(-15, -8)},
		{S(-10, -26), S(-31, 50), S(-11, 27), S(-31, 64)},
		{S(-27, 18), S(-40, 40), S(-42, 45), S(-47, 74)},
		{S(29, -37), S(5, -11), S(-5, 19), S(-11, 66)},
		{S(-13, -34), S(-78, 25), S(-27, 31), S(-50, 76)},
		{S(12, -28), S(3, 0), S(25, 9), S(28, 21)},
	},
	{ // King
		{S(243, -37), S(176, 41), S(99, 90), S(112, 73)},
		{S(263, 14), S(215, 44), S(106, 108), S(40, 132)},
		{S(119, 66), S(157, 79), S(61, 119), S(57, 138)},
		{S(21, 74), S(105, 85), S(44, 132), S(-4, 154)},
		{S(37, 89), S(148, 112), S(92, 141), S(124, 128)},
		{S(137, 96), S(344, 103), S(321, 127), S(232, 99)},
		{S(63, 100), S(166, 116), S(156, 145), S(243, 109)},
		{S(38, 15), S(187, 54), S(163, 106), S(37, 86)},
	},
}

// Pawns Square scores
var pawnScores = [7][8]Score{
	{},
	{S(-18, 24), S(6, -5), S(-34, 31), S(-7, 8), S(-31, 16), S(-5, 25), S(-8, 12), S(-42, 11)},
	{S(-27, 2), S(-24, -15), S(-27, 9), S(-26, 8), S(-36, 15), S(2, -11), S(-30, -11), S(-47, 2)},
	{S(-2, 18), S(19, 0), S(-7, -4), S(48, -18), S(35, 3), S(19, -8), S(21, 0), S(-26, 11)},
	{S(29, 18), S(34, 19), S(34, -4), S(67, -18), S(72, -30), S(28, -4), S(64, -11), S(-11, 32)},
	{S(14, 57), S(40, 41), S(79, 9), S(69, 2), S(80, -7), S(108, 12), S(14, 50), S(-14, 41)},
	{S(-7, 95), S(11, 95), S(5, 54), S(3, 56), S(122, 41), S(-15, 67), S(2, 67), S(-100, 118)},
}

var pawnsConnected = [8][4]Score{
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
	{S(25, -34), S(8, 6), S(-4, -11), S(12, 24)},
	{S(40, -1), S(24, 15), S(22, 11), S(43, 18)},
	{S(-10, 6), S(-2, 14), S(24, 18), S(14, 16)},
	{S(22, 3), S(19, 37), S(43, 16), S(24, 38)},
	{S(18, 82), S(64, 82), S(109, 80), S(106, 69)},
	{S(11, 86), S(212, -1), S(232, 31), S(312, 60)},
	{S(0, 0), S(0, 0), S(0, 0), S(0, 0)},
}

var mobilityBonus = [...][32]Score{
	{S(-60, -169), S(-47, -97), S(-26, -47), S(-20, -17), S(-2, -15), S(6, 3), // Knights
		S(14, -7), S(31, -21), S(35, -39)},
	{S(-43, -105), S(-13, -85), S(12, -63), S(16, -30), S(31, -5), S(43, 19), // Bishops
		S(45, 22), S(50, 22), S(43, 45), S(54, 38), S(96, 2), S(132, 5),
		S(73, 34), S(96, 12)},
	{S(-40, -50), S(-60, -42), S(-44, 15), S(-27, 48), S(-18, 64), S(-20, 93), // Rooks
		S(-8, 103), S(6, 106), S(16, 105), S(29, 115), S(27, 124), S(51, 114),
		S(61, 115), S(90, 105), S(134, 85)},
	{S(-31, -28), S(-69, -11), S(-5, -204), S(-14, -171), S(0, -27), S(0, -33), // Queens
		S(14, -30), S(19, -7), S(29, 28), S(24, 32), S(38, 53), S(32, 66),
		S(48, 54), S(35, 87), S(29, 92), S(28, 109), S(44, 100), S(38, 103),
		S(54, 98), S(54, 116), S(95, 80), S(95, 80), S(100, 67), S(89, 48),
		S(99, 29), S(41, 53), S(-10, -8), S(8, 5)},
}

var passedFriendlyDistance = [8]Score{
	S(0, 0), S(8, 45), S(-24, 22), S(-5, -26),
	S(-18, -33), S(-26, -33), S(-10, -34), S(-40, -21),
}

var passedEnemyDistance = [8]Score{
	S(0, 0), S(-81, -128), S(66, -47), S(32, 9),
	S(-14, 48), S(-27, 63), S(-20, 64), S(-17, 70),
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
var passedRank = [7]Score{S(0, 0), S(3, -52), S(3, -15), S(-17, 56), S(8, 122), S(72, 244), S(160, 357)}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{S(-10, 28), S(-24, 31), S(-20, 8), S(-34, -21),
	S(-26, -15), S(-1, -7), S(-1, 34), S(-26, 21),
}

var isolated = S(2, -23)
var doubled = S(-24, -56)
var backward = S(15, -5)
var backwardOpen = S(-10, -14)

var bishopPair = S(11, 83)
var bishopRammedPawns = S(-37, -17)

var bishopOutpostUndefendedBonus = S(83, 2)
var bishopOutpostDefendedBonus = S(85, 19)

var knightOutpostUndefendedBonus = S(56, -7)
var knightOutpostDefendedBonus = S(82, 48)

var minorBehindPawn = S(11, 60)

var tempo = S(29, 35)

// Rook on semiopen, open file
var rookOnFile = [2]Score{S(-23, 57), S(8, 29)}

var kingDefenders = [12]Score{
	S(-114, 16), S(-68, 15), S(-40, -4), S(-31, -5),
	S(-17, -1), S(-2, 0), S(24, 2), S(40, 0),
	S(57, 0), S(60, 0), S(16, 0), S(16, 0),
}

var kingShelter = [2][8][8]Score{
	{{S(-42, 5), S(2, -13), S(0, 11), S(48, -11),
		S(5, -24), S(3, 2), S(8, -14), S(-43, 19)},
		{S(31, -1), S(51, -18), S(-5, -2), S(-2, 9),
			S(-34, -2), S(12, -11), S(28, -56), S(-43, 6)},
		{S(15, 12), S(5, 3), S(-21, 6), S(-27, 11),
			S(-50, 0), S(-21, 2), S(-11, -8), S(-17, -4)},
		{S(-23, 31), S(-1, 9), S(-30, -8), S(-2, -2),
			S(-1, -31), S(-1, -18), S(18, -59), S(-21, -7)},
		{S(0, 9), S(-1, 3), S(-30, 0), S(-33, 6),
			S(-20, -11), S(-33, 3), S(-34, -1), S(-46, 6)},
		{S(72, -26), S(25, -20), S(-33, -13), S(8, -11),
			S(8, -30), S(-28, 6), S(58, -34), S(-28, 8)},
		{S(34, -13), S(-27, -4), S(-30, -11), S(-1, 0),
			S(-15, -14), S(2, -5), S(-10, -24), S(-15, 5)},
		{S(-31, -1), S(-26, -20), S(-11, -4), S(0, 18),
			S(0, 12), S(-26, 24), S(-42, -10), S(-102, 54)}},
	{{S(0, -2), S(-60, -24), S(-26, -7), S(-113, -31),
		S(-7, -24), S(-47, -20), S(-111, 2), S(-97, 28)},
		{S(9, 34), S(3, -21), S(-28, -7), S(-8, -5),
			S(-4, -1), S(18, -56), S(3, -23), S(-82, 29)},
		{S(21, 45), S(61, -7), S(18, -7), S(21, -13),
			S(21, 5), S(-30, -13), S(80, -20), S(-40, 18)},
		{S(9, 35), S(-44, 28), S(-27, 14), S(-28, 6),
			S(-33, 28), S(-102, 45), S(-33, -1), S(-63, 11)},
		{S(3, 79), S(16, 6), S(2, 2), S(5, -1),
			S(-24, 11), S(-10, -14), S(-1, -17), S(-37, 9)},
		{S(100, -14), S(31, -20), S(-20, 2), S(-7, -20),
			S(-5, -8), S(-39, -17), S(25, -40), S(-47, 3)},
		{S(2, 16), S(-4, -10), S(2, -27), S(-18, -11),
			S(-39, -18), S(-47, -21), S(-18, -18), S(-110, 29)},
		{S(2, 0), S(-1, -43), S(-7, -18), S(-33, -11),
			S(-18, -4), S(5, -5), S(-43, -36), S(-101, 34)}},
}

var kingStorm = [2][4][8]Score{
	{{S(14, 0), S(32, 3), S(32, 0), S(-10, 22),
		S(5, 24), S(16, 2), S(12, 24), S(8, -21)},
		{S(12, -1), S(34, -4), S(43, 0), S(-13, 28),
			S(5, 11), S(60, -20), S(43, -14), S(5, -10)},
		{S(16, 28), S(21, 14), S(22, 15), S(-10, 29),
			S(-40, 32), S(0, 21), S(34, -21), S(5, -2)},
		{S(35, 15), S(-2, 8), S(8, 0), S(-10, 2),
			S(-8, 25), S(11, 12), S(14, 9), S(-1, 3)}},
	{{S(0, 0), S(12, 22), S(-2, 11), S(31, -7),
		S(22, 18), S(-13, 19), S(11, 64), S(16, -30)},
		{S(0, 0), S(11, -40), S(-8, -8), S(85, -15),
			S(70, -21), S(-27, 2), S(3, 28), S(8, -28)},
		{S(0, 0), S(-105, 0), S(-47, -1), S(19, 3),
			S(3, 2), S(-4, -15), S(102, -72), S(-15, -4)},
		{S(0, 0), S(0, -30), S(18, -24), S(-11, 2),
			S(-23, 3), S(11, -28), S(-4, 2), S(-4, 24)}},
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

var kingSafetyAttacksWeights = [King + 1]int16{0, 3, 3, -2, 3, 0}
var kingSafetyAttackValue int16 = 122
var kingSafetyWeakSquares int16 = 19
var kingSafetyFriendlyPawns int16 = -23
var kingSafetyNoEnemyQueens int16 = 61
var kingSafetySafeQueenCheck int16 = 76
var kingSafetySafeRookCheck int16 = 107
var kingSafetySafeBishopCheck int16 = 97
var kingSafetySafeKnightCheck int16 = 142
var kingSafetyAdjustment int16 = -45

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
