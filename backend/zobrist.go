package backend

import (
	"math/rand"
)

var zobrist [6][2][64]uint64
var zobristEpSquare [64]uint64
var zobristFlags [16]uint64
var zobristColor uint64

func initZobrist() {
	var r = rand.New(rand.NewSource(0))
	for y := Pawn; y <= King; y++ {
		for x := Black; x <= White; x++ {
			for z := A1; z <= H8; z++ {
				zobrist[y][x][z] = r.Uint64()
			}
		}
	}
	for y := A4; y <= H5; y++ {
		zobristEpSquare[y] = r.Uint64()
	}
	for y := 0; y < 16; y++ {
		zobristFlags[y] = r.Uint64()
	}
	zobristColor = r.Uint64()
}

func HashPosition(pos *Position) {
	pos.Key = 0
	pos.PawnKey = 0
	var fromId int
	var fromBB uint64

	for piece := Pawn; piece <= King; piece++ {
		for colour := Black; colour <= White; colour++ {
			for fromBB = pos.Pieces[piece] & pos.Colours[colour]; fromBB != 0; fromBB &= (fromBB - 1) {
				fromId = BitScan(fromBB)
				pos.Key ^= zobrist[piece][colour][fromId]
				if piece == Pawn || piece == King {
					pos.PawnKey ^= zobrist[piece][colour][fromId]
				}
			}
		}
	}

	pos.Key ^= zobristFlags[pos.Flags]
	if pos.SideToMove == White {
		pos.Key ^= zobristColor
		pos.PawnKey ^= zobristColor
	}
	pos.Key ^= zobristEpSquare[pos.EpSquare]
}

func init() {
	initZobrist()
}
