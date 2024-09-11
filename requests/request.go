package requests

import (
	"context"
	"math/rand"
	"net/url"

	"github.com/aykhans/dodo/config"
	customerrors "github.com/aykhans/dodo/custom_errors"
	"github.com/valyala/fasthttp"
)

// getRequests generates a list of HTTP requests based on the provided parameters.
//
// Parameters:
//   - ctx: The context to control cancellation and deadlines.
//   - URL: The base URL for the requests.
//   - Headers: A map of headers to include in each request.
//   - Cookies: A map of cookies to include in each request.
//   - Params: A map of query parameters to include in each request.
//   - Method: The HTTP method to use for the requests (e.g., GET, POST).
//   - Bodies: A list of request bodies to cycle through for each request.
//   - RequestCount: The number of requests to generate.
//
// Returns:
//   - A list of fasthttp.Request objects based on the provided parameters.
//   - An error if the context is canceled.
func getRequests(
	ctx context.Context,
	URL *url.URL,
	Headers map[string][]string,
	Cookies map[string][]string,
	Params map[string][]string,
	Method string,
	Bodies []string,
	RequestCount uint,
) ([]*fasthttp.Request, error) {
	requests := make([]*fasthttp.Request, 0, RequestCount)

	bodiesLen := len(Bodies)
	getBody := func() string { return "" }
	if bodiesLen == 1 {
		getBody = func() string { return Bodies[0] }
	} else if bodiesLen > 1 {
		currentIndex := 0
		stopIndex := bodiesLen - 1

		getBody = func() string {
			body := Bodies[currentIndex%bodiesLen]
			if currentIndex == stopIndex {
				currentIndex = rand.Intn(bodiesLen)
				stopIndex = currentIndex - 1
			} else {
				currentIndex = (currentIndex + 1) % bodiesLen
			}
			return body
		}
	}
	getHeaders := getKeyValueSetFunc(Headers)
	getCookies := getKeyValueSetFunc(Cookies)
	getParams := getKeyValueSetFunc(Params)

	for range RequestCount {
		if ctx.Err() != nil {
			return nil, customerrors.ErrInterrupt
		}
		request := newRequest(
			URL,
			getHeaders(),
			getCookies(),
			getParams(),
			Method,
			getBody(),
		)
		requests = append(requests, request)
	}

	return requests, nil
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
](keyValueSet KeyValueSet) func() KeyValue {
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
			currentIndex := 0
			stopIndex := valuesLen - 1

			getKeyValue = func() string {
				value := values[currentIndex%valuesLen]
				if currentIndex == stopIndex {
					currentIndex = rand.Intn(valuesLen)
					stopIndex = currentIndex - 1
				} else {
					currentIndex = (currentIndex + 1) % valuesLen
				}
				return value
			}

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
