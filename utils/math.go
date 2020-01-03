package utils

func NearestPowerOfTwo(input int) uint64 {
	res := uint64(1)
	for res <= uint64(input) {
		res <<= 1
	}
	res >>= 1
	return res
}

func Abs(input int) int {
	if input < 0 {
		return -input
	}
	return input
}

func Max(l, r int) int {
	if l <= r {
		return r
	}
	return l
}

func Min(l, r int) int {
	if l <= r {
		return l
	}
	return r
}

func BoolToInt(input bool) (res int) {
	if input {
		res = 1
	} else {
		res = 0
	}
	return
}
