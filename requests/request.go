package requests

import (
	"context"
	"math/rand"
	"net/url"
	"time"

	"github.com/aykhans/dodo/config"
	customerrors "github.com/aykhans/dodo/custom_errors"
	"github.com/aykhans/dodo/utils"
	"github.com/valyala/fasthttp"
)

type RequestGeneratorFunc func() *fasthttp.Request

// Request represents an HTTP request to be sent using the fasthttp client.
// It isn't thread-safe and should be used by a single goroutine.
type Request struct {
	getClient  ClientGeneratorFunc
	getRequest RequestGeneratorFunc
}

// Send sends the HTTP request using the fasthttp client with a specified timeout.
// It returns the HTTP response or an error if the request fails or times out.
func (r *Request) Send(ctx context.Context, timeout time.Duration) (*fasthttp.Response, error) {
	client := r.getClient()
	request := r.getRequest()
	defer fasthttp.ReleaseRequest(request)

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

// newRequest creates a new Request instance based on the provided configuration and clients.
// It initializes a random number generator using the current time and a unique identifier (uid).
// Depending on the number of clients provided, it sets up a function to select the appropriate client.
// It also sets up a function to generate the request based on the provided configuration.
func newRequest(
	requestConfig config.RequestConfig,
	clients []*fasthttp.HostClient,
	uid int64,
) *Request {
	localRand := rand.New(rand.NewSource(time.Now().UnixNano() + uid))

	clientsCount := len(clients)
	if clientsCount < 1 {
		panic("no clients")
	}

	getClient := ClientGeneratorFunc(nil)
	if clientsCount == 1 {
		getClient = getSharedClientFuncSingle(clients[0])
	} else {
		getClient = getSharedClientFuncMultiple(clients, localRand)
	}

	getRequest := getRequestGeneratorFunc(
		requestConfig.URL,
		requestConfig.Headers,
		requestConfig.Cookies,
		requestConfig.Params,
		requestConfig.Method,
		requestConfig.Body,
		localRand,
	)

	requests := &Request{
		getClient:  getClient,
		getRequest: getRequest,
	}

	return requests
}

// getRequestGeneratorFunc returns a RequestGeneratorFunc which generates HTTP requests
// with the specified parameters.
// The function uses a local random number generator to select bodies, headers, cookies, and parameters
// if multiple options are provided.
func getRequestGeneratorFunc(
	URL *url.URL,
	Headers map[string][]string,
	Cookies map[string][]string,
	Params map[string][]string,
	Method string,
	Bodies []string,
	localRand *rand.Rand,
) RequestGeneratorFunc {
	bodiesLen := len(Bodies)
	getBody := func() string { return "" }
	if bodiesLen == 1 {
		getBody = func() string { return Bodies[0] }
	} else if bodiesLen > 1 {
		getBody = utils.RandomValueCycle(Bodies, localRand)
	}
	getHeaders := getKeyValueSetFunc(Headers, localRand)
	getCookies := getKeyValueSetFunc(Cookies, localRand)
	getParams := getKeyValueSetFunc(Params, localRand)

	return func() *fasthttp.Request {
		return newFasthttpRequest(
			URL,
			getHeaders(),
			getCookies(),
			getParams(),
			Method,
			getBody(),
		)
	}
}

// newFasthttpRequest creates a new fasthttp.Request object with the provided parameters.
// It sets the request URI, host header, headers, cookies, params, method, and body.
func newFasthttpRequest(
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

// getKeyValueSetFunc generates a function that returns a map of key-value pairs based on the provided key-value set.
// The generated function will either return fixed values or random values depending on the input.
//
// Returns:
//   - A function that returns a map of key-value pairs. If the input map contains multiple values for a key,
//     the returned function will generate random values for that key. If the input map contains a single value
//     for a key, the returned function will always return that value. If the input map is empty for a key,
//     the returned function will generate an empty string for that key.
func getKeyValueSetFunc[
	KeyValueSet map[string][]string,
	KeyValue map[string]string,
](keyValueSet KeyValueSet, localRand *rand.Rand) func() KeyValue {
	getKeyValueSlice := []map[string]func() string{}
	isRandom := false
	for key, values := range keyValueSet {
		valuesLen := len(values)

		// if values is empty, return a function that generates empty string
		// if values has only one element, return a function that generates that element
		// if values has more than one element, return a function that generates a random element
		getKeyValue := func() string { return "" }
		if valuesLen == 1 {
			getKeyValue = func() string { return values[0] }
		} else if valuesLen > 1 {
			getKeyValue = utils.RandomValueCycle(values, localRand)
			isRandom = true
		}

		getKeyValueSlice = append(
			getKeyValueSlice,
			map[string]func() string{key: getKeyValue},
		)
	}

	// if isRandom is true, return a function that generates random values,
	// otherwise return a function that generates fixed values to avoid unnecessary random number generation
	if isRandom {
		return func() KeyValue {
			keyValues := make(KeyValue, len(getKeyValueSlice))
			for _, keyValue := range getKeyValueSlice {
				for key, value := range keyValue {
					keyValues[key] = value()
				}
			}
			return keyValues
		}
	} else {
		keyValues := make(KeyValue, len(getKeyValueSlice))
		for _, keyValue := range getKeyValueSlice {
			for key, value := range keyValue {
				keyValues[key] = value()
			}
		}
		return func() KeyValue { return keyValues }
	}
}
