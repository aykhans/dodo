package requests

import (
	"net/url"

	"github.com/aykhans/dodo/config"
	"github.com/valyala/fasthttp"
)

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
