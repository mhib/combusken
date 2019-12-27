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

func (first *Score) Add(second Score) {
	first.Middle += second.Middle
	first.End += second.End
}

func (first *Score) Subtract(second Score) {
	first.Middle -= second.Middle
	first.End -= second.End
}

func addScore(first, second Score) Score {
	return Score{
		Middle: first.Middle + second.Middle,
		End:    first.End + second.End,
	}
}

var PawnValue = Score{103, 117}
var KnightValue = Score{510, 420}
var BishopValue = Score{472, 416}
var RookValue = Score{655, 681}
var QueenValue = Score{1435, 1322}

// Piece Square Values
var pieceScores = [King + 1][8][4]Score{
	{},
	{ // knight
		{Score{-59, -40}, Score{-12, -54}, Score{-39, -33}, Score{-10, -24}},
		{Score{-21, -56}, Score{-27, -36}, Score{-13, -31}, Score{-3, -22}},
		{Score{-21, -40}, Score{-1, -32}, Score{-1, -17}, Score{5, -2}},
		{Score{-16, -29}, Score{13, -26}, Score{11, 0}, Score{11, 2}},
		{Score{-2, -34}, Score{7, -21}, Score{14, 2}, Score{32, 3}},
		{Score{-37, -52}, Score{27, -43}, Score{-11, -1}, Score{28, -5}},
		{Score{-78, -51}, Score{-44, -28}, Score{56, -57}, Score{-3, -23}},
		{Score{-211, -70}, Score{-76, -78}, Score{-130, -34}, Score{1, -57}},
	},
	{ // Bishop
		{Score{-9, -11}, Score{6, -4}, Score{7, -9}, Score{7, -4}},
		{Score{3, -25}, Score{42, -31}, Score{26, -18}, Score{6, -6}},
		{Score{12, -19}, Score{29, -12}, Score{26, -1}, Score{12, 5}},
		{Score{0, -23}, Score{4, -19}, Score{9, -2}, Score{29, -1}},
		{Score{-21, -13}, Score{2, -14}, Score{-2, 0}, Score{25, 2}},
		{Score{-56, -1}, Score{-16, -13}, Score{-8, -1}, Score{-18, -8}},
		{Score{-75, 2}, Score{20, -9}, Score{-14, 2}, Score{2, -14}},
		{Score{-1, -28}, Score{-36, -19}, Score{-113, -5}, Score{-92, -4}},
	},
	{ // Rook
		{Score{-5, -17}, Score{-12, -4}, Score{12, -13}, Score{14, -15}},
		{Score{-42, 0}, Score{-6, -17}, Score{-9, -10}, Score{2, -11}},
		{Score{-37, -7}, Score{-14, -8}, Score{-3, -17}, Score{-11, -13}},
		{Score{-39, 3}, Score{-9, -3}, Score{-16, 1}, Score{-12, -2}},
		{Score{-38, 8}, Score{-23, 2}, Score{11, 6}, Score{-6, 0}},
		{Score{-25, 4}, Score{16, 1}, Score{20, -4}, Score{-8, 3}},
		{Score{7, 10}, Score{6, 15}, Score{48, 1}, Score{60, -7}},
		{Score{6, 12}, Score{8, 10}, Score{-21, 18}, Score{14, 11}},
	},
	{ // Queen
		{Score{0, -58}, Score{16, -72}, Score{15, -60}, Score{33, -68}},
		{Score{-3, -41}, Score{7, -50}, Score{31, -56}, Score{27, -43}},
		{Score{-6, -14}, Score{20, -30}, Score{-1, 5}, Score{2, -6}},
		{Score{-2, -18}, Score{-23, 34}, Score{-7, 18}, Score{-20, 44}},
		{Score{-20, 12}, Score{-28, 27}, Score{-29, 31}, Score{-33, 51}},
		{Score{20, -26}, Score{3, -8}, Score{-4, 13}, Score{-8, 45}},
		{Score{-9, -24}, Score{-54, 17}, Score{-19, 21}, Score{-35, 52}},
		{Score{8, -20}, Score{2, 0}, Score{17, 6}, Score{19, 14}},
	},
	{ // King
		{Score{192, -13}, Score{175, 21}, Score{99, 65}, Score{102, 53}},
		{Score{180, 23}, Score{140, 46}, Score{69, 80}, Score{39, 91}},
		{Score{82, 49}, Score{108, 56}, Score{42, 84}, Score{39, 91}},
		{Score{14, 49}, Score{72, 58}, Score{30, 90}, Score{-3, 100}},
		{Score{25, 61}, Score{102, 77}, Score{63, 97}, Score{85, 88}},
		{Score{94, 66}, Score{237, 71}, Score{221, 87}, Score{160, 68}},
		{Score{43, 69}, Score{114, 80}, Score{107, 100}, Score{167, 75}},
		{Score{26, 10}, Score{129, 37}, Score{112, 73}, Score{25, 59}},
	},
}

