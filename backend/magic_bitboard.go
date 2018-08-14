package backend

// Names as in https://stackoverflow.com/a/30862064
// Pretty much everything as in this answer, but index right shift is done by constant value(bishopShift, rookShift)
// Pseudo-random number generation from https://github.com/goutham/magic-bits

import (
	"math/rand"
	"sort"
	"sync"
)

const (
	MAX_ROOK_BITS        = 12
	MAX_BISHOP_BITS      = 9
	bishopShift     uint = 64 - MAX_BISHOP_BITS
	rookShift       uint = 64 - MAX_ROOK_BITS
)

var (
	rookMoveBoard                      [64][1 << 12]uint64
	bishopMoveBoard                    [64][1 << 9]uint64
	bishopBlockerMask, rookBlockerMask [64]uint64
	rookMagicIndex, bishopMagicIndex   [64]uint64
)

func generateRookBlockerMask(mask uint64) uint64 {
	res := uint64(0)
	file := getFile(mask)
	rank := getRank(mask)
	res |= file
	res |= rank
	res ^= mask
	if file != FILE_A_BB {
		res &= ^FILE_A_BB
	}
	if file != FILE_H_BB {
		res &= ^FILE_H_BB
	}
	if rank != RANK_1_BB {
		res &= ^RANK_1_BB
	}
	if rank != RANK_8_BB {
		res &= ^RANK_8_BB
	}
	return res
}

func generateBishopBlockerMask(mask uint64) uint64 {
	res := uint64(0)
	tmpMask := mask
	for tmpMask&FILE_H_BB == 0 && tmpMask&RANK_8_BB == 0 {
		res |= tmpMask
		tmpMask = NorthEast(tmpMask)
	}
	tmpMask = mask
	for tmpMask&FILE_H_BB == 0 && tmpMask&RANK_1_BB == 0 {
		res |= tmpMask
		tmpMask = SouthEast(tmpMask)
	}
	tmpMask = mask
	for tmpMask&FILE_A_BB == 0 && tmpMask&RANK_1_BB == 0 {
		res |= tmpMask
		tmpMask = SouthWest(tmpMask)
	}
	tmpMask = mask
	for tmpMask&FILE_A_BB == 0 && tmpMask&RANK_8_BB == 0 {
		res |= tmpMask
		tmpMask = NorthWest(tmpMask)
	}
	res &= ^mask
	return res
}

func combinations(x uint64) []uint64 {
	if x == 0 {
		return []uint64{0}
	}
	right_hand_bit := x & -x
	tmp := combinations(x & ^right_hand_bit)
	res := append([]uint64{}, tmp...)
	for _, el := range tmp {
		res = append(res, el|right_hand_bit)
	}
	return res
}

func sortedCombinations(x uint64) []uint64 {
	res := combinations(x)
	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
	return res
}

func initRookBlockerBoard() (rookBlockerBoard [][]uint64) {
	for _, val := range rookBlockerMask {
		rookBlockerBoard = append(rookBlockerBoard, sortedCombinations(val))
	}
	return
}

func initBishopBlockerBoard() (bishopBlockerBoard [][]uint64) {
	for _, val := range bishopBlockerMask {
		bishopBlockerBoard = append(bishopBlockerBoard, sortedCombinations(val))
	}
	return
}

func initRookMoveBoard(rookBlockerBoard [][]uint64) {
	for y, position := range rookBlockerBoard {
		for x, board := range position {
			rookMoveBoard[y][x] = generateRookMoveBoard(y, board)
		}
	}
}

