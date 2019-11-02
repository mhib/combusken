package backend

type EvaledMove struct {
	Move
	Value int
}

func (pos *Position) GenerateAllMoves(buffer []EvaledMove) []EvaledMove {
	var size = 0
	var fromBB, from, toBB, toMask uint64
	var fromId, toId int
	ourOccupation := pos.Colours[pos.SideToMove]
	theirOccupation := pos.Colours[pos.SideToMove^1]
	allOccupation := ourOccupation | theirOccupation

	if pos.SideToMove == White {
		for fromBB = pos.Pieces[Pawn] & pos.Colours[White]; fromBB > 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			from = SquareMask[uint(fromId)]
			if from&RANK_7_BB != 0 {
				toMask = from << 8
				if allOccupation&toMask == 0 {
					buffer[size].Move = NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 1, 1))
					size++
					buffer[size].Move = NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 0, 1))
					size++
					buffer[size].Move = NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 1, 0))
					size++
					buffer[size].Move = NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 0, 0))
					size++
				}
				for toBB = PawnAttacks[White][fromId] & pos.Colours[Black]; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					toMask = SquareMask[uint(toId)]
					captureType := pos.TypeOnSquare(toMask)
					buffer[size].Move = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 1, 1))
					size++
					buffer[size].Move = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 0, 1))
					size++
					buffer[size].Move = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 1, 0))
					size++
					buffer[size].Move = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 0, 0))
					size++
				}
			} else {
				toMask = from << 8
				if allOccupation&toMask == 0 {
					buffer[size].Move = NewMove(fromId, fromId+8, Pawn, None, 0)
					size++
					if from&RANK_2_BB != 0 && allOccupation&(from<<16) == 0 {
						buffer[size].Move = NewMove(fromId, fromId+16, Pawn, None, NewType(0, 0, 0, 1))
						size++
					}
				}
				for toBB = PawnAttacks[White][fromId] & pos.Colours[Black]; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					toMask = SquareMask[toId]
					captureType := pos.TypeOnSquare(SquareMask[toId])
					buffer[size].Move = NewMove(fromId, toId, Pawn, captureType, NewType(1, 0, 0, 0))
					size++
				}
			}
		}
		if pos.EpSquare != 0 {
			epBB := (SquareMask[uint(pos.EpSquare)-1] | SquareMask[uint(pos.EpSquare)] | SquareMask[uint(pos.EpSquare)+1]) & RANK_5_BB
			for fromBB = epBB & pos.Pieces[Pawn] & pos.Colours[White]; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[size].Move = NewMove(fromId, pos.EpSquare+8, Pawn, Pawn, NewType(1, 0, 0, 1))
				size++
			}
		}

		// Castling
		if allOccupation&WHITE_KING_CASTLE_BLOCK_BB == 0 && pos.Flags&WhiteKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E1, Black) && !pos.IsSquareAttacked(F1, Black) {
			buffer[size].Move = WhiteKingSideCastle
			size++
		}
		if allOccupation&WHITE_QUEEN_CASTLE_BLOCK_BB == 0 && pos.Flags&WhiteQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E1, Black) && !pos.IsSquareAttacked(D1, Black) {
			buffer[size].Move = WhiteQueenSideCastle
			size++
		}
	} else {
		for fromBB = pos.Pieces[Pawn] & pos.Colours[Black]; fromBB > 0; fromBB &= (fromBB - 1) {
			fromId = BitScan(fromBB)
			from = SquareMask[uint(fromId)]
			if from&RANK_2_BB != 0 {
				toMask = from >> 8
				if allOccupation&toMask == 0 {
					buffer[size].Move = NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 1, 1))
					size++
					buffer[size].Move = NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 1, 0))
					size++
					buffer[size].Move = NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 0, 1))
					size++
					buffer[size].Move = NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 0, 0))
					size++
				}
				for toBB = PawnAttacks[Black][fromId] & pos.Colours[White]; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					toMask = SquareMask[uint(toId)]
					captureType := pos.TypeOnSquare(toMask)
					buffer[size].Move = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 1, 1))
					size++
					buffer[size].Move = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 1, 0))
					size++
					buffer[size].Move = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 0, 1))
					size++
					buffer[size].Move = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 0, 0))
					size++
				}
			} else {
				toMask = SquareMask[fromId-8]
				if allOccupation&toMask == 0 {
					buffer[size].Move = NewMove(fromId, fromId-8, Pawn, None, 0)
					size++
					if from&RANK_7_BB != 0 && allOccupation&(from>>16) == 0 {
						buffer[size].Move = NewMove(fromId, fromId-16, Pawn, None, NewType(0, 0, 0, 1))
						size++
					}
				}
				for toBB = PawnAttacks[Black][fromId] & pos.Colours[White]; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					toMask = SquareMask[uint(toId)]
					captureType := pos.TypeOnSquare(toMask)
					buffer[size].Move = NewMove(fromId, toId, Pawn, captureType, NewType(1, 0, 0, 0))
					size++
				}
			}
		}
		if pos.EpSquare != 0 {
			epBB := (SquareMask[uint(pos.EpSquare)-1] | SquareMask[uint(pos.EpSquare)] | SquareMask[uint(pos.EpSquare)+1]) & RANK_4_BB
			for fromBB = epBB & pos.Pieces[Pawn] & pos.Colours[Black]; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[size].Move = NewMove(fromId, pos.EpSquare-8, Pawn, Pawn, NewType(1, 0, 0, 1))
				size++
			}
		}

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
		for toBB = KnightAttacks[fromId] & ^ourOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			toMask = SquareMask[uint(toId)]
			if toMask&theirOccupation != 0 {
				buffer[size].Move = NewMove(fromId, toId, Knight, pos.TypeOnSquare(toMask), NewType(1, 0, 0, 0))
				size++
			} else {
				buffer[size].Move = NewMove(fromId, toId, Knight, None, NewType(0, 0, 0, 0))
				size++
			}
		}
	}
	// end of knights

	// Kings
	from = pos.Pieces[King] & ourOccupation
	fromId = BitScan(from)
	for toBB = KingAttacks[fromId] & ^ourOccupation; toBB != 0; toBB &= (toBB - 1) {
		toId = BitScan(toBB)
		toMask = SquareMask[uint(toId)]
		if toMask&theirOccupation != 0 {
			buffer[size].Move = NewMove(fromId, toId, King, pos.TypeOnSquare(toMask), NewType(1, 0, 0, 0))
			size++
		} else {
			buffer[size].Move = NewMove(fromId, toId, King, None, NewType(0, 0, 0, 0))
			size++
		}
	}
	// end of Kings

	// Rooks
	for fromBB = pos.Pieces[Rook] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = RookAttacks(fromId, allOccupation) & ^ourOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			toMask = SquareMask[uint(toId)]
			if toMask&theirOccupation != 0 {
				buffer[size].Move = NewMove(fromId, toId, Rook, pos.TypeOnSquare(toMask), NewType(1, 0, 0, 0))
				size++
			} else {
				buffer[size].Move = NewMove(fromId, toId, Rook, None, NewType(0, 0, 0, 0))
				size++
			}
		}
	}
	// end of Rooks

	// Bishops
	for fromBB = pos.Pieces[Bishop] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = BishopAttacks(fromId, allOccupation) & ^ourOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			toMask = SquareMask[uint(toId)]
			if toMask&theirOccupation != 0 {
				buffer[size].Move = NewMove(fromId, toId, Bishop, pos.TypeOnSquare(toMask), NewType(1, 0, 0, 0))
				size++
			} else {
				buffer[size].Move = NewMove(fromId, toId, Bishop, None, NewType(0, 0, 0, 0))
				size++
			}
		}
	}
	// end of Bishops

	// Queens
	for fromBB = pos.Pieces[Queen] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = QueenAttacks(fromId, allOccupation) & ^ourOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			toMask = SquareMask[uint(toId)]
			if toMask&theirOccupation != 0 {
				buffer[size].Move = NewMove(fromId, toId, Queen, pos.TypeOnSquare(toMask), NewType(1, 0, 0, 0))
				size++
			} else {
				buffer[size].Move = NewMove(fromId, toId, Queen, None, NewType(0, 0, 0, 0))
				size++
			}
		}
	}
	// end of Queens

	return buffer[:size]
}

