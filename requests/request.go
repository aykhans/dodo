package requests

import (
	"bytes"
	"context"
	"math/rand"
	"net/url"
	"text/template"
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

type keyValueGenerator struct {
	key   func() string
	value func() string
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
	getParams := getKeyValueGeneratorFunc(params, localRand)
	getHeaders := getKeyValueGeneratorFunc(headers, localRand)
	getCookies := getKeyValueGeneratorFunc(cookies, localRand)
	getBody := getBodyValueFunc(bodies, utils.NewFuncMapGenerator(localRand), localRand)

	return func() *fasthttp.Request {
		body, contentType := getBody()
		headers := getHeaders()
		if contentType != "" {
			headers = append(headers, types.KeyValue[string, string]{
				Key:   "Content-Type",
				Value: contentType,
			})
		}

		return newFasthttpRequest(
			URL,
			getParams(),
			headers,
			getCookies(),
			method,
			body,
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
	keyValueGenerators := make([]keyValueGenerator, len(keyValueSlice))

	funcMap := *utils.NewFuncMapGenerator(localRand).GetFuncMap()

	for i, kv := range keyValueSlice {
		keyValueGenerators[i] = keyValueGenerator{
			key:   getKeyFunc(kv.Key, funcMap),
			value: getValueFunc(kv.Value, funcMap, localRand),
		}
	}

	return func() T {
		keyValues := make(T, len(keyValueGenerators))
		for i, keyValue := range keyValueGenerators {
			keyValues[i] = types.KeyValue[string, string]{
				Key:   keyValue.key(),
				Value: keyValue.value(),
			}
		}
		return keyValues
	}
}

// getKeyFunc creates a function that processes a key string through Go's template engine.
// It takes a key string and a template.FuncMap containing the available template functions.
//
// The returned function, when called, will execute the template with the given key and return
// the processed string result. If template parsing fails, the returned function will always
// return an empty string.
//
// This enables dynamic generation of keys that can include template directives and functions.
func getKeyFunc(key string, funcMap template.FuncMap) func() string {
	t, err := template.New("default").Funcs(funcMap).Parse(key)
	if err != nil {
		return func() string { return "" }
	}

	return func() string {
		var buf bytes.Buffer
		_ = t.Execute(&buf, nil)
		return buf.String()
	}
}

// getValueFunc creates a function that randomly selects and processes a value from a slice of strings
// through Go's template engine.
//
// Parameters:
//   - values: A slice of string templates that can contain template directives
//   - funcMap: A template.FuncMap containing all available template functions
//   - localRand: A random number generator for consistent randomization
//
// The returned function, when called, will:
//  1. Select a random template from the values slice
//  2. Execute the selected template
//  3. Return the processed string result
//
// If a selected template is nil (due to earlier parsing failure), the function will return an empty string.
// This enables dynamic generation of values with randomized selection from multiple templates.
func getValueFunc(
	values []string,
	funcMap template.FuncMap,
	localRand *rand.Rand,
) func() string {
	templates := make([]*template.Template, len(values))

	for i, value := range values {
		t, err := template.New("default").Funcs(funcMap).Parse(value)
		if err != nil {
			templates[i] = nil
		}
		templates[i] = t
	}

	randomTemplateFunc := utils.RandomValueCycle(templates, localRand)

	return func() string {
		if tmpl := randomTemplateFunc(); tmpl == nil {
			return ""
		} else {
			var buf bytes.Buffer
			_ = tmpl.Execute(&buf, nil)
			return buf.String()
		}
	}
}

// getBodyValueFunc creates a function that randomly selects and processes a request body from a slice of templates.
// It returns a closure that generates both the body content and the appropriate Content-Type header value.
//
// Parameters:
//   - values: A slice of string templates that can contain template directives for request bodies
//   - funcMapGenerator: Provides template functions and content type information
//   - localRand: A random number generator for consistent randomization
//
// The returned function, when called, will:
//  1. Select a random body template from the values slice
//  2. Execute the selected template with available template functions
//  3. Return both the processed body string and the appropriate Content-Type header value
//
// If the selected template is nil (due to earlier parsing failure), the function will return
// empty strings for both the body and Content-Type.
//
// This enables dynamic generation of request bodies with proper content type headers.
func getBodyValueFunc(
	values []string,
	funcMapGenerator *utils.FuncMapGenerator,
	localRand *rand.Rand,
) func() (string, string) {
	templates := make([]*template.Template, len(values))

	for i, value := range values {
		t, err := template.New("default").Funcs(*funcMapGenerator.GetFuncMap()).Parse(value)
		if err != nil {
			templates[i] = nil
		}
		templates[i] = t
	}

	randomTemplateFunc := utils.RandomValueCycle(templates, localRand)

	return func() (string, string) {
		if tmpl := randomTemplateFunc(); tmpl == nil {
			return "", ""
		} else {
			var buf bytes.Buffer
			_ = tmpl.Execute(&buf, nil)
			return buf.String(), funcMapGenerator.GetBodyDataHeader()
		}
	}
}
