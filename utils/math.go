package utils

func NearestPowerOfTwo(input int) uint64 {
	res := uint64(1)
	for int((res << 1)) <= input {
		res <<= 1
	}
	return res
}

func Abs(input int64) int64 {
	y := input >> 63
	return (input ^ y) - y
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
