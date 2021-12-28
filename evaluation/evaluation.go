package evaluation

import (
	. "github.com/mhib/combusken/chess"
	"github.com/mhib/combusken/registeel"
	. "github.com/mhib/combusken/utils"
)

const tuning = false

var T Trace

const PawnPhase = 0
const KnightPhase = 1
const BishopPhase = 1
const RookPhase = 2
const QueenPhase = 4
const TotalPhase = PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2

var PawnPsqt [16][2][64]Score   // BishopFlag, colour, Square
var Psqt [2][King + 1][64]Score // One row for every colour purposefelly left empty

var bishopExistanceTranslations [16][2][2]BishopFlag // Flag, piece colour, square colour

var PawnsConnectedSquare [2][64]Score
var pawnsConnectedMask [2][64]uint64

var passedMask [2][64]uint64

var outpustMask [2][64]uint64

var distanceBetween [64][64]int16

var adjacentFilesMask [8]uint64

var kingAreaMask [2][64]uint64

var forwardRanksMask [2][8]uint64

var forwardFileMask [2][64]uint64

// Outpost bitboards
const whiteOutpustRanks = Rank4_BB | Rank5_BB | Rank6_BB
const blackOutpustRanks = Rank5_BB | Rank4_BB | Rank3_BB

const (
	BlackBlackSquareBishopFlag = 1 << iota
	BlackWhiteSquareBishopFlag
	WhiteBlackSquareBishopFlag
	WhiteWhiteSquareBishopFlag
)

var BishopFlags [2][2]BishopFlag = [2][2]BishopFlag{{BlackBlackSquareBishopFlag, BlackWhiteSquareBishopFlag}, {WhiteBlackSquareBishopFlag, WhiteWhiteSquareBishopFlag}}
var bishopFlagPawnTranslation [16]BishopFlag

type BishopFlag uint8

func (f BishopFlag) BlackPawnPerspective() BishopFlag {
	return bishopFlagPawnTranslation[f]
}

func LoadScoresToPieceSquares() {
	for x := 0; x < 4; x++ {
		for y := 0; y < 7; y++ {
			PawnsConnectedSquare[White][y*8+x] = PawnsConnected[y][x]
			PawnsConnectedSquare[White][y*8+(7-x)] = PawnsConnected[y][x]
			PawnsConnectedSquare[Black][(7-y)*8+x] = PawnsConnected[y][x]
			PawnsConnectedSquare[Black][(7-y)*8+(7-x)] = PawnsConnected[y][x]
		}
	}
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			Psqt[White][Knight][y*8+x] = PieceScores[Knight][y][x] + KnightValue
			Psqt[Black][Knight][(7-y)*8+x] = PieceScores[Knight][y][x] + KnightValue

			Psqt[White][Bishop][y*8+x] = PieceScores[Bishop][y][x] + BishopValue
			Psqt[Black][Bishop][(7-y)*8+x] = PieceScores[Bishop][y][x] + BishopValue

			Psqt[White][Rook][y*8+x] = PieceScores[Rook][y][x] + RookValue
			Psqt[Black][Rook][(7-y)*8+x] = PieceScores[Rook][y][x] + RookValue

			Psqt[White][Queen][y*8+x] = PieceScores[Queen][y][x] + QueenValue
			Psqt[Black][Queen][(7-y)*8+x] = PieceScores[Queen][y][x] + QueenValue

			Psqt[White][King][y*8+x] = PieceScores[King][y][x]
			Psqt[Black][King][(7-y)*8+x] = PieceScores[King][y][x]
		}
	}

	for bishopFlag := 0; bishopFlag <= 15; bishopFlag++ {
		for y := 1; y < 7; y++ {
			for x := 0; x < 8; x++ {
				PawnPsqt[bishopFlag][White][y*8+x] = PawnScores[bishopFlag][y][x] + PawnValue
				PawnPsqt[bishopFlag][Black][(7-y)*8+x] = PawnScores[BishopFlag(bishopFlag).BlackPawnPerspective()][y][x] + PawnValue
			}
		}
	}
}

func init() {
	for flag := 0; flag <= 15; flag++ {
		bishopExistanceTranslations[flag][Black][Black] = BishopFlag(flag)
	}

	for flag := 0; flag <= 15; flag++ {
		var currentFlag uint16
		if flag&BlackBlackSquareBishopFlag != 0 {
			currentFlag |= BlackWhiteSquareBishopFlag
		}
		if flag&BlackWhiteSquareBishopFlag != 0 {
			currentFlag |= BlackBlackSquareBishopFlag
		}
		if flag&WhiteBlackSquareBishopFlag != 0 {
			currentFlag |= WhiteWhiteSquareBishopFlag
		}
		if flag&WhiteWhiteSquareBishopFlag != 0 {
			currentFlag |= WhiteBlackSquareBishopFlag
		}
		bishopExistanceTranslations[flag][Black][White] = BishopFlag(currentFlag)
	}

	for flag := 0; flag <= 15; flag++ {
		var currentFlag uint16
		if flag&BlackBlackSquareBishopFlag != 0 {
			currentFlag |= WhiteBlackSquareBishopFlag
		}
		if flag&BlackWhiteSquareBishopFlag != 0 {
			currentFlag |= WhiteWhiteSquareBishopFlag
		}
		if flag&WhiteBlackSquareBishopFlag != 0 {
			currentFlag |= BlackBlackSquareBishopFlag
		}
		if flag&WhiteWhiteSquareBishopFlag != 0 {
			currentFlag |= BlackWhiteSquareBishopFlag
		}
		bishopExistanceTranslations[flag][White][Black] = BishopFlag(currentFlag)
	}

	for flag := 0; flag <= 15; flag++ {
		var currentFlag uint16
		if flag&BlackBlackSquareBishopFlag != 0 {
			currentFlag |= WhiteWhiteSquareBishopFlag
		}
		if flag&BlackWhiteSquareBishopFlag != 0 {
			currentFlag |= WhiteBlackSquareBishopFlag
		}
		if flag&WhiteBlackSquareBishopFlag != 0 {
			currentFlag |= BlackWhiteSquareBishopFlag
		}
		if flag&WhiteWhiteSquareBishopFlag != 0 {
			currentFlag |= BlackBlackSquareBishopFlag
		}
		bishopExistanceTranslations[flag][White][White] = BishopFlag(currentFlag)
	}
	for flag := uint8(0); flag < 16; flag++ {
		var res BishopFlag
		if flag&WhiteWhiteSquareBishopFlag != 0 {
			res |= BlackBlackSquareBishopFlag
		}
		if flag&WhiteBlackSquareBishopFlag != 0 {
			res |= BlackWhiteSquareBishopFlag
		}
		if flag&BlackWhiteSquareBishopFlag != 0 {
			res |= WhiteBlackSquareBishopFlag
		}
		if flag&BlackBlackSquareBishopFlag != 0 {
			res |= WhiteWhiteSquareBishopFlag
		}
		bishopFlagPawnTranslation[flag] = res
	}
	LoadScoresToPieceSquares()

	// Pawn is passed if no pawn of opposite color can stop it from promoting
	for i := 8; i <= 55; i++ {
		passedMask[White][i] = 0
		for file := File(i) - 1; file <= File(i)+1; file++ {
			if file < FileA || file > FileH {
				continue
			}
			for rank := Rank(i) + 1; rank < Rank8; rank++ {
				passedMask[White][i] |= 1 << uint(rank*8+file)
			}
		}
	}
	// Outpust is similar to passed bitboard bot we do not care about pawns in same file
	for i := 8; i <= 55; i++ {
		outpustMask[White][i] = passedMask[White][i] & ^Files_BB[File(i)]
	}

	for i := 55; i >= 8; i-- {
		passedMask[Black][i] = 0
		for file := File(i) - 1; file <= File(i)+1; file++ {
			if file < FileA || file > FileH {
				continue
			}
			for rank := Rank(i) - 1; rank > Rank1; rank-- {
				passedMask[Black][i] |= 1 << uint(rank*8+file)
			}
		}
	}
	for i := 55; i >= 8; i-- {
		outpustMask[Black][i] = passedMask[Black][i] & ^Files_BB[File(i)]
	}

	for i := 8; i <= 55; i++ {
		pawnsConnectedMask[White][i] = PawnAttacks[Black][i] | PawnAttacks[Black][i+8]
		pawnsConnectedMask[Black][i] = PawnAttacks[White][i] | PawnAttacks[White][i-8]
	}

	for i := range Files_BB {
		adjacentFilesMask[i] = 0
		if i != 0 {
			adjacentFilesMask[i] |= Files_BB[i-1]
		}
		if i != 7 {
			adjacentFilesMask[i] |= Files_BB[i+1]
		}
	}

	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			distanceBetween[y][x] = int16(Max(Abs(Rank(y)-Rank(x)), Abs(File(y)-File(x))))
		}
	}

	for y := 0; y < 64; y++ {
		kingAreaMask[White][y] = KingAttacks[y] | SquareMask[y] | North(KingAttacks[y])
		kingAreaMask[Black][y] = KingAttacks[y] | SquareMask[y] | South(KingAttacks[y])
		if File(y) > FileA {
			kingAreaMask[White][y] |= West(kingAreaMask[White][y])
			kingAreaMask[Black][y] |= West(kingAreaMask[Black][y])
		}
		if File(y) < FileH {
			kingAreaMask[White][y] |= East(kingAreaMask[White][y])
			kingAreaMask[Black][y] |= East(kingAreaMask[Black][y])
		}
	}

	for rank := Rank1; rank <= Rank8; rank++ {
		for y := rank; y <= Rank8; y++ {
			forwardRanksMask[White][rank] |= Ranks_BB[y]
		}
		forwardRanksMask[Black][rank] = (^forwardRanksMask[White][rank]) | Ranks_BB[rank]
	}

	for y := 0; y < 64; y++ {
		forwardFileMask[White][y] = forwardRanksMask[White][Rank(y)] & Files_BB[File(y)] & ^SquareMask[y]
		forwardFileMask[Black][y] = forwardRanksMask[Black][Rank(y)] & Files_BB[File(y)] & ^SquareMask[y]
	}
}

