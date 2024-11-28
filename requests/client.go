package requests

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/aykhans/dodo/config"
	"github.com/aykhans/dodo/readers"
	"github.com/aykhans/dodo/utils"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type ClientGeneratorFunc func() *fasthttp.HostClient

// getClients initializes and returns a slice of fasthttp.HostClient based on the provided parameters.
// It can either return clients with proxies or a single client without proxies.
func getClients(
	ctx context.Context,
	timeout time.Duration,
	proxies []config.Proxy,
	dodosCount uint,
	maxConns uint,
	yes bool,
	noProxyCheck bool,
	URL *url.URL,
) []*fasthttp.HostClient {
	isTLS := URL.Scheme == "https"

	if proxiesLen := len(proxies); proxiesLen > 0 {
		// If noProxyCheck is true, we will return the clients without checking the proxies.
		if noProxyCheck {
			clients := make([]*fasthttp.HostClient, 0, proxiesLen)
			addr := URL.Host
			if isTLS && URL.Port() == "" {
				addr += ":443"
			}

			for _, proxy := range proxies {
				dialFunc, err := getDialFunc(&proxy, timeout)
				if err != nil {
					continue
				}

				clients = append(clients, &fasthttp.HostClient{
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
			return clients
		}

		// Else, we will check the proxies and return the active ones.
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
		if activeProxyClientsCount > 0 {
			return activeProxyClients
		}
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
	return []*fasthttp.HostClient{client}
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

// getSharedClientFuncMultiple returns a ClientGeneratorFunc that cycles through a list of fasthttp.HostClient instances.
// The function uses a local random number generator to determine the starting index and stop index for cycling through the clients.
// The returned function isn't thread-safe and should be used in a single-threaded context.
func getSharedClientFuncMultiple(clients []*fasthttp.HostClient, localRand *rand.Rand) ClientGeneratorFunc {
	return utils.RandomValueCycle(clients, localRand)
}

// getSharedClientFuncSingle returns a ClientGeneratorFunc that always returns the provided fasthttp.HostClient instance.
// This can be useful for sharing a single client instance across multiple requests.
func getSharedClientFuncSingle(client *fasthttp.HostClient) ClientGeneratorFunc {
	return func() *fasthttp.HostClient {
		return client
	}
}
