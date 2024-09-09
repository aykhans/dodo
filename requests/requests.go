package requests

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/aykhans/dodo/config"
	customerrors "github.com/aykhans/dodo/custom_errors"
	"github.com/aykhans/dodo/readers"
	"github.com/aykhans/dodo/utils"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
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
	dodosCount int,
	requestCount int,
) Responses {
	var (
		wg                  sync.WaitGroup
		streamWG            sync.WaitGroup
		requestCountPerDodo int
	)

	wg.Add(dodosCount)
	streamWG.Add(1)
	responses := make([][]Response, dodosCount)
	increase := make(chan int64, requestCount)

	streamCtx, streamCtxCancel := context.WithCancel(context.Background())
	go streamProgress(streamCtx, &streamWG, int64(requestCount), "Dodos Working🔥", increase)

	for i := 0; i < dodosCount; i++ {
		if i+1 == dodosCount {
			requestCountPerDodo = requestCount -
				(i * requestCount / dodosCount)
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

func sendRequest(
	ctx context.Context,
	request *fasthttp.Request,
	responseData *[]Response,
	increase chan<- int64,
	requestCount int,
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
			defer func() { increase <- 1 }()

			startTime := time.Now()
			response, err := clientDo(ctx, request)
			completedTime := time.Since(startTime)

			if err != nil {
				*responseData = append(*responseData, Response{
					StatusCode: 0,
					Error:      err,
					Time:       completedTime,
				})
				return
			}
			defer fasthttp.ReleaseResponse(response)

			*responseData = append(*responseData, Response{
				StatusCode: response.StatusCode(),
				Error:      nil,
				Time:       completedTime,
			})
		}()
	}
}

// getClientDoFunc returns a ClientDoFunc function that can be used to make HTTP requests.
//
// The function first checks if there are any proxies available. If there are, it retrieves the active proxy clients
// using the getActiveProxyClients function. If the context is canceled during this process, it returns nil.
// It then checks the number of active proxy clients and prompts the user to continue if there are none.
// If the user chooses to continue, it creates a fasthttp.HostClient with the appropriate settings and returns
// a ClientDoFunc function using the getSharedClientDoFunc function.
// If there is only one active proxy client, it uses that client to create the ClientDoFunc function.
// If there are multiple active proxy clients, it uses the getSharedRandomClientDoFunc function to create the ClientDoFunc function.
//
// If there are no proxies available, it creates a fasthttp.HostClient with the appropriate settings and returns
// a ClientDoFunc function using the getSharedClientDoFunc function.
func getClientDoFunc(
	ctx context.Context,
	timeout time.Duration,
	proxies []config.Proxy,
	dodosCount int,
	yes bool,
	URL *url.URL,
) ClientDoFunc {
	isTLS := URL.Scheme == "https"
	if len(proxies) > 0 {
		activeProxyClients := getActiveProxyClients(
			ctx, proxies, timeout, dodosCount, URL,
		)
		if ctx.Err() != nil {
			return nil
		}
		activeProxyClientsCount := len(activeProxyClients)
		var yesOrNoMessage string
		var yesOrNoDefault bool
		if activeProxyClientsCount == 0 {
			yesOrNoDefault = false
			yesOrNoMessage = utils.Colored(
				utils.Colors.Yellow,
				"No active proxies found. Do you want to continue?",
			)
		} else {
			yesOrNoMessage = utils.Colored(
				utils.Colors.Yellow,
				fmt.Sprintf(
					"Found %d active proxies. Do you want to continue?",
					activeProxyClientsCount,
				),
			)
		}
		if !yes {
			response := readers.CLIYesOrNoReader("\n"+yesOrNoMessage, yesOrNoDefault)
			if !response {
				utils.PrintAndExit("Exiting...")
			}
		}
		fmt.Println()
		if activeProxyClientsCount == 0 {
			client := &fasthttp.HostClient{
				IsTLS:               isTLS,
				Addr:                URL.Host,
				MaxIdleConnDuration: timeout,
				MaxConnDuration:     timeout,
				WriteTimeout:        timeout,
				ReadTimeout:         timeout,
			}
			return getSharedClientDoFunc(client, timeout)
		} else if activeProxyClientsCount == 1 {
			client := &activeProxyClients[0]
			return getSharedClientDoFunc(client, timeout)
		}
		return getSharedRandomClientDoFunc(
			activeProxyClients,
			activeProxyClientsCount,
			timeout,
		)
	}

	client := &fasthttp.HostClient{
		IsTLS:               isTLS,
		Addr:                URL.Host,
		MaxIdleConnDuration: timeout,
		MaxConnDuration:     timeout,
		WriteTimeout:        timeout,
		ReadTimeout:         timeout,
	}
	return getSharedClientDoFunc(client, timeout)
}

// getActiveProxyClients divides the proxies into slices based on the number of dodos and
// launches goroutines to find active proxy clients for each slice.
// It uses a progress tracker to monitor the progress of the search.
// Once all goroutines have completed, the function waits for them to finish and
// returns a flattened slice of active proxy clients.
func getActiveProxyClients(
	ctx context.Context,
	proxies []config.Proxy,
	timeout time.Duration,
	dodosCount int,
	URL *url.URL,
) []fasthttp.HostClient {
	activeProxyClientsArray := make([][]fasthttp.HostClient, dodosCount)
	proxiesCount := len(proxies)

	var (
		wg       sync.WaitGroup
		streamWG sync.WaitGroup
	)
	wg.Add(dodosCount)
	streamWG.Add(1)
	var proxiesSlice []config.Proxy
	increase := make(chan int64, proxiesCount)

	streamCtx, streamCtxCancel := context.WithCancel(context.Background())
	go streamProgress(streamCtx, &streamWG, int64(proxiesCount), "Searching for active proxies🌐", increase)

	for i := 0; i < dodosCount; i++ {
		if i+1 == dodosCount {
			proxiesSlice = proxies[i*proxiesCount/dodosCount:]
		} else {
			proxiesSlice = proxies[i*proxiesCount/dodosCount : (i+1)*proxiesCount/dodosCount]
		}
		go findActiveProxyClients(
			ctx,
			proxiesSlice,
			timeout,
			&activeProxyClientsArray[i],
			increase,
			URL,
			&wg,
		)
	}
	wg.Wait()
	streamCtxCancel()
	streamWG.Wait()
	return utils.Flatten(activeProxyClientsArray)
}

// findActiveProxyClients checks a list of proxies to determine which ones are active
// and appends the active ones to the provided activeProxyClients slice.
//
// Parameters:
//   - ctx: The context to control cancellation and timeout.
//   - proxies: A slice of Proxy configurations to be checked.
//   - timeout: The duration to wait for each proxy check before timing out.
//   - activeProxyClients: A pointer to a slice where active proxy clients will be appended.
//   - increase: A channel to signal the increase of checked proxies count.
//   - URL: The URL to be used for checking the proxies.
//   - wg: A WaitGroup to signal when the function is done.
//
// The function sends a GET request to each proxy using the provided URL. If the proxy
// responds with a status code of 200, it is considered active and added to the activeProxyClients slice.
// The function respects the context's cancellation and timeout settings.
func findActiveProxyClients(
	ctx context.Context,
	proxies []config.Proxy,
	timeout time.Duration,
	activeProxyClients *[]fasthttp.HostClient,
	increase chan<- int64,
	URL *url.URL,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.SetRequestURI(config.ProxyCheckURL)
	request.Header.SetMethod("GET")

	for _, proxy := range proxies {
		if ctx.Err() != nil {
			return
		}

		func() {
			defer func() { increase <- 1 }()

			response := fasthttp.AcquireResponse()
			defer fasthttp.ReleaseResponse(response)

			dialFunc, err := getDialFunc(&proxy, timeout)
			if err != nil {
				return
			}
			client := &fasthttp.Client{
				Dial: dialFunc,
			}
			defer client.CloseIdleConnections()

			ch := make(chan error)
			go func() {
				err := client.DoTimeout(request, response, timeout)
				ch <- err
			}()
			select {
			case err := <-ch:
				if err != nil {
					return
				}
				break
			case <-time.After(timeout):
				return
			case <-ctx.Done():
				return
			}

			isTLS := URL.Scheme == "https"
			addr := URL.Host
			if isTLS && URL.Port() == "" {
				addr += ":443"
			}
			if response.StatusCode() == 200 {
				*activeProxyClients = append(
					*activeProxyClients,
					fasthttp.HostClient{
						IsTLS:               isTLS,
						Addr:                addr,
						Dial:                dialFunc,
						MaxIdleConnDuration: timeout,
						MaxConnDuration:     timeout,
						WriteTimeout:        timeout,
						ReadTimeout:         timeout,
					},
				)
			}
		}()
	}
}

// getDialFunc returns a fasthttp.DialFunc based on the provided proxy configuration.
// It takes a pointer to a config.Proxy struct as input and returns a fasthttp.DialFunc and an error.
// The function parses the proxy URL, determines the scheme (socks5, socks5h, http, or https),
// and creates a dialer accordingly. If the proxy URL is invalid or the scheme is not supported,
// it returns an error.
func getDialFunc(proxy *config.Proxy, timeout time.Duration) (fasthttp.DialFunc, error) {
	parsedProxyURL, err := url.Parse(proxy.URL)
	if err != nil {
		return nil, err
	}

	var dialer fasthttp.DialFunc
	if parsedProxyURL.Scheme == "socks5" || parsedProxyURL.Scheme == "socks5h" {
		if proxy.Username != "" {
			dialer = fasthttpproxy.FasthttpSocksDialer(
				fmt.Sprintf(
					"%s://%s:%s@%s",
					parsedProxyURL.Scheme,
					proxy.Username,
					proxy.Password,
					parsedProxyURL.Host,
				),
			)
		} else {
			dialer = fasthttpproxy.FasthttpSocksDialer(
				fmt.Sprintf(
					"%s://%s",
					parsedProxyURL.Scheme,
					parsedProxyURL.Host,
				),
			)
		}
	} else if parsedProxyURL.Scheme == "http" {
		if proxy.Username != "" {
			dialer = fasthttpproxy.FasthttpHTTPDialerTimeout(
				fmt.Sprintf(
					"%s:%s@%s",
					proxy.Username, proxy.Password, parsedProxyURL.Host,
				),
				timeout,
			)
		} else {
			dialer = fasthttpproxy.FasthttpHTTPDialerTimeout(
				parsedProxyURL.Host,
				timeout,
			)
		}
	} else {
		return nil, err
	}
	return dialer, nil
}

// getSharedRandomClientDoFunc is equivalent to getSharedClientDoFunc but uses a random client from the provided slice.
func getSharedRandomClientDoFunc(
	clients []fasthttp.HostClient,
	clientsCount int,
	timeout time.Duration,
) ClientDoFunc {
	return func(ctx context.Context, request *fasthttp.Request) (*fasthttp.Response, error) {
		client := &clients[rand.Intn(clientsCount)]
		defer client.CloseIdleConnections()
		response := fasthttp.AcquireResponse()
		ch := make(chan error)
		go func() {
			err := client.DoTimeout(request, response, timeout)
			ch <- err
		}()
		select {
		case err := <-ch:
			if err != nil {
				fasthttp.ReleaseResponse(response)
				return nil, err
			}
			return response, nil
		case <-time.After(timeout):
			fasthttp.ReleaseResponse(response)
			return nil, customerrors.ErrTimeout
		case <-ctx.Done():
			return nil, customerrors.ErrInterrupt
		}
	}
}

// getSharedClientDoFunc is a function that returns a ClientDoFunc, which is a function type used for making HTTP requests using a shared client.
// It takes a client of type *fasthttp.HostClient and a timeout of type time.Duration as input parameters.
// The returned ClientDoFunc function can be used to make an HTTP request with the given client and timeout.
// It takes a context.Context and a *fasthttp.Request as input parameters and returns a *fasthttp.Response and an error.
// The function internally creates a new response using fasthttp.AcquireResponse() and a channel to handle errors.
// It then spawns a goroutine to execute the client.DoTimeout() method with the given request, response, and timeout.
// The function uses a select statement to handle three cases:
//   - If an error is received from the channel, it checks if the error is not nil. If it's not nil, it releases the response and returns nil and the error.
//     Otherwise, it returns the response and nil.
//   - If the timeout duration is reached, it releases the response and returns nil and a custom timeout error.
//   - If the context is canceled, it returns nil and a custom interrupt error.
//
// The function ensures that idle connections are closed by calling client.CloseIdleConnections() using a defer statement.
func getSharedClientDoFunc(
	client *fasthttp.HostClient,
	timeout time.Duration,
) ClientDoFunc {
	return func(ctx context.Context, request *fasthttp.Request) (*fasthttp.Response, error) {
		defer client.CloseIdleConnections()
		response := fasthttp.AcquireResponse()
		ch := make(chan error)
		go func() {
			err := client.DoTimeout(request, response, timeout)
			ch <- err
		}()
		select {
		case err := <-ch:
			if err != nil {
				fasthttp.ReleaseResponse(response)
				return nil, err
			}
			return response, nil
		case <-time.After(timeout):
			fasthttp.ReleaseResponse(response)
			return nil, customerrors.ErrTimeout
		case <-ctx.Done():
			return nil, customerrors.ErrInterrupt
		}
	}
}

// newRequest creates a new fasthttp.Request object with the provided parameters.
// It sets the request URI, host header, headers, cookies, params, method, and body.
func newRequest(
	URL *url.URL,
	Headers map[string]string,
	Cookies map[string]string,
	Params map[string]string,
	Method string,
	Body string,
) *fasthttp.Request {
	request := fasthttp.AcquireRequest()
	request.SetRequestURI(URL.Path)

	// Set the host of the request to the host header
	// If the host header is not set, the request will fail
	// If there is host header in the headers, it will be overwritten
	request.Header.Set("Host", URL.Host)
	setRequestHeaders(request, Headers)
	setRequestCookies(request, Cookies)
	setRequestParams(request, Params)
	setRequestMethod(request, Method)
	setRequestBody(request, Body)
	if URL.Scheme == "https" {
		request.URI().SetScheme("https")
	}

	return request
}

// setRequestHeaders sets the headers of the given request with the provided key-value pairs.
func setRequestHeaders(req *fasthttp.Request, headers map[string]string) {
	req.Header.Set("User-Agent", config.DefaultUserAgent)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

// setRequestCookies sets the cookies in the given request.
func setRequestCookies(req *fasthttp.Request, cookies map[string]string) {
	for key, value := range cookies {
		req.Header.SetCookie(key, value)
	}
}

// setRequestParams sets the query parameters of the given request based on the provided map of key-value pairs.
func setRequestParams(req *fasthttp.Request, params map[string]string) {
	urlParams := url.Values{}
	for key, value := range params {
		urlParams.Add(key, value)
	}
	req.URI().SetQueryString(urlParams.Encode())
}

// setRequestMethod sets the HTTP request method for the given request.
func setRequestMethod(req *fasthttp.Request, method string) {
	req.Header.SetMethod(method)
}

// setRequestBody sets the request body of the given fasthttp.Request object.
// The body parameter is a string that will be converted to a byte slice and set as the request body.
func setRequestBody(req *fasthttp.Request, body string) {
	req.SetBody([]byte(body))
}

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
