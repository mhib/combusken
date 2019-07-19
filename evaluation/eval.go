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

func addScore(first, second Score) Score {
	return Score{
		Middle: first.Middle + second.Middle,
		End:    first.End + second.End,
	}
}

var PawnValue = Score{190, 213}
var KnightValue = Score{858, 742}
var BishopValue = Score{832, 784}
var RookValue = Score{1093, 1250}
var QueenValue = Score{2413, 2452}

// Piece Square Values
var pieceScores = [7][8][4]Score{
	{},
	{},
	{ // knight
		{Score{-128, -72}, Score{-51, -74}, Score{-78, -40}, Score{-41, -16}},
		{Score{-60, -70}, Score{-39, -51}, Score{-24, -16}, Score{0, -17}},
		{Score{-59, -35}, Score{-3, -33}, Score{11, -16}, Score{3, 5}},
		{Score{-36, -24}, Score{3, 0}, Score{21, 13}, Score{15, 22}},
		{Score{0, -32}, Score{5, -8}, Score{40, 4}, Score{52, 27}},
		{Score{-11, -51}, Score{26, -38}, Score{23, 0}, Score{54, 8}},
		{Score{-67, -64}, Score{-21, -39}, Score{6, -37}, Score{0, 0}},
		{Score{-200, -98}, Score{-79, -89}, Score{-53, -53}, Score{-32, -16}},
	},
	{ // Bishop
		{Score{-44, -30}, Score{0, -19}, Score{-10, -31}, Score{-11, -3}},
		{Score{0, -37}, Score{48, -37}, Score{14, -14}, Score{3, -6}},
		{Score{-7, -18}, Score{22, -4}, Score{30, 1}, Score{14, 13}},
		{Score{0, -26}, Score{0, -3}, Score{3, 0}, Score{36, 4}},
		{Score{-7, -24}, Score{-12, -7}, Score{6, 1}, Score{31, 2}},
		{Score{-37, -22}, Score{-1, 0}, Score{0, 0}, Score{-1, 0}},
		{Score{-21, -34}, Score{0, -12}, Score{0, -5}, Score{-5, 0}},
		{Score{-37, -50}, Score{-2, -40}, Score{-12, -39}, Score{-25, -20}},
	},
	{ // Rook
		{Score{-10, -19}, Score{-25, 0}, Score{3, 0}, Score{18, -8}},
		{Score{-82, 0}, Score{-10, -7}, Score{-6, -5}, Score{0, 0}},
		{Score{-53, 0}, Score{-7, -4}, Score{-1, -7}, Score{-1, -8}},
		{Score{-22, -4}, Score{-5, 0}, Score{-4, 0}, Score{-6, 0}},
		{Score{-21, 0}, Score{-10, 0}, Score{0, 5}, Score{1, 0}},
		{Score{-1, 1}, Score{0, 1}, Score{4, 0}, Score{7, 4}},
		{Score{0, 17}, Score{9, 18}, Score{10, 21}, Score{28, 8}},
		{Score{0, 12}, Score{8, 2}, Score{0, 13}, Score{4, 7}},
	},
	{ // Queen
		{Score{-1, -68}, Score{-5, -57}, Score{-5, -47}, Score{29, -29}},
		{Score{-3, -55}, Score{-3, -31}, Score{28, -23}, Score{17, -4}},
		{Score{0, -21}, Score{14, -16}, Score{1, 0}, Score{0, 1}},
		{Score{0, -21}, Score{0, 1}, Score{0, 11}, Score{-8, 24}},
		{Score{0, 1}, Score{-2, 2}, Score{-1, 9}, Score{-3, 21}},
		{Score{0, -7}, Score{8, -5}, Score{6, 4}, Score{8, 1}},
		{Score{0, -48}, Score{-20, -25}, Score{5, 8}, Score{5, 0}},
		{Score{0, -44}, Score{0, -8}, Score{0, 1}, Score{0, 0}},
	},
	{ // King
		{Score{261, -29}, Score{331, 13}, Score{182, 79}, Score{214, 50}},
		{Score{270, 32}, Score{245, 85}, Score{106, 137}, Score{56, 151}},
		{Score{115, 73}, Score{147, 116}, Score{79, 145}, Score{0, 171}},
		{Score{20, 71}, Score{89, 115}, Score{0, 164}, Score{-3, 167}},
		{Score{37, 94}, Score{119, 154}, Score{100, 158}, Score{26, 163}},
		{Score{120, 87}, Score{159, 164}, Score{85, 183}, Score{34, 140}},
		{Score{87, 40}, Score{120, 101}, Score{64, 144}, Score{25, 141}},
		{Score{24, 0}, Score{87, 60}, Score{49, 75}, Score{0, 72}},
	},
}

