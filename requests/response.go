package requests

import (
	"os"
	"time"

	"github.com/aykhans/dodo/types"
	"github.com/aykhans/dodo/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Response struct {
	Response string
	Time     time.Duration
}

type Responses []Response

// Print prints the responses in a tabular format, including information such as
// response count, minimum time, maximum time, average time, and latency percentiles.
func (responses Responses) Print() {
	if len(responses) == 0 {
		return
	}

	mergedResponses := make(map[string]types.Durations)

	totalDurations := make(types.Durations, len(responses))
	var totalSum time.Duration
	totalCount := len(responses)

	for i, response := range responses {
		totalSum += response.Time
		totalDurations[i] = response.Time

		mergedResponses[response.Response] = append(
			mergedResponses[response.Response],
			response.Time,
		)
	}

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
		totalDurations.Sort()
		allDurationsLenAsFloat := float64(len(totalDurations) - 1)

		t.AppendRow(table.Row{
			"Total",
			totalCount,
			utils.DurationRoundBy(totalDurations[0], roundPrecision),
			utils.DurationRoundBy(totalDurations[len(totalDurations)-1], roundPrecision),
			utils.DurationRoundBy(totalSum/time.Duration(totalCount), roundPrecision), // Average
			utils.DurationRoundBy(totalDurations[int(0.90*allDurationsLenAsFloat)], roundPrecision),
			utils.DurationRoundBy(totalDurations[int(0.95*allDurationsLenAsFloat)], roundPrecision),
			utils.DurationRoundBy(totalDurations[int(0.99*allDurationsLenAsFloat)], roundPrecision),
		})
	}
	t.Render()
}
