package evaluation

import . "github.com/mhib/combusken/backend"
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

var PawnValue = Score{166, 210}
var KnightValue = Score{805, 765}
var BishopValue = Score{720, 741}
var RookValue = Score{1019, 1189}
var QueenValue = Score{2251, 2323}

var pieceScores = [7][8][4]Score{
	{},
	{},
	{ // knight
		{Score{-137, -43}, Score{-32, -90}, Score{-73, -58}, Score{-30, -40}},
		{Score{-47, -100}, Score{-68, -66}, Score{-30, -59}, Score{-2, -48}},
		{Score{-49, -68}, Score{-10, -59}, Score{-5, -40}, Score{-7, -1}},
		{Score{-33, -51}, Score{29, -47}, Score{13, 0}, Score{10, 7}},
		{Score{10, -65}, Score{7, -36}, Score{27, 2}, Score{52, 4}},
		{Score{-47, -93}, Score{61, -85}, Score{0, -7}, Score{63, -19}},
		{Score{-111, -86}, Score{-73, -46}, Score{77, -89}, Score{-10, -37}},
		{Score{-297, -108}, Score{-76, -132}, Score{-146, -78}, Score{0, -93}},
	},
	{ // Bishop
		{Score{-23, -21}, Score{5, -13}, Score{12, -18}, Score{7, -8}},
		{Score{1, -46}, Score{63, -57}, Score{41, -35}, Score{14, -15}},
		{Score{10, -36}, Score{44, -27}, Score{50, -13}, Score{22, 7}},
		{Score{0, -37}, Score{5, -27}, Score{6, -1}, Score{46, 1}},
		{Score{-26, -27}, Score{-12, -23}, Score{11, -7}, Score{49, 1}},
		{Score{-77, -8}, Score{-1, -27}, Score{12, -13}, Score{-6, -21}},
		{Score{-93, -16}, Score{26, -25}, Score{-16, -9}, Score{-2, -33}},
		{Score{-24, -42}, Score{-36, -40}, Score{-96, -35}, Score{-92, -25}},
	},
	{ // Rook
		{Score{-20, -31}, Score{-33, -9}, Score{9, -23}, Score{23, -31}},
		{Score{-87, 0}, Score{-15, -33}, Score{-31, -18}, Score{-1, -27}},
		{Score{-73, -11}, Score{-26, -13}, Score{-26, -22}, Score{-14, -28}},
		{Score{-63, 3}, Score{-21, -2}, Score{-30, 3}, Score{-21, -4}},
		{Score{-39, 9}, Score{-35, 5}, Score{20, 7}, Score{12, -4}},
		{Score{-39, 13}, Score{37, 4}, Score{39, -4}, Score{22, -2}},
		{Score{35, 11}, Score{8, 25}, Score{77, 4}, Score{97, -9}},
		{Score{7, 18}, Score{13, 14}, Score{-39, 29}, Score{21, 18}},
	},
	{ // Queen
		{Score{-13, -102}, Score{1, -98}, Score{4, -107}, Score{41, -136}},
		{Score{-7, -102}, Score{0, -85}, Score{50, -113}, Score{29, -70}},
		{Score{-1, -44}, Score{32, -58}, Score{-4, 7}, Score{-5, -3}},
		{Score{2, -29}, Score{-17, 35}, Score{-14, 37}, Score{-34, 75}},
		{Score{-4, 0}, Score{-38, 50}, Score{-22, 41}, Score{-53, 107}},
		{Score{50, -43}, Score{21, -17}, Score{17, 20}, Score{1, 76}},
		{Score{-1, -45}, Score{-87, 32}, Score{0, 16}, Score{-34, 86}},
		{Score{5, -43}, Score{0, -3}, Score{32, 8}, Score{24, 22}},
	},
	{ // King
		{Score{284, -44}, Score{337, 9}, Score{200, 66}, Score{242, 36}},
		{Score{306, 28}, Score{253, 91}, Score{151, 136}, Score{78, 160}},
		{Score{176, 76}, Score{180, 126}, Score{107, 163}, Score{49, 188}},
		{Score{25, 101}, Score{121, 137}, Score{59, 190}, Score{15, 204}},
		{Score{5, 129}, Score{127, 185}, Score{106, 198}, Score{0, 209}},
		{Score{123, 136}, Score{119, 205}, Score{117, 212}, Score{4, 184}},
		{Score{71, 133}, Score{34, 181}, Score{40, 202}, Score{26, 184}},
		{Score{51, 1}, Score{21, 98}, Score{1, 121}, Score{0, 96}},
	},
}
var pawnScores = [7][8]Score{
	{},
	{Score{-13, -11}, Score{39, -17}, Score{18, 8}, Score{39, 3}, Score{38, 7}, Score{22, 6}, Score{36, -8}, Score{-13, -18}},
	{Score{-11, -18}, Score{-14, -12}, Score{7, -3}, Score{4, -3}, Score{-2, 0}, Score{16, -10}, Score{-20, -11}, Score{-9, -23}},
	{Score{-30, -2}, Score{-22, -11}, Score{6, -18}, Score{36, -17}, Score{33, -15}, Score{12, -16}, Score{-23, -10}, Score{-30, -5}},
	{Score{-18, 22}, Score{19, -2}, Score{9, -15}, Score{52, -30}, Score{50, -39}, Score{2, -6}, Score{22, -6}, Score{-19, 19}},
	{Score{0, 79}, Score{27, 60}, Score{58, 14}, Score{54, -4}, Score{70, -25}, Score{88, 22}, Score{8, 57}, Score{9, 78}},
	{Score{-1, 155}, Score{5, 137}, Score{-2, 83}, Score{0, 72}, Score{24, 108}, Score{-28, 107}, Score{0, 122}, Score{-64, 198}},
}

