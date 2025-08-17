package types

import (
	"slices"
	"sort"
	"time"
)

type Durations []time.Duration

func (d Durations) Sort(ascending ...bool) {
	// If ascending is provided and is false, sort in descending order
	if len(ascending) > 0 && ascending[0] == false {
		sort.Slice(d, func(i, j int) bool {
			return d[i] > d[j]
		})
	} else { // Otherwise, sort in ascending order
		slices.Sort(d)
	}
}

func (d Durations) First() *time.Duration {
	return &d[0]
}

func (d Durations) Last() *time.Duration {
	return &d[len(d)-1]
}

func (d Durations) Sum() time.Duration {
	sum := time.Duration(0)
	for _, duration := range d {
		sum += duration
	}
	return sum
}

func (d Durations) Avg() time.Duration {
	return d.Sum() / time.Duration(len(d))
}
