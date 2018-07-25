package backend

import "fmt"

// Move is represented as int32\
// Bits are:
// 4 bits for move type
// 3 bits for captured piece type
// 3 bits for piece type
// 6 bits to
// 6 bits from

const (
	QuietMove              = 0
	Capture                = 1
	DoublePawnPush         = 8
	KingCastle             = 4
	QueenCastle            = 5
	EPCapture              = 9
	KnightPromotion        = 2
	BishopPromotion        = 10
	RookPromotion          = 6
	QueenPromotion         = 14
	KnightCapturePromotion = 3
	BishopCapturePromotion = 11
	RookCapturePromotion   = 7
	QueenCapturePromotion  = 15
)

type Move int32

func (m Move) From() int {
	return int(m & 0x3f)
}

func (m Move) IsCapture() bool {
	return m&(1<<18) != 0
}

func (m Move) To() int {
	return int((m >> 6) & 0x3f)
}

func (m Move) MovedPiece() int {
	return int((m >> 12) & 0x7)
}

func (m Move) CapturedPiece() int {
	return int((m >> 15) & 0x7)
}

func (m Move) Type() int {
	return int(m>>18) & 0xf
}

func NewMove(from, to, pieceType, capturedType, moveType int) Move {
	return Move(from | (to << 6) | (pieceType << 12) | (capturedType << 15) | (moveType << 18))
}

func NewType(capture, promotion, s1, s0 int) int {
	return capture | (promotion << 1) | (s1 << 2) | (s0 << 3)
}

func (m Move) Inspect() {
	fmt.Print("From: ")
	fmt.Print(m.From())
	fmt.Print(" To: ")
	fmt.Print(m.To())
	fmt.Print(" Piece: ")
	fmt.Print(m.MovedPiece())
	fmt.Print(" Captured Piece: ")
	fmt.Print(m.CapturedPiece())
	fmt.Print(" Type: ")
	fmt.Print(m.Type())
}
