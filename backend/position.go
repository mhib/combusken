package backend

import (
	"fmt"
	"github.com/mhib/combusken/utils"
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
	WhiteKingSideCastleFlag = 1 << iota
	WhiteQueenSideCastleFlag
	BlackKingSideCastleFlag
	BlackQueenSideCastleFlag
)

const (
	Black = iota
	White
)

type Position struct {
	Squares    [64]Piece
	Colours    [2]uint64
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

var InitialPosition = ParseFen(InitialPositionFen)

var rookCastleFlags [64]uint8

func init() {
	HashPosition(&InitialPosition)
	rookCastleFlags[A1] = WhiteQueenSideCastleFlag
	rookCastleFlags[H1] = WhiteKingSideCastleFlag
	rookCastleFlags[H8] = BlackKingSideCastleFlag
	rookCastleFlags[A8] = BlackQueenSideCastleFlag
}

func (pos *Position) TypeOnSquare(square int) int {
	return pos.Squares[square].Type()
}

var kingCastlingFlags = [2]uint8{BlackKingSideCastleFlag | BlackQueenSideCastleFlag, WhiteKingSideCastleFlag | WhiteQueenSideCastleFlag}

func (p *Position) MovePiece(piece int, side int, from int, to int) {
	var b = SquareMask[from] ^ SquareMask[to]
	p.Colours[side] ^= b
	p.Pieces[piece] ^= b
	p.Squares[to] = NewPiece(piece, side)
	p.Squares[from] = NoPiece
	p.Key ^= zobrist[piece][side][from] ^ zobrist[piece][side][to]

	if piece == Pawn || piece == King {
		p.PawnKey ^= zobrist[piece][side][from] ^ zobrist[piece][side][to]
	}
	if piece == King {
		p.Flags |= kingCastlingFlags[side]
	}
	if piece == Rook {
		p.Flags |= rookCastleFlags[from]
	}
}

func (p *Position) MovePieceWithoutFlags(piece int, side int, from int, to int) {
	var b = SquareMask[from] ^ SquareMask[to]
	p.Colours[side] ^= b
	p.Pieces[piece] ^= b
	p.Squares[to] = NewPiece(piece, side)
	p.Squares[from] = NoPiece
}

func (p *Position) RemovePiece(piece int, side int, square int) {
	var b = SquareMask[square]
	p.Colours[side] ^= b
	p.Pieces[piece] ^= b
	p.Squares[square] = NoPiece
	if piece == Pawn || piece == King {
		p.PawnKey ^= zobrist[piece][side][square]
	}
	if piece == King {
		p.Flags |= kingCastlingFlags[side]
	}
	if piece == Rook {
		p.Flags |= rookCastleFlags[square]
	}
}

func (p *Position) RemovePieceWithoutFlags(piece int, side int, square int) {
	var b = SquareMask[square]
	p.Colours[side] ^= b
	p.Squares[square] = NoPiece
	p.Pieces[piece] ^= b
}

func (p *Position) SetPiece(piece int, side int, square int) {
	var b = SquareMask[square]
	p.Key ^= zobrist[piece][side][square]
	p.Squares[square] = NewPiece(piece, side)
	p.Colours[side] ^= b
	p.Pieces[piece] ^= b
}

func (p *Position) SetPieceWithoutFlags(piece int, side int, square int) {
	var b = SquareMask[square]
	p.Squares[square] = NewPiece(piece, side)
	p.Colours[side] ^= b
	p.Pieces[piece] ^= b
}

func (pos *Position) MakeNullMove(undo *Undo) {
	undo.Key = pos.Key
	undo.PawnKey = pos.PawnKey
	undo.FiftyMove = pos.FiftyMove
	undo.LastMove = pos.LastMove
	undo.EpSquare = pos.EpSquare

	pos.Key ^= zobristColor ^ zobristEpSquare[pos.EpSquare]
	pos.PawnKey ^= zobristColor

	pos.FiftyMove++
	pos.LastMove = NullMove
	pos.EpSquare = 0
}

func (pos *Position) MakeMove(move Move, undo *Undo) bool {
	undo.Key = pos.Key
	undo.PawnKey = pos.PawnKey
	undo.FiftyMove = pos.FiftyMove
	undo.EpSquare = pos.EpSquare
	undo.CapturedPiece = NoPiece
	undo.LastMove = pos.LastMove
	undo.Flags = pos.Flags

	pos.Key ^= zobristEpSquare[pos.EpSquare] ^ zobristFlags[pos.Flags] ^ zobristColor
	pos.PawnKey ^= zobristColor
	pos.FiftyMove = pos.FiftyMove + 1

	movedPiece := pos.TypeOnSquare(move.From())

	pos.EpSquare = 0

	switch move.Type() {
	case NormalMove:
		if move.Special() == CaptureMove {
			pos.FiftyMove = 0
			capturedPiece := pos.TypeOnSquare(move.To())
			undo.CapturedPiece = pos.Squares[move.To()]
			pos.RemovePiece(capturedPiece, pos.SideToMove^1, move.To())
		} else if movedPiece == Pawn {
			pos.FiftyMove = 0
			if utils.Abs(int64(move.From()-move.To())) == 16 {
				pos.EpSquare = move.To()
				pos.Key ^= zobristEpSquare[move.To()]
			}
		}
		pos.MovePiece(movedPiece, pos.SideToMove, move.From(), move.To())
	case CastleMove:
		pos.MovePiece(King, pos.SideToMove, move.From(), move.To())
		switch move {
		case WhiteKingSideCastle:
			pos.MovePiece(Rook, 1, H1, F1)
		case WhiteQueenSideCastle:
			pos.MovePiece(Rook, 1, A1, D1)
		case BlackKingSideCastle:
			pos.MovePiece(Rook, 0, H8, F8)
		case BlackQueenSideCastle:
			pos.MovePiece(Rook, 0, A8, D8)
		}
	case EnpassMove:
		pos.FiftyMove = 0
		pos.RemovePiece(Pawn, pos.SideToMove^1, undo.EpSquare)
		pos.MovePiece(Pawn, pos.SideToMove, move.From(), move.To())
	case PromotionMove:
		pos.FiftyMove = 0
		pos.RemovePiece(Pawn, pos.SideToMove, move.From())
		capturedPiece := pos.TypeOnSquare(move.To())
		if capturedPiece != None {
			undo.CapturedPiece = pos.Squares[move.To()]
			pos.RemovePiece(capturedPiece, pos.SideToMove^1, move.To())
		}
		pos.SetPiece(move.PromotedPiece(), pos.SideToMove, move.To())
	}

	if pos.IsInCheck() {
		pos.SideToMove ^= 1
		return false
	}

	pos.SideToMove ^= 1
	pos.Key ^= zobristFlags[pos.Flags]
	pos.LastMove = move
	return true
}

func (pos *Position) Undo(move Move, undo *Undo) {
	if move == NullMove {
		pos.Key = undo.Key
		pos.PawnKey = undo.PawnKey
		pos.FiftyMove = undo.FiftyMove
		pos.LastMove = undo.LastMove
		pos.EpSquare = undo.EpSquare
		pos.SideToMove ^= 1
		return
	}

	movedPiece := pos.TypeOnSquare(move.To())

	switch move.Type() {
	case NormalMove:
		pos.MovePieceWithoutFlags(movedPiece, pos.SideToMove^1, move.To(), move.From())
		if move.Special() == CaptureMove {
			pos.SetPieceWithoutFlags(undo.CapturedPiece.Type(), pos.SideToMove, move.To())
		}
	case CastleMove:
		pos.MovePieceWithoutFlags(King, pos.SideToMove^1, move.To(), move.From())
		switch move {
		case WhiteKingSideCastle:
			pos.MovePieceWithoutFlags(Rook, pos.SideToMove^1, F1, H1)
		case WhiteQueenSideCastle:
			pos.MovePieceWithoutFlags(Rook, pos.SideToMove^1, D1, A1)
		case BlackKingSideCastle:
			pos.MovePieceWithoutFlags(Rook, pos.SideToMove^1, F8, H8)
		case BlackQueenSideCastle:
			pos.MovePieceWithoutFlags(Rook, pos.SideToMove^1, D8, A8)
		}
	case EnpassMove:
		pos.MovePieceWithoutFlags(Pawn, pos.SideToMove^1, move.To(), move.From())
		pos.SetPieceWithoutFlags(Pawn, pos.SideToMove, undo.EpSquare)
	case PromotionMove:
		pos.RemovePieceWithoutFlags(move.PromotedPiece(), pos.SideToMove^1, move.To())
		if undo.CapturedPiece != NoPiece {
			pos.SetPieceWithoutFlags(undo.CapturedPiece.Type(), pos.SideToMove, move.To())
		}
		pos.SetPieceWithoutFlags(Pawn, pos.SideToMove^1, move.From())
	}

	pos.Key = undo.Key
	pos.PawnKey = undo.PawnKey
	pos.FiftyMove = undo.FiftyMove
	pos.EpSquare = undo.EpSquare
	pos.LastMove = undo.LastMove
	pos.Flags = undo.Flags
	pos.SideToMove ^= 1
}

func (pos *Position) IsInCheck() bool {
	return pos.IsSquareAttacked(BitScan(pos.Colours[pos.SideToMove]&pos.Pieces[King]), pos.SideToMove^1)
}

func (pos *Position) IsSquareAttacked(square int, side int) bool {
	return PawnAttacks[side^1][square]&pos.Pieces[Pawn]&pos.Colours[side] != 0 ||
		KnightAttacks[square]&pos.Colours[side]&pos.Pieces[Knight] != 0 ||
		KingAttacks[square]&pos.Colours[side]&pos.Pieces[King] != 0 ||
		BishopAttacks(square, pos.Colours[0]|pos.Colours[1])&(pos.Pieces[Queen]|pos.Pieces[Bishop])&pos.Colours[side] != 0 ||
		RookAttacks(square, pos.Colours[0]|pos.Colours[1])&(pos.Pieces[Queen]|pos.Pieces[Rook])&pos.Colours[side] != 0
}

func (pos *Position) Print() {
	for y := 7; y >= 0; y-- {
		for x := 0; x <= 7; x++ {
			bb := uint64(1) << uint64(8*y+x)
			var char byte
			switch pos.TypeOnSquare(BitScan(bb)) {
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
	var undo Undo
	for i := range ml {
		var mv = ml[i].Move
		if strings.EqualFold(mv.String(), lan) {
			var newPosition = Position{}
			if p.MakeMove(mv, &undo) {
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
	var undo Undo
	result := make([]EvaledMove, 0)
	for _, move := range moves {
		if pos.MakeMove(move.Move, &undo) {
			result = append(result, move)
		}
		pos.Undo(move.Move, &undo)
	}
	return result
}

func (pos *Position) MakeLegalMove(move Move, undo *Undo) {
	undo.Key = pos.Key
	undo.PawnKey = pos.PawnKey
	undo.FiftyMove = pos.FiftyMove
	undo.EpSquare = pos.EpSquare
	undo.CapturedPiece = NoPiece
	undo.LastMove = pos.LastMove
	undo.Flags = pos.Flags

	pos.Key ^= zobristEpSquare[pos.EpSquare] ^ zobristFlags[pos.Flags] ^ zobristColor
	pos.PawnKey ^= zobristColor
	pos.FiftyMove = pos.FiftyMove + 1

	movedPiece := pos.TypeOnSquare(move.From())

	pos.EpSquare = 0

	switch move.Type() {
	case NormalMove:
		if move.Special() == CaptureMove {
			pos.FiftyMove = 0
			capturedPiece := pos.TypeOnSquare(move.To())
			undo.CapturedPiece = pos.Squares[move.To()]
			pos.RemovePiece(capturedPiece, pos.SideToMove^1, move.To())
		} else if movedPiece == Pawn {
			pos.FiftyMove = 0
			if utils.Abs(int64(move.From()-move.To())) == 16 {
				pos.EpSquare = move.To()
				pos.Key ^= zobristEpSquare[move.To()]
			}
		}
		pos.MovePiece(movedPiece, pos.SideToMove, move.From(), move.To())
	case CastleMove:
		pos.MovePiece(King, pos.SideToMove, move.From(), move.To())
		switch move {
		case WhiteKingSideCastle:
			pos.MovePiece(Rook, 1, H1, F1)
		case WhiteQueenSideCastle:
			pos.MovePiece(Rook, 1, A1, D1)
		case BlackKingSideCastle:
			pos.MovePiece(Rook, 0, H8, F8)
		case BlackQueenSideCastle:
			pos.MovePiece(Rook, 0, A8, D8)
		}
	case EnpassMove:
		pos.FiftyMove = 0
		pos.RemovePiece(Pawn, ^pos.SideToMove, undo.EpSquare)
		pos.MovePiece(Pawn, pos.SideToMove, move.From(), move.To())
	case PromotionMove:
		pos.FiftyMove = 0
		pos.RemovePiece(Pawn, pos.SideToMove, move.From())
		capturedPiece := pos.TypeOnSquare(move.To())
		if capturedPiece != None {
			undo.CapturedPiece = pos.Squares[move.To()]
			pos.RemovePiece(capturedPiece, pos.SideToMove^1, move.To())
		}
		pos.SetPiece(move.PromotedPiece(), pos.SideToMove, move.To())
	}

	pos.SideToMove ^= 1
	pos.Key ^= zobristFlags[pos.Flags]
	pos.LastMove = move
}

func (pos *Position) allOccupation() uint64 {
	return pos.Colours[0] | pos.Colours[1]
}
