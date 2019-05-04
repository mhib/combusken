package engine

import . "github.com/mhib/combusken/backend"

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

var PawnValue = Score{128, 213}
var KnightValue = Score{782, 865}
var BishopValue = Score{830, 918}
var RookValue = Score{1289, 1378}
var QueenValue = Score{2529, 2687}

// values from stockfish 10
var pieceScores = [...][8][4]Score{
	{},
	{},
	{ // knight
		{Score{-169, -105}, Score{-96, -74}, Score{-80, -46}, Score{-79, -18}},
		{Score{-79, -70}, Score{-39, -56}, Score{-24, -15}, Score{-9, 6}},
		{Score{-64, -38}, Score{-20, -33}, Score{4, -5}, Score{19, 27}},
		{Score{-28, -36}, Score{5, 0}, Score{41, 13}, Score{47, 34}},
		{Score{-29, -41}, Score{13, -20}, Score{42, 4}, Score{52, 35}},
		{Score{-11, -51}, Score{28, -38}, Score{63, -17}, Score{55, 19}},
		{Score{-67, -64}, Score{-21, -45}, Score{6, -37}, Score{37, 16}},
		{Score{-200, -98}, Score{-80, -89}, Score{-53, -53}, Score{-32, -16}},
	},
	{ // Bishop
		{Score{-44, -63}, Score{-4, -30}, Score{-11, -35}, Score{-28, -8}},
		{Score{-18, -38}, Score{7, -13}, Score{14, -14}, Score{3, 0}},
		{Score{-8, -18}, Score{24, 0}, Score{-3, -7}, Score{15, 13}},
		{Score{1, -26}, Score{8, -3}, Score{26, 1}, Score{37, 16}},
		{Score{-7, -24}, Score{30, -6}, Score{23, -10}, Score{28, 17}},
		{Score{-17, -26}, Score{4, 2}, Score{-1, 1}, Score{8, 16}},
		{Score{-21, -34}, Score{-19, -18}, Score{10, -7}, Score{-6, 9}},
		{Score{-48, -51}, Score{-3, -40}, Score{-12, -39}, Score{-25, -20}},
	},
	{ // Rook
		{Score{-24, -2}, Score{-13, -6}, Score{-7, -3}, Score{2, -2}},
		{Score{-18, -10}, Score{-10, -7}, Score{-5, 1}, Score{9, 0}},
		{Score{-21, 10}, Score{-7, -4}, Score{3, 2}, Score{-1, -2}},
		{Score{-13, -5}, Score{-5, 2}, Score{-4, -8}, Score{-6, 8}},
		{Score{-24, -8}, Score{-12, 5}, Score{-1, 4}, Score{6, -9}},
		{Score{-24, 3}, Score{-4, -2}, Score{4, -10}, Score{10, 7}},
		{Score{-8, 1}, Score{6, 2}, Score{10, 17}, Score{12, -8}},
		{Score{-22, 12}, Score{-24, -6}, Score{-6, 13}, Score{4, 7}},
	},
	{ // Queen
		{Score{3, -69}, Score{-5, -57}, Score{-5, -47}, Score{4, -26}},
		{Score{-3, -55}, Score{5, -31}, Score{8, -22}, Score{12, -4}},
		{Score{-3, -39}, Score{6, -18}, Score{13, -9}, Score{7, 3}},
		{Score{4, -23}, Score{5, -3}, Score{9, 13}, Score{8, 24}},
		{Score{0, -29}, Score{14, -6}, Score{12, 9}, Score{5, 21}},
		{Score{-4, -38}, Score{10, -18}, Score{6, -12}, Score{8, 1}},
		{Score{-5, -50}, Score{6, -27}, Score{10, -24}, Score{8, -8}},
		{Score{-2, -75}, Score{-2, -52}, Score{1, -43}, Score{-2, -36}},
	},
	{ // King
		{Score{272, 0}, Score{325, 41}, Score{273, 80}, Score{190, 93}},
		{Score{277, 57}, Score{305, 98}, Score{241, 138}, Score{183, 131}},
		{Score{198, 86}, Score{253, 138}, Score{168, 165}, Score{120, 173}},
		{Score{169, 103}, Score{191, 152}, Score{136, 168}, Score{108, 169}},
		{Score{145, 98}, Score{176, 166}, Score{112, 197}, Score{69, 194}},
		{Score{122, 87}, Score{159, 164}, Score{85, 174}, Score{36, 189}},
		{Score{87, 40}, Score{120, 99}, Score{64, 128}, Score{25, 141}},
		{Score{64, 5}, Score{87, 60}, Score{49, 75}, Score{0, 75}},
	},
}
var pawnScores = [...][8]Score{{},
	{Score{0, -10}, Score{-5, -3}, Score{10, 7}, Score{13, -1}, Score{21, 7}, Score{17, 6}, Score{6, 1}, Score{-3, -20}},
	{Score{-11, -6}, Score{-10, -6}, Score{15, -1}, Score{22, -1}, Score{26, -1}, Score{28, 2}, Score{4, -2}, Score{-24, -5}},
	{Score{-9, 4}, Score{-18, -5}, Score{8, -4}, Score{22, -5}, Score{33, -6}, Score{25, -13}, Score{-4, -3}, Score{-16, -7}},
	{Score{6, 18}, Score{-3, 2}, Score{-10, 2}, Score{1, -9}, Score{12, -13}, Score{6, -8}, Score{-12, 11}, Score{1, 9}},
	{Score{-6, 25}, Score{-8, 17}, Score{5, 19}, Score{11, 29}, Score{-14, 29}, Score{0, 8}, Score{-12, 4}, Score{-14, 12}},
	{Score{-10, -1}, Score{6, -6}, Score{-5, 18}, Score{-11, 22}, Score{-2, 22}, Score{-14, 17}, Score{12, 2}, Score{-1, 9}},
}