// Pawns Square scores
var pawnScores = [7][8]Score{
	{},
	{Score{-16, 3}, Score{17, -2}, Score{-6, 9}, Score{6, 5}, Score{3, 10}, Score{-6, 11}, Score{17, 0}, Score{-14, 3}},
	{Score{-12, -7}, Score{-6, -3}, Score{-1, 0}, Score{5, -3}, Score{3, -1}, Score{2, -2}, Score{-5, -6}, Score{-4, -5}},
	{Score{-18, 1}, Score{-4, 0}, Score{16, -6}, Score{25, -9}, Score{19, -7}, Score{14, -3}, Score{-6, 1}, Score{-13, 3}},
	{Score{2, 12}, Score{25, 1}, Score{16, -8}, Score{36, -14}, Score{33, -16}, Score{11, 0}, Score{29, 2}, Score{-5, 15}},
	{Score{10, 39}, Score{27, 28}, Score{51, 6}, Score{47, 1}, Score{55, -5}, Score{74, 10}, Score{9, 35}, Score{9, 43}},
	{Score{-5, 65}, Score{7, 65}, Score{3, 37}, Score{2, 38}, Score{84, 28}, Score{-11, 46}, Score{1, 46}, Score{-69, 81}},
}

var pawnsConnected = [8][4]Score{
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
	{Score{12, -23}, Score{5, 10}, Score{8, -9}, Score{1, 16}},
	{Score{6, 1}, Score{30, 4}, Score{14, 10}, Score{13, 18}},
	{Score{9, 8}, Score{22, 8}, Score{14, 9}, Score{23, 11}},
	{Score{13, 16}, Score{10, 25}, Score{28, 24}, Score{31, 21}},
	{Score{12, 58}, Score{43, 56}, Score{75, 55}, Score{73, 47}},
	{Score{7, 59}, Score{146, -1}, Score{160, 21}, Score{215, 41}},
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
}

var mobilityBonus = [...][32]Score{
	{Score{-42, -117}, Score{-32, -67}, Score{-19, -38}, Score{-20, -15}, Score{-5, -17}, Score{5, -7}, // Knights
		Score{15, -12}, Score{25, -14}, Score{39, -29}},
	{Score{-28, -73}, Score{-13, -59}, Score{8, -27}, Score{12, -7}, Score{24, 2}, Score{35, 8}, // Bishops
		Score{42, 11}, Score{45, 12}, Score{48, 16}, Score{55, 14}, Score{66, 0}, Score{91, 2},
		Score{50, 23}, Score{66, 8}},
	{Score{-28, -35}, Score{-31, -29}, Score{-17, 10}, Score{-7, 36}, Score{0, 50}, Score{4, 58}, // Rooks
		Score{10, 64}, Score{17, 66}, Score{18, 67}, Score{36, 67}, Score{37, 69}, Score{40, 73},
		Score{50, 73}, Score{61, 68}, Score{92, 58}},
	{Score{-22, -20}, Score{-48, -8}, Score{-4, -141}, Score{-10, -118}, Score{-2, -19}, Score{4, -23}, // Queens
		Score{3, -21}, Score{13, -5}, Score{17, 19}, Score{20, 22}, Score{21, 36}, Score{22, 46},
		Score{27, 37}, Score{26, 61}, Score{30, 63}, Score{30, 75}, Score{31, 69}, Score{26, 71},
		Score{37, 67}, Score{37, 80}, Score{65, 55}, Score{65, 55}, Score{69, 46}, Score{61, 33},
		Score{68, 20}, Score{28, 36}, Score{-7, -6}, Score{5, 3}},
}

