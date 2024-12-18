package customerrors

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

func OSErrorFormater(err error) error {
	errStr := err.Error()
	if strings.Contains(errStr, "no such file or directory") {
		fileName1 := strings.Index(errStr, "open")
		fileName2 := strings.LastIndex(errStr, ":")
		return NewFileNotFoundError(errStr[fileName1+5:fileName2], err)
	}
	return ErrInvalidFile
}

func shortenNamespace(namespace string) string {
	return namespace[strings.Index(namespace, ".")+1:]
}

func ValidationErrorsFormater(errs validator.ValidationErrors) error {
	errsStr := make(map[string]string)
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			errsStr[shortenNamespace(err.Namespace())] = fmt.Sprintf("Field \"%s\" is required", err.Field())
		case "gte":
			errsStr[shortenNamespace(err.Namespace())] = fmt.Sprintf("Value of \"%s\" must be greater than or equal to \"%s\"", err.Field(), err.Param())
		case "lte":
			errsStr[shortenNamespace(err.Namespace())] = fmt.Sprintf("Value of \"%s\" must be less than or equal to \"%s\"", err.Field(), err.Param())
		case "filepath":
			errsStr[shortenNamespace(err.Namespace())] = fmt.Sprintf("Invalid file path for \"%s\" field: \"%s\"", err.Field(), err.Value())
		case "http_url":
			errsStr[shortenNamespace(err.Namespace())] =
				fmt.Sprintf("Invalid url for \"%s\" field: \"%s\"", err.Field(), err.Value())
		// --------------------------------------| Custom validations |--------------------------------------
		case "http_method":
			errsStr[shortenNamespace(err.Namespace())] = fmt.Sprintf("Invalid HTTP method for \"%s\" field: \"%s\"", err.Field(), err.Value())
		case "proxy_url":
			errsStr[shortenNamespace(err.Namespace())] = fmt.Sprintf("Invalid proxy url for \"%s\" field: \"%s\" (it must be http, socks5 or socks5h)", err.Field(), err.Value())
		case "string_bool":
			errsStr[shortenNamespace(err.Namespace())] = fmt.Sprintf("Invalid value for \"%s\" field: \"%s\"", err.Field(), err.Value())
		default:
			errsStr[shortenNamespace(err.Namespace())] = fmt.Sprintf("Invalid value for \"%s\" field: \"%s\"", err.Field(), err.Value())
		}
	}
	return NewValidationErrors(errsStr, errs)
}

func RequestErrorsFormater(err error) string {
	switch e := err.(type) {
	case *url.Error:
		if netErr, ok := e.Err.(net.Error); ok && netErr.Timeout() {
			return "Timeout Error"
		}
		if strings.Contains(e.Error(), "http: ContentLength=") {
			println(e.Error())
			return "Empty Body Error"
		}
		// TODO: Add more cases
	}
	return "Unknown Error"
}
