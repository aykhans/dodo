package requests

import (
	"os"
	"time"

	. "github.com/aykhans/dodo/types"
	"github.com/aykhans/dodo/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Response struct {
	Response string
	Time     time.Duration
}

type Responses []*Response

// Print prints the responses in a tabular format, including information such as
// response count, minimum time, maximum time, average time, and latency percentiles.
func (responses Responses) Print() {
	total := struct {
		Count int
		Min   time.Duration
		Max   time.Duration
		Sum   time.Duration
		P90   time.Duration
		P95   time.Duration
		P99   time.Duration
	}{
		Count: len(responses),
		Min:   responses[0].Time,
		Max:   responses[0].Time,
	}
	mergedResponses := make(map[string]Durations)
	var allDurations Durations

	for _, response := range responses {
		if response.Time < total.Min {
			total.Min = response.Time
		}
		if response.Time > total.Max {
			total.Max = response.Time
		}
		total.Sum += response.Time

		mergedResponses[response.Response] = append(
			mergedResponses[response.Response],
			response.Time,
		)
		allDurations = append(allDurations, response.Time)
	}
	allDurations.Sort()
	allDurationsLenAsFloat := float64(len(allDurations) - 1)
	total.P90 = allDurations[int(0.90*allDurationsLenAsFloat)]
	total.P95 = allDurations[int(0.95*allDurationsLenAsFloat)]
	total.P99 = allDurations[int(0.99*allDurationsLenAsFloat)]

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 40},
	})
	t.AppendHeader(table.Row{
		"Response",
		"Count",
		"Min",
		"Max",
		"Average",
		"P90",
		"P95",
		"P99",
	})

	var roundPrecision int64 = 4
	for key, durations := range mergedResponses {
		durations.Sort()
		durationsLen := len(durations)
		durationsLenAsFloat := float64(durationsLen - 1)

		t.AppendRow(table.Row{
			key,
			durationsLen,
			utils.DurationRoundBy(*durations.First(), roundPrecision),
			utils.DurationRoundBy(*durations.Last(), roundPrecision),
			utils.DurationRoundBy(durations.Avg(), roundPrecision),
			utils.DurationRoundBy(durations[int(0.90*durationsLenAsFloat)], roundPrecision),
			utils.DurationRoundBy(durations[int(0.95*durationsLenAsFloat)], roundPrecision),
			utils.DurationRoundBy(durations[int(0.99*durationsLenAsFloat)], roundPrecision),
		})
		t.AppendSeparator()
	}

	if len(mergedResponses) > 1 {
		t.AppendRow(table.Row{
			"Total",
			total.Count,
			utils.DurationRoundBy(total.Min, roundPrecision),
			utils.DurationRoundBy(total.Max, roundPrecision),
			utils.DurationRoundBy(total.Sum/time.Duration(total.Count), roundPrecision), // Average
			utils.DurationRoundBy(total.P90, roundPrecision),
			utils.DurationRoundBy(total.P95, roundPrecision),
			utils.DurationRoundBy(total.P99, roundPrecision),
		})
	}
	t.Render()
}
