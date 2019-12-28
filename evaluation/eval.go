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

var PawnValue = Score{103, 122}
var KnightValue = Score{475, 433}
var BishopValue = Score{442, 427}
var RookValue = Score{653, 693}
var QueenValue = Score{1436, 1323}

// Piece Square Values
var pieceScores = [King + 1][8][4]Score{
	{},
	{ // knight
		{Score{-55, -38}, Score{-15, -54}, Score{-39, -33}, Score{-8, -24}},
		{Score{-17, -61}, Score{-22, -37}, Score{-12, -30}, Score{0, -22}},
		{Score{-21, -40}, Score{1, -32}, Score{-2, -16}, Score{8, -2}},
		{Score{-17, -29}, Score{13, -26}, Score{13, 1}, Score{9, 2}},
		{Score{-1, -35}, Score{7, -21}, Score{14, 4}, Score{29, 3}},
		{Score{-39, -53}, Score{28, -47}, Score{-75, 19}, Score{28, -6}},
		{Score{-77, -51}, Score{-44, -28}, Score{64, -56}, Score{-1, -20}},
		{Score{-227, -64}, Score{-75, -77}, Score{-127, -32}, Score{7, -57}},
	},
	{ // Bishop
		{Score{6, -11}, Score{16, -1}, Score{18, -7}, Score{15, -2}},
		{Score{13, -25}, Score{53, -30}, Score{34, -13}, Score{18, -4}},
		{Score{24, -19}, Score{40, -12}, Score{39, -1}, Score{23, 8}},
		{Score{-1, -24}, Score{3, -20}, Score{17, -6}, Score{20, -3}},
		{Score{-33, -14}, Score{6, -20}, Score{-28, -1}, Score{24, -6}},
		{Score{-131, 13}, Score{-50, -9}, Score{-225, 47}, Score{-50, -4}},
		{Score{-67, 6}, Score{34, -9}, Score{-6, 4}, Score{12, -10}},
		{Score{-18, -19}, Score{-31, -12}, Score{-113, -1}, Score{-93, 3}},
	},
	{ // Rook
		{Score{-1, -21}, Score{-10, -6}, Score{4, -11}, Score{10, -13}},
		{Score{-44, -1}, Score{-6, -17}, Score{-10, -11}, Score{1, -12}},
		{Score{-38, -9}, Score{-14, -8}, Score{-3, -18}, Score{-12, -14}},
		{Score{-37, 3}, Score{-8, -3}, Score{-16, 1}, Score{-13, -3}},
		{Score{-37, 8}, Score{-22, 1}, Score{10, 5}, Score{-6, 1}},
		{Score{-25, 3}, Score{16, 1}, Score{18, -3}, Score{-8, 3}},
		{Score{7, 12}, Score{6, 16}, Score{48, 3}, Score{58, -7}},
		{Score{6, 12}, Score{8, 10}, Score{-20, 18}, Score{18, 12}},
	},
	{ // Queen
		{Score{1, -58}, Score{15, -72}, Score{18, -61}, Score{36, -68}},
		{Score{-3, -41}, Score{11, -52}, Score{34, -52}, Score{30, -47}},
		{Score{-6, -13}, Score{23, -30}, Score{0, 5}, Score{4, -6}},
		{Score{-1, -18}, Score{-24, 34}, Score{-6, 18}, Score{-20, 44}},
		{Score{-20, 12}, Score{-28, 27}, Score{-32, 32}, Score{-34, 51}},
		{Score{19, -30}, Score{0, -9}, Score{-4, 15}, Score{-9, 45}},
		{Score{-9, -24}, Score{-56, 17}, Score{-17, 26}, Score{-36, 52}},
		{Score{6, -20}, Score{2, -1}, Score{17, 6}, Score{19, 14}},
	},
	{ // King
		{Score{186, -14}, Score{170, 21}, Score{95, 65}, Score{82, 55}},
		{Score{179, 21}, Score{136, 46}, Score{70, 78}, Score{39, 88}},
		{Score{81, 48}, Score{108, 54}, Score{42, 82}, Score{39, 90}},
		{Score{14, 48}, Score{73, 57}, Score{31, 88}, Score{-2, 99}},
		{Score{43, 55}, Score{102, 77}, Score{66, 96}, Score{88, 88}},
		{Score{94, 66}, Score{255, 67}, Score{221, 87}, Score{165, 64}},
		{Score{44, 69}, Score{121, 80}, Score{109, 101}, Score{166, 75}},
		{Score{29, 10}, Score{163, 29}, Score{110, 71}, Score{30, 56}},
	},
}

