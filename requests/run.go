package requests

import (
	"context"
	"sync"
	"time"

	"github.com/aykhans/dodo/config"
	customerrors "github.com/aykhans/dodo/custom_errors"
	"github.com/aykhans/dodo/utils"
	"github.com/valyala/fasthttp"
)

// Run executes the HTTP requests based on the provided request configuration.
// It checks for internet connection and returns an error if there is no connection.
// If the context is canceled while checking proxies, it returns the ErrInterrupt.
// If the context is canceled while sending requests, it returns the response objects obtained so far.
func Run(ctx context.Context, requestConfig *config.RequestConfig) (Responses, error) {
	checkConnectionCtx, checkConnectionCtxCancel := context.WithTimeout(ctx, 8*time.Second)
	if !checkConnection(checkConnectionCtx) {
		checkConnectionCtxCancel()
		return nil, customerrors.ErrNoInternet
	}
	checkConnectionCtxCancel()

	clientDoFunc := getClientDoFunc(
		ctx,
		requestConfig.Timeout,
		requestConfig.Proxies,
		requestConfig.GetValidDodosCountForProxies(),
		requestConfig.GetMaxConns(fasthttp.DefaultMaxConnsPerHost),
		requestConfig.Yes,
		requestConfig.URL,
	)
	if clientDoFunc == nil {
		return nil, customerrors.ErrInterrupt
	}

	requests, err := getRequests(
		ctx,
		requestConfig.URL,
		requestConfig.Headers,
		requestConfig.Cookies,
		requestConfig.Params,
		requestConfig.Method,
		requestConfig.Body,
		requestConfig.RequestCount,
	)
	if err != nil {
		return nil, err
	}

	responses := releaseDodos(
		ctx,
		requests,
		clientDoFunc,
		requestConfig.GetValidDodosCountForRequests(),
	)
	if ctx.Err() != nil && len(responses) == 0 {
		return nil, customerrors.ErrInterrupt
	}

	return responses, nil
}

// releaseDodos sends HTTP requests concurrently using multiple "dodos" (goroutines).
//
// Parameters:
//   - ctx: The context to control the lifecycle of the requests.
//   - requests: A slice of HTTP requests to be sent.
//   - clientDoFunc: A function to execute the HTTP requests.
//   - dodosCount: The number of dodos (goroutines) to use for sending the requests.
//
// Returns:
//   - A slice of Response objects containing the results of the requests.
//
// The function divides the requests into equal parts based on the number of dodos.
// It then sends each part concurrently using a separate goroutine.
func releaseDodos(
	ctx context.Context,
	requests []*fasthttp.Request,
	clientDoFunc ClientDoFunc,
	dodosCount uint,
) Responses {
	var (
		wg                  sync.WaitGroup
		streamWG            sync.WaitGroup
		requestCountPerDodo uint
		dodosCountInt       int  = int(dodosCount)
		totalRequestCount   uint = uint(len(requests))
		requestCount        uint = 0
		responses                = make([][]*Response, dodosCount)
		increase                 = make(chan int64, totalRequestCount)
	)

	wg.Add(dodosCountInt)
	streamWG.Add(1)
	streamCtx, streamCtxCancel := context.WithCancel(context.Background())

	go streamProgress(streamCtx, &streamWG, int64(totalRequestCount), "Dodos Working🔥", increase)

	for i := range dodosCount {
		if i+1 == dodosCount {
			requestCountPerDodo = totalRequestCount - (i * totalRequestCount / dodosCount)
		} else {
			requestCountPerDodo = ((i + 1) * totalRequestCount / dodosCount) -
				(i * totalRequestCount / dodosCount)
		}

		go sendRequest(
			ctx,
			requests[requestCount:requestCount+requestCountPerDodo],
			&responses[i],
			increase,
			clientDoFunc,
			&wg,
		)
		requestCount += requestCountPerDodo
	}
	wg.Wait()
	streamCtxCancel()
	streamWG.Wait()
	return utils.Flatten(responses)
}

// sendRequest sends multiple HTTP requests concurrently and collects their responses.
//
// Parameters:
//   - ctx: The context to control cancellation and timeout.
//   - requests: A slice of pointers to fasthttp.Request objects to be sent.
//   - responseData: A pointer to a slice of *Response objects to store the results.
//   - increase: A channel to signal the completion of each request.
//   - clientDo: A function to execute the HTTP request.
//   - wg: A wait group to synchronize the completion of the requests.
//
// The function iterates over the provided requests, sending each one using the clientDo function.
// It measures the time taken for each request and appends the response data to responseData.
// If an error occurs, it appends an error response. The function signals completion through the increase channel
// and ensures proper resource cleanup by releasing requests and responses.
func sendRequest(
	ctx context.Context,
	requests []*fasthttp.Request,
	responseData *[]*Response,
	increase chan<- int64,
	clientDo ClientDoFunc,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for _, request := range requests {
		if ctx.Err() != nil {
			return
		}

		func() {
			defer fasthttp.ReleaseRequest(request)
			startTime := time.Now()
			response, err := clientDo(ctx, request)
			completedTime := time.Since(startTime)

			if err != nil {
				if err == customerrors.ErrInterrupt {
					return
				}
				*responseData = append(*responseData, &Response{
					StatusCode: 0,
					Error:      err,
					Time:       completedTime,
				})
				increase <- 1
				return
			}
			defer fasthttp.ReleaseResponse(response)

			*responseData = append(*responseData, &Response{
				StatusCode: response.StatusCode(),
				Error:      nil,
				Time:       completedTime,
			})
			increase <- 1
		}()
	}
}