var passedFriendlyDistance = [8]Score{
	Score{0, 0}, Score{5, 23}, Score{-5, 9}, Score{-8, -12},
	Score{-14, -20}, Score{-17, -23}, Score{4, -27}, Score{-31, -15},
}

var passedEnemyDistance = [8]Score{
	Score{0, 0}, Score{-56, -72}, Score{23, -28}, Score{14, 6},
	Score{15, 27}, Score{8, 36}, Score{4, 40}, Score{-13, 48},
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
var passedRank = [7]Score{Score{0, 0}, Score{-2, -30}, Score{-7, -11}, Score{-7, 31}, Score{27, 74}, Score{39, 160}, Score{110, 249}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{Score{-6, 21}, Score{-24, 21}, Score{-29, 9}, Score{-26, -6},
	Score{-17, -5}, Score{7, -2}, Score{-10, 17}, Score{-7, 11},
}

var isolated = Score{-10, -11}
var doubled = Score{-11, -34}
var backward = Score{5, -2}
var backwardOpen = Score{-15, -7}

var bishopPair = Score{47, 60}
var bishopRammedPawns = Score{-7, -13}

var bishopOutpostUndefendedBonus = Score{40, -3}
var bishopOutpostDefendedBonus = Score{80, 2}

var knightOutpostUndefendedBonus = Score{32, -15}
var knightOutpostDefendedBonus = Score{54, 13}

var minorBehindPawn = Score{6, 27}

var tempo = Score{27, 28}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{12, 23}, Score{58, -3}}

// king can castle / king cannot castle
var trappedRook = [2]Score{Score{-52, -10}, {-104, -20}}

var kingDefenders = [12]Score{
	Score{-80, 0}, Score{-61, -4}, Score{-35, -4}, Score{-12, -4},
	Score{-1, -1}, Score{14, -1}, Score{26, 0}, Score{34, 1},
	Score{44, -1}, Score{35, 6}, Score{7, 4}, Score{7, 4},
}

var kingShelter = [2][8][8]Score{
	{{Score{-29, 3}, Score{-8, -9}, Score{0, 7}, Score{33, -8},
		Score{3, -17}, Score{2, 1}, Score{5, -10}, Score{-29, 13}},
		{Score{21, -1}, Score{34, -13}, Score{-4, -2}, Score{-2, 6},
			Score{-24, -2}, Score{8, -8}, Score{19, -39}, Score{-29, 7}},
		{Score{10, 7}, Score{1, 2}, Score{-15, 4}, Score{-19, 7},
			Score{-35, 0}, Score{-15, 1}, Score{-8, -6}, Score{-16, 1}},
		{Score{-16, 21}, Score{3, 6}, Score{-9, -6}, Score{0, -2},
			Score{6, -22}, Score{-6, -13}, Score{9, -41}, Score{-27, 3}},
		{Score{0, 6}, Score{-1, 2}, Score{-21, 0}, Score{-23, 9},
			Score{-14, -8}, Score{-23, 2}, Score{-24, -1}, Score{-28, 7}},
		{Score{46, -9}, Score{35, -17}, Score{-6, -10}, Score{-2, -7},
			Score{3, -21}, Score{1, -3}, Score{42, -30}, Score{-11, 2}},
		{Score{25, -3}, Score{-2, -5}, Score{-20, -8}, Score{-8, -3},
			Score{-11, -9}, Score{12, -4}, Score{6, -17}, Score{-35, 15}},
		{Score{-22, 0}, Score{-22, -5}, Score{-13, 3}, Score{-19, 11},
			Score{-7, 9}, Score{-17, 19}, Score{-36, 8}, Score{-58, 35}}},
	{{Score{0, -2}, Score{-42, -17}, Score{-18, -5}, Score{-78, -22},
		Score{-5, -17}, Score{-33, -14}, Score{-77, 1}, Score{-67, 19}},
		{Score{6, 23}, Score{2, -15}, Score{-20, -5}, Score{-6, -4},
			Score{-3, -1}, Score{12, -39}, Score{2, -16}, Score{-57, 19}},
		{Score{14, 31}, Score{42, -5}, Score{12, -5}, Score{14, -9},
			Score{14, 3}, Score{-21, -9}, Score{55, -14}, Score{-28, 8}},
		{Score{6, 24}, Score{-31, 19}, Score{-19, 9}, Score{-20, 4},
			Score{-23, 19}, Score{-71, 31}, Score{-23, -1}, Score{-44, 4}},
		{Score{2, 54}, Score{1, 4}, Score{-1, 0}, Score{-3, -1},
			Score{-9, 7}, Score{-2, -10}, Score{-1, -12}, Score{-35, 9}},
		{Score{69, -10}, Score{22, -10}, Score{-14, 1}, Score{-5, -14},
			Score{-4, -6}, Score{-27, -12}, Score{17, -28}, Score{-32, 5}},
		{Score{1, 11}, Score{7, -13}, Score{3, -16}, Score{-20, -8},
			Score{-17, -12}, Score{-3, -15}, Score{-1, -17}, Score{-73, 24}},
		{Score{1, 0}, Score{-3, -25}, Score{-4, -13}, Score{-23, -8},
			Score{-22, -3}, Score{3, -4}, Score{-30, -25}, Score{-62, 29}}},
}

var kingStorm = [2][4][8]Score{
	{{Score{19, 1}, Score{10, 1}, Score{17, 1}, Score{7, 7},
		Score{-1, 12}, Score{6, 9}, Score{-8, 19}, Score{0, -8}},
		{Score{14, 2}, Score{5, 6}, Score{20, -1}, Score{1, 9},
			Score{12, 4}, Score{10, 1}, Score{0, -1}, Score{1, -9}},
		{Score{15, 14}, Score{3, 8}, Score{2, 11}, Score{-8, 14},
			Score{-5, 12}, Score{1, 3}, Score{9, -11}, Score{7, -5}},
		{Score{24, 10}, Score{1, 5}, Score{6, 2}, Score{0, 3},
			Score{-5, 11}, Score{3, 8}, Score{3, 6}, Score{-4, 2}}},
	{{Score{0, 0}, Score{8, 15}, Score{-17, 7}, Score{21, -4},
		Score{16, 12}, Score{-9, 13}, Score{7, 44}, Score{11, -21}},
		{Score{0, 0}, Score{7, -28}, Score{-6, -6}, Score{58, -11},
			Score{48, -15}, Score{-19, 1}, Score{2, 19}, Score{6, -20}},
		{Score{0, 0}, Score{-73, 0}, Score{-33, -1}, Score{13, 2},
			Score{2, 1}, Score{-3, -11}, Score{70, -50}, Score{7, -3}},
		{Score{0, 0}, Score{0, -21}, Score{12, -17}, Score{-8, 0},
			Score{-5, 2}, Score{8, -20}, Score{-3, 1}, Score{-7, 16}}},
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
			whiteKnightsPos[y*8+x] = addScore(pieceScores[Knight][y][x], KnightValue)
			whiteKnightsPos[y*8+(7-x)] = addScore(pieceScores[Knight][y][x], KnightValue)
			blackKnightsPos[(7-y)*8+x] = addScore(pieceScores[Knight][y][x], KnightValue)
			blackKnightsPos[(7-y)*8+(7-x)] = addScore(pieceScores[Knight][y][x], KnightValue)

			whiteBishopsPos[y*8+x] = addScore(pieceScores[Bishop][y][x], BishopValue)
			whiteBishopsPos[y*8+(7-x)] = addScore(pieceScores[Bishop][y][x], BishopValue)
			blackBishopsPos[(7-y)*8+x] = addScore(pieceScores[Bishop][y][x], BishopValue)
			blackBishopsPos[(7-y)*8+(7-x)] = addScore(pieceScores[Bishop][y][x], BishopValue)

			whiteRooksPos[y*8+x] = addScore(pieceScores[Rook][y][x], RookValue)
			whiteRooksPos[y*8+(7-x)] = addScore(pieceScores[Rook][y][x], RookValue)
			blackRooksPos[(7-y)*8+x] = addScore(pieceScores[Rook][y][x], RookValue)
			blackRooksPos[(7-y)*8+(7-x)] = addScore(pieceScores[Rook][y][x], RookValue)

			whiteQueensPos[y*8+x] = addScore(pieceScores[Queen][y][x], QueenValue)
			whiteQueensPos[y*8+(7-x)] = addScore(pieceScores[Queen][y][x], QueenValue)
			blackQueensPos[(7-y)*8+x] = addScore(pieceScores[Queen][y][x], QueenValue)
			blackQueensPos[(7-y)*8+(7-x)] = addScore(pieceScores[Queen][y][x], QueenValue)

			whiteKingPos[y*8+x] = pieceScores[King][y][x]
			whiteKingPos[y*8+(7-x)] = pieceScores[King][y][x]
			blackKingPos[(7-y)*8+x] = pieceScores[King][y][x]
			blackKingPos[(7-y)*8+(7-x)] = pieceScores[King][y][x]
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

func evaluateKingPawns(pos *Position, pkTable PawnKingTable) (int, int) {
	if ok, midScore, endScore := pkTable.Get(pos.PawnKey); ok {
		return midScore, endScore
	}
	var fromBB uint64
	var fromId int
	whiteKingLocation := BitScan(pos.Pieces[King] & pos.Colours[White])
	blackKingLocation := BitScan(pos.Pieces[King] & pos.Colours[Black])
	midResult := 0
	endResult := 0

	// white pawns
	for fromBB = pos.Pieces[Pawn] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		midResult += int(whitePawnsPos[fromId].Middle)
		endResult += int(whitePawnsPos[fromId].End)

		// Passed bonus
		if whitePassedMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
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
		if adjacentFilesMask[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			midResult += int(isolated.Middle)
			endResult += int(isolated.End)
		}

		// Pawn is backward if there are no pawns behind it and cannot increase rank without being attacked by enemy pawn
		if blackPassedMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 &&
			PawnAttacks[White][fromId+8]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
			if FILES[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
				midResult += int(backwardOpen.Middle)
				endResult += int(backwardOpen.End)
			} else {
				midResult += int(backward.Middle)
				endResult += int(backward.End)
			}
		} else if whitePawnsConnectedMask[fromId]&(pos.Colours[White]&pos.Pieces[Pawn]) != 0 {
			midResult += int(whitePawnsConnected[fromId].Middle)
			endResult += int(whitePawnsConnected[fromId].End)
		}
	}

	// white doubled pawns
	doubledCount := PopCount(pos.Pieces[Pawn] & pos.Colours[White] & South(pos.Pieces[Pawn]&pos.Colours[White]))
	midResult += doubledCount * int(doubled.Middle)
	endResult += doubledCount * int(doubled.End)

	// black pawns
	for fromBB = pos.Pieces[Pawn] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		midResult -= int(blackPawnsPos[fromId].Middle)
		endResult -= int(blackPawnsPos[fromId].End)
		if blackPassedMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
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
		if adjacentFilesMask[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			midResult -= int(isolated.Middle)
			endResult -= int(isolated.End)
		}
		if whitePassedMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 &&
			PawnAttacks[Black][fromId-8]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
			if FILES[File(fromId)]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
				midResult -= int(backwardOpen.Middle)
				endResult -= int(backwardOpen.End)
			} else {
				midResult -= int(backward.Middle)
				endResult -= int(backward.End)
			}
		} else if blackPawnsConnectedMask[fromId]&(pos.Colours[Black]&pos.Pieces[Pawn]) != 0 {
			midResult -= int(blackPawnsConnected[fromId].Middle)
			endResult -= int(blackPawnsConnected[fromId].End)
		}
	}

	// black doubled pawns
	doubledCount = PopCount(pos.Pieces[Pawn] & pos.Colours[Black] & North(pos.Pieces[Pawn]&pos.Colours[Black]))
	midResult -= doubledCount * int(doubled.Middle)
	endResult -= doubledCount * int(doubled.End)

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
		midResult += int(kingShelter[sameFile][file][ourDist].Middle)
		endResult += int(kingShelter[sameFile][file][ourDist].End)

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		midResult += int(kingStorm[blocked][FileMirror[file]][theirDist].Middle)
		endResult += int(kingStorm[blocked][FileMirror[file]][theirDist].End)
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
		midResult -= int(kingShelter[sameFile][file][ourDist].Middle)
		endResult -= int(kingShelter[sameFile][file][ourDist].End)

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		midResult -= int(kingStorm[blocked][FileMirror[file]][theirDist].Middle)
		endResult -= int(kingStorm[blocked][FileMirror[file]][theirDist].End)
	}
	pkTable.Set(pos.PawnKey, midResult, endResult)
	return midResult, endResult
}

