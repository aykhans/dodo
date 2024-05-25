package validation

import (
	"reflect"
	"strings"

	"github.com/aykhans/dodo/config"
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
		"url_map_slice",
		func(fl validator.FieldLevel) bool {
			proxies := fl.Field().Interface().(config.ProxySlice)
			for _, proxy := range proxies {
				if _, ok := proxy["url"]; !ok {
					return false
				}
				if err := validation.Var(proxy["url"], "url"); err != nil {
					return false
				}
			}
			return true
		},
	)
	return validation
}
