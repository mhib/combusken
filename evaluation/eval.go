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

var PawnValue = Score{182, 213}
var KnightValue = Score{817, 763}
var BishopValue = Score{734, 748}
var RookValue = Score{1026, 1197}
var QueenValue = Score{2283, 2324}

var pieceScores = [7][8][4]Score{
	{},
	{},
	{ // knight
		{Score{-137, -43}, Score{-36, -92}, Score{-78, -59}, Score{-30, -43}},
		{Score{-48, -98}, Score{-67, -65}, Score{-34, -60}, Score{-7, -48}},
		{Score{-53, -69}, Score{-11, -58}, Score{-8, -39}, Score{-10, -2}},
		{Score{-39, -51}, Score{29, -47}, Score{11, 1}, Score{10, 6}},
		{Score{3, -65}, Score{3, -34}, Score{27, 2}, Score{51, 5}},
		{Score{-45, -94}, Score{62, -87}, Score{5, -9}, Score{63, -17}},
		{Score{-110, -88}, Score{-71, -43}, Score{72, -90}, Score{-10, -37}},
		{Score{-296, -113}, Score{-75, -132}, Score{-146, -78}, Score{0, -91}},
	},
	{ // Bishop
		{Score{-21, -25}, Score{8, -19}, Score{11, -21}, Score{5, -8}},
		{Score{4, -56}, Score{63, -59}, Score{36, -32}, Score{15, -14}},
		{Score{12, -36}, Score{46, -26}, Score{50, -10}, Score{22, 7}},
		{Score{-3, -38}, Score{9, -26}, Score{7, -1}, Score{47, 2}},
		{Score{-25, -26}, Score{-12, -17}, Score{13, -6}, Score{49, 1}},
		{Score{-81, -2}, Score{1, -26}, Score{1, -6}, Score{-3, -18}},
		{Score{-93, -16}, Score{24, -25}, Score{-14, -10}, Score{-2, -32}},
		{Score{-21, -41}, Score{-35, -40}, Score{-97, -31}, Score{-90, -25}},
	},
	{ // Rook
		{Score{-20, -31}, Score{-30, -8}, Score{8, -20}, Score{26, -32}},
		{Score{-89, -1}, Score{-14, -35}, Score{-27, -18}, Score{0, -25}},
		{Score{-74, -11}, Score{-25, -17}, Score{-25, -22}, Score{-17, -29}},
		{Score{-63, -1}, Score{-22, 0}, Score{-31, 3}, Score{-18, -4}},
		{Score{-44, 10}, Score{-29, 4}, Score{20, 9}, Score{21, -8}},
		{Score{-37, 16}, Score{39, 5}, Score{40, -3}, Score{26, 1}},
		{Score{31, 13}, Score{9, 24}, Score{79, 3}, Score{99, -9}},
		{Score{0, 21}, Score{13, 14}, Score{-39, 28}, Score{21, 17}},
	},
	{ // Queen
		{Score{-12, -101}, Score{0, -91}, Score{2, -100}, Score{40, -137}},
		{Score{-6, -104}, Score{-4, -75}, Score{51, -113}, Score{28, -72}},
		{Score{0, -47}, Score{33, -55}, Score{-4, 7}, Score{-6, -4}},
		{Score{1, -29}, Score{-16, 35}, Score{-13, 36}, Score{-33, 76}},
		{Score{-4, 1}, Score{-41, 50}, Score{-14, 39}, Score{-53, 107}},
		{Score{48, -43}, Score{26, -17}, Score{18, 22}, Score{2, 76}},
		{Score{0, -44}, Score{-91, 33}, Score{1, 15}, Score{-35, 86}},
		{Score{6, -43}, Score{0, -3}, Score{30, 9}, Score{24, 22}},
	},
	{ // King
		{Score{286, -44}, Score{337, 9}, Score{195, 70}, Score{236, 41}},
		{Score{306, 29}, Score{253, 87}, Score{143, 138}, Score{78, 160}},
		{Score{179, 74}, Score{173, 126}, Score{99, 163}, Score{51, 187}},
		{Score{25, 101}, Score{118, 138}, Score{67, 184}, Score{31, 201}},
		{Score{5, 131}, Score{127, 189}, Score{107, 198}, Score{2, 208}},
		{Score{119, 148}, Score{114, 220}, Score{147, 211}, Score{40, 179}},
		{Score{70, 133}, Score{35, 187}, Score{40, 207}, Score{25, 187}},
		{Score{52, 1}, Score{21, 99}, Score{2, 125}, Score{0, 97}},
	},
}
var pawnScores = [7][8]Score{
	{},
	{Score{-19, -3}, Score{29, -15}, Score{11, -7}, Score{35, -19}, Score{36, -22}, Score{12, -9}, Score{34, -11}, Score{-20, -1}},
	{Score{-12, -27}, Score{-23, -28}, Score{-3, -35}, Score{-4, -40}, Score{-8, -40}, Score{7, -43}, Score{-18, -33}, Score{-10, -31}},
	{Score{-24, -10}, Score{-25, -23}, Score{7, -46}, Score{32, -57}, Score{29, -55}, Score{10, -43}, Score{-21, -25}, Score{-21, -13}},
	{Score{-1, 36}, Score{40, 15}, Score{21, -12}, Score{66, -41}, Score{62, -50}, Score{10, 0}, Score{45, 11}, Score{-1, 37}},
	{Score{30, 199}, Score{43, 199}, Score{102, 116}, Score{118, 99}, Score{130, 88}, Score{135, 125}, Score{25, 202}, Score{22, 203}},
	{Score{136, 509}, Score{225, 471}, Score{121, 402}, Score{158, 363}, Score{280, 368}, Score{65, 426}, Score{154, 462}, Score{-3, 567}},
}

