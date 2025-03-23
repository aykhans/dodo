package requests

import (
	"context"
	"math/rand"
	"net/url"
	"time"

	"github.com/aykhans/dodo/config"
	"github.com/aykhans/dodo/types"
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
		return nil, types.ErrTimeout
	case <-ctx.Done():
		return nil, types.ErrInterrupt
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
		requestConfig.Params,
		requestConfig.Headers,
		requestConfig.Cookies,
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

// getRequestGeneratorFunc returns a RequestGeneratorFunc which generates HTTP requests with the specified parameters.
// The function uses a local random number generator to select bodies, headers, cookies, and parameters if multiple options are provided.
func getRequestGeneratorFunc(
	URL url.URL,
	params types.Params,
	headers types.Headers,
	cookies types.Cookies,
	method string,
	bodies []string,
	localRand *rand.Rand,
) RequestGeneratorFunc {
	bodiesLen := len(bodies)
	getBody := func() string { return "" }
	if bodiesLen == 1 {
		getBody = func() string { return bodies[0] }
	} else if bodiesLen > 1 {
		getBody = utils.RandomValueCycle(bodies, localRand)
	}

	getParams := getKeyValueGeneratorFunc(params, localRand)
	getHeaders := getKeyValueGeneratorFunc(headers, localRand)
	getCookies := getKeyValueGeneratorFunc(cookies, localRand)

	return func() *fasthttp.Request {
		return newFasthttpRequest(
			URL,
			getParams(),
			getHeaders(),
			getCookies(),
			method,
			getBody(),
		)
	}
}

// newFasthttpRequest creates a new fasthttp.Request object with the provided parameters.
// It sets the request URI, host header, headers, cookies, params, method, and body.
func newFasthttpRequest(
	URL url.URL,
	params []types.KeyValue[string, string],
	headers []types.KeyValue[string, string],
	cookies []types.KeyValue[string, string],
	method string,
	body string,
) *fasthttp.Request {
	request := fasthttp.AcquireRequest()
	request.SetRequestURI(URL.Path)

	// Set the host of the request to the host header
	// If the host header is not set, the request will fail
	// If there is host header in the headers, it will be overwritten
	request.Header.SetHost(URL.Host)
	setRequestParams(request, params)
	setRequestHeaders(request, headers)
	setRequestCookies(request, cookies)
	setRequestMethod(request, method)
	setRequestBody(request, body)
	if URL.Scheme == "https" {
		request.URI().SetScheme("https")
	}

	return request
}

// setRequestParams adds the query parameters of the given request based on the provided key-value pairs.
func setRequestParams(req *fasthttp.Request, params []types.KeyValue[string, string]) {
	for _, param := range params {
		req.URI().QueryArgs().Add(param.Key, param.Value)
	}
}

// setRequestHeaders adds the headers of the given request with the provided key-value pairs.
func setRequestHeaders(req *fasthttp.Request, headers []types.KeyValue[string, string]) {
	for _, header := range headers {
		req.Header.Add(header.Key, header.Value)
	}
}

// setRequestCookies adds the cookies of the given request with the provided key-value pairs.
func setRequestCookies(req *fasthttp.Request, cookies []types.KeyValue[string, string]) {
	for _, cookie := range cookies {
		req.Header.Add("Cookie", cookie.Key+"="+cookie.Value)
	}
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

// getKeyValueGeneratorFunc creates a function that generates key-value pairs for HTTP requests.
// It takes a slice of key-value pairs where each key maps to a slice of possible values,
// and a random number generator.
//
// If any key has multiple possible values, the function will randomly select one value for each
// call (using the provided random number generator). If all keys have at most one value, the
// function will always return the same set of key-value pairs for efficiency.
func getKeyValueGeneratorFunc[
	T []types.KeyValue[string, string],
](
	keyValueSlice []types.KeyValue[string, []string],
	localRand *rand.Rand,
) func() T {
	getKeyValueSlice := []map[string]func() string{}
	isRandom := false

	for _, kv := range keyValueSlice {
		valuesLen := len(kv.Value)

		getValueFunc := func() string { return "" }
		if valuesLen == 1 {
			getValueFunc = func() string { return kv.Value[0] }
		} else if valuesLen > 1 {
			getValueFunc = utils.RandomValueCycle(kv.Value, localRand)
			isRandom = true
		}

		getKeyValueSlice = append(
			getKeyValueSlice,
			map[string]func() string{kv.Key: getValueFunc},
		)
	}

	if isRandom {
		return func() T {
			keyValues := make(T, len(getKeyValueSlice))
			for i, keyValue := range getKeyValueSlice {
				for key, value := range keyValue {
					keyValues[i] = types.KeyValue[string, string]{
						Key:   key,
						Value: value(),
					}
				}
			}
			return keyValues
		}
	} else {
		keyValues := make(T, len(getKeyValueSlice))
		for i, keyValue := range getKeyValueSlice {
			for key, value := range keyValue {
				keyValues[i] = types.KeyValue[string, string]{
					Key:   key,
					Value: value(),
				}
			}
		}
		return func() T { return keyValues }
	}
}
