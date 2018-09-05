package backend

import (
	"fmt"
	"math/bits"
)

const (
	FILE_A = iota
	FILE_B
	FILE_C
	FILE_D
	FILE_E
	FILE_F
	FILE_G
	FILE_H
)

const (
	RANK_1 = iota
	RANK_2
	RANK_3
	RANK_4
	RANK_5
	RANK_6
	RANK_7
	RANK_8
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

	WHITE_KING_CASTLE_BLOCK_BB  uint64 = F1_MASK | G1_MASK
	WHITE_QUEEN_CASTLE_BLOCK_BB uint64 = B1_MASK | C1_MASK | D1_MASK
	BLACK_KING_CASTLE_BLOCK_BB  uint64 = F8_MASK | G8_MASK
	BLACK_QUEEN_CASTLE_BLOCK_BB uint64 = B8_MASK | C8_MASK | D8_MASK
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

func file(id int) int {

	return id & 7
}

func rank(id int) int {
	return id >> 3
}

const (
	A1 = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1

	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2

	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3

	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4

	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5

	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6

	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7

	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8

	NoSquare
)

const (
	A1_MASK = 1 << iota
	B1_MASK
	C1_MASK
	D1_MASK
	E1_MASK
	F1_MASK
	G1_MASK
	H1_MASK

	A2_MASK
	B2_MASK
	C2_MASK
	D2_MASK
	E2_MASK
	F2_MASK
	G2_MASK
	H2_MASK

	A3_MASK
	B3_MASK
	C3_MASK
	D3_MASK
	E3_MASK
	F3_MASK
	G3_MASK
	H3_MASK

	A4_MASK
	B4_MASK
	C4_MASK
	D4_MASK
	E4_MASK
	F4_MASK
	G4_MASK
	H4_MASK

	A5_MASK
	B5_MASK
	C5_MASK
	D5_MASK
	E5_MASK
	F5_MASK
	G5_MASK
	H5_MASK

	A6_MASK
	B6_MASK
	C6_MASK
	D6_MASK
	E6_MASK
	F6_MASK
	G6_MASK
	H6_MASK

	A7_MASK
	B7_MASK
	C7_MASK
	D7_MASK
	E7_MASK
	F7_MASK
	G7_MASK
	H7_MASK

	A8_MASK
	B8_MASK
	C8_MASK
	D8_MASK
	E8_MASK
	F8_MASK
	G8_MASK
	H8_MASK
)

var SquareMask [64]uint64
var SquareString [64]string

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

func MoreThanOne(bb uint64) bool {
	return bb != 0 && ((bb-1)&bb) != 0
}

func NorthWest(set uint64) uint64 {
	return set << 7
}

func North(set uint64) uint64 {
	return set << 8
}

func NorthEast(set uint64) uint64 {
	return set << 9
}

func East(set uint64) uint64 {
	return set << 1
}

func West(set uint64) uint64 {
	return set >> 1
}

func SouthWest(set uint64) uint64 {
	return set >> 9
}

func South(set uint64) uint64 {
	return set >> 8
}

func SouthEast(set uint64) uint64 {
	return set >> 7
}

func initArray(array *[64]uint64, method func(mask uint64) uint64) {
	for i := uint32(0); i <= 63; i++ {
		array[i] = method(uint64(1) << i)
	}
}

func kingAttacks(set uint64) uint64 {
	return NorthEast(set & ^RANK_8_BB & ^FILE_H_BB) | North(set & ^RANK_8_BB) |
		NorthWest(set & ^RANK_8_BB & ^FILE_A_BB) | East(set & ^FILE_H_BB) | West(set & ^FILE_A_BB) |
		SouthEast(set & ^RANK_1_BB & ^FILE_H_BB) | South(set & ^RANK_1_BB) | SouthWest(set & ^RANK_1_BB & ^FILE_A_BB)
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
	return NorthWest(set & ^FILE_A_BB) |
		NorthEast(set & ^FILE_H_BB)
}

func WhitePawnsMoves(set uint64, occupancy uint64) uint64 {
	return (North(set) | ((North(North(set)) & RANK_4_BB) & ^North(occupancy&RANK_3_BB))) & ^occupancy
}

func BlackPawnsAttacks(set uint64) uint64 {
	return SouthWest(set & ^FILE_A_BB) |
		SouthEast(set & ^FILE_H_BB)
}

func BlackPawnsMoves(set uint64, occupancy uint64) uint64 {
	return (South(set) | ((South(South(set)) & RANK_5_BB) & ^South(occupancy&RANK_6_BB))) & ^occupancy
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

func squareString(square int) string {
	var res string
	res += string(byte(int('a') + (square & 7)))
	res += string(byte(int('1') + square/8))
	return res
}

func init() {
	initSquareMask()
	initArray(&KingAttacks, kingAttacks)
	initArray(&KnightAttacks, KnightsAttacks)
	initArray(&WhitePawnAttacks, WhitePawnsAttacks)
	initArray(&BlackPawnAttacks, BlackPawnsAttacks)
	for i := 0; i < 64; i++ {
		SquareString[i] = squareString(i)
	}

	initArray(&rookBlockerMask, generateRookBlockerMask)
	rookBlockerBoard := initRookBlockerBoard()
	initRookMoveBoard(rookBlockerBoard)
	initRookMagicIndex(rookBlockerBoard)
	initRookAttacks(rookBlockerBoard)

	initArray(&bishopBlockerMask, generateBishopBlockerMask)
	bishopBlockerBoard := initBishopBlockerBoard()
	initBishopMoveBoard(bishopBlockerBoard)
	initBishopMagicIndex(bishopBlockerBoard)
	initBishopAttacks(bishopBlockerBoard)
}
