package requests

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aykhans/dodo/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/valyala/fasthttp"
)

type Response struct {
	StatusCode int
	Error      error
	Time       time.Duration
}

type Responses []Response

type ClientDoFunc func(ctx context.Context, request *fasthttp.Request) (*fasthttp.Response, error)

// Print prints the responses in a tabular format, including information such as
// response count, minimum time, maximum time, and average time.
func (respones *Responses) Print() {
	var (
		totalMinDuration time.Duration = (*respones)[0].Time
		totalMaxDuration time.Duration = (*respones)[0].Time
		totalDuration    time.Duration
		totalCount       int = len(*respones)
	)
	mergedResponses := make(map[string][]time.Duration)

	for _, response := range *respones {
		if response.Time < totalMinDuration {
			totalMinDuration = response.Time
		}
		if response.Time > totalMaxDuration {
			totalMaxDuration = response.Time
		}
		totalDuration += response.Time

		if response.Error != nil {
			mergedResponses[response.Error.Error()] = append(
				mergedResponses[response.Error.Error()],
				response.Time,
			)
		} else {
			mergedResponses[fmt.Sprintf("%d", response.StatusCode)] = append(
				mergedResponses[fmt.Sprintf("%d", response.StatusCode)],
				response.Time,
			)
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.SetAllowedRowLength(125)
	t.AppendHeader(table.Row{
		"Response",
		"Count",
		"Min Time",
		"Max Time",
		"Average Time",
	})
	for key, durations := range mergedResponses {
		t.AppendRow(table.Row{
			key,
			len(durations),
			utils.MinDuration(durations...),
			utils.MaxDuration(durations...),
			utils.AvgDuration(durations...),
		})
		t.AppendSeparator()
	}
	t.AppendRow(table.Row{
		"Total",
		totalCount,
		totalMinDuration,
		totalMaxDuration,
		totalDuration / time.Duration(totalCount),
	})
	t.Render()
}
