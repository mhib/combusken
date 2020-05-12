package backend

type EvaledMove struct {
	Move
	Value int32
}

func addPromotions(move Move, buffer []EvaledMove) {
	// generate for all special flags
	buffer[0].Move = move ^ Move((0xe)<<18)
	buffer[1].Move = move ^ Move((0xa)<<18)
	buffer[2].Move = move ^ Move((0x6)<<18)
	buffer[3].Move = move ^ Move((0x2)<<18)
}

var forwardByColor = [2]int{-8, +8}
var secondRank = [2]int{RANK_7, RANK_2}
var promotionBB = [2]uint64{RANK_2_BB, RANK_7_BB}
var epRankBB = [2]uint64{RANK_4_BB, RANK_5_BB}

func (pos *Position) GenerateQuiet(buffer []EvaledMove) (size uint8) {
	var fromBB, toBB, toMask uint64
	var fromId, toId int
	sideToMove := pos.SideToMove
	ourOccupation := pos.Colours[sideToMove]
	theirOccupation := pos.Colours[sideToMove^1]
	allOccupation := ourOccupation | theirOccupation
	forward := forwardByColor[sideToMove]
	for fromBB = pos.Pieces[Pawn] & ourOccupation & ^promotionBB[sideToMove]; fromBB > 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		toId = fromId + forward
		toMask = SquareMask[toId]
		if allOccupation&toMask == 0 {
			buffer[size].Move = NewMove(fromId, toId, Pawn, None, 0)
			size++

			// Double pawn push
			toId += forward
			toMask = SquareMask[toId]
			if Rank(fromId) == secondRank[sideToMove] && allOccupation&toMask == 0 {
				buffer[size].Move = NewMove(fromId, toId, Pawn, None, NewType(0, 0, 0, 1))
				size++
			}
		}
	}

	// Castling
	if pos.SideToMove == White {
		if allOccupation&WHITE_KING_CASTLE_BLOCK_BB == 0 && pos.Flags&WhiteKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E1, Black) && !pos.IsSquareAttacked(F1, Black) {
			buffer[size].Move = WhiteKingSideCastle
			size++
		}
		if allOccupation&WHITE_QUEEN_CASTLE_BLOCK_BB == 0 && pos.Flags&WhiteQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E1, Black) && !pos.IsSquareAttacked(D1, Black) {
			buffer[size].Move = WhiteQueenSideCastle
			size++
		}
	} else {
		if allOccupation&BLACK_KING_CASTLE_BLOCK_BB == 0 && pos.Flags&BlackKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E8, White) && !pos.IsSquareAttacked(F8, White) {
			buffer[size].Move = BlackKingSideCastle
			size++
		}
		if allOccupation&BLACK_QUEEN_CASTLE_BLOCK_BB == 0 && pos.Flags&BlackQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E8, White) && !pos.IsSquareAttacked(D8, White) {
			buffer[size].Move = BlackQueenSideCastle
			size++
		}

	}

	// Knights
	for fromBB = pos.Pieces[Knight] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = KnightAttacks[fromId] & ^allOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			buffer[size].Move = NewMove(fromId, toId, Knight, None, NewType(0, 0, 0, 0))
			size++
		}
	}
	// end of knights

	// Kings
	fromId = BitScan(pos.Pieces[King] & ourOccupation)
	for toBB = KingAttacks[fromId] & ^allOccupation; toBB != 0; toBB &= (toBB - 1) {
		toId = BitScan(toBB)
		buffer[size].Move = NewMove(fromId, toId, King, None, NewType(0, 0, 0, 0))
		size++
	}
	// end of Kings

	// Rooks
	for fromBB = pos.Pieces[Rook] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = RookAttacks(fromId, allOccupation) & ^allOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			buffer[size].Move = NewMove(fromId, toId, Rook, None, NewType(0, 0, 0, 0))
			size++
		}
	}
	// end of Rooks

	// Bishops
	for fromBB = pos.Pieces[Bishop] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = BishopAttacks(fromId, allOccupation) & ^allOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			buffer[size].Move = NewMove(fromId, toId, Bishop, None, NewType(0, 0, 0, 0))
			size++
		}
	}
	// end of Bishops

	// Queens
	for fromBB = pos.Pieces[Queen] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = QueenAttacks(fromId, allOccupation) & ^allOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			buffer[size].Move = NewMove(fromId, toId, Queen, None, NewType(0, 0, 0, 0))
			size++
		}
	}
	// end of Queens

	return
}

