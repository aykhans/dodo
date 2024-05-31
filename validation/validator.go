package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/net/http/httpguts"
)

// net/http/request.go/isNotToken
func isNotToken(r rune) bool {
	return !httpguts.IsTokenRune(r)
}

func NewValidator() *validator.Validate {
	validation := validator.New()
	validation.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if fld.Tag.Get("validation_name") != "" {
			return fld.Tag.Get("validation_name")
		} else {
			return fld.Tag.Get("json")
		}
	})
	validation.RegisterValidation(
		"http_method",
		func(fl validator.FieldLevel) bool {
			method := fl.Field().String()
			// net/http/request.go/validMethod
			return len(method) > 0 && strings.IndexFunc(method, isNotToken) == -1
		},
	)
	validation.RegisterValidation(
		"string_bool",
		func(fl validator.FieldLevel) bool {
			s := fl.Field().String()
			return s == "true" || s == "false" || s == ""
		},
	)
	validation.RegisterValidation(
		"proxy_url",
		func(fl validator.FieldLevel) bool {
			url := fl.Field().String()
			if url == "" {
				return false
			}
			if err := validation.Var(url, "url"); err != nil {
				return false
			}
			if !(url[:7] == "http://" ||
				url[:9] == "socks5://" ||
				url[:10] == "socks5h://") {
				return false
			}
			return true
		},
	)
	return validation
}
