package backend

import "fmt"

// Move is represented as int32\
// Bits are:
// 2 bits special
// 2 bits type
// 6 bits to
// 6 bits from

const (
	NormalMove = iota
	CastleMove
	EnpassMove
	PromotionMove
)

const (
	PromoteToKnight = iota
	PromoteToBishop
	PromoteToRook
	PromoteToQueen
)

// Does not apply to promotions
const (
	QuietMove   = 0
	CaptureMove = 1
)

type Move int32

const (
	NullMove = Move(0)
)

var WhiteKingSideCastle = NewMove(E1, G1, CastleMove, QuietMove)
var WhiteQueenSideCastle = NewMove(E1, C1, CastleMove, QuietMove)
var BlackKingSideCastle = NewMove(E8, G8, CastleMove, QuietMove)
var BlackQueenSideCastle = NewMove(E8, C8, CastleMove, QuietMove)

func (m Move) From() int {
	return int(m & 0x3f)
}

func (m Move) To() int {
	return int((m >> 6) & 0x3f)
}

func (m Move) Type() int {
	return int(m>>12) & 3
}

func (m Move) Special() int {
	return int(m >> 14)
}

func (m Move) IsPromotion() bool {
	return (m>>12)&3 == PromotionMove
}

// A bit tricky
// If move capture special is 1
// And type is either 0, 2 or 3
// In case of promotion type is always 3
// This checks that either type is EnpassMove or PromotionMove or super is not empty
func (m Move) IsCaptureOrPromotion() bool {
	return (m&(1<<13))|(m>>14) != 0
}

func (m Move) IsCastling() bool {
	return (m>>12)&3 == CastleMove
}

// This method does not check if move is a promotion
func (m Move) PromotedPiece() int {
	return Knight + int(m>>14)
}

func NewMove(from, to, moveType, special int) Move {
	return Move(from | (to << 6) | (moveType << 12) | (special << 14))
}

func (m Move) Inspect() {
	fmt.Print("From: ")
	fmt.Print(m.From())
	fmt.Print(" To: ")
	fmt.Print(m.To())
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
