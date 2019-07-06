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

var PawnValue = Score{186, 216}
var KnightValue = Score{805, 767}
var BishopValue = Score{714, 747}
var RookValue = Score{1009, 1201}
var QueenValue = Score{2251, 2325}

var pieceScores = [7][8][4]Score{
	{},
	{},
	{ // knight
		{Score{-138, -43}, Score{-40, -92}, Score{-83, -60}, Score{-32, -44}},
		{Score{-51, -100}, Score{-71, -67}, Score{-38, -60}, Score{-12, -49}},
		{Score{-56, -70}, Score{-18, -59}, Score{-10, -41}, Score{-13, -3}},
		{Score{-26, -50}, Score{31, -47}, Score{13, 1}, Score{15, 6}},
		{Score{26, -74}, Score{11, -32}, Score{43, 3}, Score{59, 7}},
		{Score{-44, -100}, Score{95, -95}, Score{29, -17}, Score{95, -27}},
		{Score{-112, -89}, Score{-75, -46}, Score{70, -90}, Score{-11, -37}},
		{Score{-296, -108}, Score{-76, -132}, Score{-147, -78}, Score{0, -92}},
	},
	{ // Bishop
		{Score{-23, -27}, Score{4, -21}, Score{5, -20}, Score{-3, -8}},
		{Score{1, -56}, Score{55, -59}, Score{31, -36}, Score{8, -17}},
		{Score{5, -38}, Score{38, -28}, Score{45, -14}, Score{16, 6}},
		{Score{1, -37}, Score{7, -26}, Score{6, 0}, Score{55, 1}},
		{Score{-9, -27}, Score{-5, -16}, Score{37, -8}, Score{73, 1}},
		{Score{-49, -10}, Score{17, -26}, Score{51, -18}, Score{13, -16}},
		{Score{-93, -16}, Score{24, -25}, Score{-15, -13}, Score{-2, -33}},
		{Score{-25, -41}, Score{-35, -40}, Score{-94, -31}, Score{-90, -25}},
	},
	{ // Rook
		{Score{-24, -32}, Score{-33, -9}, Score{9, -23}, Score{22, -32}},
		{Score{-91, -1}, Score{-15, -35}, Score{-32, -18}, Score{0, -27}},
		{Score{-73, -11}, Score{-27, -17}, Score{-26, -22}, Score{-15, -29}},
		{Score{-62, 3}, Score{-21, 0}, Score{-31, 3}, Score{-20, -5}},
		{Score{-38, 10}, Score{-35, 5}, Score{21, 7}, Score{20, -8}},
		{Score{-41, 17}, Score{37, 4}, Score{40, -4}, Score{25, -3}},
		{Score{35, 12}, Score{9, 26}, Score{77, 6}, Score{98, -9}},
		{Score{7, 18}, Score{14, 14}, Score{-39, 31}, Score{21, 18}},
	},
	{ // Queen
		{Score{-14, -102}, Score{-6, -88}, Score{2, -109}, Score{40, -137}},
		{Score{-7, -106}, Score{-2, -75}, Score{49, -113}, Score{27, -72}},
		{Score{0, -47}, Score{32, -54}, Score{-4, 7}, Score{-7, -5}},
		{Score{0, -31}, Score{-17, 35}, Score{-14, 38}, Score{-34, 75}},
		{Score{-4, 0}, Score{-38, 50}, Score{-13, 41}, Score{-45, 106}},
		{Score{50, -43}, Score{29, -17}, Score{26, 20}, Score{3, 76}},
		{Score{0, -45}, Score{-83, 32}, Score{0, 15}, Score{-33, 86}},
		{Score{6, -43}, Score{0, -5}, Score{32, 9}, Score{24, 24}},
	},
	{ // King
		{Score{284, -44}, Score{337, 9}, Score{200, 66}, Score{242, 36}},
		{Score{306, 28}, Score{253, 91}, Score{151, 136}, Score{79, 161}},
		{Score{176, 74}, Score{181, 126}, Score{107, 163}, Score{50, 188}},
		{Score{27, 101}, Score{121, 138}, Score{83, 186}, Score{32, 203}},
		{Score{5, 129}, Score{127, 189}, Score{106, 198}, Score{1, 210}},
		{Score{120, 140}, Score{82, 216}, Score{114, 217}, Score{3, 185}},
		{Score{73, 133}, Score{33, 186}, Score{40, 206}, Score{25, 189}},
		{Score{53, 1}, Score{21, 99}, Score{2, 125}, Score{0, 97}},
	},
}
var pawnScores = [7][8]Score{
	{},
	{Score{-17, -2}, Score{29, -15}, Score{11, -3}, Score{33, -17}, Score{35, -22}, Score{16, -9}, Score{34, -12}, Score{-18, 0}},
	{Score{-11, -29}, Score{-23, -28}, Score{-3, -35}, Score{-6, -40}, Score{-7, -39}, Score{9, -41}, Score{-17, -35}, Score{-6, -33}},
	{Score{-24, -10}, Score{-24, -22}, Score{7, -46}, Score{36, -53}, Score{33, -51}, Score{11, -43}, Score{-21, -21}, Score{-21, -13}},
	{Score{-1, 40}, Score{40, 16}, Score{29, -12}, Score{76, -42}, Score{70, -50}, Score{19, 0}, Score{49, 10}, Score{-1, 37}},
	{Score{46, 203}, Score{74, 198}, Score{136, 114}, Score{134, 101}, Score{146, 88}, Score{158, 124}, Score{44, 205}, Score{39, 211}},
	{Score{153, 510}, Score{227, 473}, Score{156, 400}, Score{191, 363}, Score{281, 368}, Score{129, 424}, Score{189, 460}, Score{30, 568}},
}

