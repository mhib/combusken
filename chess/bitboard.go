package chess

import (
	"fmt"
	"math/bits"
)

const (
	FileA = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

const (
	Rank1 = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)
const (
	Rank1_BB uint64 = 0x00000000000000FF
	Rank2_BB uint64 = 0x000000000000FF00
	Rank3_BB uint64 = 0x0000000000FF0000
	Rank4_BB uint64 = 0x00000000FF000000
	Rank5_BB uint64 = 0x000000FF00000000
	Rank6_BB uint64 = 0x0000FF0000000000
	Rank7_BB uint64 = 0x00FF000000000000
	Rank8_BB uint64 = 0xFF00000000000000

	FileA_BB uint64 = 0x0101010101010101
	FileB_BB uint64 = 0x0202020202020202
	FileC_BB uint64 = 0x0404040404040404
	FileD_BB uint64 = 0x0808080808080808
	FileE_BB uint64 = 0x1010101010101010
	FileF_BB uint64 = 0x2020202020202020
	FileG_BB uint64 = 0x4040404040404040
	FileH_BB uint64 = 0x8080808080808080

	WhiteSquares_BB uint64 = 0x55AA55AA55AA55AA
	BlackSquares_BB uint64 = 0xAA55AA55AA55AA55

	WhiteKingCastleBlock_BB  uint64 = F1_BB | G1_BB
	WhiteQueenCastleBlock_BB uint64 = B1_BB | C1_BB | D1_BB
	BlackKingCastleBlock_BB  uint64 = F8_BB | G8_BB
	BlackQueenCastleBlock_BB uint64 = B8_BB | C8_BB | D8_BB

	PromotionRanks_BB uint64 = Rank1_BB | Rank8_BB
	Center_BB         uint64 = (FileD_BB | FileE_BB) & (Rank4_BB | Rank5_BB)
	LongDiagonals_BB  uint64 = 0x8142241818244281

	QueenSide_BB   uint64 = FileA_BB | FileB_BB | FileC_BB | FileD_BB
	KingSide_BB    uint64 = FileE_BB | FileF_BB | FileG_BB | FileH_BB
	CenterFiles_BB uint64 = FileC_BB | FileD_BB | FileE_BB | FileF_BB
)

var Ranks_BB = [...]uint64{Rank1_BB, Rank2_BB, Rank3_BB, Rank4_BB, Rank5_BB, Rank6_BB, Rank7_BB, Rank8_BB}
var Files_BB = [...]uint64{FileA_BB, FileB_BB, FileC_BB, FileD_BB, FileE_BB, FileF_BB, FileG_BB, FileH_BB}

var KingFlank_BB = [8]uint64{QueenSide_BB ^ FileD_BB, QueenSide_BB, QueenSide_BB,
	CenterFiles_BB, CenterFiles_BB, KingSide_BB, KingSide_BB, KingSide_BB ^ FileE_BB}

func File(id int) int {
	return id & 7
}

func FileBB(id int) uint64 {
	return Files_BB[File(id)]
}

func Rank(id int) int {
	return id >> 3
}

func Colour(id int) int {
	return (File(id) ^ Rank(id)) & 1
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

var FileMirror = [8]int{FileA, FileB, FileC, FileD, FileD, FileC, FileB, FileA}

func initArray(array *[64]uint64, method func(mask uint64) uint64) {
	for i := uint32(0); i <= 63; i++ {
		array[i] = method(uint64(1) << i)
	}
}

func KingsAttacks(set uint64) uint64 {
	return NorthEast(set & ^Rank8_BB & ^FileH_BB) | North(set & ^Rank8_BB) |
		NorthWest(set & ^Rank8_BB & ^FileA_BB) | East(set & ^FileH_BB) | West(set & ^FileA_BB) |
		SouthEast(set & ^Rank1_BB & ^FileH_BB) | South(set & ^Rank1_BB) | SouthWest(set & ^Rank1_BB & ^FileA_BB)
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
	return NorthWest(set & ^FileA_BB) |
		NorthEast(set & ^FileH_BB)
}

func WhitePawnsDoubleAttacks(set uint64) uint64 {
	return NorthWest(set & ^FileA_BB) &
		NorthEast(set & ^FileH_BB)
}

func WhitePawnsMoves(set uint64, occupancy uint64) uint64 {
	return (North(set) | ((North(North(set)) & Rank4_BB) & ^North(occupancy&Rank3_BB))) & ^occupancy
}

func BlackPawnsAttacks(set uint64) uint64 {
	return SouthWest(set & ^FileA_BB) |
		SouthEast(set & ^FileH_BB)
}

func BlackPawnsDoubleAttacks(set uint64) uint64 {
	return SouthWest(set & ^FileA_BB) &
		SouthEast(set & ^FileH_BB)
}

func BlackPawnsMoves(set uint64, occupancy uint64) uint64 {
	return (South(set) | ((South(South(set)) & Rank5_BB) & ^South(occupancy&Rank6_BB))) & ^occupancy
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
