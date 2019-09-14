package backend

type EvaledMove struct {
	Move
	Value int
}

func (pos *Position) GenerateAllMoves(buffer []EvaledMove) []EvaledMove {
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
					buffer[counter].Move = NewMove(fromId, fromId+8, PromotionMove, PromoteToQueen)
					counter++
					buffer[counter].Move = NewMove(fromId, fromId+8, PromotionMove, PromoteToRook)
					counter++
					buffer[counter].Move = NewMove(fromId, fromId+8, PromotionMove, PromoteToBishop)
					counter++
					buffer[counter].Move = NewMove(fromId, fromId+8, PromotionMove, PromoteToKnight)
					counter++
				}
				for toBB = WhitePawnAttacks[fromId] & pos.Black; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					to = SquareMask[uint(toId)]
					buffer[counter].Move = NewMove(fromId, toId, PromotionMove, PromoteToQueen)
					counter++
					buffer[counter].Move = NewMove(fromId, toId, PromotionMove, PromoteToRook)
					counter++
					buffer[counter].Move = NewMove(fromId, toId, PromotionMove, PromoteToBishop)
					counter++
					buffer[counter].Move = NewMove(fromId, toId, PromotionMove, PromoteToKnight)
					counter++
				}
			} else {
				to = from << 8
				if allOccupation&to == 0 {
					buffer[counter].Move = NewMove(fromId, fromId+8, NormalMove, QuietMove)
					counter++
					if from&RANK_2_BB != 0 && allOccupation&(from<<16) == 0 {
						buffer[counter].Move = NewMove(fromId, fromId+16, NormalMove, QuietMove)
						counter++
					}
				}
				for toBB = WhitePawnAttacks[fromId] & pos.Black; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					to = SquareMask[toId]
					buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
					counter++
				}
			}
		}
		if pos.EpSquare != 0 {
			epBB := (SquareMask[uint(pos.EpSquare)-1] | SquareMask[uint(pos.EpSquare)] | SquareMask[uint(pos.EpSquare)+1]) & RANK_5_BB
			for fromBB = epBB & pos.Pawns & pos.White; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[counter].Move = NewMove(fromId, pos.EpSquare+8, EnpassMove, CaptureMove)
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
					buffer[counter].Move = NewMove(fromId, fromId-8, PromotionMove, PromoteToQueen)
					counter++
					buffer[counter].Move = NewMove(fromId, fromId-8, PromotionMove, PromoteToRook)
					counter++
					buffer[counter].Move = NewMove(fromId, fromId-8, PromotionMove, PromoteToBishop)
					counter++
					buffer[counter].Move = NewMove(fromId, fromId-8, PromotionMove, PromoteToKnight)
					counter++
				}
				for toBB = BlackPawnAttacks[fromId] & pos.White; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					to = SquareMask[uint(toId)]
					buffer[counter].Move = NewMove(fromId, toId, PromotionMove, PromoteToQueen)
					counter++
					buffer[counter].Move = NewMove(fromId, toId, PromotionMove, PromoteToRook)
					counter++
					buffer[counter].Move = NewMove(fromId, toId, PromotionMove, PromoteToBishop)
					counter++
					buffer[counter].Move = NewMove(fromId, toId, PromotionMove, PromoteToKnight)
					counter++
				}
			} else {
				to = SquareMask[fromId-8]
				if allOccupation&to == 0 {
					buffer[counter].Move = NewMove(fromId, fromId-8, NormalMove, QuietMove)
					counter++
					if from&RANK_7_BB != 0 && allOccupation&(from>>16) == 0 {
						buffer[counter].Move = NewMove(fromId, fromId-16, NormalMove, QuietMove)
						counter++
					}
				}
				for toBB = BlackPawnAttacks[fromId] & pos.White; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					to = SquareMask[uint(toId)]
					buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
					counter++
				}
			}
		}
		if pos.EpSquare != 0 {
			epBB := (SquareMask[uint(pos.EpSquare)-1] | SquareMask[uint(pos.EpSquare)] | SquareMask[uint(pos.EpSquare)+1]) & RANK_4_BB
			for fromBB = epBB & pos.Pawns & pos.Black; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[counter].Move = NewMove(fromId, pos.EpSquare-8, EnpassMove, CaptureMove)
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
				buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
				counter++
			} else {
				buffer[counter].Move = NewMove(fromId, toId, NormalMove, QuietMove)
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
			buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
			counter++
		} else {
			buffer[counter].Move = NewMove(fromId, toId, NormalMove, QuietMove)
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
				buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
				counter++
			} else {
				buffer[counter].Move = NewMove(fromId, toId, NormalMove, QuietMove)
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
				buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
				counter++
			} else {
				buffer[counter].Move = NewMove(fromId, toId, NormalMove, QuietMove)
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
				buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
				counter++
			} else {
				buffer[counter].Move = NewMove(fromId, toId, NormalMove, QuietMove)
				counter++
			}
		}
	}
	// end of Queens

	// Castling
	if pos.WhiteMove {
		if allOccupation&WHITE_KING_CASTLE_BLOCK_BB == 0 && pos.Flags&WhiteKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E1, false) && !pos.IsSquareAttacked(F1, false) {
			buffer[counter].Move = WhiteKingSideCastle
			counter++
		}
		if allOccupation&WHITE_QUEEN_CASTLE_BLOCK_BB == 0 && pos.Flags&WhiteQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E1, false) && !pos.IsSquareAttacked(D1, false) {
			buffer[counter].Move = WhiteQueenSideCastle
			counter++
		}
	} else {
		if allOccupation&BLACK_KING_CASTLE_BLOCK_BB == 0 && pos.Flags&BlackKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E8, true) && !pos.IsSquareAttacked(F8, true) {
			buffer[counter].Move = BlackKingSideCastle
			counter++
		}
		if allOccupation&BLACK_QUEEN_CASTLE_BLOCK_BB == 0 && pos.Flags&BlackQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E8, true) && !pos.IsSquareAttacked(D8, true) {
			buffer[counter].Move = BlackQueenSideCastle
			counter++
		}

	}
	// end of Castling

	return buffer[:counter]
}

