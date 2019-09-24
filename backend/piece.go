package backend

type Piece uint8

const NoPiece = Piece(0)

func (p Piece) Type() int {
	return int(p / 4)
}

func (p Piece) Colour() int {
	return int(p & 1)
}

func NewPiece(figure, colour int) Piece {
	return Piece(figure*4 + colour)
}
