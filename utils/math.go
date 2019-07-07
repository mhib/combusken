package utils

func NearestPowerOfTwo(input int) uint64 {
	res := uint64(1)
	for int((res << 1)) <= input {
		res <<= 1
	}
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
