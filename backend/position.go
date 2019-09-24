package backend

import (
	"fmt"
	"github.com/mhib/combusken/utils"
	"strconv"
	"strings"
)

const (
	None = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

const (
	WhiteKingSideCastleFlag = 1 << iota
	WhiteQueenSideCastleFlag
	BlackKingSideCastleFlag
	BlackQueenSideCastleFlag
)

type Position struct {
	Pieces                                                      [64]Piece
	Pawns, Knights, Bishops, Rooks, Queens, Kings, White, Black uint64
	Key                                                         uint64
	PawnKey                                                     uint64
	EpSquare                                                    int
	FiftyMove                                                   int
	LastMove                                                    Move
	WhiteMove                                                   bool
	Flags                                                       uint8
}

func (pos *Position) Inspect() string {
	var sb strings.Builder
	sb.WriteString(strconv.FormatUint(pos.Pawns, 16))
	sb.WriteString("-")
	sb.WriteString(strconv.FormatUint(pos.Knights, 16))
	sb.WriteString("-")
	sb.WriteString(strconv.FormatUint(pos.Bishops, 16))
	sb.WriteString("-")
	sb.WriteString(strconv.FormatUint(pos.Rooks, 16))
	sb.WriteString("-")
	sb.WriteString(strconv.FormatUint(pos.Queens, 16))
	sb.WriteString("-")
	sb.WriteString(strconv.FormatUint(pos.Kings, 16))
	sb.WriteString("-")
	sb.WriteString(strconv.FormatUint(pos.White, 16))
	sb.WriteString("-")
	sb.WriteString(strconv.FormatUint(pos.Black, 16))
	return sb.String()
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
	return pos.Pieces[square].Type()
}

var kingCastlingFlags = [2]uint8{BlackKingSideCastleFlag | BlackQueenSideCastleFlag, WhiteKingSideCastleFlag | WhiteQueenSideCastleFlag}

func (p *Position) MovePiece(piece int, side bool, from int, to int) {
	var b = SquareMask[from] ^ SquareMask[to]
	var intSide = 0
	if side {
		p.White ^= b
		intSide = 1
	} else {
		p.Black ^= b
	}
	p.Pieces[to] = NewPiece(piece, intSide)
	p.Pieces[from] = NoPiece
	switch piece {
	case Pawn:
		p.Pawns ^= b
		p.Key ^= zobrist[0][intSide][from] ^ zobrist[0][intSide][to]
		p.PawnKey ^= zobrist[0][intSide][from] ^ zobrist[0][intSide][to]
	case Knight:
		p.Knights ^= b
		p.Key ^= zobrist[1][intSide][from] ^ zobrist[1][intSide][to]
	case Bishop:
		p.Bishops ^= b
		p.Key ^= zobrist[2][intSide][from] ^ zobrist[2][intSide][to]
	case Rook:
		p.Rooks ^= b
		p.Key ^= zobrist[3][intSide][from] ^ zobrist[3][intSide][to]
		p.Flags |= rookCastleFlags[from]
	case Queen:
		p.Queens ^= b
		p.Key ^= zobrist[4][intSide][from] ^ zobrist[4][intSide][to]
	case King:
		p.Kings ^= b
		p.Key ^= zobrist[5][intSide][from] ^ zobrist[5][intSide][to]
		p.PawnKey ^= zobrist[5][intSide][from] ^ zobrist[5][intSide][to]
		p.Flags |= kingCastlingFlags[intSide]
	}
}

func (p *Position) MovePieceWithoutFlags(piece int, side bool, from int, to int) {
	var b = SquareMask[from] ^ SquareMask[to]
	var intSide = 0
	if side {
		p.White ^= b
		intSide = 1
	} else {
		p.Black ^= b
	}
	p.Pieces[to] = NewPiece(piece, intSide)
	p.Pieces[from] = NoPiece
	switch piece {
	case Pawn:
		p.Pawns ^= b
	case Knight:
		p.Knights ^= b
	case Bishop:
		p.Bishops ^= b
	case Rook:
		p.Rooks ^= b
	case Queen:
		p.Queens ^= b
	case King:
		p.Kings ^= b
	}
}

func (p *Position) RemovePiece(piece int, side bool, square int) {
	var b = SquareMask[square]
	var intSide = 0
	if side {
		p.White ^= b
		intSide = 1
	} else {
		p.Black ^= b
	}
	p.Pieces[square] = NoPiece
	switch piece {
	case Pawn:
		p.Pawns ^= b
		p.Key ^= zobrist[0][intSide][square]
		p.PawnKey ^= zobrist[0][intSide][square]
	case Knight:
		p.Knights ^= b
		p.Key ^= zobrist[1][intSide][square]
	case Bishop:
		p.Bishops ^= b
		p.Key ^= zobrist[2][intSide][square]
	case Rook:
		p.Rooks ^= b
		p.Key ^= zobrist[3][intSide][square]
		p.Flags |= rookCastleFlags[square]
	case Queen:
		p.Queens ^= b
		p.Key ^= zobrist[4][intSide][square]
	case King:
		p.Kings ^= b
		p.Key ^= zobrist[5][intSide][square]
		p.PawnKey ^= zobrist[5][intSide][square]
	}
}

func (p *Position) RemovePieceWithoutFlags(piece int, side bool, square int) {
	var b = SquareMask[square]
	if side {
		p.White ^= b
	} else {
		p.Black ^= b
	}
	p.Pieces[square] = NoPiece
	switch piece {
	case Pawn:
		p.Pawns ^= b
	case Knight:
		p.Knights ^= b
	case Bishop:
		p.Bishops ^= b
	case Rook:
		p.Rooks ^= b
	case Queen:
		p.Queens ^= b
	case King:
		p.Kings ^= b
	}
}

func (p *Position) SetPiece(piece int, side bool, square int) {
	var b = SquareMask[square]
	var intSide = 0
	if side {
		p.White ^= b
		intSide = 1
	} else {
		p.Black ^= b
	}
	p.Pieces[square] = NewPiece(piece, intSide)
	switch piece {
	case Pawn:
		p.Pawns ^= b
		p.Key ^= zobrist[0][intSide][square]
		p.PawnKey ^= zobrist[0][intSide][square]
	case Knight:
		p.Knights ^= b
		p.Key ^= zobrist[1][intSide][square]
	case Bishop:
		p.Bishops ^= b
		p.Key ^= zobrist[2][intSide][square]
	case Rook:
		p.Rooks ^= b
		p.Key ^= zobrist[3][intSide][square]
		p.Flags |= rookCastleFlags[square]
	case Queen:
		p.Queens ^= b
		p.Key ^= zobrist[4][intSide][square]
	case King:
		p.Kings ^= b
		p.Key ^= zobrist[5][intSide][square]
		p.PawnKey ^= zobrist[5][intSide][square]
	}
}

func (p *Position) SetPieceWithoutFlags(piece int, side bool, square int) {
	var b = SquareMask[square]
	var intSide = 0
	if side {
		p.White ^= b
		intSide = 1
	} else {
		p.Black ^= b
	}
	p.Pieces[square] = NewPiece(piece, intSide)
	switch piece {
	case Pawn:
		p.Pawns ^= b
	case Knight:
		p.Knights ^= b
	case Bishop:
		p.Bishops ^= b
	case Rook:
		p.Rooks ^= b
	case Queen:
		p.Queens ^= b
	case King:
		p.Kings ^= b
	}
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
			undo.CapturedPiece = pos.Pieces[move.To()]
			pos.RemovePiece(capturedPiece, !pos.WhiteMove, move.To())
		} else if movedPiece == Pawn {
			pos.FiftyMove = 0
			if utils.Abs(int64(move.From()-move.To())) == 16 {
				pos.EpSquare = move.To()
				pos.Key ^= zobristEpSquare[move.To()]
			}
		}
		pos.MovePiece(movedPiece, pos.WhiteMove, move.From(), move.To())
	case CastleMove:
		pos.MovePiece(King, pos.WhiteMove, move.From(), move.To())
		switch move {
		case WhiteKingSideCastle:
			pos.MovePiece(Rook, true, H1, F1)
		case WhiteQueenSideCastle:
			pos.MovePiece(Rook, true, A1, D1)
		case BlackKingSideCastle:
			pos.MovePiece(Rook, false, H8, F8)
		case BlackQueenSideCastle:
			pos.MovePiece(Rook, false, A8, D8)
		}
	case EnpassMove:
		pos.FiftyMove = 0
		pos.RemovePiece(Pawn, !pos.WhiteMove, undo.EpSquare)
		pos.MovePiece(Pawn, pos.WhiteMove, move.From(), move.To())
	case PromotionMove:
		pos.FiftyMove = 0
		pos.RemovePiece(Pawn, pos.WhiteMove, move.From())
		capturedPiece := pos.TypeOnSquare(move.To())
		if capturedPiece != None {
			undo.CapturedPiece = pos.Pieces[move.To()]
			pos.RemovePiece(capturedPiece, !pos.WhiteMove, move.To())
		}
		pos.SetPiece(move.PromotedPiece(), pos.WhiteMove, move.To())
	}

	if pos.IsInCheck() {
		pos.WhiteMove = !pos.WhiteMove
		return false
	}

	pos.WhiteMove = !pos.WhiteMove
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
		pos.WhiteMove = !pos.WhiteMove
		return
	}

	movedPiece := pos.TypeOnSquare(move.To())

	switch move.Type() {
	case NormalMove:
		pos.MovePieceWithoutFlags(movedPiece, !pos.WhiteMove, move.To(), move.From())
		if move.Special() == CaptureMove {
			if undo.CapturedPiece == NoPiece {
				fmt.Println(move == pos.LastMove)
			}
			pos.SetPieceWithoutFlags(undo.CapturedPiece.Type(), pos.WhiteMove, move.To())
		}
	case CastleMove:
		switch move {
		case WhiteKingSideCastle:
			pos.White ^= (1 << G1) ^ (1 << E1)
			pos.Kings ^= (1 << G1) ^ (1 << E1)
			pos.Pieces[G1] = NoPiece
			pos.Pieces[E1] = NewPiece(King, 1)
			pos.White ^= (1 << H1) ^ (1 << F1)
			pos.Rooks ^= (1 << H1) ^ (1 << F1)
			pos.Pieces[F1] = NoPiece
			pos.Pieces[H1] = NewPiece(Rook, 1)
		case WhiteQueenSideCastle:
			pos.White ^= (1 << C1) ^ (1 << E1)
			pos.Kings ^= (1 << C1) ^ (1 << E1)
			pos.Pieces[C1] = NoPiece
			pos.Pieces[E1] = NewPiece(King, 1)

			pos.White ^= (1 << A1) ^ (1 << D1)
			pos.Rooks ^= (1 << A1) ^ (1 << D1)
			pos.Pieces[D1] = NoPiece
			pos.Pieces[A1] = NewPiece(Rook, 1)
		case BlackKingSideCastle:
			pos.Black ^= (1 << G8) ^ (1 << E8)
			pos.Kings ^= (1 << G8) ^ (1 << E8)
			pos.Pieces[G8] = NoPiece
			pos.Pieces[E8] = NewPiece(King, 0)

			pos.Black ^= (1 << H8) ^ (1 << F8)
			pos.Rooks ^= (1 << H8) ^ (1 << F8)
			pos.Pieces[F8] = NoPiece
			pos.Pieces[H8] = NewPiece(Rook, 0)
		case BlackQueenSideCastle:
			pos.Black ^= (1 << C8) ^ (1 << E8)
			pos.Kings ^= (1 << C8) ^ (1 << E8)
			pos.Pieces[C8] = NoPiece
			pos.Pieces[E8] = NewPiece(King, 0)

			pos.Black ^= (1 << A8) ^ (1 << D8)
			pos.Rooks ^= (1 << A8) ^ (1 << D8)
			pos.Pieces[D8] = NoPiece
			pos.Pieces[A8] = NewPiece(Rook, 0)
		}
	case EnpassMove:
		pos.MovePieceWithoutFlags(Pawn, !pos.WhiteMove, move.To(), move.From())
		pos.SetPieceWithoutFlags(Pawn, pos.WhiteMove, undo.EpSquare)
	case PromotionMove:
		pos.RemovePieceWithoutFlags(move.PromotedPiece(), !pos.WhiteMove, move.To())
		if undo.CapturedPiece != NoPiece {
			pos.SetPieceWithoutFlags(undo.CapturedPiece.Type(), pos.WhiteMove, move.To())
		}
		pos.SetPieceWithoutFlags(Pawn, !pos.WhiteMove, move.From())
	}

	pos.Key = undo.Key
	pos.PawnKey = undo.PawnKey
	pos.FiftyMove = undo.FiftyMove
	pos.EpSquare = undo.EpSquare
	pos.LastMove = undo.LastMove
	pos.Flags = undo.Flags
	pos.WhiteMove = !pos.WhiteMove
}