func (pos *Position) GenerateNoisy(buffer []EvaledMove) (size uint8) {
	var fromBB, toBB uint64
	var fromId, toId, what int

	sideToMove := pos.SideToMove
	ourOccupation := pos.Colours[sideToMove]
	theirOccupation := pos.Colours[sideToMove^1]
	allOccupation := ourOccupation | theirOccupation

	// PAWNS
	forward := forwardByColor[sideToMove]
	if pos.EpSquare != 0 {
		fromBB = (SquareMask[uint(pos.EpSquare)-1] | SquareMask[uint(pos.EpSquare)+1]) &
			epRankBB[sideToMove] & pos.Pieces[Pawn] & ourOccupation
		for ; fromBB > 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			buffer[size].Move = NewMove(fromId, pos.EpSquare+forward, Pawn, Pawn, NewType(1, 0, 0, 1))
			size++
		}
	}
	if sideToMove == White {
		fromBB = BlackPawnsAttacks(theirOccupation) | RANK_7_BB
	} else {
		fromBB = WhitePawnsAttacks(theirOccupation) | RANK_2_BB
	}
	for fromBB &= pos.Pieces[Pawn] & ourOccupation; fromBB != 0; fromBB &= fromBB - 1 {
		fromId = BitScan(fromBB)
		if Rank(fromId) == secondRank[sideToMove^1] {
			if SquareMask[fromId+forward]&allOccupation == 0 {
				toId = fromId + forward
				addPromotions(NewMove(fromId, toId, Pawn, None, 0), buffer[size:])
				size += 4
			}
			for toBB = PawnAttacks[sideToMove][fromId] & theirOccupation; toBB > 0; toBB &= (toBB - 1) {
				toId = BitScan(toBB)
				what = pos.TypeOnSquare(SquareMask[uint(toId)])
				addPromotions(NewMove(fromId, toId, Pawn, what, 1), buffer[size:])
				size += 4
			}
		} else {
			for toBB = PawnAttacks[sideToMove][fromId] & theirOccupation; toBB > 0; toBB &= (toBB - 1) {
				toId = BitScan(toBB)
				what = pos.TypeOnSquare(SquareMask[uint(toId)])
				buffer[size].Move = NewMove(fromId, toId, Pawn, what, NewType(1, 0, 0, 0))
				size++
			}
		}
	}
	// end of pawns

	// Knights
	for fromBB = pos.Pieces[Knight] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = KnightAttacks[fromId] & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			what = pos.TypeOnSquare(SquareMask[uint(toId)])
			buffer[size].Move = NewMove(fromId, toId, Knight, what, NewType(1, 0, 0, 0))
			size++
		}
	}
	// end of knights

	// Bishops
	for fromBB = pos.Pieces[Bishop] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = BishopAttacks(fromId, allOccupation) & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			what = pos.TypeOnSquare(SquareMask[uint(toId)])
			buffer[size].Move = NewMove(fromId, toId, Bishop, what, NewType(1, 0, 0, 0))
			size++
		}
	}
	// end of Bishops

	// Rooks
	for fromBB = pos.Pieces[Rook] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = RookAttacks(fromId, allOccupation) & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			what = pos.TypeOnSquare(SquareMask[uint(toId)])
			buffer[size].Move = NewMove(fromId, toId, Rook, what, NewType(1, 0, 0, 0))
			size++
		}
	}
	// end of Rooks

	// Queens
	for fromBB = pos.Pieces[Queen] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = QueenAttacks(fromId, allOccupation) & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			what = pos.TypeOnSquare(SquareMask[uint(toId)])
			buffer[size].Move = NewMove(fromId, toId, Queen, what, NewType(1, 0, 0, 0))
			size++
		}
	}
	// end of Queens

	// Kings
	fromBB = pos.Pieces[King] & ourOccupation
	fromId = BitScan(fromBB)
	for toBB = KingAttacks[fromId] & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
		toId = BitScan(toBB)
		what = pos.TypeOnSquare(SquareMask[uint(toId)])
		buffer[size].Move = NewMove(fromId, toId, King, what, NewType(1, 0, 0, 0))
		size++
	}
	// end of Kings

	return
}
