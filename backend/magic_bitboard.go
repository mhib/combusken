package backend

import "unsafe"

const (
	MAX_ROOK_BITS   = 12
	MAX_BISHOP_BITS = 9
)

type Magic struct {
	blockerMask uint64
	magicIndex  uint64
	offset      unsafe.Pointer
	shift       uint64
}

func (magic Magic) attacks(occupancy uint64) uint64 {
	return *((*uint64)(unsafe.Pointer(uintptr(magic.offset) + uintptr((magic.blockerMask&occupancy*magic.magicIndex)>>magic.shift)<<3)))
}

var (
	bishopMagics, rookMagics [64]Magic
	bishopAttacks            [0x1480]uint64
	rookAttacks              [0x19000]uint64
)
var rookMagicIdx = [64]uint64{
	0xA180022080400230, 0x0040100040022000, 0x0080088020001002, 0x0080080280841000,
	0x4200042010460008, 0x04800A0003040080, 0x0400110082041008, 0x008000A041000880,
	0x10138001A080C010, 0x0000804008200480, 0x00010011012000C0, 0x0022004128102200,
	0x000200081201200C, 0x202A001048460004, 0x0081000100420004, 0x4000800380004500,
	0x0000208002904001, 0x0090004040026008, 0x0208808010002001, 0x2002020020704940,
	0x8048010008110005, 0x6820808004002200, 0x0A80040008023011, 0x00B1460000811044,
	0x4204400080008EA0, 0xB002400180200184, 0x2020200080100380, 0x0010080080100080,
	0x2204080080800400, 0x0000A40080360080, 0x02040604002810B1, 0x008C218600004104,
	0x8180004000402000, 0x488C402000401001, 0x4018A00080801004, 0x1230002105001008,
	0x8904800800800400, 0x0042000C42003810, 0x008408110400B012, 0x0018086182000401,
	0x2240088020C28000, 0x001001201040C004, 0x0A02008010420020, 0x0010003009010060,
	0x0004008008008014, 0x0080020004008080, 0x0282020001008080, 0x50000181204A0004,
	0x48FFFE99FECFAA00, 0x48FFFE99FECFAA00, 0x497FFFADFF9C2E00, 0x613FFFDDFFCE9200,
	0xFFFFFFE9FFE7CE00, 0xFFFFFFF5FFF3E600, 0x0010301802830400, 0x510FFFF5F63C96A0,
	0xEBFFFFB9FF9FC526, 0x61FFFEDDFEEDAEAE, 0x53BFFFEDFFDEB1A2, 0x127FFFB9FFDFB5F6,
	0x411FFFDDFFDBF4D6, 0x0801000804000603, 0x0003FFEF27EEBE74, 0x7645FFFECBFEA79E,
}

var bishopMagicIdx = [64]uint64{
	0xFFEDF9FD7CFCFFFF, 0xFC0962854A77F576, 0x5822022042000000, 0x2CA804A100200020,
	0x0204042200000900, 0x2002121024000002, 0xFC0A66C64A7EF576, 0x7FFDFDFCBD79FFFF,
	0xFC0846A64A34FFF6, 0xFC087A874A3CF7F6, 0x1001080204002100, 0x1810080489021800,
	0x0062040420010A00, 0x5028043004300020, 0xFC0864AE59B4FF76, 0x3C0860AF4B35FF76,
	0x73C01AF56CF4CFFB, 0x41A01CFAD64AAFFC, 0x040C0422080A0598, 0x4228020082004050,
	0x0200800400E00100, 0x020B001230021040, 0x7C0C028F5B34FF76, 0xFC0A028E5AB4DF76,
	0x0020208050A42180, 0x001004804B280200, 0x2048020024040010, 0x0102C04004010200,
	0x020408204C002010, 0x02411100020080C1, 0x102A008084042100, 0x0941030000A09846,
	0x0244100800400200, 0x4000901010080696, 0x0000280404180020, 0x0800042008240100,
	0x0220008400088020, 0x04020182000904C9, 0x0023010400020600, 0x0041040020110302,
	0xDCEFD9B54BFCC09F, 0xF95FFA765AFD602B, 0x1401210240484800, 0x0022244208010080,
	0x1105040104000210, 0x2040088800C40081, 0x43FF9A5CF4CA0C01, 0x4BFFCD8E7C587601,
	0xFC0FF2865334F576, 0xFC0BF6CE5924F576, 0x80000B0401040402, 0x0020004821880A00,
	0x8200002022440100, 0x0009431801010068, 0xC3FFB7DC36CA8C89, 0xC3FF8A54F4CA2C89,
	0xFFFFFCFCFD79EDFF, 0xFC0863FCCB147576, 0x040C000022013020, 0x2000104000420600,
	0x0400000260142410, 0x0800633408100500, 0xFC087E8E4BB2F736, 0x43FF9E4EF4CA2C89,
}

