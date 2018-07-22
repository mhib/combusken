package backend

import (
	"fmt"
	"math/bits"
)

const (
	RANK_1_BB uint64 = 0x00000000000000FF
	RANK_2_BB uint64 = 0x000000000000FF00
	RANK_3_BB uint64 = 0x0000000000FF0000
	RANK_4_BB uint64 = 0x00000000FF000000
	RANK_5_BB uint64 = 0x000000FF00000000
	RANK_6_BB uint64 = 0x0000FF0000000000
	RANK_7_BB uint64 = 0x00FF000000000000
	RANK_8_BB uint64 = 0xFF00000000000000

	FILE_A_BB uint64 = 0x0101010101010101
	FILE_B_BB uint64 = 0x0202020202020202
	FILE_C_BB uint64 = 0x0404040404040404
	FILE_D_BB uint64 = 0x0808080808080808
	FILE_E_BB uint64 = 0x1010101010101010
	FILE_F_BB uint64 = 0x2020202020202020
	FILE_G_BB uint64 = 0x4040404040404040
	FILE_H_BB uint64 = 0x8080808080808080
)

var RANKS = [...]uint64{RANK_1_BB, RANK_2_BB, RANK_3_BB, RANK_4_BB, RANK_5_BB, RANK_6_BB, RANK_7_BB, RANK_8_BB}
var FILES = [...]uint64{FILE_A_BB, FILE_B_BB, FILE_C_BB, FILE_D_BB, FILE_E_BB, FILE_F_BB, FILE_G_BB, FILE_H_BB}

func getRank(mask uint64) uint64 {
	for _, el := range RANKS {
		if mask&el != 0 {
			return el
		}
	}
	return 0
}

func getFile(mask uint64) uint64 {
	for _, el := range FILES {
		if mask&el != 0 {
			return el
		}
	}
	return 0
}

const (
	A1 = iota
	A2
	A3
	A4
	A5
	A6
	A7
	A8

	B1
	B2
	B3
	B4
	B5
	B6
	B7
	B8

	C1
	C2
	C3
	C4
	C5
	C6
	C7
	C8

	D1
	D2
	D3
	D4
	D5
	D6
	D7
	D8

	E1
	E2
	E3
	E4
	E5
	E6
	E7
	E8

	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8

	G1
	G2
	G3
	G4
	G5
	G6
	G7
	G8

	H1
	H2
	H3
	H4
	H5
	H6
	H7
	H8

	NoSquare
)

var SquareMask [64]uint64

func initSquareMask() {
	for i := uint(0); i < 64; i++ {
		SquareMask[i] = 1 << i
	}
}

func PopCount(set uint64) int {
	return bits.OnesCount64(set)
}

var (
	WhitePawnAttacks, BlackPawnAttacks [64]uint64
	KnightAttacks, KingAttacks         [64]uint64
)

var index64 = [64]int{0, 47, 1, 56, 48, 27, 2, 60,
	57, 49, 41, 37, 28, 16, 3, 61,
	54, 58, 35, 52, 50, 42, 21, 44,
	38, 32, 29, 23, 17, 11, 4, 62,
	46, 55, 26, 59, 40, 36, 15, 53,
	34, 51, 20, 43, 31, 22, 10, 45,
	25, 39, 14, 33, 19, 30, 9, 24,
	13, 18, 8, 12, 7, 6, 5, 63,
}

const debruijn64 uint64 = 0x03f79d71b4cb0a89

func BitScan(bb uint64) int {
	return index64[((bb^(bb-1))*debruijn64)>>58]
}

func northWest(set uint64) uint64 {
	return set << 7
}

func north(set uint64) uint64 {
	return set << 8
}

func northEast(set uint64) uint64 {
	return set << 9
}

func east(set uint64) uint64 {
	return set << 1
}

func west(set uint64) uint64 {
	return set >> 1
}

func southWest(set uint64) uint64 {
	return set >> 9
}

func south(set uint64) uint64 {
	return set >> 8
}

func southEast(set uint64) uint64 {
	return set >> 7
}

func initArray(array *[64]uint64, method func(mask uint64) uint64) {
	for i := uint32(0); i <= 63; i++ {
		array[i] = method(uint64(1) << i)
	}
}