// Pawns Square scores
var pawnScores = [7][8]Score{
	{},
	{Score{-6, -11}, Score{13, -2}, Score{3, 14}, Score{21, 0}, Score{21, 0}, Score{6, 13}, Score{15, 1}, Score{-9, -12}},
	{Score{-3, -15}, Score{-10, -2}, Score{7, -1}, Score{15, -5}, Score{12, 0}, Score{16, -6}, Score{-8, -6}, Score{-6, -17}},
	{Score{-21, 0}, Score{-18, 0}, Score{5, -12}, Score{36, -16}, Score{33, -14}, Score{10, -11}, Score{-17, 0}, Score{-21, -2}},
	{Score{-4, 21}, Score{1, 14}, Score{0, -1}, Score{37, -11}, Score{27, -14}, Score{3, 0}, Score{2, 11}, Score{0, 16}},
	{Score{2, 61}, Score{0, 49}, Score{5, 27}, Score{0, 1}, Score{2, 0}, Score{16, 45}, Score{0, 62}, Score{2, 72}},
	{Score{-1, 111}, Score{0, 89}, Score{0, 43}, Score{0, 22}, Score{0, 38}, Score{0, 70}, Score{0, 82}, Score{0, 130}},
}

var pawnsConnected = [8][4]Score{
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
	{Score{4, -18}, Score{14, -1}, Score{24, 1}, Score{4, 17}},
	{Score{19, 0}, Score{30, 5}, Score{19, 10}, Score{31, 34}},
	{Score{8, 1}, Score{16, 5}, Score{15, 11}, Score{33, 16}},
	{Score{-2, 24}, Score{17, 16}, Score{32, 37}, Score{61, 22}},
	{Score{59, 29}, Score{55, 60}, Score{72, 58}, Score{88, 61}},
	{Score{106, 2}, Score{201, 9}, Score{224, 31}, Score{236, 51}},
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
}

var mobilityBonus = [...][32]Score{
	{
		Score{-38, -81}, Score{-35, -56}, Score{-13, -39}, Score{-13, -19}, Score{3, -8}, Score{11, 19}, // Knights
		Score{24, 12}, Score{40, 19}, Score{49, 0},
	},
	{
		Score{-70, -59}, Score{-41, -53}, Score{-4, -38}, Score{3, -8}, Score{25, -1}, Score{46, 4}, // Bishops
		Score{55, 18}, Score{62, 17}, Score{63, 28}, Score{59, 26}, Score{66, 14}, Score{49, 32},
		Score{54, 51}, Score{33, 32},
	},
	{
		Score{-133, -76}, Score{-39, -20}, Score{-13, 0}, Score{-13, 38}, Score{-8, 66}, Score{0, 84}, // Rooks
		Score{1, 112}, Score{19, 117}, Score{23, 131}, Score{44, 134}, Score{40, 145}, Score{54, 154},
		Score{62, 159}, Score{48, 160}, Score{53, 151},
	},
	{
		Score{-39, -36}, Score{-3, -14}, Score{-1, 0}, Score{-1, 0}, Score{-2, 0}, Score{-2, -10}, // Queens
		Score{16, -3}, Score{29, 0}, Score{41, 0}, Score{52, 4}, Score{50, 56}, Score{55, 61},
		Score{60, 88}, Score{56, 103}, Score{67, 113}, Score{61, 124}, Score{70, 133}, Score{73, 134},
		Score{79, 140}, Score{88, 143}, Score{88, 148}, Score{99, 166}, Score{102, 169}, Score{101, 164},
		Score{105, 183}, Score{102, 155}, Score{111, 203}, Score{47, 178},
	},
}

