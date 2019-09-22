package evaluation

import . "github.com/mhib/combusken/backend"

var SEEValues = []int{0, 100, 450, 450, 675, 1300, Mate / 2}

// Returns true if see non-negative
func SeeSign(pos *Position, move Move) bool {
	return SeeAbove(pos, move, 0)
}

// based on laser implementation
func SeeAbove(pos *Position, move Move, cutoff int) bool {
	// Special case for ep and castling
	if move.Type() == EnpassMove || move.IsCastling() {
		return cutoff <= 0
	}
	lastPiece := pos.TypeOnSquare(move.From())
	capturedValue := SEEValues[pos.TypeOnSquare(move.To())]
	if move.IsPromotion() {
		lastPiece = move.PromotedPiece()
		capturedValue += SEEValues[lastPiece] - SEEValues[Pawn]
	}
	value := capturedValue - cutoff
	// return when free piece is not enough
	if value < 0 {
		return false
	}

	value -= SEEValues[lastPiece]
	// return when after recapture it is still good enough
	if value >= 0 {
		return true
	}

	to := move.To()
	occ := (pos.White ^ pos.Black ^ SquareMask[move.From()]) | SquareMask[to]
	side := !pos.WhiteMove
	for {
		nextVictim, from := getLeastValuableAttacker(pos, to, side, occ)
		if nextVictim == None {
			break
		}
		// Last capture with king was illegal, as there were opposide side attackers
		if lastPiece == King {
			side = !side
			break
		}
		occ ^= SquareMask[from]
		side = !side
		value = -value - 1 - SEEValues[nextVictim]
		lastPiece = nextVictim
		// lastPiece belonged to `side`
		// if after the recapture of lastPiece opponents score is positive then `side` loses
		if value >= 0 {
			break
		}
	}
	return side != pos.WhiteMove
}

func getLeastValuableAttacker(pos *Position, to int, side bool, occupancy uint64) (piece, from int) {
	from = NoSquare
	var sideOccupancy uint64

	if side {
		sideOccupancy = occupancy & pos.White
		if attacks := BlackPawnAttacks[to] & (pos.Pawns & sideOccupancy); attacks != 0 {
			return Pawn, BitScan(attacks)
		}
	} else {
		sideOccupancy = occupancy & pos.Black
		if attacks := WhitePawnAttacks[to] & (pos.Pawns & sideOccupancy); attacks != 0 {
			return Pawn, BitScan(attacks)
		}
	}

	if attacks := KnightAttacks[to] & (sideOccupancy & pos.Knights); attacks != 0 {
		return Knight, BitScan(attacks)
	}

	if attacks := BishopAttacks(to, occupancy) & (sideOccupancy & pos.Bishops); attacks != 0 {
		return Bishop, BitScan(attacks)
	}

	if attacks := RookAttacks(to, occupancy) & (sideOccupancy & pos.Rooks); attacks != 0 {
		return Rook, BitScan(attacks)
	}

	if attacks := QueenAttacks(to, occupancy) & (sideOccupancy & pos.Queens); attacks != 0 {
		return Queen, BitScan(attacks)
	}

	if attacks := KingAttacks[to] & (sideOccupancy & pos.Kings); attacks != 0 {
		return King, BitScan(attacks)
	}

	return
}
