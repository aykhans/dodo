package utils

func IsNilOrZero[T comparable](value *T) bool {
	if value == nil {
		return true
	}

	var zero T
	if *value == zero {
		return true
	}

	return false
}