func kingAttacks(set uint64) uint64 {
	return northEast(set & ^RANK_8_BB & ^FILE_H_BB) | north(set & ^RANK_8_BB) |
		northWest(set & ^RANK_8_BB & ^FILE_A_BB) | east(set & ^FILE_H_BB) | west(set & ^FILE_A_BB) |
		southEast(set & ^RANK_1_BB & ^FILE_H_BB) | south(set & ^RANK_1_BB) | southWest(set & ^RANK_1_BB & ^FILE_A_BB)
}

func InspectBB(bb uint64) {
	for y := 7; y >= 0; y-- {
		for x := 0; x <= 7; x++ {
			if bb&(uint64(1)<<uint64(8*y+x)) != 0 {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
			fmt.Print(" ")
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func KnightsAttacks(set uint64) uint64 {
	l1 := set >> 1 & 0x7f7f7f7f7f7f7f7f
	l2 := set >> 2 & 0x3f3f3f3f3f3f3f3f
	r1 := set << 1 & 0xfefefefefefefefe
	r2 := set << 2 & 0xfcfcfcfcfcfcfcfc
	h1 := l1 | r1
	h2 := l2 | r2
	return (h1 << 16) | (h1 >> 16) | (h2 << 8) | (h2 >> 8)
}

func WhitePawnsAttacks(set uint64) uint64 {
	return northWest(set & ^FILE_A_BB) |
		northEast(set & ^FILE_H_BB)
}

func WhitePawnsMoves(set uint64, occupancy uint64) uint64 {
	return (north(set) | ((north(north(set)) & RANK_4_BB) & ^north(occupancy&RANK_3_BB))) & ^occupancy
}

func BlackPawnsAttacks(set uint64) uint64 {
	return southWest(set & ^FILE_A_BB) |
		southEast(set & ^FILE_H_BB)
}

func BlackPawnsMoves(set uint64, occupancy uint64) uint64 {
	return (south(set) | ((south(south(set)) & RANK_5_BB) & ^south(occupancy&RANK_6_BB))) & ^occupancy
}

func QueenAttacks(square int, occupancy uint64) uint64 {
	return RookAttacks(square, occupancy) | BishopAttacks(square, occupancy)
}

func RookAttacks(square int, occupancy uint64) uint64 {
	return rookMoveBoard[square][(rookBlockerMask[square]&occupancy)*rookMagicIndex[square]>>rookShift]
}

func BishopAttacks(square int, occupancy uint64) uint64 {
	return bishopMoveBoard[square][(bishopBlockerMask[square]&occupancy)*bishopMagicIndex[square]>>bishopShift]
}

func RooksAttacks(set uint64, occupancy uint64) uint64 {
	res := uint64(0)
	for set > 0 {
		position := BitScan(set)
		set &= set - 1

		res |= RookAttacks(int(position), occupancy)
	}
	return res
}

func BishopsAttacks(set uint64, occupancy uint64) uint64 {
	res := uint64(0)
	for set > 0 {
		position := BitScan(set)
		set &= set - 1

		res |= BishopAttacks(int(position), occupancy)
	}
	return res
}

func QueensAttacks(set uint64, occupancy uint64) uint64 {
	res := uint64(0)
	for set > 0 {
		position := BitScan(set)
		set &= set - 1

		res |= QueenAttacks(int(position), occupancy)
	}
	return res
}

func InitBB() {
	initSquareMask()
	initArray(&KingAttacks, kingAttacks)
	initArray(&KnightAttacks, KnightsAttacks)
	initArray(&WhitePawnAttacks, WhitePawnsAttacks)
	initArray(&BlackPawnAttacks, BlackPawnsAttacks)

	initArray(&rookBlockerMask, generateRookBlockerMask)
	initRookBlockerBoard()
	initRookMoveBoard()
	initRookMagicIndex()
	initRookAttacks()

	initArray(&bishopBlockerMask, generateBishopBlockerMask)
	initBishopBlockerBoard()
	initBishopMoveBoard()
	initBishopMagicIndex()
	initBishopAttacks()
}