var pawnsConnected = [8][4]Score{
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
	{Score{19, -46}, Score{17, 0}, Score{24, -11}, Score{7, 16}},
	{Score{21, 7}, Score{69, 2}, Score{23, 20}, Score{51, 44}},
	{Score{22, 15}, Score{48, 11}, Score{31, 27}, Score{57, 22}},
	{Score{8, 21}, Score{8, 26}, Score{41, 39}, Score{57, 36}},
	{Score{-42, 66}, Score{35, 41}, Score{72, 69}, Score{97, 78}},
	{Score{0, 263}, Score{127, 51}, Score{147, 0}, Score{0, 128}},
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
}

var mobilityBonus = [...][32]Score{
	{Score{-77, -192}, Score{-57, -138}, Score{-37, -82}, Score{-38, -43}, Score{-13, -37}, Score{3, -14}, // Knights
		Score{21, -24}, Score{39, -22}, Score{60, -50}},
	{Score{-60, -175}, Score{-30, -122}, Score{8, -63}, Score{14, -20}, Score{34, -2}, Score{54, 15}, // Bishops
		Score{65, 23}, Score{71, 29}, Score{76, 37}, Score{80, 34}, Score{96, 17}, Score{132, 16},
		Score{89, 39}, Score{73, 20}},
	{Score{-51, -70}, Score{-58, -65}, Score{-34, 1}, Score{-22, 45}, Score{-12, 77}, Score{-9, 100}, // Rooks
		Score{2, 118}, Score{13, 124}, Score{16, 121}, Score{42, 127}, Score{54, 130}, Score{69, 132},
		Score{79, 136}, Score{90, 132}, Score{200, 97}},
	{Score{-39, -36}, Score{-21, -15}, Score{-25, 0}, Score{-46, -4}, Score{-28, -1}, Score{-7, -48}, // Queens
		Score{-4, -17}, Score{21, -2}, Score{32, 34}, Score{44, 20}, Score{49, 58}, Score{50, 79},
		Score{55, 73}, Score{61, 110}, Score{61, 124}, Score{65, 125}, Score{68, 129}, Score{58, 135},
		Score{99, 114}, Score{83, 148}, Score{101, 144}, Score{122, 124}, Score{126, 108}, Score{139, 95},
		Score{107, 89}, Score{67, 90}, Score{4, 0}, Score{39, 63}},
}

var passedFriendlyDistance = [8]Score{
	Score{0, 0}, Score{-49, -7}, Score{-124, 7}, Score{-82, -11},
	Score{-69, -21}, Score{-29, -37}, Score{21, -65}, Score{-37, -45},
}

var passedEnemyDistance = [8]Score{
	Score{0, 0}, Score{19, 1}, Score{38, -2}, Score{57, 1},
	Score{49, 2}, Score{41, 6}, Score{48, 3}, Score{18, 22},
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
var passedRank = [7]Score{Score{0, 0}, Score{63, 100}, Score{28, 122}, Score{50, 124}, Score{16, 137}, Score{39, 142}, Score{-4, 163}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{Score{-91, -4}, Score{-79, -1}, Score{-42, -24}, Score{-28, -24},
	Score{-50, -22}, Score{-57, -14}, Score{-61, -3}, Score{-62, -14},
}

var isolated = Score{-21, -16}
var doubled = Score{-26, -41}
var backward = Score{21, -8}
var backwardOpen = Score{-27, -2}

var bishopPair = Score{106, 105}
var bishopRammedPawns = Score{-10, -16}

var bishopOutpostUndefendedBonus = Score{36, -9}
var bishopOutpostDefendedBonus = Score{92, -4}

var knightOutpostUndefendedBonus = Score{34, -23}
var knightOutpostDefendedBonus = Score{75, 23}

var minorBehindPawn = Score{9, 52}

var tempo = Score{39, 50}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{22, 39}, Score{86, 2}}

// this bonus only improves midScore
var pawnShieldBonus = [...]Score{Score{11, 0}, Score{-11, 0}} // score for every pawn

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