// Pawns Square scores
var pawnScores = [7][8]Score{
	{},
	{Score{-19, 4}, Score{20, -2}, Score{-9, 10}, Score{3, 5}, Score{1, 9}, Score{-9, 13}, Score{19, -2}, Score{-20, 4}},
	{Score{-16, -4}, Score{-8, -3}, Score{1, -2}, Score{1, -6}, Score{2, -2}, Score{2, -4}, Score{-9, -6}, Score{-10, -4}},
	{Score{-20, 4}, Score{-9, 1}, Score{12, -7}, Score{21, -12}, Score{21, -8}, Score{14, -5}, Score{-10, 3}, Score{-16, 3}},
	{Score{0, 14}, Score{25, -1}, Score{14, -6}, Score{35, -11}, Score{35, -14}, Score{11, 0}, Score{29, 2}, Score{-8, 16}},
	{Score{10, 38}, Score{31, 27}, Score{51, 6}, Score{48, 0}, Score{55, -5}, Score{75, 10}, Score{12, 35}, Score{9, 42}},
	{Score{-7, 65}, Score{42, 62}, Score{3, 38}, Score{4, 38}, Score{92, 28}, Score{-8, 48}, Score{32, 47}, Score{-70, 80}},
}

var pawnsConnected = [8][4]Score{
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
	{Score{11, -20}, Score{5, 6}, Score{6, -8}, Score{3, 16}},
	{Score{7, 0}, Score{30, 2}, Score{10, 8}, Score{15, 20}},
	{Score{9, 5}, Score{23, 6}, Score{16, 11}, Score{21, 14}},
	{Score{13, 16}, Score{9, 25}, Score{29, 23}, Score{32, 19}},
	{Score{11, 58}, Score{41, 60}, Score{75, 55}, Score{65, 49}},
	{Score{45, 58}, Score{146, -1}, Score{166, 22}, Score{282, 57}},
	{Score{0, 0}, Score{0, 0}, Score{0, 0}, Score{0, 0}},
}

var mobilityBonus = [...][32]Score{
	{Score{-45, -121}, Score{-34, -66}, Score{-23, -36}, Score{-22, -12}, Score{-7, -13}, Score{4, -4}, // Knights
		Score{14, -9}, Score{26, -13}, Score{38, -28}},
	{Score{-29, -73}, Score{-13, -59}, Score{9, -27}, Score{14, -8}, Score{27, 0}, Score{36, 8}, // Bishops
		Score{42, 11}, Score{46, 12}, Score{49, 14}, Score{51, 15}, Score{67, 3}, Score{91, 3},
		Score{46, 28}, Score{81, 11}},
	{Score{-23, -35}, Score{-27, -29}, Score{-14, 18}, Score{-10, 35}, Score{-3, 48}, Score{1, 57}, // Rooks
		Score{4, 67}, Score{17, 66}, Score{16, 66}, Score{33, 67}, Score{37, 71}, Score{37, 73},
		Score{50, 72}, Score{61, 68}, Score{90, 57}},
	{Score{-22, -20}, Score{-54, -8}, Score{-3, -205}, Score{-10, -180}, Score{-5, -21}, Score{2, -24}, // Queens
		Score{1, -29}, Score{12, -5}, Score{15, 19}, Score{20, 21}, Score{22, 37}, Score{20, 46},
		Score{28, 38}, Score{27, 61}, Score{30, 64}, Score{28, 75}, Score{32, 70}, Score{27, 71},
		Score{37, 66}, Score{39, 80}, Score{65, 55}, Score{61, 55}, Score{69, 46}, Score{60, 32},
		Score{68, 16}, Score{31, 31}, Score{0, -7}, Score{4, -1}},
}