func (pos *Position) GenerateAllCaptures(buffer []EvaledMove) []EvaledMove {
	var fromBB, toBB uint64
	var fromId, toId, what int

	ourOccupation := pos.Colours[pos.SideToMove]
	theirOccupation := pos.Colours[pos.SideToMove^1]
	allOccupation := ourOccupation | theirOccupation

	var size = 0

	// PAWNS
	if pos.SideToMove == White {
		if pos.EpSquare != 0 {
			epBB := (SquareMask[uint(pos.EpSquare)-1] | SquareMask[uint(pos.EpSquare)] | SquareMask[uint(pos.EpSquare)+1]) & RANK_5_BB
			for fromBB = epBB & pos.Pieces[Pawn] & pos.Colours[White]; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[size] = EvaledMove{NewMove(fromId, pos.EpSquare+8, Pawn, Pawn, NewType(1, 0, 0, 1)), 0}
				size++
			}
		}
		for fromBB = (BlackPawnsAttacks(theirOccupation) | RANK_7_BB) & pos.Pieces[Pawn] & pos.Colours[White]; fromBB != 0; fromBB &= fromBB - 1 {
			fromId = BitScan(fromBB)
			if Rank(fromId) == RANK_7 {
				if SquareMask[fromId+8]&allOccupation == 0 {
					buffer[size] = EvaledMove{NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 1, 1)), 0}
					size++
				}
				if File(fromId) > FILE_A && (SquareMask[fromId+7]&theirOccupation) != 0 {
					what = pos.TypeOnSquare(SquareMask[fromId+7])
					buffer[size] = EvaledMove{NewMove(fromId, fromId+7, Pawn, what, NewType(1, 1, 1, 1)), 0}
					size++
				}
				if File(fromId) < FILE_H && (SquareMask[fromId+9]&theirOccupation) != 0 {
					what = pos.TypeOnSquare(SquareMask[fromId+9])
					buffer[size] = EvaledMove{NewMove(fromId, fromId+9, Pawn, what, NewType(1, 1, 1, 1)), 0}
					size++
				}
			} else {
				if File(fromId) > FILE_A && (SquareMask[fromId+7]&theirOccupation) != 0 {
					what = pos.TypeOnSquare(SquareMask[fromId+7])
					buffer[size] = EvaledMove{NewMove(fromId, fromId+7, Pawn, what, NewType(1, 0, 0, 0)), 0}
					size++
				}
				if File(fromId) < FILE_H && (SquareMask[fromId+9]&theirOccupation) != 0 {
					what = pos.TypeOnSquare(SquareMask[fromId+9])
					buffer[size] = EvaledMove{NewMove(fromId, fromId+9, Pawn, what, NewType(1, 0, 0, 0)), 0}
					size++
				}
			}
		}
	} else {
		if pos.EpSquare != 0 {
			epBB := (SquareMask[uint(pos.EpSquare)-1] | SquareMask[uint(pos.EpSquare)] | SquareMask[uint(pos.EpSquare)+1]) & RANK_4_BB
			for fromBB = epBB & pos.Pieces[Pawn] & pos.Colours[Black]; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[size] = EvaledMove{NewMove(fromId, pos.EpSquare-8, Pawn, Pawn, NewType(1, 0, 0, 1)), 0}
				size++
			}
		}
		for fromBB = (WhitePawnsAttacks(theirOccupation) | RANK_2_BB) & pos.Pieces[Pawn] & pos.Colours[Black]; fromBB != 0; fromBB &= fromBB - 1 {
			fromId = BitScan(fromBB)
			if Rank(fromId) == RANK_2 {
				if SquareMask[fromId-8]&allOccupation == 0 {
					buffer[size] = EvaledMove{NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 1, 1)), 0}
					size++
				}
				if File(fromId) > FILE_A && (SquareMask[fromId-9]&theirOccupation) != 0 {
					what = pos.TypeOnSquare(SquareMask[fromId-9])
					buffer[size] = EvaledMove{NewMove(fromId, fromId-9, Pawn, what, NewType(1, 1, 1, 1)), 0}
					size++
				}
				if File(fromId) < FILE_H && (SquareMask[fromId-7]&theirOccupation) != 0 {
					what = pos.TypeOnSquare(SquareMask[fromId-7])
					buffer[size] = EvaledMove{NewMove(fromId, fromId-7, Pawn, what, NewType(1, 1, 1, 1)), 0}
					size++
				}
			} else {
				if File(fromId) > FILE_A && (SquareMask[fromId-9]&theirOccupation) != 0 {
					what = pos.TypeOnSquare(SquareMask[fromId-9])
					buffer[size] = EvaledMove{NewMove(fromId, fromId-9, Pawn, what, NewType(1, 0, 0, 0)), 0}
					size++
				}
				if File(fromId) < FILE_H && (SquareMask[fromId-7]&theirOccupation) != 0 {
					what = pos.TypeOnSquare(SquareMask[fromId-7])
					buffer[size] = EvaledMove{NewMove(fromId, fromId-7, Pawn, what, NewType(1, 0, 0, 0)), 0}
					size++
				}
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
			buffer[size] = EvaledMove{NewMove(fromId, toId, Knight, what, NewType(1, 0, 0, 0)), 0}
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
			buffer[size] = EvaledMove{NewMove(fromId, toId, Bishop, what, NewType(1, 0, 0, 0)), 0}
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
			buffer[size] = EvaledMove{NewMove(fromId, toId, Rook, what, NewType(1, 0, 0, 0)), 0}
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
			buffer[size] = EvaledMove{NewMove(fromId, toId, Queen, what, NewType(1, 0, 0, 0)), 0}
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
		buffer[size] = EvaledMove{NewMove(fromId, toId, King, what, NewType(1, 0, 0, 0)), 0}
		size++
	}
	// end of Kings

	return buffer[:size]
}