var mobilityBonus = [...][32]Score{
	{Score{-62, -81}, Score{-53, -56}, Score{-12, -30}, Score{-4, -14}, Score{3, 8}, Score{13, 15}, // Knights
		Score{22, 23}, Score{28, 27}, Score{33, 33}},
	{Score{-48, -59}, Score{-20, -23}, Score{16, -3}, Score{26, 13}, Score{38, 24}, Score{51, 42}, // Bishops
		Score{55, 54}, Score{63, 57}, Score{63, 65}, Score{68, 73}, Score{81, 78}, Score{81, 86},
		Score{91, 88}, Score{98, 97}},
	{Score{-58, -76}, Score{-27, -18}, Score{-15, 28}, Score{-10, 55}, Score{-5, 69}, Score{-2, 82}, // Rooks
		Score{9, 112}, Score{16, 118}, Score{30, 132}, Score{29, 142}, Score{32, 155}, Score{38, 165},
		Score{46, 166}, Score{48, 169}, Score{58, 171}},
	{Score{-39, -36}, Score{-21, -15}, Score{3, 8}, Score{3, 18}, Score{14, 34}, Score{22, 54}, // Queens
		Score{28, 61}, Score{41, 73}, Score{43, 79}, Score{48, 92}, Score{56, 94}, Score{60, 104},
		Score{60, 113}, Score{66, 120}, Score{67, 123}, Score{70, 126}, Score{71, 133}, Score{73, 136},
		Score{79, 140}, Score{88, 143}, Score{88, 148}, Score{99, 166}, Score{102, 170}, Score{102, 175},
		Score{106, 184}, Score{109, 191}, Score{113, 206}, Score{116, 212}},
}

var blackPawnsPos [64]Score
var whitePawnsPos [64]Score

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

// PassedRank[Rank] contains a bonus according to the rank of a passed pawn
var passedRank = [7]Score{Score{0, 0}, Score{5, 18}, Score{12, 23}, Score{10, 31}, Score{57, 62}, Score{163, 167}, Score{271, 250}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{Score{-1, 7}, Score{0, 9}, Score{-9, -8}, Score{-30, -14},
	Score{-30, -14}, Score{-9, -8}, Score{0, 9}, Score{-1, 7},
}

const bishopPair = 55

var tempo = Score{25, 12}

var doubled = Score{11, 56}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{18, 7}, Score{44, 20}}

