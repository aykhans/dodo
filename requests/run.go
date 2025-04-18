package requests

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/aykhans/dodo/config"
	"github.com/aykhans/dodo/types"
	"github.com/aykhans/dodo/utils"
	"github.com/valyala/fasthttp"
)

// Run executes the main logic for processing requests based on the provided configuration.
// It initializes clients based on the request configuration and releases the dodos.
// If the context is canceled and no responses are collected, it returns an interrupt error.
//
// Parameters:
//   - ctx: The context for managing request lifecycle and cancellation.
//   - requestConfig: The configuration for the request, including timeout, proxies, and other settings.
func Run(ctx context.Context, requestConfig *config.RequestConfig) (Responses, error) {
	if requestConfig.Duration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, requestConfig.Duration)
		defer cancel()
	}

	clients := getClients(
		ctx,
		requestConfig.Timeout,
		requestConfig.Proxies,
		requestConfig.GetMaxConns(fasthttp.DefaultMaxConnsPerHost),
		requestConfig.URL,
	)
	if clients == nil {
		return nil, types.ErrInterrupt
	}

	responses := releaseDodos(ctx, requestConfig, clients)
	if ctx.Err() != nil && len(responses) == 0 {
		return nil, types.ErrInterrupt
	}

	return responses, nil
}

// releaseDodos sends requests concurrently using multiple dodos (goroutines) and returns the aggregated responses.
//
// The function performs the following steps:
//  1. Initializes wait groups and other necessary variables.
//  2. Starts a goroutine to stream progress updates.
//  3. Distributes the total request count among the dodos.
//  4. Starts a goroutine for each dodo to send requests concurrently.
//  5. Waits for all dodos to complete their requests.
//  6. Cancels the progress streaming context and waits for the progress goroutine to finish.
//  7. Flattens and returns the aggregated responses.
func releaseDodos(
	ctx context.Context,
	requestConfig *config.RequestConfig,
	clients []*fasthttp.HostClient,
) Responses {
	var (
		wg                  sync.WaitGroup
		streamWG            sync.WaitGroup
		requestCountPerDodo uint
		dodosCount          = requestConfig.GetValidDodosCountForRequests()
		responses           = make([][]*Response, dodosCount)
		increase            = make(chan int64, requestConfig.RequestCount)
	)

	wg.Add(int(dodosCount))
	streamWG.Add(1)
	streamCtx, streamCtxCancel := context.WithCancel(context.Background())

	go streamProgress(streamCtx, &streamWG, requestConfig.RequestCount, "Dodos Working🔥", increase)

	if requestConfig.RequestCount == 0 {
		for i := range dodosCount {
			go sendRequest(
				ctx,
				newRequest(*requestConfig, clients, int64(i)),
				requestConfig.Timeout,
				&responses[i],
				increase,
				&wg,
			)
		}
	} else {
		for i := range dodosCount {
			if i+1 == dodosCount {
				requestCountPerDodo = requestConfig.RequestCount - (i * requestConfig.RequestCount / dodosCount)
			} else {
				requestCountPerDodo = ((i + 1) * requestConfig.RequestCount / dodosCount) -
					(i * requestConfig.RequestCount / dodosCount)
			}

			go sendRequestByCount(
				ctx,
				newRequest(*requestConfig, clients, int64(i)),
				requestConfig.Timeout,
				requestCountPerDodo,
				&responses[i],
				increase,
				&wg,
			)
		}
	}

	wg.Wait()
	streamCtxCancel()
	streamWG.Wait()
	return utils.Flatten(responses)
}

// sendRequestByCount sends a specified number of HTTP requests concurrently with a given timeout.
// It appends the responses to the provided responseData slice and sends the count of completed requests
// to the increase channel. The function terminates early if the context is canceled or if a custom
// interrupt error is encountered.
func sendRequestByCount(
	ctx context.Context,
	request *Request,
	timeout time.Duration,
	requestCount uint,
	responseData *[]*Response,
	increase chan<- int64,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for range requestCount {
		if ctx.Err() != nil {
			return
		}

		func() {
			startTime := time.Now()
			response, err := request.Send(ctx, timeout)
			completedTime := time.Since(startTime)
			if response != nil {
				defer fasthttp.ReleaseResponse(response)
			}

			if err != nil {
				if err == types.ErrInterrupt {
					return
				}
				*responseData = append(*responseData, &Response{
					Response: err.Error(),
					Time:     completedTime,
				})
				increase <- 1
				return
			}

			*responseData = append(*responseData, &Response{
				Response: strconv.Itoa(response.StatusCode()),
				Time:     completedTime,
			})
			increase <- 1
		}()
	}
}

// sendRequest continuously sends HTTP requests until the context is canceled.
// It records the response status code or error message along with the response time,
// and signals each completed request through the increase channel.
func sendRequest(
	ctx context.Context,
	request *Request,
	timeout time.Duration,
	responseData *[]*Response,
	increase chan<- int64,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		if ctx.Err() != nil {
			return
		}

		func() {
			startTime := time.Now()
			response, err := request.Send(ctx, timeout)
			completedTime := time.Since(startTime)
			if response != nil {
				defer fasthttp.ReleaseResponse(response)
			}

			if err != nil {
				if err == types.ErrInterrupt {
					return
				}
				*responseData = append(*responseData, &Response{
					Response: err.Error(),
					Time:     completedTime,
				})
				increase <- 1
				return
			}

			*responseData = append(*responseData, &Response{
				Response: strconv.Itoa(response.StatusCode()),
				Time:     completedTime,
			})
			increase <- 1
		}()
	}
}
