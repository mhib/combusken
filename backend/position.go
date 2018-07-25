package backend

import (
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
	LastMove                                                    Move
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
	0xffff, 0xffff000000000000, 0, 0, true, 0}

func (pos *Position) GenerateAllMoves(buffer []Move) []Move {
	var counter = 0
	allOccupation := pos.White | pos.Black
	var fromBB, from, toBB, to, ourOccupation, theirOccupation uint64
	var fromId, toId int
	if pos.WhiteMove {
		ourOccupation = pos.White
		theirOccupation = pos.Black
		for fromBB = pos.Pawns & pos.White; fromBB > 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			from = SquareMask[uint(fromId)]
			if from&RANK_7_BB != 0 {
				to = from << 8
				if allOccupation&to == 0 {
					buffer[counter] = NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 0, 0))
					counter++
					buffer[counter] = NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 0, 1))
					counter++
					buffer[counter] = NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 1, 0))
					counter++
					buffer[counter] = NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 1, 1))
					counter++
				}
				for toBB = WhitePawnAttacks[fromId] & pos.Black; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					to = SquareMask[uint(toId)]
					captureType := pos.TypeOnSquare(to)
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 0, 0))
					counter++
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 0, 1))
					counter++
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 1, 0))
					counter++
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 1, 1))
					counter++
				}
			} else {
				to = from << 8
				if allOccupation&to == 0 {
					buffer[counter] = NewMove(fromId, fromId+8, Pawn, None, 0)
					counter++
					if from&RANK_2_BB != 0 && allOccupation&(from<<16) == 0 {
						buffer[counter] = NewMove(fromId, fromId+16, Pawn, None, NewType(0, 0, 0, 1))
						counter++
					}
				}
				for toBB = WhitePawnAttacks[fromId] & pos.Black; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					to = SquareMask[toId]
					captureType := pos.TypeOnSquare(SquareMask[toId])
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 0, 0, 0))
					counter++
				}
			}
		}
		if pos.EpSquare != 0 {
			epSquareBB := SquareMask[uint(pos.EpSquare)]
			epBB := (west(epSquareBB) | east(epSquareBB)) & RANK_5_BB
			for fromBB = epBB & pos.Pawns & pos.White; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[counter] = NewMove(fromId, pos.EpSquare+8, Pawn, Pawn, NewType(1, 0, 0, 1))
				counter++
			}
		}
	} else {
		ourOccupation = pos.Black
		theirOccupation = pos.White
		for fromBB = pos.Pawns & pos.Black; fromBB > 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			from = SquareMask[uint(fromId)]
			if from&RANK_2_BB != 0 {
				to = from >> 8
				if allOccupation&to == 0 {
					buffer[counter] = NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 0, 0))
					counter++
					buffer[counter] = NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 0, 1))
					counter++
					buffer[counter] = NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 1, 0))
					counter++
					buffer[counter] = NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 1, 1))
					counter++
				}
				for toBB = BlackPawnAttacks[fromId] & pos.White; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					to = SquareMask[uint(toId)]
					captureType := pos.TypeOnSquare(to)
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 0, 0))
					counter++
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 0, 1))
					counter++
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 1, 0))
					counter++
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 1, 1))
					counter++
				}
			} else {
				to = SquareMask[fromId-8]
				if allOccupation&to == 0 {
					buffer[counter] = NewMove(fromId, fromId-8, Pawn, None, 0)
					counter++
					if from&RANK_7_BB != 0 && allOccupation&(from>>16) == 0 {
						buffer[counter] = NewMove(fromId, fromId-16, Pawn, None, NewType(0, 0, 0, 1))
						counter++
					}
				}
				for toBB = BlackPawnAttacks[fromId] & pos.White; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					to = SquareMask[uint(toId)]
					captureType := pos.TypeOnSquare(to)
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 0, 0, 0))
					counter++
				}
			}
		}
		if pos.EpSquare != 0 {
			epSquareBB := SquareMask[uint(pos.EpSquare)]
			epBB := (west(epSquareBB) | east(epSquareBB)) & RANK_4_BB
			for fromBB = epBB & pos.Pawns & pos.Black; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[counter] = NewMove(fromId, pos.EpSquare-8, Pawn, Pawn, NewType(1, 0, 0, 1))
				counter++
			}
		}
	}
	// Knights
	for fromBB = pos.Knights & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = KnightAttacks[fromId] & (^ourOccupation); toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			to = SquareMask[uint(toId)]
			if to&theirOccupation != 0 {
				buffer[counter] = NewMove(fromId, toId, Knight, pos.TypeOnSquare(to), NewType(1, 0, 0, 0))
				counter++
			} else {
				buffer[counter] = NewMove(fromId, toId, Knight, None, NewType(0, 0, 0, 0))
				counter++
			}
		}
	}
	// end of knights

	// Kings
	from = pos.Kings & ourOccupation
	fromId = BitScan(from)
	for toBB = KingAttacks[fromId] & (^ourOccupation); toBB != 0; toBB &= (toBB - 1) {
		toId = BitScan(toBB)
		to = SquareMask[uint(toId)]
		if to&theirOccupation != 0 {
			buffer[counter] = NewMove(fromId, toId, King, pos.TypeOnSquare(to), NewType(1, 0, 0, 0))
			counter++
		} else {
			buffer[counter] = NewMove(fromId, toId, King, None, NewType(0, 0, 0, 0))
			counter++
		}
	}
	// end of Kings

	// Rooks
	for fromBB = pos.Rooks & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = RookAttacks(fromId, allOccupation) & ^ourOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			to = SquareMask[uint(toId)]
			if to&theirOccupation != 0 {
				buffer[counter] = NewMove(fromId, toId, Rook, pos.TypeOnSquare(to), NewType(1, 0, 0, 0))
				counter++
			} else {
				buffer[counter] = NewMove(fromId, toId, Rook, None, NewType(0, 0, 0, 0))
				counter++
			}
		}
	}
	// end of Rooks

	// Bishops
	for fromBB = pos.Bishops & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = BishopAttacks(fromId, allOccupation) & ^ourOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			to = SquareMask[uint(toId)]
			if to&theirOccupation != 0 {
				buffer[counter] = NewMove(fromId, toId, Bishop, pos.TypeOnSquare(to), NewType(1, 0, 0, 0))
				counter++
			} else {
				buffer[counter] = NewMove(fromId, toId, Bishop, None, NewType(0, 0, 0, 0))
				counter++
			}
		}
	}
	// end of Bishops

	// Queens
	for fromBB = pos.Queens & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = QueenAttacks(fromId, allOccupation) & ^ourOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			to = SquareMask[uint(toId)]
			if to&theirOccupation != 0 {
				buffer[counter] = NewMove(fromId, toId, Queen, pos.TypeOnSquare(to), NewType(1, 0, 0, 0))
				counter++
			} else {
				buffer[counter] = NewMove(fromId, toId, Queen, None, NewType(0, 0, 0, 0))
				counter++
			}
		}
	}
	// end of Queens

	// Castling
	if pos.WhiteMove {
		if allOccupation&WHITE_KING_CASTLE_BLOCK_BB == 0 && pos.Flags&WhiteKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E1_MASK, false) && !pos.IsSquareAttacked(F1_MASK, false) {
			buffer[counter] = WhiteKingSideCastle
			counter++
		}
		if allOccupation&WHITE_QUEEN_CASTLE_BLOCK_BB == 0 && pos.Flags&WhiteQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E1_MASK, false) && !pos.IsSquareAttacked(D1_MASK, false) {
			buffer[counter] = WhiteQueenSideCastle
			counter++
		}
	} else {
		if allOccupation&BLACK_KING_CASTLE_BLOCK_BB == 0 && pos.Flags&BlackKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E8_MASK, true) && !pos.IsSquareAttacked(F8_MASK, true) {
			buffer[counter] = BlackKingSideCastle
			counter++
		}
		if allOccupation&BLACK_QUEEN_CASTLE_BLOCK_BB == 0 && pos.Flags&BlackQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E8_MASK, true) && !pos.IsSquareAttacked(D1_MASK, true) {
			buffer[counter] = BlackQueenSideCastle
			counter++
		}

	}
	// end of Castling

	return buffer[:counter]
}

