package chess

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
	ColourNone
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
	Checkers   uint64
	SideToMove int
	EpSquare   int
	FiftyMove  int
	LastMove   Move
	Flags      uint8
}

var InitialPosition Position = ParseFen(InitialPositionFen)

var castleFlags [64]uint8

func init() {
	castleFlags[A1] = WhiteQueenSideCastleFlag
	castleFlags[H1] = WhiteKingSideCastleFlag
	castleFlags[E1] = WhiteQueenSideCastleFlag | WhiteKingSideCastleFlag
	castleFlags[H8] = BlackKingSideCastleFlag
	castleFlags[A8] = BlackQueenSideCastleFlag
	castleFlags[E8] = BlackQueenSideCastleFlag | BlackKingSideCastleFlag
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

func (p *Position) MovePiece(piece, side, from, to int) {
	var b = SquareMask[from] ^ SquareMask[to]
	p.Colours[side] ^= b
	p.Pieces[piece] ^= b
	p.Key ^= zobrist[piece][side][from] ^ zobrist[piece][side][to]
	p.Flags |= castleFlags[from] | castleFlags[to]
	if piece == King || piece == Pawn {
		p.PawnKey ^= zobrist[piece][side][from] ^ zobrist[piece][side][to]
	}
}

func (p *Position) TogglePiece(piece, side, square int) {
	var b = SquareMask[square]
	p.Colours[side] ^= b
	p.Pieces[piece] ^= b
	p.Key ^= zobrist[piece][side][square]
	p.Flags |= castleFlags[square]
	if piece == Pawn {
		p.PawnKey ^= zobrist[Pawn][side][square]
	}
}

func (pos *Position) MakeNullMove(res *Position) {
	res.Colours[Black] = pos.Colours[Black]
	res.Colours[White] = pos.Colours[White]
	res.Pieces[Pawn] = pos.Pieces[Pawn]
	res.Pieces[Knight] = pos.Pieces[Knight]
	res.Pieces[Bishop] = pos.Pieces[Bishop]
	res.Pieces[Rook] = pos.Pieces[Rook]
	res.Pieces[Queen] = pos.Pieces[Queen]
	res.Pieces[King] = pos.Pieces[King]
	res.SideToMove = pos.SideToMove ^ 1
	res.Flags = pos.Flags
	res.Checkers = 0
	res.Key = pos.Key ^ zobristColor ^ zobristEpSquare[pos.EpSquare]
	res.PawnKey = pos.PawnKey

	res.FiftyMove = pos.FiftyMove + 1
	res.LastMove = NullMove
	res.EpSquare = 0
}

func (pos *Position) MakeMove(move Move, res *Position) bool {
	res.Colours[Black] = pos.Colours[Black]
	res.Colours[White] = pos.Colours[White]
	res.Pieces[Pawn] = pos.Pieces[Pawn]
	res.Pieces[Knight] = pos.Pieces[Knight]
	res.Pieces[Bishop] = pos.Pieces[Bishop]
	res.Pieces[Rook] = pos.Pieces[Rook]
	res.Pieces[Queen] = pos.Pieces[Queen]
	res.Pieces[King] = pos.Pieces[King]
	res.SideToMove = pos.SideToMove
	res.Flags = pos.Flags
	res.Key = pos.Key ^ zobristColor ^ zobristEpSquare[pos.EpSquare] ^ zobristFlags[pos.Flags]
	res.PawnKey = pos.PawnKey

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
	res.Checkers = res.AllSquareAttackers(BitScan(res.Colours[res.SideToMove]&res.Pieces[King]), res.SideToMove^1)
	res.LastMove = move
	return true
}

func (pos *Position) IsInCheck() bool {
	return pos.Checkers != 0
}

func (pos *Position) AllSquareAttackers(square, side int) uint64 {
	theirOccupation := pos.Colours[side]
	return (PawnAttacks[side^1][square] & pos.Pieces[Pawn] & theirOccupation) |
		(KnightAttacks[square] & theirOccupation & pos.Pieces[Knight]) |
		(KingAttacks[square] & pos.Pieces[King] & theirOccupation) |
		(BishopAttacks(square, pos.Colours[Black]|pos.Colours[White]) & (pos.Pieces[Bishop] | pos.Pieces[Queen]) & theirOccupation) |
		(RookAttacks(square, pos.Colours[Black]|pos.Colours[White]) & (pos.Pieces[Queen] | pos.Pieces[Rook]) & theirOccupation)
}

func (pos *Position) IsSquareAttacked(square, side int) bool {
	theirOccupation := pos.Colours[side]
	theirBishops := (pos.Pieces[Bishop] | pos.Pieces[Queen]) & theirOccupation
	theirRooks := (pos.Pieces[Rook] | pos.Pieces[Queen]) & theirOccupation
	return PawnAttacks[side^1][square]&pos.Pieces[Pawn]&theirOccupation != 0 ||
		KnightAttacks[square]&theirOccupation&pos.Pieces[Knight] != 0 ||
		KingAttacks[square]&pos.Pieces[King]&theirOccupation != 0 ||
		(theirBishops != 0 && BishopAttacks(square, pos.Colours[Black]|pos.Colours[White])&theirBishops != 0) ||
		(theirRooks != 0 && RookAttacks(square, pos.Colours[Black]|pos.Colours[White])&theirRooks != 0)
}

func (pos *Position) Print() {
	pieceChar := "pnbrqk."
	for y := 7; y >= 0; y-- {
		for x := 0; x <= 7; x++ {
			bb := SquareMask[8*y+x]
			char := pieceChar[pos.TypeOnSquare(bb)]
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
	noisySize := GenerateNoisy(p, buffer[:])
	quietsSize := GenerateQuiet(p, buffer[noisySize:])
	for i := range buffer[:noisySize+quietsSize] {
		var mv = buffer[i].Move
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

func (pos *Position) MakeLegalMove(move Move, res *Position) {
	res.Colours[Black] = pos.Colours[Black]
	res.Colours[White] = pos.Colours[White]
	res.Pieces[Pawn] = pos.Pieces[Pawn]
	res.Pieces[Knight] = pos.Pieces[Knight]
	res.Pieces[Bishop] = pos.Pieces[Bishop]
	res.Pieces[Rook] = pos.Pieces[Rook]
	res.Pieces[Queen] = pos.Pieces[Queen]
	res.Pieces[King] = pos.Pieces[King]
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

	res.Key ^= zobristFlags[res.Flags]
	res.SideToMove = pos.SideToMove ^ 1
	res.Checkers = res.AllSquareAttackers(BitScan(res.Colours[res.SideToMove]&res.Pieces[King]), res.SideToMove^1)
	res.LastMove = move
}

func (pos *Position) IsMovePseudoLegal(move Move) bool {
	we := pos.Colours[pos.SideToMove]
	them := pos.Colours[pos.SideToMove^1]
	occupancy := we | them
	fromMask := SquareMask[move.From()]
	toMask := SquareMask[move.To()]

	if move == NullMove || (we&fromMask) == 0 || (we&toMask) != 0 ||
		(move.IsCapture() && (move.CapturedPiece() >= King || (move.Type() != EPCapture && pos.Pieces[move.CapturedPiece()]&them&toMask == 0))) ||
		(!move.IsCapture() && toMask&them != 0) {
		return false
	}

	switch move.MovedPiece() {
	case Knight:
		return move.IsNormal() && (pos.Pieces[Knight]&fromMask) != 0 && KnightAttacks[move.From()]&toMask != 0
	case Bishop:
		return move.IsNormal() && (pos.Pieces[Bishop]&fromMask) != 0 && BishopAttacks(move.From(), occupancy)&toMask != 0
	case Rook:
		return move.IsNormal() && (pos.Pieces[Rook]&fromMask) != 0 && RookAttacks(move.From(), occupancy)&toMask != 0
	case Queen:
		return move.IsNormal() && (pos.Pieces[Queen]&fromMask) != 0 && QueenAttacks(move.From(), occupancy)&toMask != 0
	case Pawn:
		if (pos.Pieces[Pawn] & fromMask) == 0 {
			return false
		}
		var attacks, forward uint64
		if pos.SideToMove == White {
			if move.Type() == EPCapture {
				return pos.EpSquare != 0 && (SquareMask[uint(pos.EpSquare)-1]|SquareMask[uint(pos.EpSquare)+1])&Rank5_BB&fromMask != 0
			}
			attacks = PawnAttacks[White][move.From()]
			forward = North(fromMask) & ^occupancy
		} else {
			if move.Type() == EPCapture {
				return pos.EpSquare != 0 && (SquareMask[uint(pos.EpSquare)-1]|SquareMask[uint(pos.EpSquare)+1])&Rank4_BB&fromMask != 0
			}
			attacks = PawnAttacks[Black][move.From()]
			forward = South(fromMask) & ^occupancy
		}
		if move.IsPromotion() {
			return PromotionRanks_BB&((attacks&them)|forward) != 0 && move.PromotedPiece() <= Queen
		}

		// Invalid move type as promotions and EPCapture were checked
		if move.IsCapture() && move.Type() != Capture {
			return false
		}

		// Double pawn push
		if forward != 0 && pos.SideToMove == White && fromMask&Rank2_BB != 0 && North(forward)&occupancy == 0 {
			forward |= North(forward)
		} else if forward != 0 && pos.SideToMove == Black && fromMask&Rank7_BB != 0 && South(forward)&occupancy == 0 {
			forward |= South(forward)
		}
		return (^PromotionRanks_BB)&((attacks&them)|forward)&toMask != 0
	case King:
		if (pos.Pieces[King] & fromMask) == 0 {
			return false
		}
		if move.IsNormal() {
			return KingAttacks[move.From()]&toMask != 0
		}
		if pos.SideToMove == White {
			if move == WhiteKingSideCastle {
				return occupancy&WhiteKingCastleBlock_BB == 0 && pos.Flags&WhiteKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E1, Black) && !pos.IsSquareAttacked(F1, Black)
			} else if move == WhiteQueenSideCastle {
				return occupancy&WhiteQueenCastleBlock_BB == 0 && pos.Flags&WhiteQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E1, Black) && !pos.IsSquareAttacked(D1, Black)
			} else {
				return false
			}
		} else {
			if move == BlackKingSideCastle {
				return occupancy&BlackKingCastleBlock_BB == 0 && pos.Flags&BlackKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E8, White) && !pos.IsSquareAttacked(F8, White)
			} else if move == BlackQueenSideCastleFlag {
				return occupancy&BlackQueenCastleBlock_BB == 0 && pos.Flags&BlackQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E8, White) && !pos.IsSquareAttacked(D8, White)
			} else {
				return false
			}
		}
	}

	return false
}
