package utils

func Max(l, r int) int {
	if l > r {
		return l
	}
	return r
}

func Min(l, r int) int {
	if l < r {
		return l
	}
	return r
}

func Abs(l int) int {
	if l < 0 {
		return -l
	}
	return l
}

func NearestPowerOfTwo(input int) uint64 {
	res := uint64(1)
	for int((res << 1)) <= input {
		res <<= 1
	}
	return res
}