func generateRookBlockerMask(mask uint64) uint64 {
	res := uint64(0)
	square := BitScan(mask)
	file := File(square)
	rank := Rank(square)
	res |= FILES[file] | RANKS[rank]
	res ^= mask

	if file != FILE_A {
		res &= ^FILE_A_BB
	}
	if file != FILE_H {
		res &= ^FILE_H_BB
	}
	if rank != RANK_1 {
		res &= ^RANK_1_BB
	}
	if rank != RANK_8 {
		res &= ^RANK_8_BB
	}
	return res
}

func generateBishopBlockerMask(mask uint64) uint64 {
	res := uint64(0)
	for tmpMask := mask; tmpMask&FILE_H_BB == 0 && tmpMask&RANK_8_BB == 0; tmpMask = NorthEast(tmpMask) {
		res |= tmpMask
	}
	for tmpMask := mask; tmpMask&FILE_H_BB == 0 && tmpMask&RANK_1_BB == 0; tmpMask = SouthEast(tmpMask) {
		res |= tmpMask
	}
	for tmpMask := mask; tmpMask&FILE_A_BB == 0 && tmpMask&RANK_1_BB == 0; tmpMask = SouthWest(tmpMask) {
		res |= tmpMask
	}
	for tmpMask := mask; tmpMask&FILE_A_BB == 0 && tmpMask&RANK_8_BB == 0; tmpMask = NorthWest(tmpMask) {
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

func initRookBlockerBoard(rookBlockerMask *[64]uint64) (rookBlockerBoard [][]uint64) {
	for _, val := range *rookBlockerMask {
		rookBlockerBoard = append(rookBlockerBoard, combinations(val))
	}
	return
}

func initBishopBlockerBoard(bishopBlockerMask *[64]uint64) (bishopBlockerBoard [][]uint64) {
	for _, val := range *bishopBlockerMask {
		bishopBlockerBoard = append(bishopBlockerBoard, combinations(val))
	}
	return
}

func initRookMoveBoard(rookBlockerMask *[64]uint64, rookBlockerBoard [][]uint64) [64][1 << MAX_ROOK_BITS]uint64 {
	var rookMoveBoard [64][1 << MAX_ROOK_BITS]uint64
	for y, boards := range rookBlockerBoard {
		for x, board := range boards {
			rookMoveBoard[y][x] = generateRookMoveBoard(rookBlockerMask[y], y, board)
		}
	}
	return rookMoveBoard
}

func generateRookMoveBoard(blockerMask uint64, idx int, board uint64) (res uint64) {
	mask := uint64(1) << uint64(idx)

	if File(idx) != FILE_A {
		for tmpMask := West(mask); ; tmpMask = West(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}
	if File(idx) != FILE_H {
		for tmpMask := East(mask); ; tmpMask = East(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}
	if Rank(idx) != RANK_8 {
		for tmpMask := North(mask); ; tmpMask = North(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}
	if Rank(idx) != RANK_1 {
		for tmpMask := South(mask); ; tmpMask = South(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}

	return
}

func generateBishopMoveBoard(blockerMask uint64, idx int, board uint64) (res uint64) {
	mask := uint64(1) << uint64(idx)

	if mask&FILE_H_BB == 0 && mask&RANK_8_BB == 0 {
		for tmpMask := NorthEast(mask); ; tmpMask = NorthEast(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}
	if mask&FILE_H_BB == 0 && mask&RANK_1_BB == 0 {
		for tmpMask := SouthEast(mask); ; tmpMask = SouthEast(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}
	if mask&FILE_A_BB == 0 && mask&RANK_1_BB == 0 {
		for tmpMask := SouthWest(mask); ; tmpMask = SouthWest(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}
	if mask&FILE_A_BB == 0 && mask&RANK_8_BB == 0 {
		for tmpMask := NorthWest(mask); ; tmpMask = NorthWest(tmpMask) {
			res |= tmpMask
			if blockerMask&tmpMask == 0 || board&tmpMask > 0 {
				break
			}
		}
	}

	return
}

func initBishopMoveBoard(blockerMask *[64]uint64, bishopBlockerBoard [][]uint64) [64][1 << MAX_BISHOP_BITS]uint64 {
	var bishopMoveBoard [64][1 << MAX_BISHOP_BITS]uint64
	for y, boards := range bishopBlockerBoard {
		for x, board := range boards {
			bishopMoveBoard[y][x] = generateBishopMoveBoard(blockerMask[y], y, board)
		}
	}
	return bishopMoveBoard
}

func initRookAttacks(blockerMask *[64]uint64, rookBlockerBoard [][]uint64, rookMoveBoard *[64][1 << MAX_ROOK_BITS]uint64) {
	rookMagics[0].offset = unsafe.Pointer(&rookAttacks[0])
	for idx := range rookMagicIdx {
		magic := &rookMagics[idx]
		magic.blockerMask = blockerMask[idx]
		magic.magicIndex = rookMagicIdx[idx]
		magic.shift = uint64(64 - PopCount(blockerMask[idx]))
		for innerIdx, el := range rookBlockerBoard[idx] {
			mult := uintptr(el*magic.magicIndex>>magic.shift) * 8
			*((*uint64)(unsafe.Pointer(uintptr(magic.offset) + mult))) = rookMoveBoard[idx][innerIdx]
		}
		if idx != H8 {
			rookMagics[idx+1].offset = unsafe.Pointer(uintptr(magic.offset) + uintptr(8<<PopCount(blockerMask[idx])))
		}
	}
}

func initBishopAttacks(blockerMask *[64]uint64, bishopBlockerBoard [][]uint64, bishopMoveBoard *[64][1 << MAX_BISHOP_BITS]uint64) {
	bishopMagics[0].offset = unsafe.Pointer(&bishopAttacks[0])
	for idx := range bishopMagicIdx {
		magic := &bishopMagics[idx]
		magic.blockerMask = blockerMask[idx]
		magic.magicIndex = bishopMagicIdx[idx]
		magic.shift = uint64(64 - PopCount(blockerMask[idx]))
		for innerIdx, el := range bishopBlockerBoard[idx] {
			mult := uintptr(el*magic.magicIndex>>magic.shift) * 8
			*((*uint64)(unsafe.Pointer(uintptr(magic.offset) + mult))) = bishopMoveBoard[idx][innerIdx]
		}
		if idx != H8 {
			bishopMagics[idx+1].offset = unsafe.Pointer(uintptr(magic.offset) + uintptr(8<<PopCount(blockerMask[idx])))
		}
	}
}

func init() {
	var rookBlockerMask [64]uint64
	initArray(&rookBlockerMask, generateRookBlockerMask)
	rookBlockerBoard := initRookBlockerBoard(&rookBlockerMask)
	rookMoveBoard := initRookMoveBoard(&rookBlockerMask, rookBlockerBoard)
	initRookAttacks(&rookBlockerMask, rookBlockerBoard, &rookMoveBoard)

	var bishopBlockerMask [64]uint64
	initArray(&bishopBlockerMask, generateBishopBlockerMask)
	bishopBlockerBoard := initBishopBlockerBoard(&bishopBlockerMask)
	bishopMoveBoard := initBishopMoveBoard(&bishopBlockerMask, bishopBlockerBoard)
	initBishopAttacks(&bishopBlockerMask, bishopBlockerBoard, &bishopMoveBoard)
}
