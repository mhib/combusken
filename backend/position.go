package backend

import (
	"fmt"
	"strings"
)

const (
	Pawn = iota
	Knight
	Bishop
	Rook
	Queen
	King
	None
)

const (
	Black = iota
	White
)

const (
	WhiteKingSideCastleFlag = 1 << iota
	WhiteQueenSideCastleFlag
	BlackKingSideCastleFlag
	BlackQueenSideCastleFlag
)

type Position struct {
	Colours    [White + 1]uint64
	Pieces     [King + 1]uint64
	Key        uint64
	PawnKey    uint64
	SideToMove int
	EpSquare   int
	FiftyMove  int
	LastMove   Move
	Flags      uint8
}

const maxMoves = 256

var InitialPosition Position = ParseFen(InitialPositionFen)

var rookCastleFlags [64]uint8

func init() {
	HashPosition(&InitialPosition)
	rookCastleFlags[A1] = WhiteQueenSideCastleFlag
	rookCastleFlags[H1] = WhiteKingSideCastleFlag
	rookCastleFlags[H8] = BlackKingSideCastleFlag
	rookCastleFlags[A8] = BlackQueenSideCastleFlag
}

func (pos *Position) TypeOnSquare(squareBB uint64) int {
	if squareBB&pos.Pieces[Pawn] != 0 {
		return Pawn
	} else if squareBB&pos.Pieces[Knight] != 0 {
		return Knight
	} else if squareBB&pos.Pieces[Bishop] != 0 {
		return Bishop
	} else if squareBB&pos.Pieces[Rook] != 0 {
		return Rook
	} else if squareBB&pos.Pieces[Queen] != 0 {
		return Queen
	} else if squareBB&pos.Pieces[King] != 0 {
		return King
	}
	return None
}

var kingCastleFlags = [2]uint8{BlackKingSideCastleFlag | BlackQueenSideCastleFlag, WhiteKingSideCastleFlag | WhiteQueenSideCastleFlag}

func (p *Position) MovePiece(piece, side, from, to int) {
	var b = SquareMask[from] ^ SquareMask[to]
	p.Colours[side] ^= b
	p.Pieces[piece] ^= b
	p.Key ^= zobrist[piece][side][from] ^ zobrist[piece][side][to]
	switch piece {
	case King:
		p.Flags |= kingCastleFlags[side]
		fallthrough
	case Pawn:
		p.PawnKey ^= zobrist[piece][side][from] ^ zobrist[piece][side][to]
	case Rook:
		p.Flags |= rookCastleFlags[from]
	}
}

func (p *Position) TogglePiece(piece, side, square int) {
	var b = SquareMask[square]
	p.Colours[side] ^= b
	p.Pieces[piece] ^= b
	switch piece {
	// Commented out as this function should not be called with King
	//case King:
	//fallthrough
	case Pawn:
		p.PawnKey ^= zobrist[Pawn][side][square]
	case Rook:
		p.Flags |= rookCastleFlags[square]
	}
}

func (pos *Position) MakeNullMove(res *Position) {
	copy(res.Colours[:], pos.Colours[:])
	copy(res.Pieces[:], pos.Pieces[:])
	res.SideToMove = pos.SideToMove ^ 1
	res.Flags = pos.Flags
	res.Key = pos.Key ^ zobristColor ^ zobristEpSquare[pos.EpSquare]
	res.PawnKey = pos.PawnKey ^ zobristColor

	res.FiftyMove = pos.FiftyMove + 1
	res.LastMove = NullMove
	res.EpSquare = 0
}

func (pos *Position) MakeMove(move Move, res *Position) bool {
	copy(res.Colours[:], pos.Colours[:])
	copy(res.Pieces[:], pos.Pieces[:])
	res.SideToMove = pos.SideToMove
	res.Flags = pos.Flags
	res.Key = pos.Key ^ zobristColor ^ zobristEpSquare[pos.EpSquare] ^ zobristFlags[pos.Flags]
	res.PawnKey = pos.PawnKey ^ zobristColor

	if move.MovedPiece() == Pawn || move.IsCapture() {
		res.FiftyMove = 0
	} else {
		res.FiftyMove = pos.FiftyMove + 1
	}

	res.EpSquare = 0

	if !move.IsPromotion() {
		res.MovePiece(move.MovedPiece(), pos.SideToMove, move.From(), move.To())
		switch move.Type() {
		case DoublePawnPush:
			res.EpSquare = move.To()
			res.Key ^= zobristEpSquare[move.To()]
		case Capture:
			res.TogglePiece(move.CapturedPiece(), pos.SideToMove^1, move.To())
		case KingCastle:
			if pos.SideToMove == White {
				res.MovePiece(Rook, White, H1, F1)
			} else {
				res.MovePiece(Rook, Black, H8, F8)
			}
		case QueenCastle:
			if pos.SideToMove == White {
				res.MovePiece(Rook, White, A1, D1)
			} else {
				res.MovePiece(Rook, Black, A8, D8)
			}
		case EPCapture:
			res.TogglePiece(Pawn, pos.SideToMove^1, pos.EpSquare)
		}
	} else {
		res.TogglePiece(Pawn, pos.SideToMove, move.From())
		res.TogglePiece(move.PromotedPiece(), pos.SideToMove, move.To())
		if move.IsCapture() {
			res.TogglePiece(move.CapturedPiece(), pos.SideToMove^1, move.To())
		}
	}

	// IsInCheck inlined
	// Replace when Go will be better at inlining
	if res.IsSquareAttacked(BitScan(res.Colours[res.SideToMove]&res.Pieces[King]), res.SideToMove^1) {
		return false
	}

	res.Key ^= zobristFlags[res.Flags]
	res.SideToMove = pos.SideToMove ^ 1
	res.LastMove = move
	return true
}

