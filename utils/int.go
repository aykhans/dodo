package utils

type Number interface {
	int | int8 | int16 | int32 | int64
}

func NumLen[T Number](n T) T {
	if n < 0 {
		n = -n
	}
	if n == 0 {
		return 1
	}

	var count T = 0
	for n > 0 {
		n /= 10
		count++
	}
	return count
}
