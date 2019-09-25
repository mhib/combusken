package backend

type Piece uint8

const NoPiece = Piece(None * 4)

func (s Piece) Type() int {
	return int(s / 4)
}

func (s Piece) Colour() int {
	return int(s & 1)
}

func NewPiece(figure, colour int) Piece {
	return Piece(figure*4 + colour)
}
