package utils

import "time"

func MinDuration(durations ...time.Duration) time.Duration {
	min := durations[0]
	for _, d := range durations {
		if d < min {
			min = d
		}
	}
	return min
}

func MaxDuration(durations ...time.Duration) time.Duration {
	max := durations[0]
	for _, d := range durations {
		if d > max {
			max = d
		}
	}
	return max
}

func AvgDuration(durations ...time.Duration) time.Duration {
	total := time.Duration(0)
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}
