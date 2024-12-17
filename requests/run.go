package requests

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/aykhans/dodo/config"
	customerrors "github.com/aykhans/dodo/custom_errors"
	"github.com/aykhans/dodo/utils"
	"github.com/valyala/fasthttp"
)

// Run executes the main logic for processing requests based on the provided configuration.
// It first checks for an internet connection with a timeout context. If no connection is found,
// it returns an error. Then, it initializes clients based on the request configuration and
// releases the dodos. If the context is canceled and no responses are collected, it returns an interrupt error.
//
// Parameters:
//   - ctx: The context for managing request lifecycle and cancellation.
//   - requestConfig: The configuration for the request, including timeout, proxies, and other settings.
//
// Returns:
//   - Responses: A collection of responses from the executed requests.
//   - error: An error if the operation fails, such as no internet connection or an interrupt.
func Run(ctx context.Context, requestConfig *config.RequestConfig) (Responses, error) {
	checkConnectionCtx, checkConnectionCtxCancel := context.WithTimeout(ctx, 8*time.Second)
	if !checkConnection(checkConnectionCtx) {
		checkConnectionCtxCancel()
		return nil, customerrors.ErrNoInternet
	}
	checkConnectionCtxCancel()

	clients := getClients(
		ctx,
		requestConfig.Timeout,
		requestConfig.Proxies,
		requestConfig.GetValidDodosCountForProxies(),
		requestConfig.GetMaxConns(fasthttp.DefaultMaxConnsPerHost),
		requestConfig.Yes,
		requestConfig.NoProxyCheck,
		requestConfig.URL,
	)
	if clients == nil {
		return nil, customerrors.ErrInterrupt
	}

	responses := releaseDodos(ctx, requestConfig, clients)
	if ctx.Err() != nil && len(responses) == 0 {
		return nil, customerrors.ErrInterrupt
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
		dodosCount          uint = requestConfig.GetValidDodosCountForRequests()
		dodosCountInt       int  = int(dodosCount)
		requestCount        uint = uint(requestConfig.RequestCount)
		responses                = make([][]*Response, dodosCount)
		increase                 = make(chan int64, requestCount)
	)

	wg.Add(dodosCountInt)
	streamWG.Add(1)
	streamCtx, streamCtxCancel := context.WithCancel(context.Background())

	go streamProgress(streamCtx, &streamWG, int64(requestCount), "Dodos Working🔥", increase)

	for i := range dodosCount {
		if i+1 == dodosCount {
			requestCountPerDodo = requestCount - (i * requestCount / dodosCount)
		} else {
			requestCountPerDodo = ((i + 1) * requestCount / dodosCount) -
				(i * requestCount / dodosCount)
		}

		go sendRequest(
			ctx,
			newRequest(*requestConfig, clients, int64(i)),
			requestConfig.Timeout,
			requestCountPerDodo,
			&responses[i],
			increase,
			&wg,
		)
	}
	wg.Wait()
	streamCtxCancel()
	streamWG.Wait()
	return utils.Flatten(responses)
}

// sendRequest sends a specified number of HTTP requests concurrently with a given timeout.
// It appends the responses to the provided responseData slice and sends the count of completed requests
// to the increase channel. The function terminates early if the context is canceled or if a custom
// interrupt error is encountered.
func sendRequest(
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
				if err == customerrors.ErrInterrupt {
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
