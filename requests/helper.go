package requests

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
)

// streamProgress streams the progress of a task to the console using a progress bar.
// It listens for increments on the provided channel and updates the progress bar accordingly.
//
// The function will stop and mark the progress as errored if the context is cancelled.
// It will also stop and mark the progress as done when the total number of increments is reached.
func streamProgress(
	ctx context.Context,
	wg *sync.WaitGroup,
	total uint,
	message string,
	increase <-chan int64,
) {
	defer wg.Done()
	pw := progress.NewWriter()
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetStyle(progress.StyleBlocks)
	pw.SetTrackerLength(40)
	pw.SetUpdateFrequency(time.Millisecond * 250)
	if total == 0 {
		pw.Style().Visibility.Percentage = false
	}
	go pw.Render()
	dodosTracker := progress.Tracker{
		Message: message,
		Total:   int64(total),
	}
	pw.AppendTracker(&dodosTracker)

	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err == context.Canceled || err == context.DeadlineExceeded {
				dodosTracker.MarkAsDone()
			} else {
				dodosTracker.MarkAsErrored()
			}
			time.Sleep(time.Millisecond * 300)
			fmt.Printf("\r")
			return

		case value := <-increase:
			dodosTracker.Increment(value)
		}
	}
}