var pawnsConnected = [8][4]Score{
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
	{Score{17, -48}, Score{17, 0}, Score{22, -13}, Score{7, 16}},
	{Score{19, 5}, Score{65, 3}, Score{23, 20}, Score{50, 36}},
	{Score{22, 14}, Score{48, 11}, Score{31, 25}, Score{52, 27}},
	{Score{8, 21}, Score{7, 22}, Score{43, 39}, Score{53, 40}},
	{Score{-33, 66}, Score{34, 32}, Score{73, 54}, Score{87, 78}},
	{Score{0, 261}, Score{127, 51}, Score{147, 0}, Score{0, 123}},
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
}

var mobilityBonus = [...][32]Score{
	{Score{-78, -192}, Score{-57, -137}, Score{-38, -79}, Score{-37, -42}, Score{-12, -38}, Score{4, -15}, // Knights
		Score{21, -23}, Score{39, -22}, Score{62, -50}},
	{Score{-62, -159}, Score{-31, -106}, Score{8, -54}, Score{14, -12}, Score{35, 5}, Score{55, 16}, // Bishops
		Score{69, 21}, Score{72, 28}, Score{77, 37}, Score{85, 31}, Score{101, 17}, Score{134, 18},
		Score{87, 39}, Score{73, 22}},
	{Score{-50, -64}, Score{-61, -62}, Score{-35, 1}, Score{-23, 48}, Score{-14, 77}, Score{-9, 100}, // Rooks
		Score{2, 118}, Score{13, 124}, Score{16, 123}, Score{44, 126}, Score{54, 130}, Score{69, 132},
		Score{88, 132}, Score{101, 128}, Score{200, 99}},
	{Score{-39, -36}, Score{-21, -15}, Score{-30, 0}, Score{-49, -4}, Score{-28, -1}, Score{-9, -46}, // Queens
		Score{-5, -19}, Score{19, -3}, Score{30, 31}, Score{43, 18}, Score{47, 57}, Score{47, 71},
		Score{55, 72}, Score{60, 110}, Score{66, 108}, Score{65, 123}, Score{68, 129}, Score{59, 135},
		Score{94, 114}, Score{84, 150}, Score{119, 144}, Score{125, 124}, Score{138, 109}, Score{140, 98},
		Score{108, 91}, Score{67, 94}, Score{4, 0}, Score{38, 63}},
}

