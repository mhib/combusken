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

	WHITE_SQUARES uint64 = 0x55AA55AA55AA55AA
	BLACK_SQUARES uint64 = 0xAA55AA55AA55AA55

	WHITE_KING_CASTLE_BLOCK_BB  uint64 = F1_BB | G1_BB
	WHITE_QUEEN_CASTLE_BLOCK_BB uint64 = B1_BB | C1_BB | D1_BB
	BLACK_KING_CASTLE_BLOCK_BB  uint64 = F8_BB | G8_BB
	BLACK_QUEEN_CASTLE_BLOCK_BB uint64 = B8_BB | C8_BB | D8_BB

	PROMOTION_RANKS uint64 = RANK_1_BB | RANK_8_BB
	CENTER          uint64 = (FILE_D_BB | FILE_E_BB) & (RANK_4_BB | RANK_5_BB)
	LONG_DIAGONALS  uint64 = 0x8142241818244281
)

var RANKS = [...]uint64{RANK_1_BB, RANK_2_BB, RANK_3_BB, RANK_4_BB, RANK_5_BB, RANK_6_BB, RANK_7_BB, RANK_8_BB}
var FILES = [...]uint64{FILE_A_BB, FILE_B_BB, FILE_C_BB, FILE_D_BB, FILE_E_BB, FILE_F_BB, FILE_G_BB, FILE_H_BB}

func File(id int) int {
	return id & 7
}

func FileBB(id int) uint64 {
	return FILES[File(id)]
}

func Rank(id int) int {
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
	A1_BB = 1 << iota
	B1_BB
	C1_BB
	D1_BB
	E1_BB
	F1_BB
	G1_BB
	H1_BB

	A2_BB
	B2_BB
	C2_BB
	D2_BB
	E2_BB
	F2_BB
	G2_BB
	H2_BB

	A3_BB
	B3_BB
	C3_BB
	D3_BB
	E3_BB
	F3_BB
	G3_BB
	H3_BB

	A4_BB
	B4_BB
	C4_BB
	D4_BB
	E4_BB
	F4_BB
	G4_BB
	H4_BB

	A5_BB
	B5_BB
	C5_BB
	D5_BB
	E5_BB
	F5_BB
	G5_BB
	H5_BB

	A6_BB
	B6_BB
	C6_BB
	D6_BB
	E6_BB
	F6_BB
	G6_BB
	H6_BB

	A7_BB
	B7_BB
	C7_BB
	D7_BB
	E7_BB
	F7_BB
	G7_BB
	H7_BB

	A8_BB
	B8_BB
	C8_BB
	D8_BB
	E8_BB
	F8_BB
	G8_BB
	H8_BB
)

var SquareMask [64]uint64
var SquareString [64]string

func PopCount(set uint64) int {
	return bits.OnesCount64(set)
}

var (
	PawnAttacks                [2][64]uint64
	KnightAttacks, KingAttacks [64]uint64
)

// Least significant bit
func BitScan(bb uint64) int {
	return bits.TrailingZeros64(bb)
}

func MostSignificantBit(bb uint64) int {
	return bits.LeadingZeros64(bb)
}

func MoreThanOne(bb uint64) bool {
	return bb != 0 && ((bb-1)&bb) != 0
}

func OnlyOne(bb uint64) bool {
	return bb != 0 && ((bb-1)&bb) == 0
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

var FileMirror = [8]int{FILE_A, FILE_B, FILE_C, FILE_D, FILE_D, FILE_C, FILE_B, FILE_A}

func initArray(array *[64]uint64, method func(mask uint64) uint64) {
	for i := uint32(0); i <= 63; i++ {
		array[i] = method(uint64(1) << i)
	}
}

func KingsAttacks(set uint64) uint64 {
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

func squareString(square int) (res string) {
	res += string(byte(int('a') + File(square)))
	res += string(byte(int('1') + Rank(square)))
	return
}

func init() {
	for i := uint(0); i < 64; i++ {
		SquareMask[i] = 1 << i
	}
	initArray(&KingAttacks, KingsAttacks)
	initArray(&KnightAttacks, KnightsAttacks)
	initArray(&PawnAttacks[White], WhitePawnsAttacks)
	initArray(&PawnAttacks[Black], BlackPawnsAttacks)
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
