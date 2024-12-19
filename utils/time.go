package utils

import "time"

func DurationRoundBy(duration time.Duration, n int64) time.Duration {
	if durationLen := NumLen(duration.Nanoseconds()); durationLen > n {
		roundNum := 1
		for range durationLen - n {
			roundNum *= 10
		}
		return duration.Round(time.Duration(roundNum))
	}
	return duration
}
