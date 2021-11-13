package evaluation

import (
	. "github.com/mhib/combusken/chess"
	. "github.com/mhib/combusken/utils"
)

var SEEValues = [...]int{100, 450, 450, 675, 1300, Mate / 2, 0}

// Returns true if see non-negative
func SeeSign(pos *Position, move Move) bool {
	return SeeAbove(pos, move, 0)
}

// based on laser implementation
func SeeAbove(pos *Position, move Move, cutoff int) bool {
	// Special case for ep and castling
	if move.Type() == EPCapture || move.IsCastling() {
		return cutoff <= 0
	}
	lastPiece := move.MovedPiece()
	capturedValue := SEEValues[move.CapturedPiece()]
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
	occ := (pos.Colours[Black] ^ pos.Colours[White] ^ SquareMask[move.From()]) | SquareMask[to]
	side := pos.SideToMove ^ 1
	for {
		nextVictim, from := getLeastValuableAttacker(pos, to, side, occ)
		if nextVictim == None {
			break
		}
		// Last capture with king was illegal, as there were opposide side attackers
		if lastPiece == King {
			side ^= 1
			break
		}
		occ ^= SquareMask[from]
		side ^= 1
		value = -value - 1 - SEEValues[nextVictim]
		lastPiece = nextVictim
		// lastPiece belonged to `side`
		// if after the recapture of lastPiece opponents score is positive then `side` loses
		if value >= 0 {
			break
		}
	}
	return side != pos.SideToMove
}

func getLeastValuableAttacker(pos *Position, to int, side int, occupancy uint64) (int, int) {
	sideOccupancy := pos.Colours[side] & occupancy

	if attacks := PawnAttacks[side^1][to] & (pos.Pieces[Pawn] & sideOccupancy); attacks != 0 {
		return Pawn, BitScan(attacks)
	}

	if attacks := KnightAttacks[to] & (sideOccupancy & pos.Pieces[Knight]); attacks != 0 {
		return Knight, BitScan(attacks)
	}

	if attacks := BishopAttacks(to, occupancy) & (sideOccupancy & pos.Pieces[Bishop]); attacks != 0 {
		return Bishop, BitScan(attacks)
	}

	if attacks := RookAttacks(to, occupancy) & (sideOccupancy & pos.Pieces[Rook]); attacks != 0 {
		return Rook, BitScan(attacks)
	}

	if attacks := QueenAttacks(to, occupancy) & (sideOccupancy & pos.Pieces[Queen]); attacks != 0 {
		return Queen, BitScan(attacks)
	}

	if attacks := KingAttacks[to] & (sideOccupancy & pos.Pieces[King]); attacks != 0 {
		return King, BitScan(attacks)
	}

	return None, NoSquare
}