var pawnsConnected = [8][4]Score{
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
	{Score{7, -26}, Score{17, 3}, Score{17, -7}, Score{7, 24}},
	{Score{13, 3}, Score{57, -3}, Score{17, 7}, Score{44, 30}},
	{Score{32, 1}, Score{44, 5}, Score{30, 14}, Score{47, 10}},
	{Score{25, 20}, Score{29, 24}, Score{44, 37}, Score{56, 30}},
	{Score{0, 89}, Score{52, 76}, Score{89, 83}, Score{132, 44}},
	{Score{0, 262}, Score{130, 7}, Score{150, 0}, Score{0, 46}},
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
}

var mobilityBonus = [...][32]Score{
	{Score{-76, -192}, Score{-55, -136}, Score{-35, -82}, Score{-38, -43}, Score{-13, -37}, Score{3, -14}, // Knights
		Score{21, -24}, Score{39, -22}, Score{60, -50}},
	{Score{-60, -173}, Score{-30, -120}, Score{8, -63}, Score{14, -20}, Score{34, -2}, Score{54, 15}, // Bishops
		Score{65, 23}, Score{71, 28}, Score{76, 33}, Score{80, 34}, Score{96, 17}, Score{132, 17},
		Score{88, 40}, Score{73, 20}},
	{Score{-49, -70}, Score{-60, -65}, Score{-35, 0}, Score{-24, 46}, Score{-14, 77}, Score{-9, 100}, // Rooks
		Score{2, 118}, Score{13, 124}, Score{16, 121}, Score{42, 127}, Score{54, 130}, Score{69, 132},
		Score{79, 136}, Score{90, 132}, Score{200, 97}},
	{Score{-39, -36}, Score{-21, -15}, Score{-26, 0}, Score{-47, -4}, Score{-28, -2}, Score{-9, -48}, // Queens
		Score{-4, -18}, Score{20, -3}, Score{31, 31}, Score{44, 20}, Score{49, 58}, Score{48, 77},
		Score{55, 73}, Score{61, 110}, Score{61, 116}, Score{65, 125}, Score{68, 127}, Score{58, 134},
		Score{95, 114}, Score{87, 148}, Score{109, 143}, Score{124, 125}, Score{128, 111}, Score{139, 100},
		Score{109, 89}, Score{67, 96}, Score{4, 0}, Score{39, 63}},
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
var passedRank = [7]Score{Score{0, 0}, Score{25, -14}, Score{32, -1}, Score{9, 55}, Score{56, 110}, Score{74, 234}, Score{124, 415}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{Score{-32, 54}, Score{-72, 55}, Score{-60, 21}, Score{-54, -5},
	Score{-15, -11}, Score{63, -16}, Score{31, 21}, Score{22, 10},
}

var isolated = Score{-22, -18}
var doubled = Score{-20, -55}
var backward = Score{7, -5}
var backwardOpen = Score{-30, -12}

var bishopPair = Score{102, 107}

var bishopOutpostUndefendedBonus = Score{60, -13}
var bishopOutpostDefendedBonus = Score{124, -5}

var knightOutpostUndefendedBonus = Score{50, -27}
var knightOutpostDefendedBonus = Score{91, 20}

var minorBehindPawn = Score{7, 50}

var tempo = Score{43, 50}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{23, 39}, Score{86, 1}}

// this bonus only improves midScore
var pawnShieldBonus = [...]Score{Score{11, 0}, Score{-11, 0}} // score for every pawn

var blackPassedMask [64]uint64
var whitePassedMask [64]uint64

var whiteOutpostMask [64]uint64
var blackOutpostMask [64]uint64

var adjacentFilesMask [8]uint64

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
		if adjacentFilesMask[File(fromId)]&(pos.Pawns&pos.White) == 0 {
			midResult += int(isolated.Middle)
			endResult += int(isolated.End)
		}
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

	// king
	fromId = BitScan(pos.Kings & pos.White)
	midResult += int(whiteKingPos[fromId].Middle)
	endResult += int(whiteKingPos[fromId].End)

	// shield
	if (pos.Kings&pos.White)&whiteKingKingSide != 0 {
		midResult += PopCount(pos.White&pos.Pawns&whiteKingKingSideShield1) * int(pawnShieldBonus[0].Middle)
		midResult += PopCount(pos.White&pos.Pawns&whiteKingKingSideShield2) * int(pawnShieldBonus[1].Middle)
	}
	if (pos.Kings&pos.White)&whiteKingQueenSide != 0 {
		midResult += PopCount(pos.White&pos.Pawns&whiteKingQueenSideShield1) * int(pawnShieldBonus[0].Middle)
		midResult += PopCount(pos.White&pos.Pawns&whiteKingQueenSideShield2) * int(pawnShieldBonus[1].Middle)
	}

	// black pawns
	for fromBB = pos.Pawns & pos.Black; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		if blackPassedMask[fromId]&(pos.Pawns&pos.White) == 0 {
			midResult -= int(passedRank[7-Rank(fromId)].Middle + passedFile[File(fromId)].Middle)
			endResult -= int(passedRank[7-Rank(fromId)].End + passedFile[File(fromId)].End)
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

	fromId = BitScan(pos.Kings & pos.Black)
	midResult -= int(blackKingPos[fromId].Middle)
	endResult -= int(blackKingPos[fromId].End)
	// shield
	if (pos.Kings&pos.Black)&blackKingKingSide != 0 {
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingKingSideShield1) * int(pawnShieldBonus[0].Middle)
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingKingSideShield2) * int(pawnShieldBonus[1].Middle)
	}
	if (pos.Kings&pos.Black)&blackKingQueenSide != 0 {
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingQueenSideShield1) * int(pawnShieldBonus[0].Middle)
		midResult -= PopCount(pos.Black&pos.Pawns&blackKingQueenSideShield2) * int(pawnShieldBonus[1].Middle)
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