var passedFriendlyDistance = [8]Score{
	Score{0, 0}, Score{4, 28}, Score{-3, 12}, Score{-3, -8},
	Score{-14, -18}, Score{-15, -19}, Score{6, -27}, Score{-31, -15},
}

var passedEnemyDistance = [8]Score{
	Score{0, 0}, Score{-56, -77}, Score{25, -31}, Score{17, 4},
	Score{10, 26}, Score{2, 37}, Score{0, 42}, Score{-13, 47},
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
var passedRank = [7]Score{Score{0, 0}, Score{1, -36}, Score{-7, -14}, Score{-10, 28}, Score{24, 72}, Score{39, 160}, Score{110, 253}}

// PassedFile[File] contains a bonus according to the file of a passed pawn
var passedFile = [8]Score{Score{-4, 22}, Score{-21, 21}, Score{-24, 9}, Score{-26, -6},
	Score{-20, -3}, Score{7, -2}, Score{-10, 17}, Score{-6, 10},
}

var isolated = Score{-10, -11}
var doubled = Score{-12, -37}
var backward = Score{3, -1}
var backwardOpen = Score{-17, -4}

var bishopPair = Score{38, 63}
var bishopRammedPawns = Score{-10, -13}

var bishopOutpostUndefendedBonus = Score{74, -1}
var bishopOutpostDefendedBonus = Score{104, 5}

var knightOutpostUndefendedBonus = Score{35, -14}
var knightOutpostDefendedBonus = Score{56, 16}

var minorBehindPawn = Score{4, 28}

var tempo = Score{33, 29}

// Rook on semiopen, open file
var rookOnFile = [2]Score{Score{10, 24}, Score{55, -3}}

// king can castle / king cannot castle
var trappedRook = [2]Score{Score{-1, -17}, Score{-66, -4}}

var kingDefenders = [12]Score{
	Score{-79, 0}, Score{-61, -4}, Score{-34, -3}, Score{-10, -2},
	Score{1, 1}, Score{17, 3}, Score{29, 4}, Score{38, 4},
	Score{49, 2}, Score{42, 7}, Score{7, 4}, Score{7, 4},
}

var kingShelter = [2][8][8]Score{
	{{Score{-33, 5}, Score{-9, -11}, Score{-8, 8}, Score{38, -10},
		Score{3, -20}, Score{2, 1}, Score{6, -10}, Score{-29, 13}},
		{Score{21, -1}, Score{33, -14}, Score{-5, -2}, Score{-7, 6},
			Score{-28, -2}, Score{8, -9}, Score{19, -39}, Score{-29, 7}},
		{Score{13, 9}, Score{9, 2}, Score{-17, 4}, Score{-18, 8},
			Score{-32, 1}, Score{-6, -1}, Score{0, -5}, Score{-18, 3}},
		{Score{-14, 21}, Score{0, 7}, Score{-9, -6}, Score{-4, -3},
			Score{1, -23}, Score{-6, -14}, Score{5, -39}, Score{-29, 3}},
		{Score{0, 7}, Score{-1, 2}, Score{-21, 0}, Score{-23, 10},
			Score{-13, -7}, Score{-15, 1}, Score{-9, -4}, Score{-27, 8}},
		{Score{46, -9}, Score{29, -17}, Score{-7, -9}, Score{-6, -8},
			Score{1, -18}, Score{0, -3}, Score{41, -28}, Score{-13, 4}},
		{Score{27, -2}, Score{-1, -5}, Score{-20, -9}, Score{-10, -3},
			Score{-11, -9}, Score{15, -5}, Score{10, -21}, Score{-36, 16}},
		{Score{-22, 2}, Score{-19, -5}, Score{-6, 4}, Score{-15, 12},
			Score{-5, 9}, Score{-13, 17}, Score{-32, 7}, Score{-57, 35}}},
	{{Score{-1, -2}, Score{-44, -17}, Score{-22, -4}, Score{-80, -22},
		Score{-5, -17}, Score{-33, -30}, Score{-75, 0}, Score{-67, 19}},
		{Score{8, 26}, Score{3, -17}, Score{-19, -5}, Score{-4, -3},
			Score{-3, -2}, Score{11, -41}, Score{2, -16}, Score{-61, 20}},
		{Score{16, 34}, Score{42, -5}, Score{13, -4}, Score{15, -12},
			Score{16, 3}, Score{-22, -10}, Score{59, -16}, Score{-27, 10}},
		{Score{9, 24}, Score{-32, 18}, Score{-18, 9}, Score{-15, -1},
			Score{-24, 21}, Score{-71, 31}, Score{-21, -2}, Score{-44, 2}},
		{Score{2, 47}, Score{1, 4}, Score{1, 0}, Score{-7, -1},
			Score{-9, 11}, Score{0, -10}, Score{0, -12}, Score{-34, 9}},
		{Score{72, -19}, Score{27, -13}, Score{-12, -1}, Score{-5, -14},
			Score{-7, -6}, Score{-24, -15}, Score{25, -32}, Score{-33, 5}},
		{Score{1, 10}, Score{8, -15}, Score{3, -14}, Score{-19, -9},
			Score{-17, -11}, Score{-1, -16}, Score{3, -18}, Score{-71, 24}},
		{Score{2, 0}, Score{-4, -25}, Score{-4, -13}, Score{-23, -8},
			Score{-20, -3}, Score{7, -5}, Score{-28, -25}, Score{-62, 31}}},
}

var kingStorm = [2][4][8]Score{
	{{Score{19, 1}, Score{7, 2}, Score{14, 2}, Score{6, 6},
		Score{1, 13}, Score{5, 9}, Score{-8, 16}, Score{1, -12}},
		{Score{14, 2}, Score{4, 5}, Score{19, -2}, Score{2, 8},
			Score{10, 4}, Score{8, 1}, Score{3, -2}, Score{3, -11}},
		{Score{15, 12}, Score{-1, 9}, Score{2, 7}, Score{-9, 14},
			Score{-5, 9}, Score{5, 2}, Score{10, -11}, Score{5, -5}},
		{Score{19, 12}, Score{2, 4}, Score{10, 1}, Score{-1, 5},
			Score{-6, 9}, Score{9, 8}, Score{4, 7}, Score{-1, 0}}},
	{{Score{0, 0}, Score{8, 14}, Score{-16, 7}, Score{21, -4},
		Score{15, 12}, Score{-10, 14}, Score{11, 44}, Score{12, -28}},
		{Score{0, 0}, Score{9, -33}, Score{-15, -6}, Score{70, -16},
			Score{47, -17}, Score{-30, 3}, Score{1, 21}, Score{4, -19}},
		{Score{0, 0}, Score{-74, -3}, Score{-32, 0}, Score{14, 1},
			Score{1, -1}, Score{7, -17}, Score{69, -48}, Score{2, -5}},
		{Score{0, 0}, Score{0, -21}, Score{12, -17}, Score{-8, -1},
			Score{-4, 2}, Score{10, -18}, Score{-4, 1}, Score{-6, 16}}},
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

var kingSafetyAttacksWeights = [King + 1]int16{0, -3, -1, -1, 8, 0}
var kingSafetyAttackValue int16 = 109
var kingSafetyWeakSquares int16 = 19
var kingSafetyFriendlyPawns int16 = -7
var kingSafetyNoEnemyQueens int16 = 60
var kingSafetySafeQueenCheck int16 = 76
var kingSafetySafeRookCheck int16 = 109
var kingSafetySafeBishopCheck int16 = 97
var kingSafetySafeKnightCheck int16 = 143
var kingSafetyAdjustment int16 = -74

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