func evaluateKingPawns(pos *Position) Score {
	if !tuning {
		if ok, score := GlobalPawnKingTable.Get(pos.PawnKey); ok {
			return score
		}
	}
	var fromBB uint64
	var fromId int
	whitePawns := pos.Pieces[Pawn] & pos.Colours[White]
	blackPawns := pos.Pieces[Pawn] & pos.Colours[Black]
	var kingLocation [2]int
	kingLocation[White] = BitScan(pos.Pieces[King] & pos.Colours[White])
	kingLocation[Black] = BitScan(pos.Pieces[King] & pos.Colours[Black])
	score := ScoreZero

	// white pawns
	for fromBB = whitePawns; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		neighbours := adjacentFilesMask[File(fromId)] & whitePawns

		// Isolated pawn penalty
		if neighbours == 0 {
			score += Isolated
			if tuning {
				T.Isolated++
			}
		}

		// Pawn is backward if there are no pawns behind it and cannot increase rank without being attacked by enemy pawn
		if passedMask[Black][fromId]&whitePawns == 0 &&
			PawnAttacks[White][fromId+8]&blackPawns != 0 {
			if Files_BB[File(fromId)]&blackPawns == 0 {
				score += BackwardOpen
				if tuning {
					T.BackwardOpen++
				}
			} else {
				score += Backward
				if tuning {
					T.Backward++
				}
			}
		} else if pawnsConnectedMask[White][fromId]&whitePawns != 0 {
			score += PawnsConnectedSquare[White][fromId]
			if tuning {
				T.PawnsConnected[Rank(fromId)][FileMirror[File(fromId)]]++
			}
		}

		// Note that Passed has its own stacked evaluation
		if forwardFileMask[White][fromId]&whitePawns != 0 && passedMask[White][fromId]&blackPawns != 0 {
			friendlyBlockers := passedMask[White][fromId] & blackPawns
			isDoubled := BoolToInt(SquareMask[fromId+8]&whitePawns != 0)
			canBeTraded := BoolToInt(friendlyBlockers & ^(forwardFileMask[White][fromId]&whitePawns) != 0 ||
				(friendlyBlockers != 0 && (neighbours != 0 || PawnAttacks[White][fromId]&blackPawns != 0)))
			score += StackedPawns[isDoubled][canBeTraded][File(fromId)]
			if tuning {
				T.StackedPawns[isDoubled][canBeTraded][File(fromId)]++
			}
		}
	}

	// black pawns
	for fromBB = blackPawns; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		neighbours := adjacentFilesMask[File(fromId)] & blackPawns
		if neighbours == 0 {
			score -= Isolated
			if tuning {
				T.Isolated--
			}
		}
		if passedMask[White][fromId]&blackPawns == 0 &&
			PawnAttacks[Black][fromId-8]&whitePawns != 0 {
			if Files_BB[File(fromId)]&whitePawns == 0 {
				score -= BackwardOpen
				if tuning {
					T.BackwardOpen--
				}
			} else {
				score -= Backward
				if tuning {
					T.Backward--
				}
			}
		} else if pawnsConnectedMask[Black][fromId]&blackPawns != 0 {
			score -= PawnsConnectedSquare[Black][fromId]
			if tuning {
				T.PawnsConnected[7-Rank(fromId)][FileMirror[File(fromId)]]--
			}
		}
		// Note that Passed has its own stacked evaluation
		if forwardFileMask[Black][fromId]&blackPawns != 0 && passedMask[Black][fromId]&whitePawns != 0 {
			friendlyBlockers := passedMask[Black][fromId] & blackPawns
			isDoubled := BoolToInt(SquareMask[fromId-8]&blackPawns != 0)
			canBeTraded := BoolToInt(friendlyBlockers & ^(forwardFileMask[Black][fromId]&blackPawns) != 0 ||
				(friendlyBlockers != 0 && (neighbours != 0 || PawnAttacks[Black][fromId]&whitePawns != 0)))
			score -= StackedPawns[isDoubled][canBeTraded][File(fromId)]
			if tuning {
				T.StackedPawns[isDoubled][canBeTraded][File(fromId)]--
			}
		}
	}

	// White king storm shelter
	for file := Max(File(kingLocation[White])-1, FileA); file <= Min(File(kingLocation[White])+1, FileH); file++ {
		ours := pos.Pieces[Pawn] & Files_BB[file] & pos.Colours[White] & forwardRanksMask[White][Rank(kingLocation[White])]
		var ourDist int
		if ours == 0 {
			ourDist = 7
		} else {
			ourDist = Abs(Rank(kingLocation[White]) - Rank(BitScan(ours)))
		}
		theirs := pos.Pieces[Pawn] & Files_BB[file] & pos.Colours[Black] & forwardRanksMask[White][Rank(kingLocation[White])]
		var theirDist int
		if theirs == 0 {
			theirDist = 7
		} else {
			theirDist = Abs(Rank(kingLocation[White]) - Rank(BitScan(theirs)))
		}
		sameFile := BoolToInt(file == File(kingLocation[White]))
		score += KingShelter[sameFile][file][ourDist]
		if tuning {
			T.KingShelter[sameFile][file][ourDist]++
		}

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		score += KingStorm[blocked][file][theirDist]

		if tuning {
			T.KingStorm[blocked][file][theirDist]++
		}
	}
	if KingFlank_BB[File(kingLocation[White])]&pos.Pieces[Pawn] == 0 {
		score += KingOnPawnlessFlank
		if tuning {
			T.KingOnPawnlessFlank++
		}
	}

	// Black king storm / shelter
	for file := Max(File(kingLocation[Black])-1, FileA); file <= Min(File(kingLocation[Black])+1, FileH); file++ {
		ours := pos.Pieces[Pawn] & Files_BB[file] & pos.Colours[Black] & forwardRanksMask[Black][Rank(kingLocation[Black])]
		var ourDist int
		if ours == 0 {
			ourDist = 7
		} else {
			ourDist = Abs(Rank(kingLocation[Black]) - Rank(MostSignificantBit(ours)))
		}
		theirs := pos.Pieces[Pawn] & Files_BB[file] & pos.Colours[White] & forwardRanksMask[Black][Rank(kingLocation[Black])]
		var theirDist int
		if theirs == 0 {
			theirDist = 7
		} else {
			theirDist = Abs(Rank(kingLocation[Black]) - Rank(MostSignificantBit(theirs)))
		}
		sameFile := BoolToInt(file == File(kingLocation[Black]))
		score -= KingShelter[sameFile][file][ourDist]
		if tuning {
			T.KingShelter[sameFile][file][ourDist]--
		}

		blocked := BoolToInt(ourDist != 7 && ourDist == theirDist-1)
		score -= KingStorm[blocked][file][theirDist]
		if tuning {
			T.KingStorm[blocked][file][theirDist]--
		}
	}
	if KingFlank_BB[File(kingLocation[Black])]&pos.Pieces[Pawn] == 0 {
		score -= KingOnPawnlessFlank
		if tuning {
			T.KingOnPawnlessFlank--
		}
	}
	if !tuning {
		GlobalPawnKingTable.Set(pos.PawnKey, score)
	}
	return score
}