var passedFriendlyDistance = [8]Score{
	Score{0, 0}, Score{-57, -11}, Score{-124, 3}, Score{-66, -19},
	Score{-69, -23}, Score{-33, -37}, Score{13, -63}, Score{-37, -44},
}

var passedEnemyDistance = [8]Score{
	Score{0, 0}, Score{3, -3}, Score{22, -4}, Score{41, -1},
	Score{48, -1}, Score{41, 1}, Score{48, -1}, Score{19, 21},
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
var passedRank = [7]Score{Score{0, 0}, Score{63, 102}, Score{20, 120}, Score{18, 126}, Score{0, 135}, Score{39, 137}, Score{8, 157}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{Score{-75, -14}, Score{-75, -5}, Score{-42, -25}, Score{19, -43},
	Score{-25, -28}, Score{-57, -15}, Score{-53, -8}, Score{-54, -18},
}

var isolated = Score{-21, -16}
var doubled = Score{-26, -43}
var backward = Score{19, -2}
var backwardOpen = Score{-19, -6}

var bishopPair = Score{102, 107}
var bishopRammedPawns = Score{-17, -20}

var bishopOutpostUndefendedBonus = Score{60, -14}
var bishopOutpostDefendedBonus = Score{125, -2}

var knightOutpostUndefendedBonus = Score{51, -24}
var knightOutpostDefendedBonus = Score{91, 24}

var minorBehindPawn = Score{6, 51}

var tempo = Score{43, 50}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{26, 35}, Score{94, -2}}

// this bonus only improves midScore
var pawnShieldBonus = [...]Score{Score{10, 0}, Score{-11, 0}} // score for every pawn

var blackPassedMask [64]uint64
var whitePassedMask [64]uint64

var whiteOutpostMask [64]uint64
var blackOutpostMask [64]uint64

var adjacentFilesMask [8]uint64
var distanceBetween [64][64]int16

const whiteKingKingSide = F1 | G1 | H1
const whiteKingKingSideShield1 = (whiteKingKingSide << 8)  // one rank up
const whiteKingKingSideShield2 = (whiteKingKingSide << 16) // two ranks up
const whiteKingQueenSide = A1 | B1 | C1
const whiteKingQueenSideShield1 = (whiteKingQueenSide << 8)  // one rank up
const whiteKingQueenSideShield2 = (whiteKingQueenSide << 16) // two ranks up
const whiteOutpustRanks = RANK_4_BB | RANK_5_BB | RANK_6_BB
const blackOutpustRanks = RANK_5_BB | RANK_4_BB | RANK_3_BB

const blackKingKingSide = F8 | G8 | H8
const blackKingKingSideShield1 = (blackKingKingSide >> 8)  // one rank down
const blackKingKingSideShield2 = (blackKingKingSide >> 16) // two ranks down
const blackKingQueenSide = A8 | B8 | C8
const blackKingQueenSideShield1 = (blackKingQueenSide >> 8)  // one rank down
const blackKingQueenSideShield2 = (blackKingQueenSide >> 16) // two ranks down

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
	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {
			distanceBetween[x][y] = int16(Max((Abs(Rank(x) - Rank(y))), Abs(File(x)-File(y))))
		}
	}
}

func init() {
	loadScoresToPieceSquares()
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
}

// CounterGO's version
func IsLateEndGame(pos *Position) bool {
	if pos.WhiteMove {
		return ((pos.Rooks|pos.Queens)&pos.White) == 0 && !MoreThanOne((pos.Knights|pos.Bishops)&pos.White)

	} else {
		return ((pos.Rooks|pos.Queens)&pos.Black) == 0 && !MoreThanOne((pos.Knights|pos.Bishops)&pos.Black)
	}
}

