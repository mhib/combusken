package backend

// Names as in https://stackoverflow.com/a/30862064
// Pretty much everything as in this answer, but index right shift is done by constant values(bishopShift, rookShift)
// Pseudo-random number generation from https://github.com/goutham/magic-bits

import (
	"math/rand"
)

const (
	MaxRookBits        = 12
	MaxBishopBits      = 9
	bishopShift   uint = 64 - MaxBishopBits
	rookShift     uint = 64 - MaxRookBits
)

var (
	rookMoveBoard                      [64][1 << MaxRookBits]uint64
	bishopMoveBoard                    [64][1 << MaxBishopBits]uint64
	bishopBlockerMask, rookBlockerMask [64]uint64
	rookMagicIndex, bishopMagicIndex   [64]uint64
)

func generateRookBlockerMask(mask uint64) uint64 {
	res := uint64(0)
	square := BitScan(mask)
	file := File(square)
	rank := Rank(square)
	res |= Files_BB[file] | Ranks_BB[rank]
	res ^= mask

	if file != FileA {
		res &= ^FileA_BB
	}
	if file != FileH {
		res &= ^FileH_BB
	}
	if rank != Rank1 {
		res &= ^Rank1_BB
	}
	if rank != Rank8 {
		res &= ^Rank8_BB
	}
	return res
}

func generateBishopBlockerMask(mask uint64) uint64 {
	res := uint64(0)
	for tmpMask := mask; tmpMask&FileH_BB == 0 && tmpMask&Rank8_BB == 0; tmpMask = NorthEast(tmpMask) {
		res |= tmpMask
	}
	for tmpMask := mask; tmpMask&FileH_BB == 0 && tmpMask&Rank1_BB == 0; tmpMask = SouthEast(tmpMask) {
		res |= tmpMask
	}
	for tmpMask := mask; tmpMask&FileA_BB == 0 && tmpMask&Rank1_BB == 0; tmpMask = SouthWest(tmpMask) {
		res |= tmpMask
	}
	for tmpMask := mask; tmpMask&FileA_BB == 0 && tmpMask&Rank8_BB == 0; tmpMask = NorthWest(tmpMask) {
		res |= tmpMask
	}
	res &= ^mask
	return res
}

func combinations(x uint64) []uint64 {
	if x == 0 {
		return []uint64{0}
	}
	rightHandBit := x & -x
	recursive := combinations(x & ^rightHandBit)
	res := append([]uint64{}, recursive...)
	for _, el := range recursive {
		res = append(res, el|rightHandBit)
	}
	return res
}

func initRookBlockerBoard() (rookBlockerBoard [][]uint64) {
	for _, val := range rookBlockerMask {
		rookBlockerBoard = append(rookBlockerBoard, combinations(val))
	}
	return
}

func initBishopBlockerBoard() (bishopBlockerBoard [][]uint64) {
	for _, val := range bishopBlockerMask {
		bishopBlockerBoard = append(bishopBlockerBoard, combinations(val))
	}
	return
}

func initRookMoveBoard(rookBlockerBoard [][]uint64) {
	for y, boards := range rookBlockerBoard {
		for x, board := range boards {
			rookMoveBoard[y][x] = generateRookMoveBoard(y, board)
		}
	}
}

func generateRookMoveBoard(idx int, board uint64) (res uint64) {
	mask := uint64(1) << uint64(idx)
	blockerMask := rookBlockerMask[idx]

	if File(idx) != FileA {
		for tmpMask := West(mask); ; tmpMask = West(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}
	if File(idx) != FileH {
		for tmpMask := East(mask); ; tmpMask = East(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}
	if Rank(idx) != Rank8 {
		for tmpMask := North(mask); ; tmpMask = North(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}
	if Rank(idx) != Rank1 {
		for tmpMask := South(mask); ; tmpMask = South(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}

	return
}

func generateBishopMoveBoard(idx int, board uint64) (res uint64) {
	mask := uint64(1) << uint64(idx)
	blockerMask := bishopBlockerMask[idx]

	if mask&FileH_BB == 0 && mask&Rank8_BB == 0 {
		for tmpMask := NorthEast(mask); ; tmpMask = NorthEast(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}
	if mask&FileH_BB == 0 && mask&Rank1_BB == 0 {
		for tmpMask := SouthEast(mask); ; tmpMask = SouthEast(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}
	if mask&FileA_BB == 0 && mask&Rank1_BB == 0 {
		for tmpMask := SouthWest(mask); ; tmpMask = SouthWest(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}
	if mask&FileA_BB == 0 && mask&Rank8_BB == 0 {
		for tmpMask := NorthWest(mask); ; tmpMask = NorthWest(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}

	return
}

func initBishopMoveBoard(bishopBlockerBoard [][]uint64) {
	for y, boards := range bishopBlockerBoard {
		for x, board := range boards {
			bishopMoveBoard[y][x] = generateBishopMoveBoard(y, board)
		}
	}
}

func initRookMagicIndex(rookBlockerBoard [][]uint64) {
	for idx := range rookBlockerBoard {
		rookMagicIndex[idx] = findMagic(rookBlockerBoard[idx], rookMoveBoard[idx][:], rookShift)
	}
}

func initBishopMagicIndex(bishopBlockerBoard [][]uint64) {
	for idx := range bishopBlockerBoard {
		bishopMagicIndex[idx] = findMagic(bishopBlockerBoard[idx], bishopMoveBoard[idx][:], bishopShift)
	}
}

func u64rand() uint64 {
	return (uint64(0xFFFF&rand.Uint32()) << 48) |
		(uint64(0xFFFF&rand.Uint32()) << 32) |
		(uint64(0xFFFF&rand.Uint32()) << 16) |
		uint64(0xFFFF&rand.Uint32())
}

func biasedRandom() uint64 {
	return u64rand() & u64rand() & u64rand()
}

func findMagic(array []uint64, cmpArray []uint64, bits uint) uint64 {
	for {
		magic := biasedRandom()
		others := make(map[uint64]int)
		unique := true
		for idx, el := range array {
			mult := uint64(el*magic) >> bits
			if x, found := others[mult]; found {
				if cmpArray[x] != cmpArray[idx] {
					unique = false
					break
				}
			}
			others[mult] = idx
		}
		if unique {
			return magic
		}
	}
}

func initRookAttacks(rookBlockerBoard [][]uint64) {
	var rookAttacks [64][1 << 12]uint64
	for idx, magic := range rookMagicIndex {
		for innerIdx, el := range rookBlockerBoard[idx] {
			mult := uint64(el*magic) >> rookShift
			rookAttacks[idx][mult] = rookMoveBoard[idx][innerIdx]
		}
	}
	copy(rookMoveBoard[:], rookAttacks[:])
}

func initBishopAttacks(bishopBlockerBoard [][]uint64) {
	var bishopAttacks [64][1 << 9]uint64
	for idx, magic := range bishopMagicIndex {
		for innerIdx, el := range bishopBlockerBoard[idx] {
			mult := uint64(el*magic) >> bishopShift
			bishopAttacks[idx][mult] = bishopMoveBoard[idx][innerIdx]
		}
	}
	copy(bishopMoveBoard[:], bishopAttacks[:])
}