func (pos *Position) TypeOnSquare(squareBB uint64) int {
	if squareBB&pos.Pawns != 0 {
		return Pawn
	}
	if squareBB&pos.Rooks != 0 {
		return Rook
	}
	if squareBB&pos.Knights != 0 {
		return Knight
	}
	if squareBB&pos.Bishops != 0 {
		return Bishop
	}
	if squareBB&pos.Queens != 0 {
		return Queen
	}
	if squareBB&pos.Kings != 0 {
		return King
	}
	return None
}

func (p *Position) MovePiece(piece int, side bool, from int, to int) {
	var b = SquareMask[from] ^ SquareMask[to]
	if side {
		p.White ^= b
	} else {
		p.Black ^= b
	}
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

func (p *Position) TogglePiece(piece int, side bool, square int) {
	var b = SquareMask[square]
	if side {
		p.White ^= b
	} else {
		p.Black ^= b
	}
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

	res.EpSquare = 0

	switch move {
	case WhiteKingSideCastle:
		res.MovePiece(King, true, E1, G1)
		res.MovePiece(Rook, true, H1, F1)
	case WhiteQueenSideCastle:
		res.MovePiece(King, true, E1, C1)
		res.MovePiece(Rook, true, A1, D1)
	case BlackKingSideCastle:
		res.MovePiece(King, false, E8, G8)
		res.MovePiece(Rook, false, H8, F8)
	case BlackQueenSideCastle:
		res.MovePiece(King, false, E8, C8)
		res.MovePiece(Rook, false, A8, D8)
	default:
		if move.Type() == DoublePawnPush {
			res.EpSquare = move.To()
		} else if move.Type() == 1 {
			res.TogglePiece(move.CapturedPiece(), !pos.WhiteMove, move.To())
		} else if move.Type() == EPCapture {
			if pos.WhiteMove {
				res.TogglePiece(Pawn, false, move.To()-8)
			} else {
				res.TogglePiece(Pawn, true, move.To()+8)
			}
		}
		res.MovePiece(move.MovedPiece(), pos.WhiteMove, move.From(), move.To())
	}
	if !res.IsValid() {
		return false
	}
	res.WhiteMove = !pos.WhiteMove
	res.LastMove = move
	return true
}

func (pos *Position) IsValid() bool {
	if pos.WhiteMove {
		return !pos.IsSquareAttacked(pos.White&pos.Kings, !pos.WhiteMove)
	} else {
		return !pos.IsSquareAttacked(pos.Black&pos.Kings, !pos.WhiteMove)
	}
}

func (pos *Position) IsSquareAttacked(squareBB uint64, side bool) bool {
	var ourOccupancy, attackedSquares uint64
	allOccupation := pos.White | pos.Black
	if side {
		ourOccupancy = pos.White
		attackedSquares = WhitePawnsAttacks(pos.Pawns & pos.White)
	} else {
		ourOccupancy = pos.Black
		attackedSquares = BlackPawnsAttacks(pos.Pawns & pos.Black)
	}
	if attackedSquares&squareBB != 0 {
		return true
	}
	if KnightsAttacks(ourOccupancy&pos.Knights)&squareBB != 0 {
		return true
	}
	if BishopsAttacks(ourOccupancy&pos.Bishops, allOccupation)&squareBB != 0 {
		return true
	}
	if RooksAttacks(ourOccupancy&pos.Rooks, allOccupation)&squareBB != 0 {
		return true
	}
	if QueensAttacks(ourOccupancy&pos.Queens, allOccupation)&squareBB != 0 {
		return true
	}
	if kingAttacks(ourOccupancy&pos.Kings)&squareBB != 0 {
		return true
	}
	return false
}