// this bonus only improves midScore
var pawnShieldBonus = [...]int{15, 7} // score for every pawn

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

func init() {

	for x := 0; x < 4; x++ {
		for y := 0; y < 8; y++ {
			whiteKnightsPos[y*8+x] = pieceScores[2][y][x]
			whiteKnightsPos[y*8+(7-x)] = pieceScores[2][y][x]
			blackKnightsPos[(7-y)*8+x] = pieceScores[2][y][x]
			blackKnightsPos[(7-y)*8+(7-x)] = pieceScores[2][y][x]

			whiteBishopsPos[y*8+x] = pieceScores[3][y][x]
			whiteBishopsPos[y*8+(7-x)] = pieceScores[3][y][x]
			blackBishopsPos[(7-y)*8+x] = pieceScores[3][y][x]
			blackBishopsPos[(7-y)*8+(7-x)] = pieceScores[3][y][x]

			whiteRooksPos[y*8+x] = pieceScores[4][y][x]
			whiteRooksPos[y*8+(7-x)] = pieceScores[4][y][x]
			blackRooksPos[(7-y)*8+x] = pieceScores[4][y][x]
			blackRooksPos[(7-y)*8+(7-x)] = pieceScores[4][y][x]

			whiteQueensPos[y*8+x] = pieceScores[5][y][x]
			whiteQueensPos[y*8+(7-x)] = pieceScores[5][y][x]
			blackQueensPos[(7-y)*8+x] = pieceScores[5][y][x]
			blackQueensPos[(7-y)*8+(7-x)] = pieceScores[5][y][x]

			whiteKingPos[y*8+x] = pieceScores[6][y][x]
			whiteKingPos[y*8+(7-x)] = pieceScores[6][y][x]
			blackKingPos[(7-y)*8+x] = pieceScores[6][y][x]
			blackKingPos[(7-y)*8+(7-x)] = pieceScores[6][y][x]
		}
	}

	for y := 1; y < 7; y++ {
		for x := 0; x < 8; x++ {
			whitePawnsPos[y*8+x] = pawnScores[y][x]
			blackPawnsPos[(7-y)*8+(7-x)] = pawnScores[y][x]
		}
	}

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
}

