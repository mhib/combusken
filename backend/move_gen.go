package backend

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
					buffer[counter] = NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 1, 1))
					counter++
					buffer[counter] = NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 1, 0))
					counter++
					buffer[counter] = NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 0, 1))
					counter++
					buffer[counter] = NewMove(fromId, fromId+8, Pawn, None, NewType(0, 1, 0, 0))
					counter++
				}
				for toBB = WhitePawnAttacks[fromId] & pos.Black; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					to = SquareMask[uint(toId)]
					captureType := pos.TypeOnSquare(to)
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 1, 1))
					counter++
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 1, 0))
					counter++
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 0, 1))
					counter++
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 0, 0))
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
			epBB := (West(epSquareBB) | East(epSquareBB)) & RANK_5_BB
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
					buffer[counter] = NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 1, 1))
					counter++
					buffer[counter] = NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 1, 0))
					counter++
					buffer[counter] = NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 0, 1))
					counter++
					buffer[counter] = NewMove(fromId, fromId-8, Pawn, None, NewType(0, 1, 0, 0))
					counter++
				}
				for toBB = BlackPawnAttacks[fromId] & pos.White; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					to = SquareMask[uint(toId)]
					captureType := pos.TypeOnSquare(to)
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 1, 1))
					counter++
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 1, 0))
					counter++
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 0, 1))
					counter++
					buffer[counter] = NewMove(fromId, toId, Pawn, captureType, NewType(1, 1, 0, 0))
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
			epBB := (West(epSquareBB) | East(epSquareBB)) & RANK_4_BB
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
		if allOccupation&BLACK_QUEEN_CASTLE_BLOCK_BB == 0 && pos.Flags&BlackQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E8_MASK, true) && !pos.IsSquareAttacked(D8_MASK, true) {
			buffer[counter] = BlackQueenSideCastle
			counter++
		}

	}
	// end of Castling

	return buffer[:counter]
}