func Evaluate(pos *Position, pkTable PawnKingTable) int {
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

	midResult, endResult := evaluateKingPawns(pos, pkTable)

	// white knights
	for fromBB = pos.Pieces[Knight] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= knightPhase
		fromId = BitScan(fromBB)

		attacks = KnightAttacks[fromId]
		mobility := PopCount(whiteMobilityArea & attacks)
		midResult += int(whiteKnightsPos[fromId].Middle)
		endResult += int(whiteKnightsPos[fromId].End)
		midResult += int(mobilityBonus[0][mobility].Middle)
		endResult += int(mobilityBonus[0][mobility].End)

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Knight] |= attacks

		if (pos.Pieces[Pawn]>>8)&SquareMask[fromId] != 0 {
			midResult += int(minorBehindPawn.Middle)
			endResult += int(minorBehindPawn.End)
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && whiteOutpostMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			if PawnAttacks[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
				midResult += int(knightOutpostDefendedBonus.Middle)
				endResult += int(knightOutpostDefendedBonus.End)
			} else {
				midResult += int(knightOutpostUndefendedBonus.Middle)
				endResult += int(knightOutpostUndefendedBonus.End)
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
		midResult -= int(blackKnightsPos[fromId].Middle)
		endResult -= int(blackKnightsPos[fromId].End)
		midResult -= int(mobilityBonus[0][mobility].Middle)
		endResult -= int(mobilityBonus[0][mobility].End)

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Knight] |= attacks

		if (pos.Pieces[Pawn]<<8)&SquareMask[fromId] != 0 {
			midResult -= int(minorBehindPawn.Middle)
			endResult -= int(minorBehindPawn.End)
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && blackOutpostMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			if PawnAttacks[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
				midResult -= int(knightOutpostDefendedBonus.Middle)
				endResult -= int(knightOutpostDefendedBonus.End)
			} else {
				midResult -= int(knightOutpostUndefendedBonus.Middle)
				endResult -= int(knightOutpostUndefendedBonus.End)
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
		midResult += int(mobilityBonus[1][mobility].Middle)
		endResult += int(mobilityBonus[1][mobility].End)
		midResult += int(whiteBishopsPos[fromId].Middle)
		endResult += int(whiteBishopsPos[fromId].End)

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Bishop] |= attacks

		if (pos.Pieces[Pawn]>>8)&SquareMask[fromId] != 0 {
			midResult += int(minorBehindPawn.Middle)
			endResult += int(minorBehindPawn.End)
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && whiteOutpostMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			if PawnAttacks[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
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
		if attacks&blackKingArea != 0 {
			whiteKingAttacksCount += int16(PopCount(attacks & blackKingArea))
			whiteKingAttackersCount++
			whiteKingAttackersWeight += kingSafetyAttacksWeights[Bishop]
		}
	}

	// Bishop pair bonus
	// It is not checked if bishops have opposite colors, but that is almost always the case
	if MoreThanOne(pos.Pieces[Bishop] & pos.Colours[White]) {
		midResult += int(bishopPair.Middle)
		endResult += int(bishopPair.End)
	}

	// black bishops
	blackRammedPawns := North(pos.Pieces[Pawn]&pos.Colours[White]) & (pos.Pieces[Pawn] & pos.Colours[Black])
	for fromBB = pos.Pieces[Bishop] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= bishopPhase
		fromId = BitScan(fromBB)

		attacks = BishopAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		midResult -= int(mobilityBonus[1][mobility].Middle)
		endResult -= int(mobilityBonus[1][mobility].End)
		midResult -= int(blackBishopsPos[fromId].Middle)
		endResult -= int(blackBishopsPos[fromId].End)

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Bishop] |= attacks

		if (pos.Pieces[Pawn]<<8)&SquareMask[fromId] != 0 {
			midResult -= int(minorBehindPawn.Middle)
			endResult -= int(minorBehindPawn.End)
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && blackOutpostMask[fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			if PawnAttacks[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
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
		if attacks&whiteKingArea != 0 {
			blackKingAttacksCount += int16(PopCount(attacks & whiteKingArea))
			blackKingAttackersCount++
			blackKingAttackersWeight += kingSafetyAttacksWeights[Bishop]
		}
	}

	if MoreThanOne(pos.Pieces[Bishop] & pos.Colours[Black]) {
		midResult -= int(bishopPair.Middle)
		endResult -= int(bishopPair.End)
	}

	// white rooks
	for fromBB = pos.Pieces[Rook] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= rookPhase
		fromId = BitScan(fromBB)

		attacks = RookAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		midResult += int(mobilityBonus[2][mobility].Middle)
		endResult += int(mobilityBonus[2][mobility].End)
		midResult += int(whiteRooksPos[fromId].Middle)
		endResult += int(whiteRooksPos[fromId].End)

		whiteAttackedByTwo |= whiteAttacked & attacks
		whiteAttacked |= attacks
		whiteAttackedBy[Rook] |= attacks

		if pos.Pieces[Pawn]&FILES[File(fromId)] == 0 {
			midResult += int(rookOnFile[1].Middle)
			endResult += int(rookOnFile[1].End)
		} else if (pos.Pieces[Pawn]&pos.Colours[White])&FILES[File(fromId)] == 0 {
			midResult += int(rookOnFile[0].Middle)
			endResult += int(rookOnFile[0].End)
		} else if mobility < 3 {
			kingFile := File(whiteKingLocation)
			if (kingFile < FILE_E) == (File(fromId) < kingFile) {
				cannotCastle := BoolToInt(^pos.Flags&(WhiteKingSideCastleFlag|WhiteQueenSideCastleFlag) == 0)
				midResult += int(trappedRook[cannotCastle].Middle)
				endResult += int(trappedRook[cannotCastle].End)
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
		midResult -= int(mobilityBonus[2][mobility].Middle)
		endResult -= int(mobilityBonus[2][mobility].End)
		midResult -= int(blackRooksPos[fromId].Middle)
		endResult -= int(blackRooksPos[fromId].End)

		blackAttackedByTwo |= blackAttacked & attacks
		blackAttacked |= attacks
		blackAttackedBy[Rook] |= attacks

		if pos.Pieces[Pawn]&FILES[File(fromId)] == 0 {
			midResult -= int(rookOnFile[1].Middle)
			endResult -= int(rookOnFile[1].End)
		} else if (pos.Pieces[Pawn]&pos.Colours[Black])&FILES[File(fromId)] == 0 {
			midResult -= int(rookOnFile[0].Middle)
			endResult -= int(rookOnFile[0].End)
		} else if mobility < 3 {
			kingFile := File(blackKingLocation)
			if (kingFile < FILE_E) == (File(fromId) < kingFile) {
				cannotCastle := BoolToInt(^pos.Flags&(BlackKingSideCastleFlag|BlackQueenSideCastleFlag) == 0)
				midResult -= int(trappedRook[cannotCastle].Middle)
				endResult -= int(trappedRook[cannotCastle].End)
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
		midResult += int(mobilityBonus[3][mobility].Middle)
		endResult += int(mobilityBonus[3][mobility].End)
		midResult += int(whiteQueensPos[fromId].Middle)
		endResult += int(whiteQueensPos[fromId].End)

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
		midResult -= int(mobilityBonus[3][mobility].Middle)
		endResult -= int(mobilityBonus[3][mobility].End)
		midResult -= int(blackQueensPos[fromId].Middle)
		endResult -= int(blackQueensPos[fromId].End)

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
		midResult += int(tempo.Middle)
		endResult += int(tempo.End)
	} else {
		midResult -= int(tempo.Middle)
		endResult -= int(tempo.End)
	}

	if phase < 0 {
		phase = 0
	}

	// white king
	whiteKingDefenders := PopCount(
		(pos.Pieces[Pawn] | pos.Pieces[Bishop] | pos.Pieces[Knight]) & pos.Colours[White] & whiteKingAreaMask[whiteKingLocation],
	)
	midResult += int(whiteKingPos[whiteKingLocation].Middle)
	endResult += int(whiteKingPos[whiteKingLocation].End)
	midResult += int(kingDefenders[whiteKingDefenders].Middle)
	midResult += int(kingDefenders[whiteKingDefenders].End)
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
			midResult -= count * count / 720
			endResult -= count / 20
		}
	}

	// black king
	blackKingDefenders := PopCount(
		(pos.Pieces[Pawn] | pos.Pieces[Bishop] | pos.Pieces[Knight]) & pos.Colours[Black] & blackKingAreaMask[blackKingLocation],
	)
	midResult -= int(blackKingPos[blackKingLocation].Middle)
	endResult -= int(blackKingPos[blackKingLocation].End)
	midResult -= int(kingDefenders[blackKingDefenders].Middle)
	midResult -= int(kingDefenders[blackKingDefenders].End)
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
			midResult += count * count / 720
			endResult += count / 20
		}
	}

	scale := scaleFactor(pos, endResult)

	// tapering eval
	phase = (phase*256 + (totalPhase / 2)) / totalPhase
	result := (midResult*(256-phase) + (endResult * phase * scale / SCALE_NORMAL)) / 256

	if pos.SideToMove == White {
		return result
	}
	return -result
}

const SCALE_NORMAL = 1
const SCALE_DRAW = 0

func scaleFactor(pos *Position, endResult int) int {
	if (endResult > 0 && PopCount(pos.Colours[White]) == 2 && (pos.Colours[White]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) ||
		(endResult < 0 && PopCount(pos.Colours[Black]) == 2 && (pos.Colours[Black]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0) {
		return SCALE_DRAW
	}
	return SCALE_NORMAL
}