func (pos *Position) IsInCheck() bool {
	return pos.IsSquareAttacked(BitScan(pos.Colours[pos.SideToMove]&pos.Pieces[King]), pos.SideToMove^1)
}

func (pos *Position) IsSquareAttacked(square, side int) bool {
	theirOccupation := pos.Colours[side]
	return PawnAttacks[side^1][square]&pos.Pieces[Pawn]&theirOccupation != 0 ||
		KnightAttacks[square]&theirOccupation&pos.Pieces[Knight] != 0 ||
		KingAttacks[square]&pos.Pieces[King]&theirOccupation != 0 ||
		BishopAttacks(square, pos.Colours[Black]|pos.Colours[White])&(pos.Pieces[Bishop]|pos.Pieces[Queen])&theirOccupation != 0 ||
		RookAttacks(square, pos.Colours[Black]|pos.Colours[White])&(pos.Pieces[Queen]|pos.Pieces[Rook])&theirOccupation != 0
}

func (pos *Position) Print() {
	for y := 7; y >= 0; y-- {
		for x := 0; x <= 7; x++ {
			bb := uint64(1) << uint64(8*y+x)
			var char byte
			switch pos.TypeOnSquare(bb) {
			case Pawn:
				char = 'p'
			case Knight:
				char = 'n'
			case Bishop:
				char = 'b'
			case Rook:
				char = 'r'
			case Queen:
				char = 'q'
			case King:
				char = 'k'
			default:
				char = '.'
			}
			if pos.Colours[White]&bb != 0 {
				fmt.Print(strings.ToUpper(string(char)))
			} else {
				fmt.Print(string(char))
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func (p *Position) MakeMoveLAN(lan string) (Position, bool) {
	var buffer [256]EvaledMove
	var ml = p.GenerateAllMoves(buffer[:])
	for i := range ml {
		var mv = ml[i].Move
		if strings.EqualFold(mv.String(), lan) {
			var newPosition = Position{}
			if p.MakeMove(mv, &newPosition) {
				return newPosition, true
			} else {
				return Position{}, false
			}
		}
	}
	return Position{}, false
}

func (pos *Position) GenerateAllLegalMoves() []EvaledMove {
	var buffer [256]EvaledMove
	var moves = pos.GenerateAllMoves(buffer[:])
	var child Position
	result := make([]EvaledMove, 0)
	for _, move := range moves {
		if pos.MakeMove(move.Move, &child) {
			result = append(result, move)
		}
	}
	return result
}

func (pos *Position) MakeLegalMove(move Move, res *Position) {
	copy(res.Colours[:], pos.Colours[:])
	copy(res.Pieces[:], pos.Pieces[:])
	res.SideToMove = pos.SideToMove
	res.Flags = pos.Flags
	res.Key = pos.Key ^ zobristColor ^ zobristEpSquare[pos.EpSquare] ^ zobristFlags[pos.Flags]
	res.PawnKey = pos.PawnKey ^ zobristColor

	if move.MovedPiece() == Pawn || move.IsCapture() {
		res.FiftyMove = 0
	} else {
		res.FiftyMove = pos.FiftyMove + 1
	}

	res.EpSquare = 0

	if !move.IsPromotion() {
		res.MovePiece(move.MovedPiece(), pos.SideToMove, move.From(), move.To())
		switch move.Type() {
		case DoublePawnPush:
			res.EpSquare = move.To()
			res.Key ^= zobristEpSquare[move.To()]
		case Capture:
			res.TogglePiece(move.CapturedPiece(), pos.SideToMove^1, move.To())
			if move.CapturedPiece() == Rook {
				res.Flags |= rookCastleFlags[move.To()]
			}
		case KingCastle:
			if pos.SideToMove == White {
				res.MovePiece(Rook, White, H1, F1)
			} else {
				res.MovePiece(Rook, Black, H8, F8)
			}
		case QueenCastle:
			if pos.SideToMove == White {
				res.MovePiece(Rook, White, A1, D1)
			} else {
				res.MovePiece(Rook, Black, A8, D8)
			}
		case EPCapture:
			res.TogglePiece(Pawn, pos.SideToMove^1, pos.EpSquare)
		}
	} else {
		res.TogglePiece(Pawn, pos.SideToMove, move.From())
		res.TogglePiece(move.PromotedPiece(), pos.SideToMove, move.To())
		if move.IsCapture() {
			res.TogglePiece(move.CapturedPiece(), pos.SideToMove^1, move.To())
			if move.CapturedPiece() == Rook {
				res.Flags |= rookCastleFlags[move.To()]
			}
		}
	}

	res.Key ^= zobristFlags[res.Flags]
	res.SideToMove = pos.SideToMove ^ 1
	res.LastMove = move
}
