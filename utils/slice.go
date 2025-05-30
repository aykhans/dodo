package utils

import "math/rand"

func Flatten[T any](nested [][]*T) []*T {
	flattened := make([]*T, 0)
	for _, n := range nested {
		flattened = append(flattened, n...)
	}
	return flattened
}

// RandomValueCycle returns a function that cycles through the provided slice of values
// in a random order. Each call to the returned function will yield a value from the slice.
// The order of values is determined by the provided random number generator.
//
// The returned function will cycle through the values in a random order until all values
// have been returned at least once. After all values have been returned, the function will
// reset and start cycling through the values in a random order again.
// The returned function isn't thread-safe and should be used in a single-threaded context.
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
