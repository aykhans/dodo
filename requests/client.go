package requests

import (
	"context"
	"errors"
	"math/rand"
	"net/url"
	"time"

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
	proxies []url.URL,
	maxConns uint,
	URL url.URL,
) []*fasthttp.HostClient {
	isTLS := URL.Scheme == "https"

	if proxiesLen := len(proxies); proxiesLen > 0 {
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

// getDialFunc returns the appropriate fasthttp.DialFunc based on the provided proxy URL scheme.
// It supports SOCKS5 ('socks5' or 'socks5h') and HTTP ('http') proxy schemes.
// For HTTP proxies, the timeout parameter determines connection timeouts.
// Returns an error if the proxy scheme is unsupported.
func getDialFunc(proxy *url.URL, timeout time.Duration) (fasthttp.DialFunc, error) {
	var dialer fasthttp.DialFunc

	switch proxy.Scheme {
	case "socks5", "socks5h":
		dialer = fasthttpproxy.FasthttpSocksDialerDualStack(proxy.String())
	case "http":
		dialer = fasthttpproxy.FasthttpHTTPDialerDualStackTimeout(proxy.String(), timeout)
	default:
		return nil, errors.New("unsupported proxy scheme")
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
