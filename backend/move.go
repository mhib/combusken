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
	QueenCastle            = 12
	EPCapture              = 9
	KnightPromotion        = 2
	BishopPromotion        = 6
	RookPromotion          = 10
	QueenPromotion         = 14
	KnightCapturePromotion = 3
	BishopCapturePromotion = 7
	RookCapturePromotion   = 11
	QueenCapturePromotion  = 15
)

type Move int32

const (
	NullMove = Move(0)
)

var WhiteKingSideCastle = NewMove(E1, G1, King, None, KingCastle)
var WhiteQueenSideCastle = NewMove(E1, C1, King, None, QueenCastle)
var BlackKingSideCastle = NewMove(E8, G8, King, None, KingCastle)
var BlackQueenSideCastle = NewMove(E8, C8, King, None, QueenCastle)

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
	return int(m >> 18)
}

func (m Move) Special() int {
	return int(m >> 20)
}

func (m Move) IsPromotion() bool {
	return m&(1<<19) != 0
}

func (m Move) IsCaptureOrPromotion() bool {
	return m&((1<<19)|(1<<18)) != 0
}

func (m Move) IsCastling() bool {
	t := m >> 18 // Type() inlined
	return t&3 == 0 && t&(1<<2) != 0
}

// This method does not check if move is a promotion
func (m Move) PromotedPiece() int {
	return Knight + int(m>>20)
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
	fmt.Println("")
}

func (m Move) String() string {
	if m == 0 {
		return "0000"
	}
	var promo = ""
	if m.IsPromotion() {
		promo = string("nbrq"[m.Special()])
	}
	return SquareString[m.From()] + SquareString[m.To()] + promo
}
