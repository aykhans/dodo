package requests

import (
	"os"
	"time"

	. "github.com/aykhans/dodo/types"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Response struct {
	Response string
	Time     time.Duration
}

type Responses []*Response

// Print prints the responses in a tabular format, including information such as
// response count, minimum time, maximum time, average time, and latency percentiles.
func (respones Responses) Print() {
	total := struct {
		Count int
		Min   time.Duration
		Max   time.Duration
		Sum   time.Duration
		P90   time.Duration
		P95   time.Duration
		P99   time.Duration
	}{
		Count: len(respones),
		Min:   respones[0].Time,
		Max:   respones[0].Time,
	}
	mergedResponses := make(map[string]Durations)
	var allDurations Durations

	for _, response := range respones {
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
		"Min Time",
		"Max Time",
		"Average Time",
		"P90",
		"P95",
		"P99",
	})

	for key, durations := range mergedResponses {
		durations.Sort()
		durationsLen := len(durations)
		durationsLenAsFloat := float64(durationsLen - 1)

		t.AppendRow(table.Row{
			key,
			durationsLen,
			durations.First(),
			durations.Last(),
			durations.Avg(),
			durations[int(0.90*durationsLenAsFloat)],
			durations[int(0.95*durationsLenAsFloat)],
			durations[int(0.99*durationsLenAsFloat)],
		})
		t.AppendSeparator()
	}

	if len(mergedResponses) > 1 {
		t.AppendRow(table.Row{
			"Total",
			total.Count,
			total.Min,
			total.Max,
			total.Sum / time.Duration(total.Count), // Average
			total.P90,
			total.P95,
			total.P99,
		})
	}
	t.Render()
}
