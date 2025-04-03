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
func RandomValueCycle[Value any](values []Value, localRand *rand.Rand) func() Value {
	var (
		clientsCount = len(values)
		currentIndex = localRand.Intn(clientsCount)
		stopIndex    = currentIndex
	)

	return func() Value {
		client := values[currentIndex]
		currentIndex++
		if currentIndex == clientsCount {
			currentIndex = 0
		}
		if currentIndex == stopIndex {
			currentIndex = localRand.Intn(clientsCount)
			stopIndex = currentIndex
		}

		return client
	}
}
