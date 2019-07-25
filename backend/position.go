package backend

import (
	"fmt"
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
	Pawns, Knights, Bishops, Rooks, Queens, Kings, White, Black uint64
	Flags                                                       int
	EpSquare                                                    int
	WhiteMove                                                   bool
	FiftyMove                                                   int32
	LastMove                                                    Move
	Key                                                         uint64
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

var InitialPosition Position = Position{
	0xff00000000ff00, 0x4200000000000042, 0x2400000000000024,
	0x8100000000000081, 0x800000000000008, 0x1000000000000010,
	0xffff, 0xffff000000000000, 0, 0, true, 0, 0, 0}

func init() {
	HashPosition(&InitialPosition)
}

func (pos *Position) TypeOnSquare(squareBB uint64) int {
	if squareBB&pos.Pawns != 0 {
		return Pawn
	} else if squareBB&pos.Knights != 0 {
		return Knight
	} else if squareBB&pos.Bishops != 0 {
		return Bishop
	} else if squareBB&pos.Rooks != 0 {
		return Rook
	} else if squareBB&pos.Queens != 0 {
		return Queen
	} else if squareBB&pos.Kings != 0 {
		return King
	}
	return None
}

func (p *Position) MovePiece(piece int, side bool, from int, to int) {
	var b = SquareMask[from] ^ SquareMask[to]
	var intSide = 0
	if side {
		p.White ^= b
	} else {
		p.Black ^= b
		intSide = 1
	}
	switch piece {
	case Pawn:
		p.Pawns ^= b
		p.Key ^= zobrist[0][intSide][from] ^ zobrist[0][intSide][to]
	case Knight:
		p.Knights ^= b
		p.Key ^= zobrist[1][intSide][from] ^ zobrist[1][intSide][to]
	case Bishop:
		p.Bishops ^= b
		p.Key ^= zobrist[2][intSide][from] ^ zobrist[2][intSide][to]
	case Rook:
		p.Rooks ^= b
		p.Key ^= zobrist[3][intSide][from] ^ zobrist[3][intSide][to]
		if side {
			if from == H1 {
				p.Flags |= WhiteKingSideCastleFlag
			} else if from == A1 {
				p.Flags |= WhiteQueenSideCastleFlag
			}
		} else {
			if from == H8 {
				p.Flags |= BlackKingSideCastleFlag
			} else if from == A8 {
				p.Flags |= BlackQueenSideCastleFlag
			}
		}
	case Queen:
		p.Queens ^= b
		p.Key ^= zobrist[4][intSide][from] ^ zobrist[4][intSide][to]
	case King:
		p.Kings ^= b
		p.Key ^= zobrist[5][intSide][from] ^ zobrist[5][intSide][to]
		if side {
			p.Flags |= WhiteKingSideCastleFlag | WhiteQueenSideCastleFlag
		} else {
			p.Flags |= BlackKingSideCastleFlag | BlackQueenSideCastleFlag
		}
	}
}

func (p *Position) TogglePiece(piece int, side bool, square int) {
	var b = SquareMask[square]
	var intSide = 0
	if side {
		p.White ^= b
	} else {
		p.Black ^= b
		intSide = 1
	}
	switch piece {
	case Pawn:
		p.Pawns ^= b
		p.Key ^= zobrist[0][intSide][square]
	case Knight:
		p.Knights ^= b
		p.Key ^= zobrist[1][intSide][square]
	case Bishop:
		p.Bishops ^= b
		p.Key ^= zobrist[2][intSide][square]
	case Rook:
		p.Rooks ^= b
		p.Key ^= zobrist[3][intSide][square]
	case Queen:
		p.Queens ^= b
		p.Key ^= zobrist[4][intSide][square]
	case King:
		p.Kings ^= b
		p.Key ^= zobrist[5][intSide][square]
	}
}

func (pos *Position) MakeNullMove(res *Position) {
	res.WhiteMove = !pos.WhiteMove
	res.Pawns = pos.Pawns
	res.Knights = pos.Knights
	res.Bishops = pos.Bishops
	res.Rooks = pos.Rooks
	res.Kings = pos.Kings
	res.Queens = pos.Queens
	res.White = pos.White
	res.Black = pos.Black
	res.Flags = pos.Flags
	res.Key = pos.Key ^ zobristColor ^ zobristEpSquare[pos.EpSquare]

	res.FiftyMove = pos.FiftyMove
	res.LastMove = NullMove
	res.EpSquare = 0
}

func (pos *Position) MakeMove(move Move, res *Position) bool {
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

	if move.MovedPiece() == Pawn || move.IsCapture() {
		res.FiftyMove = 0
	} else {
		res.FiftyMove = pos.FiftyMove + 1
	}

	res.EpSquare = 0

	if !move.IsPromotion() {
		res.MovePiece(move.MovedPiece(), pos.WhiteMove, move.From(), move.To())
		switch move.Type() {
		case DoublePawnPush:
			res.EpSquare = move.To()
			res.Key ^= zobristEpSquare[move.To()]
		case Capture:
			res.TogglePiece(move.CapturedPiece(), !pos.WhiteMove, move.To())
			if move.CapturedPiece() == Rook {
				switch move.To() {
				case A1:
					res.Flags |= WhiteQueenSideCastleFlag
				case H1:
					res.Flags |= WhiteKingSideCastleFlag
				case H8:
					res.Flags |= BlackKingSideCastleFlag
				case A8:
					res.Flags |= BlackQueenSideCastleFlag
				}
			}
		case KingCastle:
			if pos.WhiteMove {
				res.MovePiece(Rook, true, H1, F1)
			} else {
				res.MovePiece(Rook, false, H8, F8)
			}
		case QueenCastle:
			if pos.WhiteMove {
				res.MovePiece(Rook, true, A1, D1)
			} else {
				res.MovePiece(Rook, false, A8, D8)
			}
		case EPCapture:
			res.TogglePiece(Pawn, !pos.WhiteMove, pos.EpSquare)
		}
	} else {
		res.TogglePiece(Pawn, pos.WhiteMove, move.From())
		switch move.Type() {
		case QueenPromotion:
			res.TogglePiece(Queen, pos.WhiteMove, move.To())
		case QueenCapturePromotion:
			res.TogglePiece(Queen, pos.WhiteMove, move.To())
			res.TogglePiece(move.CapturedPiece(), !pos.WhiteMove, move.To())
		case RookPromotion:
			res.TogglePiece(Rook, pos.WhiteMove, move.To())
		case RookCapturePromotion:
			res.TogglePiece(Rook, pos.WhiteMove, move.To())
			res.TogglePiece(move.CapturedPiece(), !pos.WhiteMove, move.To())
		case KnightPromotion:
			res.TogglePiece(Knight, pos.WhiteMove, move.To())
		case BishopPromotion:
			res.TogglePiece(Bishop, pos.WhiteMove, move.To())
		case KnightCapturePromotion:
			res.TogglePiece(Knight, pos.WhiteMove, move.To())
			res.TogglePiece(move.CapturedPiece(), !pos.WhiteMove, move.To())
		case BishopCapturePromotion:
			res.TogglePiece(Bishop, pos.WhiteMove, move.To())
			res.TogglePiece(move.CapturedPiece(), !pos.WhiteMove, move.To())
		}
	}

	if res.IsInCheck() {
		return false
	}

	res.Key ^= zobristFlags[res.Flags]
	res.WhiteMove = !pos.WhiteMove
	res.LastMove = move
	return true
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

func (pos *Position) IntSide() int {
	if pos.WhiteMove {
		return 1
	}
	return 0
}

func (pos *Position) IsMovePseudoLegal(move Move) bool {
	var we, them uint64
	if pos.WhiteMove {
		we = pos.White
		them = pos.Black
	} else {
		we = pos.Black
		them = pos.White
	}
	occupancy := we | them
	fromMask := SquareMask[move.From()]
	toMask := SquareMask[move.To()]

	if move == NullMove || (we&fromMask) == 0 {
		return false
	}

	switch move.MovedPiece() {
	case Knight:
		return move.IsNormal() && (pos.Knights&fromMask) != 0 && (KnightAttacks[move.From()] & ^we)&toMask != 0
	case Bishop:
		return move.IsNormal() && (pos.Bishops&fromMask) != 0 && (BishopAttacks(move.From(), occupancy) & ^we)&toMask != 0
	case Rook:
		return move.IsNormal() && (pos.Rooks&fromMask) != 0 && (RookAttacks(move.From(), occupancy) & ^we)&toMask != 0
	case Queen:
		return move.IsNormal() && (pos.Queens&fromMask) != 0 && (QueenAttacks(move.From(), occupancy) & ^we)&toMask != 0
	case Pawn:
		if (pos.Pawns&fromMask) == 0 || move.IsCastling() {
			return false
		}
		var attacks, forward uint64
		if pos.WhiteMove {
			if move.Type() == EPCapture && pos.EpSquare != 0 {
				return (SquareMask[uint(pos.EpSquare)-1]|SquareMask[uint(pos.EpSquare)]|SquareMask[uint(pos.EpSquare)+1])&RANK_5_BB&fromMask != 0
			}
			attacks = WhitePawnAttacks[move.From()]
			forward = North(fromMask) & ^occupancy
		} else {
			if move.Type() == EPCapture && pos.EpSquare != 0 {
				return (SquareMask[uint(pos.EpSquare)-1]|SquareMask[uint(pos.EpSquare)]|SquareMask[uint(pos.EpSquare)+1])&RANK_4_BB&fromMask != 0
			}
			attacks = BlackPawnAttacks[move.From()]
			forward = South(fromMask) & ^occupancy
		}
		if move.IsPromotion() {
			return PROMOTION_RANKS&((attacks&them)|forward) != 0
		}

		// Double pawn push
		if forward != 0 && pos.WhiteMove && fromMask&RANK_2_BB != 0 && North(forward)&occupancy == 0 {
			forward |= North(forward)
		} else if forward != 0 && !pos.WhiteMove && fromMask&RANK_7_BB != 0 && South(forward)&occupancy == 0 {
			forward |= South(forward)
		}
		return (^PROMOTION_RANKS)&((attacks&them)|forward)&toMask != 0
	case King:
		if (pos.Kings & fromMask) == 0 {
			return false
		}
		if move.IsNormal() {
			return move.IsNormal() && (KingAttacks[move.From()] & ^we)&toMask != 0
		}
		if pos.WhiteMove {
			if move == WhiteKingSideCastle {
				return occupancy&WHITE_KING_CASTLE_BLOCK_BB == 0 && pos.Flags&WhiteKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E1, false) && !pos.IsSquareAttacked(F1, false)
			} else if move == WhiteQueenSideCastle {
				return occupancy&WHITE_QUEEN_CASTLE_BLOCK_BB == 0 && pos.Flags&WhiteQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E1, false) && !pos.IsSquareAttacked(D1, false)
			} else {
				return false
			}
		} else {
			if move == BlackKingSideCastle {
				return occupancy&BLACK_KING_CASTLE_BLOCK_BB == 0 && pos.Flags&BlackKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E8, true) && !pos.IsSquareAttacked(F8, true)
			} else if move == BlackQueenSideCastleFlag {
				return occupancy&BLACK_QUEEN_CASTLE_BLOCK_BB == 0 && pos.Flags&BlackQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E8, true) && !pos.IsSquareAttacked(D8, true)
			} else {
				return false
			}
		}
	}

	return false
}
