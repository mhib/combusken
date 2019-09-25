package backend

type EvaledMove struct {
	Move
	Value int
}

func (pos *Position) GenerateAllMoves(buffer []EvaledMove) []EvaledMove {
	var counter = 0
	ourOccupation := pos.Colours[pos.SideToMove]
	theirOccupation := pos.Colours[pos.SideToMove^1]
	allOccupation := ourOccupation | theirOccupation
	var fromBB, from, toBB, to uint64
	var fromId, toId int
	if pos.SideToMove == White {
		for fromBB = pos.Pieces[Pawn] & ourOccupation; fromBB > 0; fromBB &= (fromBB - 1) {
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
				for toBB = PawnAttacks[White][fromId] & theirOccupation; toBB > 0; toBB &= (toBB - 1) {
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
				for toBB = PawnAttacks[White][fromId] & theirOccupation; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					to = SquareMask[toId]
					buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
					counter++
				}
			}
		}
		if pos.EpSquare != 0 {
			epBB := (SquareMask[uint(pos.EpSquare)-1] | SquareMask[uint(pos.EpSquare)] | SquareMask[uint(pos.EpSquare)+1]) & RANK_5_BB
			for fromBB = epBB & pos.Pieces[Pawn] & ourOccupation; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[counter].Move = NewMove(fromId, pos.EpSquare+8, EnpassMove, CaptureMove)
				counter++
			}
		}
	} else {
		for fromBB = pos.Pieces[Pawn] & ourOccupation; fromBB > 0; fromBB &= (fromBB - 1) {
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
				for toBB = PawnAttacks[Black][fromId] & theirOccupation; toBB > 0; toBB &= (toBB - 1) {
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
				for toBB = PawnAttacks[Black][fromId] & theirOccupation; toBB > 0; toBB &= (toBB - 1) {
					toId = BitScan(toBB)
					to = SquareMask[uint(toId)]
					buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
					counter++
				}
			}
		}
		if pos.EpSquare != 0 {
			epBB := (SquareMask[uint(pos.EpSquare)-1] | SquareMask[uint(pos.EpSquare)] | SquareMask[uint(pos.EpSquare)+1]) & RANK_4_BB
			for fromBB = epBB & pos.Pieces[Pawn] & ourOccupation; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[counter].Move = NewMove(fromId, pos.EpSquare-8, EnpassMove, CaptureMove)
				counter++
			}
		}
	}

	// Knights
	for fromBB = pos.Pieces[Knight] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
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
	from = pos.Pieces[King] & ourOccupation
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
	for fromBB = pos.Pieces[Rook] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
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
	for fromBB = pos.Pieces[Bishop] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
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
	for fromBB = pos.Pieces[Queen] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
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
	if pos.SideToMove == White {
		if allOccupation&WHITE_KING_CASTLE_BLOCK_BB == 0 && pos.Flags&WhiteKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E1, Black) && !pos.IsSquareAttacked(F1, Black) {
			buffer[counter].Move = WhiteKingSideCastle
			counter++
		}
		if allOccupation&WHITE_QUEEN_CASTLE_BLOCK_BB == 0 && pos.Flags&WhiteQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E1, Black) && !pos.IsSquareAttacked(D1, Black) {
			buffer[counter].Move = WhiteQueenSideCastle
			counter++
		}
	} else {
		if allOccupation&BLACK_KING_CASTLE_BLOCK_BB == 0 && pos.Flags&BlackKingSideCastleFlag == 0 && !pos.IsSquareAttacked(E8, White) && !pos.IsSquareAttacked(F8, White) {
			buffer[counter].Move = BlackKingSideCastle
			counter++
		}
		if allOccupation&BLACK_QUEEN_CASTLE_BLOCK_BB == 0 && pos.Flags&BlackQueenSideCastleFlag == 0 && !pos.IsSquareAttacked(E8, White) && !pos.IsSquareAttacked(D8, White) {
			buffer[counter].Move = BlackQueenSideCastle
			counter++
		}

	}
	// end of Castling

	return buffer[:counter]
}

func (pos *Position) GenerateAllCaptures(buffer []EvaledMove) []EvaledMove {
	var fromBB, toBB uint64
	var fromId, toId int

	ourOccupation := pos.Colours[pos.SideToMove]
	theirOccupation := pos.Colours[pos.SideToMove^1]
	allOccupation := ourOccupation | theirOccupation

	var counter = 0

	// PAWNS
	if pos.SideToMove == White {
		if pos.EpSquare != 0 {
			epBB := (SquareMask[uint(pos.EpSquare)-1] | SquareMask[uint(pos.EpSquare)] | SquareMask[uint(pos.EpSquare)+1]) & RANK_5_BB
			for fromBB = epBB & pos.Pieces[Pawn] & ourOccupation; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[counter].Move = NewMove(fromId, pos.EpSquare+8, EnpassMove, CaptureMove)
				counter++
			}
		}
		for fromBB = (BlackPawnsAttacks(theirOccupation) | RANK_7_BB) & pos.Pieces[Pawn] & ourOccupation; fromBB != 0; fromBB &= fromBB - 1 {
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
			for fromBB = epBB & pos.Pieces[Pawn] & ourOccupation; fromBB > 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				buffer[counter].Move = NewMove(fromId, pos.EpSquare-8, EnpassMove, CaptureMove)
				counter++
			}
		}
		for fromBB = (WhitePawnsAttacks(theirOccupation) | RANK_2_BB) & pos.Pieces[Pawn] & ourOccupation; fromBB != 0; fromBB &= fromBB - 1 {
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
	for fromBB = pos.Pieces[Knight] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = KnightAttacks[fromId] & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
			counter++
		}
	}
	// end of knights

	// Bishops
	for fromBB = pos.Pieces[Bishop] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = BishopAttacks(fromId, allOccupation) & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
			counter++
		}
	}
	// end of Bishops

	// Rooks
	for fromBB = pos.Pieces[Rook] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = RookAttacks(fromId, allOccupation) & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
			counter++
		}
	}
	// end of Rooks

	// Queens
	for fromBB = pos.Pieces[Queen] & ourOccupation; fromBB != 0; fromBB &= (fromBB - 1) {
		fromId = BitScan(fromBB)
		for toBB = QueenAttacks(fromId, allOccupation) & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
			toId = BitScan(toBB)
			buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
			counter++
		}
	}
	// end of Queens

	// Kings
	fromBB = pos.Pieces[King] & ourOccupation
	fromId = BitScan(fromBB)
	for toBB = KingAttacks[fromId] & theirOccupation; toBB != 0; toBB &= (toBB - 1) {
		toId = BitScan(toBB)
		buffer[counter].Move = NewMove(fromId, toId, NormalMove, CaptureMove)
		counter++
	}
	// end of Kings

	return buffer[:counter]
}
