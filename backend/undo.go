package backend

type Undo struct {
	Key           uint64
	PawnKey       uint64
	FiftyMove     int
	EpSquare      int
	LastMove      Move
	CapturedPiece Piece
	Flags         uint8
}
