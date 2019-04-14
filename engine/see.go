package engine

import . "github.com/mhib/combusken/backend"

var SEEValues = []int{0, 100, 357, 377, 712, 12534, 20000} // From zurichess

// Returns true if see non-negative
func seeSign(pos *Position, move Move) bool {
	return seeAbove(pos, move, 0)
}

// based on laser implementation
func seeAbove(pos *Position, move Move, cutoff int) bool {
	if move.Type() == EPCapture || move.IsCastling() {
		return cutoff <= 0
	}
	capturedValue := SEEValues[move.CapturedPiece()]
	if move.IsPromotion() {
		capturedValue += SEEValues[move.Special()+Knight] - SEEValues[Pawn]
	}
	value := capturedValue - cutoff
	// return when free piece is not enough
	if value < 0 {
		return false
	}

	lastPiece := move.MovedPiece()
	if lastPiece == King {
		value -= Mate / 2
	} else {
		value -= SEEValues[lastPiece]
	}
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
		// lastPiece belonged to side
		// if after the recapture of lastPiece opponents score is positive then side loses
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