func (pos *Position) IsInCheck() bool {
	if pos.WhiteMove {
		return pos.IsSquareAttacked(BitScan(pos.White&pos.Kings), false)
	} else {
		return pos.IsSquareAttacked(BitScan(pos.Black&pos.Kings), true)
	}
}

func (pos *Position) IsSquareAttacked(square int, side bool) bool {
	var theirOccupancy, attackedSquares uint64
	if side {
		theirOccupancy = pos.White
		attackedSquares = BlackPawnAttacks[square] & pos.Pawns & theirOccupancy
	} else {
		theirOccupancy = pos.Black
		attackedSquares = WhitePawnAttacks[square] & pos.Pawns & theirOccupancy
	}
	if attackedSquares != 0 {
		return true
	}
	if KnightAttacks[square]&theirOccupancy&pos.Knights != 0 {
		return true
	}
	if KingAttacks[square]&pos.Kings&theirOccupancy != 0 {
		return true
	}
	allOccupation := pos.White | pos.Black
	if BishopAttacks(square, allOccupation)&(pos.Queens|pos.Bishops)&theirOccupancy != 0 {
		return true
	}
	if RookAttacks(square, allOccupation)&(pos.Queens|pos.Rooks)&theirOccupancy != 0 {
		return true
	}
	return false
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
			if pos.White&bb != 0 {
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

func (pos *Position) IntSide() (res int) {
	if pos.WhiteMove {
		res = 1
	} else {
		res = 0
	}
	return
}

func (pos *Position) MakeLegalMove(move Move, res *Position) {
	res.WhiteMove = pos.WhiteMove
	res.Pawns = pos.Pawns
	res.Knights = pos.Knights
	res.Bishops = pos.Bishops
	res.Rooks = pos.Rooks
	res.Kings = pos.Kings
	res.Queens = pos.Queens
	res.White = pos.White
	res.Black = pos.Black
	res.Flags = pos.Flags
	res.Key = pos.Key ^ zobristColor ^ zobristEpSquare[pos.EpSquare] ^ zobristFlags[pos.Flags]
	res.PawnKey = pos.PawnKey ^ zobristColor
	copy(res.Pieces[:], pos.Pieces[:])

	movedPiece := pos.TypeOnSquare(move.From())

	res.FiftyMove = pos.FiftyMove + 1

	res.EpSquare = 0

	switch move.Type() {
	case NormalMove:
		res.MovePiece(movedPiece, pos.WhiteMove, move.From(), move.To())
		if move.Special() == CaptureMove {
			res.FiftyMove = 0
			capturedPiece := pos.TypeOnSquare(move.To())
			res.RemovePiece(capturedPiece, !pos.WhiteMove, move.To())
		} else if movedPiece == Pawn {
			res.FiftyMove = 0
			if move.Special() == QuietMove && utils.Abs(int64(move.From()-move.To())) == 16 {
				res.EpSquare = move.To()
				res.Key ^= zobristEpSquare[move.To()]
			}
		}
	case CastleMove:
		res.MovePiece(King, pos.WhiteMove, move.From(), move.To())
		switch move {
		case WhiteKingSideCastle:
			res.MovePiece(Rook, true, H1, F1)
		case WhiteQueenSideCastle:
			res.MovePiece(Rook, true, A1, D1)
		case BlackKingSideCastle:
			res.MovePiece(Rook, false, H8, F8)
		case BlackQueenSideCastle:
			res.MovePiece(Rook, false, A8, D8)
		}
	case EnpassMove:
		res.FiftyMove = 0
		res.MovePiece(Pawn, pos.WhiteMove, move.From(), move.To())
		res.RemovePiece(Pawn, !pos.WhiteMove, pos.EpSquare)
	case PromotionMove:
		res.FiftyMove = 0
		res.RemovePiece(Pawn, pos.WhiteMove, move.From())
		capturedPiece := pos.TypeOnSquare(move.To())
		if capturedPiece != None {
			res.RemovePiece(capturedPiece, !pos.WhiteMove, move.To())
		}
		res.SetPiece(move.PromotedPiece(), pos.WhiteMove, move.To())
	}

	res.Key ^= zobristFlags[res.Flags]
	res.WhiteMove = !pos.WhiteMove
	res.LastMove = move
}
