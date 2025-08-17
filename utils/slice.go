package utils

import "math/rand"

func Flatten[T any](nested [][]T) []T {
	flattened := make([]T, 0)
	for _, n := range nested {
		flattened = append(flattened, n...)
	}
	return flattened
}

// RandomValueCycle returns a function that cycles through the provided values in a pseudo-random order.
// Each value in the input slice will be returned before any value is repeated.
// If the input slice is empty, the returned function will always return the zero value of type T.
// If the input slice contains only one element, that element is always returned.
// This function is not thread-safe and should not be called concurrently.
func RandomValueCycle[T any](values []T, localRand *rand.Rand) func() T {
	switch valuesLen := len(values); valuesLen {
	case 0:
		var zero T
		return func() T { return zero }
	case 1:
		return func() T { return values[0] }
	default:
		currentIndex := localRand.Intn(valuesLen)
		stopIndex := currentIndex
		return func() T {
			value := values[currentIndex]
			currentIndex++
			if currentIndex == valuesLen {
				currentIndex = 0
			}
			if currentIndex == stopIndex {
				currentIndex = localRand.Intn(valuesLen)
				stopIndex = currentIndex
			}

			return value
		}
	}
}
