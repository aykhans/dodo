package requests

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/valyala/fasthttp"
)

// streamProgress streams the progress of a task to the console using a progress bar.
// It listens for increments on the provided channel and updates the progress bar accordingly.
//
// The function will stop and mark the progress as errored if the context is cancelled.
// It will also stop and mark the progress as done when the total number of increments is reached.
func streamProgress(
	ctx context.Context,
	wg *sync.WaitGroup,
	total int64,
	message string,
	increase <-chan int64,
) {
	defer wg.Done()
	pw := progress.NewWriter()
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetStyle(progress.StyleBlocks)
	pw.SetTrackerLength(40)
	pw.SetUpdateFrequency(time.Millisecond * 250)
	go pw.Render()
	dodosTracker := progress.Tracker{
		Message: message,
		Total:   total,
	}
	pw.AppendTracker(&dodosTracker)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\r")
			dodosTracker.MarkAsErrored()
			time.Sleep(time.Millisecond * 300)
			pw.Stop()
			return

		case value := <-increase:
			dodosTracker.Increment(value)
		}
	}
}

// checkConnection checks the internet connection by making requests to different websites.
// It returns true if the connection is successful, otherwise false.
func checkConnection(ctx context.Context) bool {
	ch := make(chan bool)
	go func() {
		_, _, err := fasthttp.Get(nil, "https://www.google.com")
		if err != nil {
			_, _, err = fasthttp.Get(nil, "https://www.bing.com")
			if err != nil {
				_, _, err = fasthttp.Get(nil, "https://www.yahoo.com")
				ch <- err == nil
			}
			ch <- true
		}
		ch <- true
	}()

	select {
	case <-ctx.Done():
		return false
	case res := <-ch:
		return res
	}
}