// CounterGO's version
func isLateEndGame(pos *Position) bool {
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
		midResult += int(PawnValue.Middle + whitePawnsPos[fromId].Middle)
		endResult += int(PawnValue.End + whitePawnsPos[fromId].End)
		phase -= pawnPhase
	}

	// white doubled pawns
	doubledCount := PopCount(pos.Pawns & pos.White & South(pos.Pawns&pos.White))
	midResult -= doubledCount * int(doubled.Middle)
	endResult -= doubledCount * int(doubled.End)

	// white knights

	for fromBB = pos.Knights & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(whiteMobilityArea & KnightAttacks[fromId])
		midResult += int(KnightValue.Middle + whiteKnightsPos[fromId].Middle)
		endResult += int(KnightValue.End + whiteKnightsPos[fromId].End)
		midResult += int(mobilityBonus[0][mobility].Middle)
		endResult += int(mobilityBonus[0][mobility].End)
		phase -= knightPhase
	}

	// white bishops
	for fromBB = pos.Bishops & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(whiteMobilityArea & BishopAttacks(fromId, allOccupation))
		midResult += int(mobilityBonus[1][mobility].Middle)
		endResult += int(mobilityBonus[1][mobility].End)
		midResult += int(BishopValue.Middle + whiteBishopsPos[fromId].Middle)
		endResult += int(BishopValue.End + whiteBishopsPos[fromId].End)
		phase -= bishopPhase
	}
	// bishop pair bonus
	if MoreThanOne(pos.Bishops & pos.White) {
		midResult += bishopPair
		endResult += bishopPair
	}

	// white rooks
	for fromBB = pos.Rooks & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(whiteMobilityArea & RookAttacks(fromId, allOccupation))
		midResult += int(mobilityBonus[2][mobility].Middle)
		endResult += int(mobilityBonus[2][mobility].End)
		midResult += int(RookValue.Middle) + int(whiteRooksPos[fromId].Middle)
		endResult += int(RookValue.End) + int(whiteRooksPos[fromId].End)
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
		midResult += int(QueenValue.Middle + whiteQueensPos[fromId].Middle)
		endResult += int(QueenValue.End + whiteQueensPos[fromId].End)
		phase -= queenPhase
	}

	// king
	fromId = BitScan(pos.Kings & pos.White)
	midResult += int(whiteKingPos[fromId].Middle)
	endResult += int(whiteKingPos[fromId].End)

	// shield
	if (pos.Kings&pos.White)&whiteKingKingSide != 0 {
		midResult += PopCount(pos.White&pos.Pawns&whiteKingKingSideShield1) * pawnShieldBonus[0]
		midResult += PopCount(pos.White&pos.Pawns&whiteKingKingSideShield2) * pawnShieldBonus[1]
	}
	if (pos.Kings&pos.White)&whiteKingQueenSide != 0 {
		midResult += PopCount(pos.White&pos.Pawns&whiteKingQueenSideShield1) * pawnShieldBonus[0]
		midResult += PopCount(pos.White&pos.Pawns&whiteKingQueenSideShield2) * pawnShieldBonus[1]
	}

	// black pawns
	for fromBB = pos.Pawns & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		if blackPassedMask[fromId]&(pos.Pawns&pos.White) == 0 {
			midResult -= int(passedRank[7-Rank(fromId)].Middle + passedFile[File(fromId)].Middle)
			endResult -= int(passedRank[7-Rank(fromId)].End + passedFile[File(fromId)].End)
		}
		midResult -= int(PawnValue.Middle + blackPawnsPos[fromId].Middle)
		endResult -= int(PawnValue.End + blackPawnsPos[fromId].End)
		phase -= pawnPhase
	}

	// black doubled pawns
	doubledCount = PopCount(pos.Pawns & pos.Black & North(pos.Pawns&pos.Black))
	midResult += doubledCount * int(doubled.Middle)
	endResult += doubledCount * int(doubled.End)

	// black knights
	for fromBB = pos.Knights & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(blackMobilityArea & KnightAttacks[fromId])
		midResult -= int(KnightValue.Middle + blackKnightsPos[fromId].Middle)
		endResult -= int(KnightValue.End + blackKnightsPos[fromId].End)
		midResult -= int(mobilityBonus[0][mobility].Middle)
		endResult -= int(mobilityBonus[0][mobility].End)
		phase -= knightPhase
	}

	// black bishops
	for fromBB = pos.Bishops & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(blackMobilityArea & BishopAttacks(fromId, allOccupation))
		midResult -= int(mobilityBonus[1][mobility].Middle)
		endResult -= int(mobilityBonus[1][mobility].End)
		midResult -= int(BishopValue.Middle + blackBishopsPos[fromId].Middle)
		endResult -= int(BishopValue.End + blackBishopsPos[fromId].End)
		phase -= bishopPhase
	}

	if MoreThanOne(pos.Bishops & pos.Black) {
		midResult -= bishopPair
		endResult -= bishopPair
	}

	// black rooks
	for fromBB = pos.Rooks & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		mobility := PopCount(blackMobilityArea & RookAttacks(fromId, allOccupation))
		midResult -= int(mobilityBonus[2][mobility].Middle)
		endResult -= int(mobilityBonus[2][mobility].End)
		midResult -= int(RookValue.Middle + blackRooksPos[fromId].Middle)
		endResult -= int(RookValue.End + blackRooksPos[fromId].End)
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
		midResult -= int(QueenValue.Middle)
		endResult -= int(QueenValue.End)
		phase -= queenPhase
	}

	fromId = BitScan(pos.Kings & pos.Black)
	midResult -= int(blackKingPos[fromId].Middle)
	endResult -= int(blackKingPos[fromId].End)
	// shield
	if (pos.Kings&pos.Black)&blackKingKingSide != 0 {
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingKingSideShield1) * pawnShieldBonus[0]
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingKingSideShield2) * pawnShieldBonus[1]
	}
	if (pos.Kings&pos.Black)&blackKingQueenSide != 0 {
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingQueenSideShield1) * pawnShieldBonus[0]
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingQueenSideShield2) * pawnShieldBonus[1]
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