type EvaluationContext struct {
	registeel.RegisteelNetwork
	contempt Score
}

func (ec *EvaluationContext) SetContempt(contempt Score) {
	ec.contempt = contempt
}

func (ec *EvaluationContext) Evaluate(pos *Position) int {
	var fromId int
	var fromBB uint64
	var attacks uint64

	var attacked [2]uint64
	var attackedBy [2][King + 1]uint64
	var attackedByTwo [2]uint64
	var kingAttacksCount [2]int16
	var kingAttackersCount [2]int16
	var kingAttackersWeight [2]Score

	var bishopFlag BishopFlag

	var kingLocation [2]int
	var kingArea [2]uint64

	phase := TotalPhase
	whiteMobilityArea := ^((pos.Pieces[Pawn] & pos.Colours[White]) | (BlackPawnsAttacks(pos.Pieces[Pawn] & pos.Colours[Black])))
	blackMobilityArea := ^((pos.Pieces[Pawn] & pos.Colours[Black]) | (WhitePawnsAttacks(pos.Pieces[Pawn] & pos.Colours[White])))
	allOccupation := pos.Colours[White] | pos.Colours[Black]

	kingLocation[White] = BitScan(pos.Pieces[King] & pos.Colours[White])
	attacks = KingAttacks[kingLocation[White]]
	attacked[White] |= attacks
	attackedBy[White][King] |= attacks
	kingArea[White] = kingAreaMask[White][kingLocation[White]]

	kingLocation[Black] = BitScan(pos.Pieces[King] & pos.Colours[Black])
	attacks = KingAttacks[kingLocation[Black]]
	attacked[Black] |= attacks
	attackedBy[Black][King] |= attacks
	kingArea[Black] = kingAreaMask[Black][kingLocation[Black]]

	// white pawns
	whitePawns := pos.Pieces[Pawn] & pos.Colours[White]
	attacks = WhitePawnsAttacks(whitePawns)
	attackedByTwo[White] |= attacked[White] & attacks
	attackedByTwo[White] |= WhitePawnsDoubleAttacks(whitePawns)
	attacked[White] |= attacks
	attackedBy[White][Pawn] |= attacks
	kingAttacksCount[White] += int16(PopCount(attacks & kingArea[Black]))

	// black pawns
	blackPawns := pos.Pieces[Pawn] & pos.Colours[Black]
	attacks = BlackPawnsAttacks(blackPawns)
	attackedByTwo[Black] |= attacked[Black] & attacks
	attackedByTwo[Black] |= BlackPawnsDoubleAttacks(blackPawns)
	attacked[Black] |= attacks
	attackedBy[Black][Pawn] |= attacks
	kingAttacksCount[Black] += int16(PopCount(attacks & kingArea[White]))

	score := ec.contempt + evaluateKingPawns(pos)

	// white knights
	for fromBB = pos.Pieces[Knight] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= KnightPhase
		fromId = BitScan(fromBB)

		attacks = KnightAttacks[fromId]
		mobility := PopCount(whiteMobilityArea & attacks)
		score += Psqt[White][Knight][fromId]
		score += MobilityBonus[0][mobility]
		if tuning {
			T.KnightValue++
			T.PieceScores[Knight][Rank(fromId)][File(fromId)]++
			T.MobilityBonus[0][mobility]++
		}

		attackedByTwo[White] |= attacked[White] & attacks
		attacked[White] |= attacks
		attackedBy[White][Knight] |= attacks

		if (pos.Pieces[Pawn]>>8)&SquareMask[fromId] != 0 {
			score += MinorBehindPawn
			if tuning {
				T.MinorBehindPawn++
			}
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && outpustMask[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			if PawnAttacks[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
				score += KnightOutpostDefendedBonus
				if tuning {
					T.KnightOutpostDefendedBonus++
				}
			} else {
				score += KnightOutpostUndefendedBonus
				if tuning {
					T.KnightOutpostUndefendedBonus++
				}
			}
		}

		kingDistance := Min(int(distanceBetween[fromId][kingLocation[White]]), int(distanceBetween[fromId][kingLocation[Black]]))
		if kingDistance >= 4 {
			score += DistantKnight[kingDistance-4]
			if tuning {
				T.DistantKnight[kingDistance-4]++
			}
		}
		if attacks&kingArea[Black] != 0 {
			kingAttacksCount[White] += int16(PopCount(attacks & kingArea[Black]))
			kingAttackersCount[White]++
			kingAttackersWeight[White] += KingSafetyAttacksWeights[Knight]
			if tuning {
				T.KingSafetyAttacksWeights[Black][Knight]++
			}
		}
	}

	// black knights
	for fromBB = pos.Pieces[Knight] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= KnightPhase
		fromId = BitScan(fromBB)

		attacks = KnightAttacks[fromId]
		mobility := PopCount(blackMobilityArea & attacks)
		score -= Psqt[Black][Knight][fromId]
		score -= MobilityBonus[0][mobility]
		if tuning {
			T.KnightValue--
			T.PieceScores[Knight][7-Rank(fromId)][File(fromId)]--
			T.MobilityBonus[0][mobility]--
		}

		attackedByTwo[Black] |= attacked[Black] & attacks
		attacked[Black] |= attacks
		attackedBy[Black][Knight] |= attacks

		if (pos.Pieces[Pawn]<<8)&SquareMask[fromId] != 0 {
			score -= MinorBehindPawn
			if tuning {
				T.MinorBehindPawn--
			}
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && outpustMask[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			if PawnAttacks[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
				score -= KnightOutpostDefendedBonus
				if tuning {
					T.KnightOutpostDefendedBonus--
				}
			} else {
				score -= KnightOutpostUndefendedBonus
				if tuning {
					T.KnightOutpostUndefendedBonus--
				}
			}
		}
		kingDistance := Min(int(distanceBetween[fromId][kingLocation[White]]), int(distanceBetween[fromId][kingLocation[Black]]))
		if kingDistance >= 4 {
			score -= DistantKnight[kingDistance-4]
			if tuning {
				T.DistantKnight[kingDistance-4]--
			}
		}
		if attacks&kingArea[White] != 0 {
			kingAttacksCount[Black] += int16(PopCount(attacks & kingArea[White]))
			kingAttackersCount[Black]++
			kingAttackersWeight[Black] += KingSafetyAttacksWeights[Knight]

			if tuning {
				T.KingSafetyAttacksWeights[White][Knight]++
			}
		}
	}

	// white bishops
	whiteRammedPawns := South(pos.Pieces[Pawn]&pos.Colours[Black]) & (pos.Pieces[Pawn] & pos.Colours[White])
	for fromBB = pos.Pieces[Bishop] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= BishopPhase
		fromId = BitScan(fromBB)

		bishopFlag |= BishopFlags[White][Colour(fromId)]
		attacks = BishopAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		score += MobilityBonus[1][mobility]
		score += Psqt[White][Bishop][fromId]
		if tuning {
			T.BishopValue++
			T.PieceScores[Bishop][Rank(fromId)][File(fromId)]++
			T.MobilityBonus[1][mobility]++
		}

		attackedByTwo[White] |= attacked[White] & attacks
		attacked[White] |= attacks
		attackedBy[White][Bishop] |= attacks

		if (pos.Pieces[Pawn]>>8)&SquareMask[fromId] != 0 {
			score += MinorBehindPawn
			if tuning {
				T.MinorBehindPawn++
			}
		}
		if (LongDiagonals_BB&SquareMask[fromId]) != 0 && (MoreThanOne(BishopAttacks(fromId, pos.Pieces[Pawn]) & Center_BB)) {
			score += LongDiagonalBishop
			if tuning {
				T.LongDiagonalBishop++
			}
		}
		if SquareMask[fromId]&whiteOutpustRanks != 0 && outpustMask[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			if PawnAttacks[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) != 0 {
				score += BishopOutpostDefendedBonus
				if tuning {
					T.BishopOutpostDefendedBonus++
				}
			} else {
				score += BishopOutpostUndefendedBonus
				if tuning {
					T.BishopOutpostUndefendedBonus++
				}
			}
		}

		// Bishop is worth less if there are friendly rammed pawns of its color
		var rammedCount Score
		if SquareMask[fromId]&WhiteSquares_BB != 0 {
			rammedCount = Score(PopCount(whiteRammedPawns & WhiteSquares_BB))
		} else {
			rammedCount = Score(PopCount(whiteRammedPawns & BlackSquares_BB))
		}
		score += BishopRammedPawns * rammedCount
		if tuning {
			T.BishopRammedPawns += int(rammedCount)
		}

		kingDistance := Min(int(distanceBetween[fromId][kingLocation[White]]), int(distanceBetween[fromId][kingLocation[Black]]))
		if kingDistance >= 4 {
			score += DistantBishop[kingDistance-4]
			if tuning {
				T.DistantBishop[kingDistance-4]++
			}
		}
		if attacks&kingArea[Black] != 0 {
			kingAttacksCount[White] += int16(PopCount(attacks & kingArea[Black]))
			kingAttackersCount[White]++
			kingAttackersWeight[White] += KingSafetyAttacksWeights[Bishop]
			if tuning {
				T.KingSafetyAttacksWeights[Black][Bishop]++
			}
		}
	}

	// Bishop pair bonus
	// It is not checked if bishops have opposite colors, but that is almost always the case
	if MoreThanOne(pos.Pieces[Bishop] & pos.Colours[White]) {
		score += BishopPair
		if tuning {
			T.BishopPair++
		}
	}

	// black bishops
	blackRammedPawns := North(pos.Pieces[Pawn]&pos.Colours[White]) & (pos.Pieces[Pawn] & pos.Colours[Black])
	for fromBB = pos.Pieces[Bishop] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= BishopPhase
		fromId = BitScan(fromBB)

		bishopFlag |= BishopFlags[Black][Colour(fromId)]
		attacks = BishopAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		score -= MobilityBonus[1][mobility]
		score -= Psqt[Black][Bishop][fromId]
		if tuning {
			T.BishopValue--
			T.PieceScores[Bishop][7-Rank(fromId)][File(fromId)]--
			T.MobilityBonus[1][mobility]--
		}

		attackedByTwo[Black] |= attacked[Black] & attacks
		attacked[Black] |= attacks
		attackedBy[Black][Bishop] |= attacks

		if (pos.Pieces[Pawn]<<8)&SquareMask[fromId] != 0 {
			score -= MinorBehindPawn
			if tuning {
				T.MinorBehindPawn--
			}
		}
		if (LongDiagonals_BB&SquareMask[fromId]) != 0 && (MoreThanOne(BishopAttacks(fromId, pos.Pieces[Pawn]) & Center_BB)) {
			score -= LongDiagonalBishop
			if tuning {
				T.LongDiagonalBishop--
			}
		}
		if SquareMask[fromId]&blackOutpustRanks != 0 && outpustMask[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			if PawnAttacks[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) != 0 {
				score -= BishopOutpostDefendedBonus
				if tuning {
					T.BishopOutpostDefendedBonus--
				}
			} else {
				score -= BishopOutpostUndefendedBonus
				if tuning {
					T.BishopOutpostUndefendedBonus--
				}
			}
		}
		var rammedCount Score
		if SquareMask[fromId]&WhiteSquares_BB != 0 {
			rammedCount = Score(PopCount(blackRammedPawns & WhiteSquares_BB))
		} else {
			rammedCount = Score(PopCount(blackRammedPawns & BlackSquares_BB))
		}
		score -= BishopRammedPawns * rammedCount
		if tuning {
			T.BishopRammedPawns -= int(rammedCount)
		}

		kingDistance := Min(int(distanceBetween[fromId][kingLocation[White]]), int(distanceBetween[fromId][kingLocation[Black]]))
		if kingDistance >= 4 {
			score -= DistantBishop[kingDistance-4]
			if tuning {
				T.DistantBishop[kingDistance-4]--
			}
		}
		if attacks&kingArea[White] != 0 {
			kingAttacksCount[Black] += int16(PopCount(attacks & kingArea[White]))
			kingAttackersCount[Black]++
			kingAttackersWeight[Black] += KingSafetyAttacksWeights[Bishop]
			if tuning {
				T.KingSafetyAttacksWeights[White][Bishop]++
			}
		}
	}

	if MoreThanOne(pos.Pieces[Bishop] & pos.Colours[Black]) {
		score -= BishopPair

		if tuning {
			T.BishopPair--
		}
	}

	// white rooks
	for fromBB = pos.Pieces[Rook] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= RookPhase
		fromId = BitScan(fromBB)

		attacks = RookAttacks(fromId, allOccupation)
		mobility := PopCount(whiteMobilityArea & attacks)
		score += MobilityBonus[2][mobility]
		score += Psqt[White][Rook][fromId]
		score += RookBishopExistence[bishopExistanceTranslations[bishopFlag][White][Colour(fromId)]]

		if tuning {
			T.RookValue++
			T.PieceScores[Rook][Rank(fromId)][File(fromId)]++
			T.MobilityBonus[2][mobility]++
			T.RookBishopExistence[bishopExistanceTranslations[bishopFlag][White][Colour(fromId)]]++
		}

		attackedByTwo[White] |= attacked[White] & attacks
		attacked[White] |= attacks
		attackedBy[White][Rook] |= attacks

		if pos.Pieces[Pawn]&Files_BB[File(fromId)] == 0 {
			score += RookOnFile[1]
			if tuning {
				T.RookOnFile[1]++
			}
		} else if (pos.Pieces[Pawn]&pos.Colours[White])&Files_BB[File(fromId)] == 0 {
			score += RookOnFile[0]
			if tuning {
				T.RookOnFile[0]++
			}
		}

		if FileBB(fromId)&pos.Pieces[Queen] != 0 {
			score += RookOnQueenFile
			if tuning {
				T.RookOnQueenFile++
			}
		}

		if mobility <= 3 {
			kingFile := File(kingLocation[White])
			if (kingFile <= FileE) == (File(fromId) < kingFile) && pos.Flags|WhiteQueenSideCastleFlag|WhiteKingSideCastleFlag == pos.Flags {
				score += TrappedRook
				if tuning {
					T.TrappedRook++
				}
			}
		}

		if attacks&kingArea[Black] != 0 {
			kingAttacksCount[White] += int16(PopCount(attacks & kingArea[Black]))
			kingAttackersCount[White]++
			kingAttackersWeight[White] += KingSafetyAttacksWeights[Rook]

			if tuning {
				T.KingSafetyAttacksWeights[Black][Rook]++
			}
		}
	}

	// black rooks
	for fromBB = pos.Pieces[Rook] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= RookPhase
		fromId = BitScan(fromBB)

		attacks = RookAttacks(fromId, allOccupation)
		mobility := PopCount(blackMobilityArea & attacks)
		score -= MobilityBonus[2][mobility]
		score -= Psqt[Black][Rook][fromId]
		score -= RookBishopExistence[bishopExistanceTranslations[bishopFlag][Black][Colour(fromId)]]

		if tuning {
			T.RookValue--
			T.PieceScores[Rook][7-Rank(fromId)][File(fromId)]--
			T.MobilityBonus[2][mobility]--
			T.RookBishopExistence[bishopExistanceTranslations[bishopFlag][Black][Colour(fromId)]]--
		}

		attackedByTwo[Black] |= attacked[Black] & attacks
		attacked[Black] |= attacks
		attackedBy[Black][Rook] |= attacks

		if pos.Pieces[Pawn]&Files_BB[File(fromId)] == 0 {
			score -= RookOnFile[1]
			if tuning {
				T.RookOnFile[1]--
			}
		} else if (pos.Pieces[Pawn]&pos.Colours[Black])&Files_BB[File(fromId)] == 0 {
			score -= RookOnFile[0]
			if tuning {
				T.RookOnFile[0]--
			}
		}

		if FileBB(fromId)&pos.Pieces[Queen] != 0 {
			score -= RookOnQueenFile
			if tuning {
				T.RookOnQueenFile--
			}
		}

		if mobility <= 3 {
			kingFile := File(kingLocation[Black])
			if (kingFile <= FileE) == (File(fromId) < kingFile) && pos.Flags|BlackQueenSideCastleFlag|BlackKingSideCastleFlag == pos.Flags {
				score -= TrappedRook
				if tuning {
					T.TrappedRook--
				}
			}
		}

		if attacks&kingArea[White] != 0 {
			kingAttacksCount[Black] += int16(PopCount(attacks & kingArea[White]))
			kingAttackersCount[Black]++
			kingAttackersWeight[Black] += KingSafetyAttacksWeights[Rook]
			if tuning {
				T.KingSafetyAttacksWeights[White][Rook]++
			}
		}
	}

	//white queens
	for fromBB = pos.Pieces[Queen] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= QueenPhase
		fromId = BitScan(fromBB)

		rookAttacks := RookAttacks(fromId, allOccupation)
		bishopAttacks := BishopAttacks(fromId, allOccupation)

		attacks = rookAttacks | bishopAttacks
		mobility := PopCount(whiteMobilityArea & attacks)
		score += MobilityBonus[3][mobility]
		score += Psqt[White][Queen][fromId]
		score += QueenBishopExistence[bishopExistanceTranslations[bishopFlag][White][Colour(fromId)]]

		if tuning {
			T.QueenValue++
			T.PieceScores[Queen][Rank(fromId)][File(fromId)]++
			T.MobilityBonus[3][mobility]++
			T.QueenBishopExistence[bishopExistanceTranslations[bishopFlag][White][Colour(fromId)]]++
		}

		pinningRooks := pos.Pieces[Rook] & pos.Colours[Black] & ^rookAttacks
		pinningBishops := pos.Pieces[Bishop] & pos.Colours[Black] & ^bishopAttacks
		if (pinningRooks != 0 && RookAttacks(fromId, allOccupation & ^rookAttacks)&pinningRooks != 0) ||
			(pinningBishops != 0 && BishopAttacks(fromId, allOccupation & ^bishopAttacks)&pinningBishops != 0) {
			score += QueenPinned

			if tuning {
				T.QueenPinned++
			}
		}

		attackedByTwo[White] |= attacked[White] & attacks
		attacked[White] |= attacks
		attackedBy[White][Queen] |= attacks

		if attacks&kingArea[Black] != 0 {
			kingAttacksCount[White] += int16(PopCount(attacks & kingArea[Black]))
			kingAttackersCount[White]++
			kingAttackersWeight[White] += KingSafetyAttacksWeights[Queen]
			if tuning {
				T.KingSafetyAttacksWeights[Black][Queen]++
			}
		}
	}

	// black queens
	for fromBB = pos.Pieces[Queen] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		phase -= QueenPhase
		fromId = BitScan(fromBB)

		rookAttacks := RookAttacks(fromId, allOccupation)
		bishopAttacks := BishopAttacks(fromId, allOccupation)

		attacks = rookAttacks | bishopAttacks
		mobility := PopCount(blackMobilityArea & attacks)
		score -= MobilityBonus[3][mobility]
		score -= Psqt[Black][Queen][fromId]
		score -= QueenBishopExistence[bishopExistanceTranslations[bishopFlag][Black][Colour(fromId)]]

		if tuning {
			T.QueenValue--
			T.PieceScores[Queen][7-Rank(fromId)][File(fromId)]--
			T.MobilityBonus[3][mobility]--
			T.QueenBishopExistence[bishopExistanceTranslations[bishopFlag][Black][Colour(fromId)]]--
		}

		pinningRooks := pos.Pieces[Rook] & pos.Colours[White] & ^rookAttacks
		pinningBishops := pos.Pieces[Bishop] & pos.Colours[White] & ^bishopAttacks
		if (pinningRooks != 0 && RookAttacks(fromId, allOccupation & ^rookAttacks)&pinningRooks != 0) ||
			(pinningBishops != 0 && BishopAttacks(fromId, allOccupation & ^bishopAttacks)&pinningBishops != 0) {
			score -= QueenPinned

			if tuning {
				T.QueenPinned--
			}
		}

		attackedByTwo[Black] |= attacked[Black] & attacks
		attacked[Black] |= attacks
		attackedBy[Black][Queen] |= attacks
		if attacks&kingArea[White] != 0 {
			kingAttacksCount[Black] += int16(PopCount(attacks & kingArea[White]))
			kingAttackersCount[Black]++
			kingAttackersWeight[Black] += KingSafetyAttacksWeights[Queen]
			if tuning {
				T.KingSafetyAttacksWeights[White][Queen]++
			}
		}
	}

	// white pawns
	for fromBB = pos.Pieces[Pawn] & pos.Colours[White]; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		score += PawnPsqt[bishopFlag][White][fromId]
		if tuning {
			T.PawnValue++
			T.PawnScores[bishopFlag][Rank(fromId)][File(fromId)]++
		}

		// Passed bonus
		if passedMask[White][fromId]&(pos.Pieces[Pawn]&pos.Colours[Black]) == 0 {
			rookOrQueenBehind := forwardFileMask[Black][fromId]&
				// Try to not to count situation when a stacked pawn is between current pawn and a major piece
				(^pos.Pieces[Pawn])&
				(pos.Pieces[Rook]|pos.Pieces[Queen])&
				pos.Colours[White] != 0

			// Bonus is calculated based on ability to push, rank, file, distance from friendly and enemy king
			advance := North(SquareMask[fromId])
			canAdvance := BoolToInt(allOccupation&advance == 0)
			safeAdvance := BoolToInt(attacked[Black]&advance == 0)
			advanceDefended := BoolToInt(rookOrQueenBehind || (attacked[White]&advance) != 0)
			score +=
				PassedRank[canAdvance][safeAdvance][advanceDefended][Rank(fromId)] +
					PassedFile[File(fromId)] +
					PassedFriendlyDistance[distanceBetween[kingLocation[White]][fromId]] +
					PassedEnemyDistance[distanceBetween[kingLocation[Black]][fromId]]

			if tuning {
				T.PassedRank[canAdvance][safeAdvance][advanceDefended][Rank(fromId)]++
				T.PassedFile[File(fromId)]++
				T.PassedFriendlyDistance[distanceBetween[kingLocation[White]][fromId]]++
				T.PassedEnemyDistance[distanceBetween[kingLocation[Black]][fromId]]++
			}

			push := forwardFileMask[White][fromId]

			stacked := pos.Pieces[Pawn]&pos.Colours[White]&push != 0

			if stacked {
				score += PassedStacked[Rank(fromId)]
				if tuning {
					T.PassedStacked[Rank(fromId)]++
				}
			}
			// Rank seventh's push == advance so it is already calculated
			if !stacked && Rank(fromId) != Rank7 {
				if push&(attacked[Black]|pos.Colours[Black]) == 0 {
					if rookOrQueenBehind || (push&attacked[White]) == push {
						score += PassedPushUncontestedDefended[Rank(fromId)]
						if tuning {
							T.PassedPushUncontestedDefended[Rank(fromId)]++
						}
					} else {
						score += PassedUncontested[Rank(fromId)]
						if tuning {
							T.PassedUncontested[Rank(fromId)]++
						}
					}
				} else if rookOrQueenBehind || (push&attacked[White]) == push {
					score += PassedPushDefended[Rank(fromId)]
					if tuning {
						T.PassedPushDefended[Rank(fromId)]++
					}
				}
			}
		}
	}

	{
		safeWhitePawns := ((^attacked[Black]) | attacked[White]) & pos.Pieces[Pawn] & pos.Colours[White]
		blackPiecesAttackedByPawns := WhitePawnsAttacks(safeWhitePawns) & pos.Colours[Black] & (^pos.Pieces[Pawn])
		if blackPiecesAttackedByPawns > 0 {
			knightsAttacked := PopCount(blackPiecesAttackedByPawns & pos.Pieces[Knight])
			bishopsAttacked := PopCount(blackPiecesAttackedByPawns & pos.Pieces[Bishop])
			rooksAttacked := PopCount(blackPiecesAttackedByPawns & pos.Pieces[Rook])
			queensAttacked := PopCount(blackPiecesAttackedByPawns & pos.Pieces[Queen])
			kingAttacked := BoolToInt(blackPiecesAttackedByPawns&pos.Pieces[King] != 0)
			score += AttackedBySafePawn[0]*Score(knightsAttacked) +
				AttackedBySafePawn[1]*Score(bishopsAttacked) +
				AttackedBySafePawn[2]*Score(rooksAttacked) +
				AttackedBySafePawn[3]*Score(queensAttacked) +
				AttackedBySafePawn[4]*Score(kingAttacked)
			if tuning {
				T.AttackedBySafePawn[0] += knightsAttacked
				T.AttackedBySafePawn[1] += bishopsAttacked
				T.AttackedBySafePawn[2] += rooksAttacked
				T.AttackedBySafePawn[3] += queensAttacked
				T.AttackedBySafePawn[4] += kingAttacked
			}
		}
	}

	{
		safeBlackPawns := pos.Pieces[Pawn] & pos.Colours[Black]
		whitePiecesAttackedByPawns := BlackPawnsAttacks(safeBlackPawns) & pos.Colours[White] & (^pos.Pieces[Pawn])
		if whitePiecesAttackedByPawns > 0 {
			knightsAttacked := PopCount(whitePiecesAttackedByPawns & pos.Pieces[Knight])
			bishopsAttacked := PopCount(whitePiecesAttackedByPawns & pos.Pieces[Bishop])
			rooksAttacked := PopCount(whitePiecesAttackedByPawns & pos.Pieces[Rook])
			queensAttacked := PopCount(whitePiecesAttackedByPawns & pos.Pieces[Queen])
			kingAttacked := BoolToInt(whitePiecesAttackedByPawns&pos.Pieces[King] != 0)
			score -= AttackedBySafePawn[0]*Score(knightsAttacked) +
				AttackedBySafePawn[1]*Score(bishopsAttacked) +
				AttackedBySafePawn[2]*Score(rooksAttacked) +
				AttackedBySafePawn[3]*Score(queensAttacked) +
				AttackedBySafePawn[4]*Score(kingAttacked)
			if tuning {
				T.AttackedBySafePawn[0] -= knightsAttacked
				T.AttackedBySafePawn[1] -= bishopsAttacked
				T.AttackedBySafePawn[2] -= rooksAttacked
				T.AttackedBySafePawn[3] -= queensAttacked
				T.AttackedBySafePawn[4] -= kingAttacked
			}
		}
	}

	// black pawns
	for fromBB = pos.Pieces[Pawn] & pos.Colours[Black]; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)

		score -= PawnPsqt[bishopFlag][Black][fromId]

		if tuning {
			T.PawnValue--
			T.PawnScores[bishopFlag.BlackPawnPerspective()][7-Rank(fromId)][File(fromId)]--
		}

		if passedMask[Black][fromId]&(pos.Pieces[Pawn]&pos.Colours[White]) == 0 {
			rookOrQueenBehind := forwardFileMask[White][fromId]&
				// Try to not to count situation when a stacked pawn is between current pawn and a major piece
				(^pos.Pieces[Pawn])&
				(pos.Pieces[Rook]|pos.Pieces[Queen])&
				pos.Colours[Black] != 0

			advance := South(SquareMask[fromId])
			canAdvance := BoolToInt(allOccupation&advance == 0)
			safeAdvance := BoolToInt(attacked[White]&advance == 0)
			advanceDefended := BoolToInt(rookOrQueenBehind || (attacked[Black]&advance) != 0)
			score -=
				PassedRank[canAdvance][safeAdvance][advanceDefended][7-Rank(fromId)] +
					PassedFile[File(fromId)] +
					PassedFriendlyDistance[distanceBetween[kingLocation[Black]][fromId]] +
					PassedEnemyDistance[distanceBetween[kingLocation[White]][fromId]]
			if tuning {
				T.PassedRank[canAdvance][safeAdvance][advanceDefended][7-Rank(fromId)]--
				T.PassedFile[File(fromId)]--
				T.PassedFriendlyDistance[distanceBetween[kingLocation[Black]][fromId]]--
				T.PassedEnemyDistance[distanceBetween[kingLocation[White]][fromId]]--
			}

			push := forwardFileMask[Black][fromId]
			stacked := pos.Pieces[Pawn]&pos.Colours[Black]&push != 0
			if stacked {
				score -= PassedStacked[7-Rank(fromId)]
				if tuning {
					T.PassedStacked[7-Rank(fromId)]--
				}
			}

			if !stacked && Rank(fromId) != Rank2 {
				if push&(attacked[White]|pos.Colours[White]) == 0 {
					if rookOrQueenBehind || (push&attacked[Black]) == push {
						score -= PassedPushUncontestedDefended[7-Rank(fromId)]
						if tuning {
							T.PassedPushUncontestedDefended[7-Rank(fromId)]--
						}
					} else {
						score -= PassedUncontested[7-Rank(fromId)]
						if tuning {
							T.PassedUncontested[7-Rank(fromId)]--
						}
					}
				} else if rookOrQueenBehind || (push&attacked[Black]) == push {
					score -= PassedPushDefended[7-Rank(fromId)]
					if tuning {
						T.PassedPushDefended[7-Rank(fromId)]--
					}
				}
			}
		}
	}

	if phase < 0 {
		phase = 0
	}

	// white king
	whiteKingDefenders := PopCount(
		(pos.Pieces[Pawn] | pos.Pieces[Bishop] | pos.Pieces[Knight]) & pos.Colours[White] & kingAreaMask[White][kingLocation[White]],
	)
	score += Psqt[White][King][kingLocation[White]]
	score += KingDefenders[whiteKingDefenders]
	score += KingBishopExistence[bishopExistanceTranslations[bishopFlag][White][Colour(kingLocation[White])]]
	if tuning {
		T.PieceScores[King][Rank(kingLocation[White])][File(kingLocation[White])]++
		T.KingDefenders[whiteKingDefenders]++
		T.KingBishopExistence[bishopExistanceTranslations[bishopFlag][White][Colour(kingLocation[White])]]++
	}

	// Weak squares are attacked by the enemy, defended no more
	// than once and only defended by our Queens or our King
	// Idea from Ethereal
	weakForWhite := attacked[Black] & ^attackedByTwo[White] & (^attacked[White] | attackedBy[White][Queen] | attackedBy[White][King])
	if int(kingAttackersCount[Black]) > 1-PopCount(pos.Colours[Black]&pos.Pieces[Queen]) {
		safe := ^pos.Colours[Black] & (^attacked[White] | (weakForWhite & attackedByTwo[Black]))

		knightThreats := KnightAttacks[kingLocation[White]]
		bishopThreats := BishopAttacks(kingLocation[White], allOccupation)
		rookThreats := RookAttacks(kingLocation[White], allOccupation)
		queenThreats := bishopThreats | rookThreats

		knightChecks := knightThreats & safe & attackedBy[Black][Knight]
		bishopChecks := bishopThreats & safe & attackedBy[Black][Bishop]
		rookChecks := rookThreats & safe & attackedBy[Black][Rook]
		queenChecks := queenThreats & safe & attackedBy[Black][Queen]
		safetyScore := Score(kingAttackersCount[Black]) * kingAttackersWeight[Black]
		numerator := 9 * int(kingAttackersCount[Black])
		denumarator := int(PopCount(kingArea[White]))
		safetyScore += S(int16(int(KingSafetyAttackValue.Middle())*numerator/denumarator), int16(int(KingSafetyAttackValue.End())*numerator/denumarator))
		safetyScore += KingSafetyWeakSquares * Score(PopCount(kingArea[White]&weakForWhite))
		safetyScore += KingSafetyFriendlyPawns * Score(PopCount(pos.Colours[White]&pos.Pieces[Pawn]&kingArea[White] & ^weakForWhite))
		safetyScore += KingSafetyNoEnemyQueens * Score(BoolToInt(pos.Colours[Black]&pos.Pieces[Queen] == 0))
		safetyScore += KingSafetySafeQueenCheck * Score(PopCount(queenChecks))
		safetyScore += KingSafetySafeRookCheck * Score(PopCount(rookChecks))
		safetyScore += KingSafetySafeBishopCheck * Score(PopCount(bishopChecks))
		safetyScore += KingSafetySafeKnightCheck * Score(PopCount(knightChecks))
		safetyScore += KingSafetyAdjustment
		middle := int(safetyScore.Middle())
		end := int(safetyScore.End())
		score += S(
			int16((-middle*Max(middle, 0))/720),
			-int16(Max(end, 0)/20),
		)

		if tuning {
			for piece := Knight; piece <= Queen; piece++ {
				T.KingSafetyAttacksWeights[White][piece] *= int(kingAttackersCount[Black])
			}
			if numerator > 0 {
				T.KingSafetyAttackValueNumerator[White] = numerator
				T.KingSafetyAttackValueDenumerator[White] = denumarator
			}
			T.KingSafetyWeakSquares[White] = PopCount(kingArea[White] & weakForWhite)
			T.KingSafetyFriendlyPawns[White] = PopCount(pos.Colours[White] & pos.Pieces[Pawn] & kingArea[White] & ^weakForWhite)
			T.KingSafetyNoEnemyQueens[White] = BoolToInt(pos.Colours[Black]&pos.Pieces[Queen] == 0)
			T.KingSafetySafeQueenCheck[White] = PopCount(queenChecks)
			T.KingSafetySafeRookCheck[White] = PopCount(rookChecks)
			T.KingSafetySafeBishopCheck[White] = PopCount(bishopChecks)
			T.KingSafetySafeKnightCheck[White] = PopCount(knightChecks)
			T.KingSafetyAdjustment[White] = 1
		}
	}

	// black king
	blackKingDefenders := PopCount(
		(pos.Pieces[Pawn] | pos.Pieces[Bishop] | pos.Pieces[Knight]) & pos.Colours[Black] & kingAreaMask[Black][kingLocation[Black]],
	)
	score -= Psqt[Black][King][kingLocation[Black]]
	score -= KingDefenders[blackKingDefenders]
	score -= KingBishopExistence[bishopExistanceTranslations[bishopFlag][Black][Colour(kingLocation[Black])]]
	if tuning {
		T.PieceScores[King][7-Rank(kingLocation[Black])][File(kingLocation[Black])]--
		T.KingDefenders[blackKingDefenders]--
		T.KingBishopExistence[bishopExistanceTranslations[bishopFlag][Black][Colour(kingLocation[Black])]]--
	}

	// Weak squares are attacked by the enemy, defended no more
	// than once and only defended by our Queens or our King
	// Idea from Ethereal
	weakForBlack := attacked[White] & ^attackedByTwo[Black] & (^attacked[Black] | attackedBy[Black][Queen] | attackedBy[Black][King])
	if int(kingAttackersCount[White]) > 1-PopCount(pos.Colours[White]&pos.Pieces[Queen]) {
		safe := ^pos.Colours[White] & (^attacked[Black] | (weakForBlack & attackedByTwo[White]))

		knightThreats := KnightAttacks[kingLocation[Black]]
		bishopThreats := BishopAttacks(kingLocation[Black], allOccupation)
		rookThreats := RookAttacks(kingLocation[Black], allOccupation)
		queenThreats := bishopThreats | rookThreats

		knightChecks := knightThreats & safe & attackedBy[White][Knight]
		bishopChecks := bishopThreats & safe & attackedBy[White][Bishop]
		rookChecks := rookThreats & safe & attackedBy[White][Rook]
		queenChecks := queenThreats & safe & attackedBy[White][Queen]

		safetyScore := Score(kingAttackersCount[White]) * kingAttackersWeight[White]
		numerator := 9 * int(kingAttackersCount[White])
		denumerator := PopCount(kingArea[Black])
		safetyScore += S(int16(int(KingSafetyAttackValue.Middle())*numerator/denumerator), int16(int(KingSafetyAttackValue.End())*numerator/denumerator))
		safetyScore += KingSafetyWeakSquares * Score(PopCount(kingArea[Black]&weakForBlack))
		safetyScore += KingSafetyFriendlyPawns * Score(PopCount(pos.Colours[Black]&pos.Pieces[Pawn]&kingArea[Black] & ^weakForBlack))
		safetyScore += KingSafetyNoEnemyQueens * Score(BoolToInt(pos.Colours[White]&pos.Pieces[Queen] == 0))
		safetyScore += KingSafetySafeQueenCheck * Score(PopCount(queenChecks))
		safetyScore += KingSafetySafeRookCheck * Score(PopCount(rookChecks))
		safetyScore += KingSafetySafeBishopCheck * Score(PopCount(bishopChecks))
		safetyScore += KingSafetySafeKnightCheck * Score(PopCount(knightChecks))
		safetyScore += KingSafetyAdjustment
		middle := int(safetyScore.Middle())
		end := int(safetyScore.End())
		score -= S(
			int16((-middle*Max(middle, 0))/720),
			-int16(Max(end, 0)/20),
		)

		if tuning {
			for piece := Knight; piece <= Queen; piece++ {
				T.KingSafetyAttacksWeights[Black][piece] *= int(kingAttackersCount[White])
			}
			if numerator > 0 {
				T.KingSafetyAttackValueNumerator[Black] = numerator
				T.KingSafetyAttackValueDenumerator[Black] = denumerator
			}
			T.KingSafetyWeakSquares[Black] = PopCount(kingArea[Black] & weakForBlack)
			T.KingSafetyFriendlyPawns[Black] = PopCount(pos.Colours[Black] & pos.Pieces[Pawn] & kingArea[Black] & ^weakForBlack)
			T.KingSafetyNoEnemyQueens[Black] = BoolToInt(pos.Colours[White]&pos.Pieces[Queen] == 0)
			T.KingSafetySafeQueenCheck[Black] = PopCount(queenChecks)
			T.KingSafetySafeRookCheck[Black] = PopCount(rookChecks)
			T.KingSafetySafeBishopCheck[Black] = PopCount(bishopChecks)
			T.KingSafetySafeKnightCheck[Black] = PopCount(knightChecks)
			T.KingSafetyAdjustment[Black] = 1
		}

	}

	// White threats
	blackStronglyProtected := attackedBy[Black][Pawn] | (attackedByTwo[Black] & ^attackedByTwo[White])
	blackDefended := pos.Colours[Black] & ^pos.Pieces[Pawn] & blackStronglyProtected
	if ((pos.Colours[Black] & weakForBlack) | blackDefended) != 0 {
		for fromBB = pos.Colours[Black] & (blackDefended | weakForBlack) & (attackedBy[White][Knight] | attackedBy[White][Bishop]) & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score += ThreatByMinor[threatenedPiece]
			if tuning {
				T.ThreatByMinor[threatenedPiece]++
			}
		}

		for fromBB = pos.Colours[Black] & (blackDefended | weakForBlack) & attackedBy[White][Rook] & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) & ^pos.Pieces[Pawn] {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score += ThreatByRook[threatenedPiece]
			if tuning {
				T.ThreatByRook[threatenedPiece]++
			}
		}

		if weakForBlack&pos.Colours[Black]&attackedBy[White][King] != 0 {
			score += ThreatByKing
			if tuning {
				T.ThreatByKing++
			}
		}

		// Bonus if enemy has a hanging piece
		score += Hanging *
			Score(PopCount((pos.Colours[Black] & ^pos.Pieces[Pawn] & attackedByTwo[White])&weakForBlack))

		if tuning {
			T.Hanging += PopCount((pos.Colours[Black] & ^pos.Pieces[Pawn] & attackedByTwo[White]) & weakForBlack)
		}

	}

	// Black threats
	whiteStronglyProtected := attackedBy[White][Pawn] | (attackedByTwo[White] & ^attackedByTwo[Black])
	whiteDefended := pos.Colours[White] & ^pos.Pieces[Pawn] & whiteStronglyProtected
	if ((pos.Colours[White] & weakForWhite) | whiteDefended) != 0 {
		for fromBB = pos.Colours[White] & (whiteDefended | weakForWhite) & (attackedBy[Black][Knight] | attackedBy[Black][Bishop]) & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score -= ThreatByMinor[threatenedPiece]
			if tuning {
				T.ThreatByMinor[threatenedPiece]--
			}
		}

		for fromBB = pos.Colours[White] & (whiteDefended | weakForWhite) & attackedBy[Black][Rook] & ^pos.Pieces[Pawn]; fromBB != 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			threatenedPiece := pos.TypeOnSquare(SquareMask[fromId])
			score -= ThreatByRook[threatenedPiece]
			if tuning {
				T.ThreatByRook[threatenedPiece]--
			}
		}

		if weakForWhite&pos.Colours[White]&attackedBy[Black][King] != 0 {
			score -= ThreatByKing
			if tuning {
				T.ThreatByKing--
			}
		}

		// Bonus if enemy has a hanging piece
		score -= Hanging *
			Score(PopCount(pos.Colours[White] & ^pos.Pieces[Pawn] & attackedByTwo[Black] & weakForWhite))

		if tuning {
			T.Hanging -= PopCount(pos.Colours[White] & ^pos.Pieces[Pawn] & attackedByTwo[Black] & weakForWhite)
		}
	}

	{
		sign := BoolToInt(score.End() > 0) - BoolToInt(score.End() < 0)
		pawnsOnBothFlanks := BoolToInt((pos.Pieces[Pawn]&KingSide_BB != 0) && (pos.Pieces[Pawn]&QueenSide_BB != 0))
		pawnEndgame := BoolToInt(pos.Pieces[Knight]|pos.Pieces[Bishop]|pos.Pieces[Rook]|pos.Pieces[Queen] == 0)
		infiltration := BoolToInt(Rank(kingLocation[White]) > Rank4 || Rank(kingLocation[Black]) < Rank5)

		complexity := ComplexityTotalPawns*Score(PopCount(pos.Pieces[Pawn])) +
			ComplexityPawnBothFlanks*Score(pawnsOnBothFlanks) +
			ComplexityPawnEndgame*Score(pawnEndgame) +
			ComplexityInfiltration*Score(infiltration) +
			ComplexityAdjustment

		if tuning {
			T.ComplexityTotalPawns = PopCount(pos.Pieces[Pawn])
			T.ComplexityPawnBothFlanks = pawnsOnBothFlanks
			T.ComplexityPawnEndgame = pawnEndgame
			T.ComplexityInfiltration = infiltration
			T.ComplexityAdjustment = 1

			T.BeforeComplexity = score
			T.Complexity = complexity
		}
		score += S(0, int16(sign*Max(int(complexity.End()), -Abs(int(score.End())))))
	}

	scale := ScaleNormal
	{
		winning := BoolToInt(score.End() > 0)
		switch PopCount(pos.Colours[winning]) {
		case 2:
			if (pos.Colours[winning] & (pos.Pieces[Bishop] | pos.Pieces[Knight])) != 0 {
				scale = ScaleDraw
			} else if (pos.Colours[winning]&pos.Pieces[Rook]) != 0 && (pos.Colours[winning^1]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0 {
				scale = ScaleDraw
			}
		case 3:
			if OnlyOne(pos.Colours[winning]&pos.Pieces[Rook]) && (pos.Colours[winning]&(pos.Pieces[Bishop]|pos.Pieces[Knight])) != 0 && (pos.Colours[winning^1]&pos.Pieces[Rook]) != 0 {
				scale = ScaleDraw
			}
		default:
			if OnlyOne(pos.Colours[Black]&pos.Pieces[Bishop]) &&
				OnlyOne(pos.Colours[White]&pos.Pieces[Bishop]) &&
				OnlyOne(pos.Pieces[Bishop]&WhiteSquares_BB) &&
				(pos.Pieces[Knight]|pos.Pieces[Rook]|pos.Pieces[Queen]) == 0 {
				scale = ScaleHard
			}
		}

	}

	if !tuning {
		if scale != ScaleDraw {
			score += ec.CorrectEvaluation(pos)
		}
	}

	if tuning {
		T.Scale = scale
	}

	// tapering eval
	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase
	result := (int(score.Middle())*(256-phase) + (int(score.End()) * phase * scale / ScaleNormal)) / 256

	if pos.SideToMove == White {
		return result + int(Tempo)
	}
	return -result + int(Tempo)
}

const ScaleNormal = 2
const ScaleHard = 1
const ScaleDraw = 0
