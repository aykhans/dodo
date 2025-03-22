package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/aykhans/dodo/types"
	"github.com/aykhans/dodo/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	VERSION             string        = "0.6.1"
	DefaultUserAgent    string        = "Dodo/" + VERSION
	DefaultMethod       string        = "GET"
	DefaultTimeout      time.Duration = time.Second * 10
	DefaultDodosCount   uint          = 1
	DefaultRequestCount uint          = 1
	DefaultYes          bool          = false
)

var SupportedProxySchemes []string = []string{"http", "socks5", "socks5h"}

type RequestConfig struct {
	Method       string
	URL          url.URL
	Timeout      time.Duration
	DodosCount   uint
	RequestCount uint
	Yes          bool
	Params       types.Params
	Headers      types.Headers
	Cookies      types.Cookies
	Body         types.Body
	Proxies      types.Proxies
}

func NewRequestConfig(conf *Config) *RequestConfig {
	return &RequestConfig{
		Method:       *conf.Method,
		URL:          conf.URL.URL,
		Timeout:      conf.Timeout.Duration,
		DodosCount:   *conf.DodosCount,
		RequestCount: *conf.RequestCount,
		Yes:          *conf.Yes,
		Params:       conf.Params,
		Headers:      conf.Headers,
		Cookies:      conf.Cookies,
		Body:         conf.Body,
		Proxies:      conf.Proxies,
	}
}

func (rc *RequestConfig) GetValidDodosCountForRequests() uint {
	return min(rc.DodosCount, rc.RequestCount)
}

func (rc *RequestConfig) GetMaxConns(minConns uint) uint {
	maxConns := max(
		minConns, rc.GetValidDodosCountForRequests(),
	)
	return ((maxConns * 50 / 100) + maxConns)
}

func (rc *RequestConfig) Print() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Number: 2,
			WidthMaxEnforcer: func(col string, maxLen int) string {
				lines := strings.Split(col, "\n")
				for i, line := range lines {
					if len(line) > maxLen {
						lines[i] = line[:maxLen-3] + "..."
					}
				}
				return strings.Join(lines, "\n")
			},
			WidthMax: 50},
	})

	t.AppendHeader(table.Row{"Request Configuration"})
	t.AppendRow(table.Row{"URL", rc.URL.String()})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Method", rc.Method})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Timeout", rc.Timeout})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Dodos", rc.DodosCount})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Requests", rc.RequestCount})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Params", rc.Params.String()})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Headers", rc.Headers.String()})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Cookies", rc.Cookies.String()})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Proxy", rc.Proxies.String()})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Body", rc.Body.String()})

	t.Render()
}

type Config struct {
	Method       *string           `json:"method" yaml:"method"`
	URL          *types.RequestURL `json:"url" yaml:"url"`
	Timeout      *types.Timeout    `json:"timeout" yaml:"timeout"`
	DodosCount   *uint             `json:"dodos" yaml:"dodos"`
	RequestCount *uint             `json:"requests" yaml:"requests"`
	Yes          *bool             `json:"yes" yaml:"yes"`
	Params       types.Params      `json:"params" yaml:"params"`
	Headers      types.Headers     `json:"headers" yaml:"headers"`
	Cookies      types.Cookies     `json:"cookies" yaml:"cookies"`
	Body         types.Body        `json:"body" yaml:"body"`
	Proxies      types.Proxies     `json:"proxy" yaml:"proxy"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Validate() []error {
	var errs []error
	if utils.IsNilOrZero(c.URL) {
		errs = append(errs, errors.New("request URL is required"))
	} else {
		if c.URL.Scheme == "" {
			c.URL.Scheme = "http"
		}
		if c.URL.Scheme != "http" && c.URL.Scheme != "https" {
			errs = append(errs, errors.New("request URL scheme must be http or https"))
		}

		urlParams := types.Params{}
		for key, values := range c.URL.Query() {
			for _, value := range values {
				urlParams = append(urlParams, types.KeyValue[string, []string]{
					Key:   key,
					Value: []string{value},
				})
			}
		}
		c.Params = append(urlParams, c.Params...)
		c.URL.RawQuery = ""
	}

	if utils.IsNilOrZero(c.Method) {
		errs = append(errs, errors.New("request method is required"))
	}
	if utils.IsNilOrZero(c.Timeout) {
		errs = append(errs, errors.New("request timeout must be greater than 0"))
	}
	if utils.IsNilOrZero(c.DodosCount) {
		errs = append(errs, errors.New("dodos count must be greater than 0"))
	}
	if utils.IsNilOrZero(c.RequestCount) {
		errs = append(errs, errors.New("request count must be greater than 0"))
	}

	for i, proxy := range c.Proxies {
		if proxy.String() == "" {
			errs = append(errs, fmt.Errorf("proxies[%d]: proxy cannot be empty", i))
		} else if schema := proxy.Scheme; !slices.Contains(SupportedProxySchemes, schema) {
			errs = append(errs,
				fmt.Errorf("proxies[%d]: proxy has unsupported scheme \"%s\" (supported schemes: %s)",
					i, proxy.String(), strings.Join(SupportedProxySchemes, ", "),
				),
			)
		}
	}

	return errs
}

func (config *Config) MergeConfig(newConfig *Config) {
	if newConfig.Method != nil {
		config.Method = newConfig.Method
	}
	if newConfig.URL != nil {
		config.URL = newConfig.URL
	}
	if newConfig.Timeout != nil {
		config.Timeout = newConfig.Timeout
	}
	if newConfig.DodosCount != nil {
		config.DodosCount = newConfig.DodosCount
	}
	if newConfig.RequestCount != nil {
		config.RequestCount = newConfig.RequestCount
	}
	if newConfig.Yes != nil {
		config.Yes = newConfig.Yes
	}
	if len(newConfig.Params) != 0 {
		config.Params = newConfig.Params
	}
	if len(newConfig.Headers) != 0 {
		config.Headers = newConfig.Headers
	}
	if len(newConfig.Cookies) != 0 {
		config.Cookies = newConfig.Cookies
	}
	if len(newConfig.Body) != 0 {
		config.Body = newConfig.Body
	}
	if len(newConfig.Proxies) != 0 {
		config.Proxies = newConfig.Proxies
	}
}

func (config *Config) SetDefaults() {
	if config.Method == nil {
		config.Method = utils.ToPtr(DefaultMethod)
	}
	if config.Timeout == nil {
		config.Timeout = &types.Timeout{Duration: DefaultTimeout}
	}
	if config.DodosCount == nil {
		config.DodosCount = utils.ToPtr(DefaultDodosCount)
	}
	if config.RequestCount == nil {
		config.RequestCount = utils.ToPtr(DefaultRequestCount)
	}
	if config.Yes == nil {
		config.Yes = utils.ToPtr(DefaultYes)
	}
	config.Headers.SetIfNotExists("User-Agent", DefaultUserAgent)
}