var passedFriendlyDistance = [8]Score{
	Score{0, 0}, Score{56, 33}, Score{0, 23}, Score{-40, -2},
	Score{-42, -16}, Score{-40, -12}, Score{-6, -6}, Score{0, 0},
}

var passedEnemyDistance = [8]Score{
	Score{0, 0}, Score{-61, -107}, Score{40, -18}, Score{55, 33},
	Score{40, 71}, Score{27, 90}, Score{23, 93}, Score{0, 80},
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
var passedRank = [7]Score{Score{0, 0}, Score{-2, -79}, Score{-14, -52}, Score{-11, 20}, Score{34, 96}, Score{109, 248}, Score{138, 476}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{
	Score{-35, 23}, Score{-20, 10}, Score{-42, 1}, Score{-39, -10},
	Score{-39, -10}, Score{-42, 1}, Score{-20, 10}, Score{-35, 23},
}

var isolated = Score{-33, -18}
var doubled = Score{-25, -50}
var backward = Score{6, -5}
var backwardOpen = Score{-29, -10}

var bishopPair = Score{87, 99}
var bishopRammedPawns = Score{-13, -27}

var bishopOutpostUndefendedBonus = Score{47, -8}
var bishopOutpostDefendedBonus = Score{84, 1}

var knightOutpostUndefendedBonus = Score{58, -28}
var knightOutpostDefendedBonus = Score{96, 17}

var minorBehindPawn = Score{8, 47}

var tempo = Score{45, 56}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{37, 33}, Score{112, -13}}

// this bonus only improves midScore
var pawnShieldBonus = [...]Score{Score{8, 0}, Score{0, 0}} // score for every pawn

var blackPassedMask [64]uint64
var whitePassedMask [64]uint64

var whiteOutpostMask [64]uint64
var blackOutpostMask [64]uint64

var distanceBetween [64][64]int16

var adjacentFilesMask [8]uint64

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

	allOccupation := pos.White | pos.Black
	whiteBlockedOrLowRankPawns := (South(allOccupation) | RANK_2_BB | RANK_3_BB) & (pos.Pawns & pos.White)
	blackBlockedOrLowRankPawns := (North(allOccupation) | RANK_7_BB | RANK_6_BB) & (pos.Pawns & pos.Black)
	whiteMobilityArea := ^(whiteBlockedOrLowRankPawns | BlackPawnsAttacks(pos.Pawns&pos.Black) | (pos.Kings & pos.White))
	blackMobilityArea := ^(blackBlockedOrLowRankPawns | WhitePawnsAttacks(pos.Pawns&pos.White) | (pos.Kings & pos.Black))

	// white king
	whiteKingLocation := BitScan(pos.Kings & pos.White)
	midResult += int(whiteKingPos[whiteKingLocation].Middle)
	endResult += int(whiteKingPos[whiteKingLocation].End)
	// Kingside shield bonus
	if (pos.Kings&pos.White)&whiteKingKingSide != 0 {
		midResult += PopCount(pos.White&pos.Pawns&whiteKingKingSideShield1) * int(pawnShieldBonus[0].Middle)
		midResult += PopCount(pos.White&pos.Pawns&whiteKingKingSideShield2) * int(pawnShieldBonus[1].Middle)
	}
	// Queenside shield bonus
	if (pos.Kings&pos.White)&whiteKingQueenSide != 0 {
		midResult += PopCount(pos.White&pos.Pawns&whiteKingQueenSideShield1) * int(pawnShieldBonus[0].Middle)
		midResult += PopCount(pos.White&pos.Pawns&whiteKingQueenSideShield2) * int(pawnShieldBonus[1].Middle)
	}

	// black king
	blackKingLocation := BitScan(pos.Kings & pos.Black)
	midResult -= int(blackKingPos[blackKingLocation].Middle)
	endResult -= int(blackKingPos[blackKingLocation].End)
	// Kingside shield bonus
	if (pos.Kings&pos.Black)&blackKingKingSide != 0 {
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingKingSideShield1) * int(pawnShieldBonus[0].Middle)
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingKingSideShield2) * int(pawnShieldBonus[1].Middle)
	}
	// Queenside shield bonus
	if (pos.Kings&pos.Black)&blackKingQueenSide != 0 {
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingQueenSideShield1) * int(pawnShieldBonus[0].Middle)
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingQueenSideShield2) * int(pawnShieldBonus[1].Middle)
	}

	// white pawns
	for fromBB = pos.Pawns & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
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
		midResult += int(whitePawnsPos[fromId].Middle)
		endResult += int(whitePawnsPos[fromId].End)
		phase -= pawnPhase
	}

	// white doubled pawns
	doubledCount := PopCount(pos.Pawns & pos.White & South(pos.Pawns&pos.White))
	midResult += doubledCount * int(doubled.Middle)
	endResult += doubledCount * int(doubled.End)

	// black pawns
	for fromBB = pos.Pawns & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
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
		midResult -= int(blackPawnsPos[fromId].Middle)
		endResult -= int(blackPawnsPos[fromId].End)
		phase -= pawnPhase
	}

	// black doubled pawns
	doubledCount = PopCount(pos.Pawns & pos.Black & North(pos.Pawns&pos.Black))
	midResult -= doubledCount * int(doubled.Middle)
	endResult -= doubledCount * int(doubled.End)

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
		if SquareMask[fromId]&whiteOutpustRanks != 0 && whiteOutpostMask[fromId]&(pos.Pawns&pos.Black) == 0 {
			if BlackPawnAttacks[fromId]&(pos.Pawns&pos.White) != 0 {
				midResult += int(knightOutpostDefendedBonus.Middle)
				endResult += int(knightOutpostDefendedBonus.End)
			} else {
				midResult += int(knightOutpostUndefendedBonus.Middle)
				endResult += int(knightOutpostUndefendedBonus.End)
			}

		}
		phase -= knightPhase
	}

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
		if SquareMask[fromId]&blackOutpustRanks != 0 && blackOutpostMask[fromId]&(pos.Pawns&pos.White) == 0 {
			if WhitePawnAttacks[fromId]&(pos.Pawns&pos.Black) != 0 {
				midResult -= int(knightOutpostDefendedBonus.Middle)
				endResult -= int(knightOutpostDefendedBonus.End)
			} else {
				midResult -= int(knightOutpostUndefendedBonus.Middle)
				endResult -= int(knightOutpostUndefendedBonus.End)
			}
		}
		phase -= knightPhase
	}

	// white bishops
	whiteRammedPawns := South(pos.Pawns&pos.Black) & (pos.Pawns & pos.White)
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
		phase -= bishopPhase
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
		phase -= bishopPhase
	}

	if MoreThanOne(pos.Bishops & pos.Black) {
		midResult -= int(bishopPair.Middle)
		endResult -= int(bishopPair.End)
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

	// tapering eval
	phase = (phase*256 + (totalPhase / 2)) / totalPhase
	result := ((midResult * (256 - phase)) + (endResult * phase)) / 256

	if pos.WhiteMove {
		return result
	}
	return -result
}
