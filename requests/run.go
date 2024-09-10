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

	request := newRequest(
		requestConfig.URL,
		requestConfig.Headers,
		requestConfig.Cookies,
		requestConfig.Params,
		requestConfig.Method,
		requestConfig.Body,
	)
	defer fasthttp.ReleaseRequest(request)
	responses := releaseDodos(
		ctx,
		request,
		clientDoFunc,
		requestConfig.GetValidDodosCountForRequests(),
		requestConfig.RequestCount,
	)
	if ctx.Err() != nil && len(responses) == 0 {
		return nil, customerrors.ErrInterrupt
	}

	return responses, nil
}

// releaseDodos sends multiple HTTP requests concurrently using multiple "dodos" (goroutines).
// It takes a mainRequest as the base request, timeout duration for each request, clientDoFunc for customizing the client behavior,
// dodosCount as the number of goroutines to be used, and requestCount as the total number of requests to be sent.
// It returns the responses received from all the requests.
func releaseDodos(
	ctx context.Context,
	mainRequest *fasthttp.Request,
	clientDoFunc ClientDoFunc,
	dodosCount uint,
	requestCount uint,
) Responses {
	var (
		wg                  sync.WaitGroup
		streamWG            sync.WaitGroup
		requestCountPerDodo uint
		dodosCountInt       = int(dodosCount)
	)

	wg.Add(dodosCountInt)
	streamWG.Add(1)
	responses := make([][]Response, dodosCount)
	increase := make(chan int64, requestCount)

	streamCtx, streamCtxCancel := context.WithCancel(context.Background())
	go streamProgress(streamCtx, &streamWG, int64(requestCount), "Dodos Working🔥", increase)

	for i := range dodosCount {
		if i+1 == dodosCount {
			requestCountPerDodo = requestCount - (i * requestCount / dodosCount)
		} else {
			requestCountPerDodo = ((i + 1) * requestCount / dodosCount) -
				(i * requestCount / dodosCount)
		}
		dodoSpecificRequest := &fasthttp.Request{}
		mainRequest.CopyTo(dodoSpecificRequest)

		go sendRequest(
			ctx,
			dodoSpecificRequest,
			&responses[i],
			increase,
			requestCountPerDodo,
			clientDoFunc,
			&wg,
		)
	}
	wg.Wait()
	streamCtxCancel()
	streamWG.Wait()
	return utils.Flatten(responses)
}

// sendRequest sends an HTTP request using the provided clientDo function and handles the response.
//
// Parameters:
//   - ctx: The context to control cancellation and timeout.
//   - request: The HTTP request to be sent.
//   - responseData: A slice to store the response data.
//   - increase: A channel to signal the completion of a request.
//   - requestCount: The number of requests to be sent.
//   - clientDo: A function to execute the HTTP request.
//   - wg: A wait group to signal the completion of the function.
//
// The function sends the specified number of requests, handles errors, and appends the response data
// to the responseData slice.
func sendRequest(
	ctx context.Context,
	request *fasthttp.Request,
	responseData *[]Response,
	increase chan<- int64,
	requestCount uint,
	clientDo ClientDoFunc,
	wg *sync.WaitGroup,
) {
	defer fasthttp.ReleaseRequest(request)
	defer wg.Done()

	for range requestCount {
		if ctx.Err() != nil {
			return
		}

		func() {
			startTime := time.Now()
			response, err := clientDo(ctx, request)
			completedTime := time.Since(startTime)

			if err != nil {
				if err == customerrors.ErrInterrupt {
					return
				}
				*responseData = append(*responseData, Response{
					StatusCode: 0,
					Error:      err,
					Time:       completedTime,
				})
				increase <- 1
				return
			}
			defer fasthttp.ReleaseResponse(response)

			*responseData = append(*responseData, Response{
				StatusCode: response.StatusCode(),
				Error:      nil,
				Time:       completedTime,
			})
			increase <- 1
		}()
	}
}
