package evaluation

import . "github.com/mhib/combusken/backend"

type Trace struct {
	PawnValue                     int
	KnightValue                   int
	BishopValue                   int
	RookValue                     int
	QueenValue                    int
	PawnScores                    [16][7][8]int
	PieceScores                   [King + 1][8][8]int
	PawnsConnected                [7][4]int
	MobilityBonus                 [4][32]int
	PassedFriendlyDistance        [8]int
	PassedEnemyDistance           [8]int
	PassedRank                    [2][2][2][7]int
	PassedFile                    [8]int
	PassedStacked                 [7]int
	PassedUncontested             [6]int
	PassedPushDefended            [6]int
	PassedPushUncontestedDefended [6]int
	Isolated                      int
	Doubled                       int
	Backward                      int
	BackwardOpen                  int
	BishopPair                    int
	BishopRammedPawns             int
	BishopOutpostUndefendedBonus  int
	BishopOutpostDefendedBonus    int
	LongDiagonalBishop            int
	DistantBishop                 [4]int
	KnightOutpostUndefendedBonus  int
	KnightOutpostDefendedBonus    int
	DistantKnight                 [4]int
	MinorBehindPawn               int
	RookOnFile                    [2]int
	RookOnQueenFile               int
	KingDefenders                 [12]int
	KingShelter                   [2][8][8]int
	KingStorm                     [2][4][8]int
	KingOnPawnlessFlank           int
	Hanging                       int
	ThreatByKing                  int
	ThreatByMinor                 [King + 1]int
	ThreatByRook                  [King + 1]int
	RookBishopExistence           [16]int
	QueenBishopExistence          [16]int
	KingBishopExistence           [16]int
	Scale                         int
}
