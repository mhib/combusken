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
var KnightValue = Score{805, 762}
var BishopValue = Score{730, 746}
var RookValue = Score{1010, 1189}
var QueenValue = Score{2251, 2309}

// values from stockfish 10
var pieceScores = [7][8][4]Score{
	{},
	{},
	{ // knight
		{Score{-135, -43}, Score{-33, -92}, Score{-82, -59}, Score{-30, -42}},
		{Score{-51, -99}, Score{-70, -67}, Score{-34, -61}, Score{-7, -47}},
		{Score{-50, -71}, Score{-16, -60}, Score{-7, -39}, Score{-13, -3}},
		{Score{-22, -50}, Score{30, -47}, Score{15, 2}, Score{19, 7}},
		{Score{26, -74}, Score{16, -32}, Score{45, 2}, Score{61, 7}},
		{Score{-44, -100}, Score{95, -95}, Score{30, -17}, Score{95, -27}},
		{Score{-112, -87}, Score{-73, -44}, Score{72, -89}, Score{-11, -37}},
		{Score{-298, -112}, Score{-76, -132}, Score{-149, -77}, Score{0, -93}},
	},
	{ // Bishop
		{Score{-23, -27}, Score{4, -19}, Score{7, -18}, Score{1, -8}},
		{Score{2, -56}, Score{56, -58}, Score{31, -35}, Score{8, -17}},
		{Score{6, -38}, Score{38, -27}, Score{45, -14}, Score{16, 6}},
		{Score{3, -37}, Score{7, -26}, Score{6, 0}, Score{53, 1}},
		{Score{-9, -27}, Score{-6, -16}, Score{29, -7}, Score{70, 1}},
		{Score{-56, -10}, Score{18, -26}, Score{48, -18}, Score{12, -16}},
		{Score{-93, -16}, Score{23, -27}, Score{-16, -13}, Score{-2, -32}},
		{Score{-24, -43}, Score{-36, -41}, Score{-101, -31}, Score{-93, -24}},
	},
	{ // Rook
		{Score{-20, -33}, Score{-31, -9}, Score{9, -23}, Score{22, -32}},
		{Score{-86, -1}, Score{-12, -36}, Score{-30, -18}, Score{0, -27}},
		{Score{-73, -10}, Score{-26, -16}, Score{-26, -22}, Score{-15, -29}},
		{Score{-62, 3}, Score{-21, -2}, Score{-29, 2}, Score{-20, -5}},
		{Score{-40, 9}, Score{-33, 3}, Score{19, 8}, Score{12, -8}},
		{Score{-38, 13}, Score{37, 3}, Score{39, -4}, Score{22, -3}},
		{Score{35, 12}, Score{9, 24}, Score{80, 1}, Score{98, -9}},
		{Score{8, 18}, Score{14, 14}, Score{-38, 29}, Score{21, 17}},
	},
	{ // Queen
		{Score{-12, -101}, Score{-2, -87}, Score{6, -108}, Score{41, -136}},
		{Score{-7, -102}, Score{-3, -79}, Score{49, -113}, Score{28, -71}},
		{Score{0, -44}, Score{32, -54}, Score{-4, 7}, Score{-5, -6}},
		{Score{1, -30}, Score{-17, 35}, Score{-14, 38}, Score{-35, 75}},
		{Score{-6, 0}, Score{-39, 50}, Score{-23, 41}, Score{-53, 106}},
		{Score{44, -43}, Score{12, -17}, Score{15, 22}, Score{2, 76}},
		{Score{0, -45}, Score{-87, 31}, Score{-1, 15}, Score{-34, 86}},
		{Score{6, -43}, Score{0, -3}, Score{29, 9}, Score{24, 23}},
	},
	{ // King
		{Score{284, -44}, Score{337, 9}, Score{200, 66}, Score{241, 36}},
		{Score{306, 28}, Score{253, 91}, Score{151, 136}, Score{79, 161}},
		{Score{176, 74}, Score{181, 126}, Score{107, 163}, Score{50, 188}},
		{Score{26, 101}, Score{120, 137}, Score{82, 184}, Score{0, 207}},
		{Score{5, 129}, Score{127, 188}, Score{103, 198}, Score{1, 208}},
		{Score{119, 138}, Score{84, 212}, Score{114, 209}, Score{13, 182}},
		{Score{73, 131}, Score{32, 186}, Score{39, 205}, Score{26, 184}},
		{Score{52, 1}, Score{21, 95}, Score{2, 124}, Score{0, 96}},
	},
}
var pawnScores = [7][8]Score{
	{},
	{Score{-13, -11}, Score{39, -17}, Score{18, 8}, Score{38, 3}, Score{38, 6}, Score{22, 6}, Score{36, -8}, Score{-13, -18}},
	{Score{-11, -18}, Score{-14, -12}, Score{7, -3}, Score{4, -3}, Score{-2, 0}, Score{18, -11}, Score{-20, -11}, Score{-9, -23}},
	{Score{-30, -2}, Score{-22, -11}, Score{6, -18}, Score{36, -17}, Score{33, -15}, Score{12, -16}, Score{-23, -11}, Score{-30, -5}},
	{Score{-20, 22}, Score{19, -2}, Score{9, -15}, Score{52, -30}, Score{50, -39}, Score{6, -6}, Score{21, -7}, Score{-19, 19}},
	{Score{0, 71}, Score{27, 57}, Score{59, 8}, Score{54, -8}, Score{70, -25}, Score{87, 16}, Score{6, 55}, Score{9, 75}},
	{Score{0, 158}, Score{5, 136}, Score{-1, 83}, Score{0, 72}, Score{22, 105}, Score{-26, 108}, Score{0, 122}, Score{-62, 196}},
}

