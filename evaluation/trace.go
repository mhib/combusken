package evaluation

import . "github.com/mhib/combusken/backend"

type Trace struct {
	PawnValue                    int
	KnightValue                  int
	BishopValue                  int
	RookValue                    int
	QueenValue                   int
	PieceScores                  [King + 1][8][4]int
	PawnScores                   [7][8]int
	PawnsConnected               [8][4]int
	MobilityBonus                [4][32]int
	PassedFriendlyDistance       [8]int
	PassedEnemyDistance          [8]int
	PassedRank                   [7]int
	PassedFile                   [8]int
	Isolated                     int
	Doubled                      int
	Backward                     int
	BackwardOpen                 int
	BishopPair                   int
	BishopRammedPawns            int
	BishopOutpostUndefendedBonus int
	BishopOutpostDefendedBonus   int
	KnightOutpostUndefendedBonus int
	KnightOutpostDefendedBonus   int
	MinorBehindPawn              int
	Tempo                        int
	RookOnFile                   [2]int
	KingDefenders                [12]int
	KingShelter                  [2][8][8]int
	KingStorm                    [2][4][8]int
}