func (pos *Position) GenerateAllCaptures(buffer []EvaledMove) []EvaledMove {
	var fromBB, toBB, ourOccupation, theirOccupation uint64
	var fromId, toId int

	allOccupation := pos.White | pos.Black

	if pos.WhiteMove {
		ourOccupation = pos.White
		theirOccupation = pos.Black
	} else {
		ourOccupation = pos.Black
		theirOccupation = pos.White
	}

	var counter = 0

	// PAWNS
	if pos.WhiteMove {
		if pos.EpSquare != 0 {
			epBB := (SquareMask[uint(pos.EpSquare)-1] | SquareMask[uint(pos.EpSquare)] | SquareMask[uint(pos.EpSquare)+1]) & RANK_5_BB
			for fromBB = epBB & pos.Pawns & pos.White; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[counter].Move = NewMove(fromId, pos.EpSquare+8, EnpassMove, CaptureMove)
				counter++
			}
		}
		for fromBB = (BlackPawnsAttacks(theirOccupation) | RANK_7_BB) & pos.Pawns & pos.White; fromBB != 0; fromBB &= fromBB - 1 {
			fromId = BitScan(fromBB)
			if Rank(fromId) == RANK_7 {
				if SquareMask[fromId+8]&allOccupation == 0 {
					buffer[counter].Move = NewMove(fromId, fromId+8, PromotionMove, PromoteToQueen)
					counter++
				}
				if File(fromId) > FILE_A && (SquareMask[fromId+7]&theirOccupation) != 0 {
					buffer[counter].Move = NewMove(fromId, fromId+7, PromotionMove, PromoteToQueen)
					counter++
				}
				if File(fromId) < FILE_H && (SquareMask[fromId+9]&theirOccupation) != 0 {
					buffer[counter].Move = NewMove(fromId, fromId+9, PromotionMove, PromoteToQueen)
					counter++
				}
			} else {
				if File(fromId) > FILE_A && (SquareMask[fromId+7]&theirOccupation) != 0 {
					buffer[counter].Move = NewMove(fromId, fromId+7, NormalMove, CaptureMove)
					counter++
				}
				if File(fromId) < FILE_H && (SquareMask[fromId+9]&theirOccupation) != 0 {
					buffer[counter].Move = NewMove(fromId, fromId+9, NormalMove, CaptureMove)
					counter++
				}
			}
		}
	} else {
		if pos.EpSquare != 0 {
			epBB := (SquareMask[uint(pos.EpSquare)-1] | SquareMask[uint(pos.EpSquare)] | SquareMask[uint(pos.EpSquare)+1]) & RANK_4_BB
			for fromBB = epBB & pos.Pawns & pos.Black; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[counter].Move = NewMove(fromId, pos.EpSquare-8, EnpassMove, CaptureMove)
				counter++
			}
		}
		for fromBB = (WhitePawnsAttacks(theirOccupation) | RANK_2_BB) & pos.Pawns & pos.Black; fromBB != 0; fromBB &= fromBB - 1 {
			fromId = BitScan(fromBB)
			if Rank(fromId) == RANK_2 {
				if SquareMask[fromId-8]&allOccupation == 0 {
					buffer[counter].Move = NewMove(fromId, fromId-8, PromotionMove, PromoteToQueen)
					counter++
				}
				if File(fromId) > FILE_A && (SquareMask[fromId-9]&theirOccupation) != 0 {
					buffer[counter].Move = NewMove(fromId, fromId-9, PromotionMove, PromoteToQueen)
					counter++
				}
				if File(fromId) < FILE_H && (SquareMask[fromId-7]&theirOccupation) != 0 {
					buffer[counter].Move = NewMove(fromId, fromId-7, PromotionMove, PromoteToQueen)
					counter++
				}
			} else {
				if File(fromId) > FILE_A && (SquareMask[fromId-9]&theirOccupation) != 0 {
					buffer[counter].Move = NewMove(fromId, fromId-9, NormalMove, CaptureMove)
					counter++
				}
				if File(fromId) < FILE_H && (SquareMask[fromId-7]&theirOccupation) != 0 {
					buffer[counter].Move = NewMove(fromId, fromId-7, NormalMove, CaptureMove)
					counter++
				}
			}
		}
	}
	// end of pawns

	// Knights
	for fromBB = pos.Knights & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = KnightAttacks[fromId] & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
			counter++
		}
	}
	// end of knights

	// Bishops
	for fromBB = pos.Bishops & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = BishopAttacks(fromId, allOccupation) & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
			counter++
		}
	}
	// end of Bishops

	// Rooks
	for fromBB = pos.Rooks & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = RookAttacks(fromId, allOccupation) & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
			counter++
		}
	}
	// end of Rooks

	// Queens
	for fromBB = pos.Queens & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = QueenAttacks(fromId, allOccupation) & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
			counter++
		}
	}
	// end of Queens

	// Kings
	fromBB = pos.Kings & ourOccupation
	fromId = BitScan(fromBB)
	for toBB = KingAttacks[fromId] & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
		toId = BitScan(toBB)
		buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
		counter++
	}
	// end of Kings

	return buffer[:counter]
}