var pawnsConnected = [8][4]Score{
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
	{Score{8, -25}, Score{17, 3}, Score{17, -7}, Score{8, 27}},
	{Score{13, 3}, Score{57, -3}, Score{18, 8}, Score{44, 30}},
	{Score{32, 1}, Score{44, 5}, Score{30, 14}, Score{47, 10}},
	{Score{25, 20}, Score{27, 23}, Score{44, 36}, Score{52, 30}},
	{Score{0, 89}, Score{52, 76}, Score{91, 84}, Score{132, 44}},
	{Score{0, 261}, Score{130, 6}, Score{153, 0}, Score{0, 45}},
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
}

var mobilityBonus = [...][32]Score{
	{Score{-77, -192}, Score{-57, -136}, Score{-37, -82}, Score{-38, -43}, Score{-13, -37}, Score{4, -18}, // Knights
		Score{21, -24}, Score{39, -22}, Score{60, -50}},
	{Score{-60, -169}, Score{-29, -119}, Score{8, -63}, Score{14, -20}, Score{34, -2}, Score{54, 15}, // Bishops
		Score{65, 23}, Score{71, 28}, Score{76, 37}, Score{80, 34}, Score{97, 17}, Score{133, 16},
		Score{89, 38}, Score{75, 20}},
	{Score{-48, -70}, Score{-60, -57}, Score{-36, 3}, Score{-22, 45}, Score{-12, 77}, Score{-9, 100}, // Rooks
		Score{2, 118}, Score{13, 124}, Score{16, 121}, Score{42, 127}, Score{54, 129}, Score{69, 132},
		Score{79, 136}, Score{90, 132}, Score{200, 99}},
	{Score{-39, -36}, Score{-21, -15}, Score{-26, 0}, Score{-46, -4}, Score{-27, -1}, Score{-7, -48}, // Queens
		Score{-4, -17}, Score{21, -3}, Score{32, 33}, Score{44, 20}, Score{48, 57}, Score{48, 80},
		Score{55, 73}, Score{61, 110}, Score{61, 122}, Score{65, 124}, Score{69, 129}, Score{58, 135},
		Score{95, 118}, Score{83, 148}, Score{109, 144}, Score{121, 124}, Score{142, 110}, Score{143, 96},
		Score{107, 92}, Score{68, 92}, Score{4, 0}, Score{39, 63}},
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

var blackPassedMask [64]uint64
var whitePassedMask [64]uint64

var adjacentFilesMask [8]uint64

// PassedRank[Rank] contains a bonus according to the rank of a passed pawn
var passedRank = [7]Score{Score{0, 0}, Score{25, -14}, Score{32, -1}, Score{9, 55}, Score{56, 110}, Score{74, 234}, Score{124, 415}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{Score{-32, 54}, Score{-72, 53}, Score{-60, 21}, Score{-54, -5},
	Score{-15, -11}, Score{63, -17}, Score{31, 21}, Score{21, 10},
}

var isolated = Score{-22, -18}
var doubled = Score{-20, -55}
var backward = Score{7, -5}
var backwardOpen = Score{-30, -12}

var bishopPair = Score{94, 107}
var bishopRammedPawns = Score{-14, -20}

var minorBehindPawn = Score{9, 52}

var tempo = Score{39, 50}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{22, 39}, Score{86, -2}}

// this bonus only improves midScore
var pawnShieldBonus = [...]Score{Score{10, 0}, Score{-11, 0}} // score for every pawn

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
