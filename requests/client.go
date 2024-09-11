package requests

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/aykhans/dodo/config"
	customerrors "github.com/aykhans/dodo/custom_errors"
	"github.com/aykhans/dodo/readers"
	"github.com/aykhans/dodo/utils"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type ClientDoFunc func(ctx context.Context, request *fasthttp.Request) (*fasthttp.Response, error)

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
	dodosCount uint,
	maxConns uint,
	yes bool,
	URL *url.URL,
) ClientDoFunc {
	isTLS := URL.Scheme == "https"

	if len(proxies) > 0 {
		activeProxyClients := getActiveProxyClients(
			ctx, proxies, timeout, dodosCount, maxConns, URL,
		)
		if ctx.Err() != nil {
			return nil
		}
		activeProxyClientsCount := uint(len(activeProxyClients))
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
				MaxConns:            int(maxConns),
				IsTLS:               isTLS,
				Addr:                URL.Host,
				MaxIdleConnDuration: timeout,
				MaxConnDuration:     timeout,
				WriteTimeout:        timeout,
				ReadTimeout:         timeout,
			}
			return getSharedClientDoFunc(client, timeout)
		} else if activeProxyClientsCount == 1 {
			client := activeProxyClients[0]
			return getSharedClientDoFunc(client, timeout)
		}
		return getSharedRandomClientDoFunc(
			activeProxyClients,
			activeProxyClientsCount,
			timeout,
		)
	}

	client := &fasthttp.HostClient{
		MaxConns:            int(maxConns),
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
	dodosCount uint,
	maxConns uint,
	URL *url.URL,
) []*fasthttp.HostClient {
	activeProxyClientsArray := make([][]*fasthttp.HostClient, dodosCount)
	proxiesCount := len(proxies)
	dodosCountInt := int(dodosCount)

	var (
		wg       sync.WaitGroup
		streamWG sync.WaitGroup
	)
	wg.Add(dodosCountInt)
	streamWG.Add(1)
	var proxiesSlice []config.Proxy
	increase := make(chan int64, proxiesCount)

	streamCtx, streamCtxCancel := context.WithCancel(context.Background())
	go streamProgress(streamCtx, &streamWG, int64(proxiesCount), "Searching for active proxies🌐", increase)

	for i := range dodosCountInt {
		if i+1 == dodosCountInt {
			proxiesSlice = proxies[i*proxiesCount/dodosCountInt:]
		} else {
			proxiesSlice = proxies[i*proxiesCount/dodosCountInt : (i+1)*proxiesCount/dodosCountInt]
		}
		go findActiveProxyClients(
			ctx,
			proxiesSlice,
			timeout,
			&activeProxyClientsArray[i],
			increase,
			maxConns,
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
	activeProxyClients *[]*fasthttp.HostClient,
	increase chan<- int64,
	maxConns uint,
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
					&fasthttp.HostClient{
						MaxConns:            int(maxConns),
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
	clients []*fasthttp.HostClient,
	clientsCount uint,
	timeout time.Duration,
) ClientDoFunc {
	clientsCountInt := int(clientsCount)
	return func(ctx context.Context, request *fasthttp.Request) (*fasthttp.Response, error) {
		client := clients[rand.Intn(clientsCountInt)]
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
