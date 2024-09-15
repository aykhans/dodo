package utils

import "math/rand"

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

func RandomValueCycle[Value any](values []Value, localRand *rand.Rand) func() Value {
	var (
		clientsCount int = len(values)
		currentIndex int = localRand.Intn(clientsCount)
		stopIndex    int = currentIndex
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