func evaluateKingPawns(pos *Position, pk PawnKingTable) Score {
	if ok, value := pk.Get(pos.PawnKey); ok {
		return value
	}
	var fromId int
	var fromBB uint64
	midResult := int16(0)
	endResult := int16(0)

	whitePassed := uint64(0)
	blackPassed := uint64(0)
	for fromBB = pos.Pawns & pos.White; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		if whitePassedMask[fromId]&(pos.Pawns&pos.Black) == 0 {
			whitePassed |= SquareMask[fromId]
		}
		if adjacentFilesMask[File(fromId)]&(pos.Pawns&pos.White) == 0 {
			midResult += isolated.Middle
			endResult += isolated.End
		}
		if blackPassedMask[fromId]&(pos.Pawns&pos.White) == 0 &&
			WhitePawnAttacks[fromId+8]&(pos.Pawns&pos.Black) != 0 {
			if FILES[File(fromId)]&(pos.Pawns&pos.Black) == 0 {
				midResult += backwardOpen.Middle
				endResult += backwardOpen.End
			} else {
				midResult += backward.Middle
				endResult += backward.End
			}
		} else if whitePawnsConnectedMask[fromId]&(pos.White&pos.Pawns) != 0 {
			midResult += whitePawnsConnected[fromId].Middle
			endResult += whitePawnsConnected[fromId].End
		}
		midResult += whitePawnsPos[fromId].Middle
		endResult += whitePawnsPos[fromId].End
	}

	// white doubled pawns
	doubledCount := int16(PopCount(pos.Pawns & pos.White & South(pos.Pawns&pos.White)))
	midResult += doubledCount * doubled.Middle
	endResult += doubledCount * doubled.End

	// king
	whiteKingPosition := BitScan(pos.Kings & pos.White)
	midResult += whiteKingPos[whiteKingPosition].Middle
	endResult += whiteKingPos[whiteKingPosition].End

	// shield
	if (pos.Kings&pos.White)&whiteKingKingSide != 0 {
		midResult += int16(PopCount(pos.White&pos.Pawns&whiteKingKingSideShield1) * int(pawnShieldBonus[0].Middle))
		midResult += int16(PopCount(pos.White&pos.Pawns&whiteKingKingSideShield2) * int(pawnShieldBonus[1].Middle))
	}
	if (pos.Kings&pos.White)&whiteKingQueenSide != 0 {
		midResult += int16(PopCount(pos.White&pos.Pawns&whiteKingQueenSideShield1) * int(pawnShieldBonus[0].Middle))
		midResult += int16(PopCount(pos.White&pos.Pawns&whiteKingQueenSideShield2) * int(pawnShieldBonus[1].Middle))
	}

	// black pawns
	for fromBB = pos.Pawns & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		if blackPassedMask[fromId]&(pos.Pawns&pos.White) == 0 {
			blackPassed |= SquareMask[fromId]
		}
		if adjacentFilesMask[File(fromId)]&(pos.Pawns&pos.Black) == 0 {
			midResult -= (isolated.Middle)
			endResult -= (isolated.End)
		}
		if whitePassedMask[fromId]&(pos.Pawns&pos.Black) == 0 &&
			BlackPawnAttacks[fromId-8]&(pos.Pawns&pos.White) != 0 {
			if FILES[File(fromId)]&(pos.Pawns&pos.White) == 0 {
				midResult -= (backwardOpen.Middle)
				endResult -= (backwardOpen.End)
			} else {
				midResult -= (backward.Middle)
				endResult -= (backward.End)
			}
		} else if blackPawnsConnectedMask[fromId]&(pos.Black&pos.Pawns) != 0 {
			midResult -= (blackPawnsConnected[fromId].Middle)
			endResult -= (blackPawnsConnected[fromId].End)
		}
		midResult -= (blackPawnsPos[fromId].Middle)
		endResult -= blackPawnsPos[fromId].End
	}

	// black doubled pawns
	doubledCount = int16(PopCount(pos.Pawns & pos.Black & North(pos.Pawns&pos.Black)))
	midResult -= doubledCount * doubled.Middle
	endResult -= doubledCount * doubled.End

	blackKingPosition := BitScan(pos.Kings & pos.Black)
	midResult -= blackKingPos[blackKingPosition].Middle
	endResult -= blackKingPos[blackKingPosition].End
	// shield
	if (pos.Kings&pos.Black)&blackKingKingSide != 0 {
		midResult -= int16(PopCount(pos.Black&pos.Pawns&blackKingKingSideShield1) * int(pawnShieldBonus[0].Middle))
		midResult -= int16(PopCount(pos.Black&pos.Pawns&blackKingKingSideShield2) * int(pawnShieldBonus[1].Middle))
	}
	if (pos.Kings&pos.Black)&blackKingQueenSide != 0 {
		midResult -= int16(PopCount(pos.Black&pos.Pawns&blackKingQueenSideShield1) * int(pawnShieldBonus[0].Middle))
		midResult -= int16(PopCount(pos.Black&pos.Pawns&blackKingQueenSideShield2) * int(pawnShieldBonus[1].Middle))
	}

	for fromBB = whitePassed; fromBB != 0; fromBB &= (fromBB - 1) {
		// Rank and file
		midResult += passedRank[Rank(fromId)].Middle + passedFile[File(fromId)].Middle
		endResult += passedRank[Rank(fromId)].End + passedFile[File(fromId)].End

		// friendly king distance
		midResult += passedFriendlyDistance[distanceBetween[fromId][whiteKingPosition]].Middle
		endResult += passedFriendlyDistance[distanceBetween[fromId][whiteKingPosition]].End

		// enemy king distance
		midResult += passedEnemyDistance[distanceBetween[fromId][blackKingPosition]].Middle
		endResult += passedEnemyDistance[distanceBetween[fromId][blackKingPosition]].End
	}

	for fromBB = blackPassed; fromBB != 0; fromBB &= (fromBB - 1) {
		// Rank and file
		midResult -= passedRank[7-Rank(fromId)].Middle + passedFile[File(fromId)].Middle
		endResult -= passedRank[7-Rank(fromId)].End + passedFile[File(fromId)].End

		// friendly king distance
		midResult -= passedFriendlyDistance[distanceBetween[fromId][blackKingPosition]].Middle
		endResult -= passedFriendlyDistance[distanceBetween[fromId][blackKingPosition]].End

		// enemy king distance
		midResult -= passedEnemyDistance[distanceBetween[fromId][whiteKingPosition]].Middle
		endResult -= passedEnemyDistance[distanceBetween[fromId][whiteKingPosition]].End
	}

	result := Score{midResult, endResult}

	pk.Set(pos.PawnKey, result)

	return result
}

func Evaluate(pos *Position, pk PawnKingTable) int {
	var fromId int
	var fromBB uint64

	phase := totalPhase
	midResult := 0
	endResult := 0
	whiteMobilityArea := ^((pos.Pawns & pos.White) | (BlackPawnsAttacks(pos.Pawns & pos.Black)))
	blackMobilityArea := ^((pos.Pawns & pos.Black) | (WhitePawnsAttacks(pos.Pawns & pos.White)))
	allOccupation := pos.White | pos.Black

	pawnKingScore := evaluateKingPawns(pos, pk)
	midResult += int(pawnKingScore.Middle)
	endResult += int(pawnKingScore.End)

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

	whiteRammedPawns := South(pos.Pawns&pos.Black) & (pos.Pawns & pos.White)
	// white bishops
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
	// bishop pair bonus
	if MoreThanOne(pos.Bishops & pos.White) {
		midResult += int(bishopPair.Middle)
		endResult += int(bishopPair.End)
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
