package utils

func Flatten[T any](nested [][]*T) []*T {
	flattened := make([]*T, 0)
	for _, n := range nested {
		flattened = append(flattened, n...)
	}
	return flattened
}

func Contains[T comparable](slice []T, item T) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}