func generateRookMoveBoard(idx int, board uint64) uint64 {
	res := uint64(0)
	mask := uint64(1) << uint64(idx)
	blockerMask := rookBlockerMask[idx]

	if mask&FILE_A_BB == 0 {
		tmpMask := West(mask)
		for blockerMask&tmpMask > 0 && board&tmpMask == 0 {
			res |= tmpMask
			tmpMask = West(tmpMask)
		}
		res |= tmpMask
	}
	if mask&FILE_H_BB == 0 {
		tmpMask := East(mask)
		for blockerMask&tmpMask > 0 && board&tmpMask == 0 {
			res |= tmpMask
			tmpMask = East(tmpMask)
		}
		res |= tmpMask
	}
	if mask&RANK_8_BB == 0 {
		tmpMask := North(mask)
		for blockerMask&tmpMask > 0 && board&tmpMask == 0 {
			res |= tmpMask
			tmpMask = North(tmpMask)
		}
		res |= tmpMask
	}
	if mask&RANK_1_BB == 0 {
		tmpMask := South(mask)
		for blockerMask&tmpMask > 0 && board&tmpMask == 0 {
			res |= tmpMask
			tmpMask = South(tmpMask)
		}
		res |= tmpMask
	}

	return res
}

func generateBishopMoveBoard(idx int, board uint64) uint64 {
	res := uint64(0)

	mask := uint64(1) << uint64(idx)
	blockerMask := bishopBlockerMask[idx]

	if mask&FILE_H_BB == 0 && mask&RANK_8_BB == 0 {
		tmpMask := NorthEast(mask)
		for blockerMask&tmpMask > 0 && board&tmpMask == 0 {
			res |= tmpMask
			tmpMask = NorthEast(tmpMask)
		}
		res |= tmpMask
	}
	if mask&FILE_H_BB == 0 && mask&RANK_1_BB == 0 {
		tmpMask := SouthEast(mask)
		for blockerMask&tmpMask > 0 && board&tmpMask == 0 {
			res |= tmpMask
			tmpMask = SouthEast(tmpMask)
		}
		res |= tmpMask
	}
	if mask&FILE_A_BB == 0 && mask&RANK_1_BB == 0 {
		tmpMask := SouthWest(mask)
		for blockerMask&tmpMask > 0 && board&tmpMask == 0 {
			res |= tmpMask
			tmpMask = SouthWest(tmpMask)
		}
		res |= tmpMask
	}
	if mask&FILE_A_BB == 0 && mask&RANK_8_BB == 0 {
		tmpMask := NorthWest(mask)
		for blockerMask&tmpMask > 0 && board&tmpMask == 0 {
			res |= tmpMask
			tmpMask = NorthWest(tmpMask)
		}
		res |= tmpMask
	}

	return res
}

func initBishopMoveBoard(bishopBlockerBoard [][]uint64) {
	for y, position := range bishopBlockerBoard {
		for x, board := range position {
			bishopMoveBoard[y][x] = generateBishopMoveBoard(y, board)
		}
	}
}

func initRookMagicIndex(rookBlockerBoard [][]uint64) {
	var wg sync.WaitGroup
	for idx, _ := range rookBlockerMask {
		wg.Add(1)
		go func(i int) {
			val := findMagic(rookBlockerBoard[i], rookMoveBoard[i][:], rookShift)
			rookMagicIndex[i] = val
			wg.Done()
		}(idx)
	}
	wg.Wait()
}

func initBishopMagicIndex(bishopBlockerBoard [][]uint64) {
	var wg sync.WaitGroup
	for idx, _ := range bishopBlockerMask {
		wg.Add(1)
		go func(i int) {
			bishopMagicIndex[i] = findMagic(bishopBlockerBoard[i], bishopMoveBoard[i][:], bishopShift)
			wg.Done()
		}(idx)
	}
	wg.Wait()
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
	rookMoveBoard = rookAttacks
}

func initBishopAttacks(bishopBlockerBoard [][]uint64) {
	var bishopAttacks [64][1 << 9]uint64
	for idx, magic := range bishopMagicIndex {
		for innerIdx, el := range bishopBlockerBoard[idx] {
			mult := uint64(el*magic) >> bishopShift
			bishopAttacks[idx][mult] = bishopMoveBoard[idx][innerIdx]
		}
	}
	bishopMoveBoard = bishopAttacks
}
