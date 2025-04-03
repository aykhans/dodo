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
	VERSION             string        = "0.6.2"
	DefaultUserAgent    string        = "Dodo/" + VERSION
	DefaultMethod       string        = "GET"
	DefaultTimeout      time.Duration = time.Second * 10
	DefaultDodosCount   uint          = 1
	DefaultRequestCount uint          = 0
	DefaultDuration     time.Duration = 0
	DefaultYes          bool          = false
)

var SupportedProxySchemes []string = []string{"http", "socks5", "socks5h"}

type RequestConfig struct {
	Method       string
	URL          url.URL
	Timeout      time.Duration
	DodosCount   uint
	RequestCount uint
	Duration     time.Duration
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
		Duration:     conf.Duration.Duration,
		Yes:          *conf.Yes,
		Params:       conf.Params,
		Headers:      conf.Headers,
		Cookies:      conf.Cookies,
		Body:         conf.Body,
		Proxies:      conf.Proxies,
	}
}

func (rc *RequestConfig) GetValidDodosCountForRequests() uint {
	if rc.RequestCount == 0 {
		return rc.DodosCount
	}
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
	if rc.RequestCount > 0 {
		t.AppendRow(table.Row{"Requests", rc.RequestCount})
	} else {
		t.AppendRow(table.Row{"Requests"})
	}
	t.AppendSeparator()
	if rc.Duration > 0 {
		t.AppendRow(table.Row{"Duration", rc.Duration})
	} else {
		t.AppendRow(table.Row{"Duration"})
	}
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
	Duration     *types.Duration   `json:"duration" yaml:"duration"`
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

func (config *Config) Validate() []error {
	var errs []error
	if utils.IsNilOrZero(config.URL) {
		errs = append(errs, errors.New("request URL is required"))
	} else {
		if config.URL.Scheme == "" {
			config.URL.Scheme = "http"
		}
		if config.URL.Scheme != "http" && config.URL.Scheme != "https" {
			errs = append(errs, errors.New("request URL scheme must be http or https"))
		}

		urlParams := types.Params{}
		for key, values := range config.URL.Query() {
			for _, value := range values {
				urlParams = append(urlParams, types.KeyValue[string, []string]{
					Key:   key,
					Value: []string{value},
				})
			}
		}
		config.Params = append(urlParams, config.Params...)
		config.URL.RawQuery = ""
	}

	if utils.IsNilOrZero(config.Method) {
		errs = append(errs, errors.New("request method is required"))
	}
	if utils.IsNilOrZero(config.Timeout) {
		errs = append(errs, errors.New("request timeout must be greater than 0"))
	}
	if utils.IsNilOrZero(config.DodosCount) {
		errs = append(errs, errors.New("dodos count must be greater than 0"))
	}
	if utils.IsNilOrZero(config.Duration) && utils.IsNilOrZero(config.RequestCount) {
		errs = append(errs, errors.New("you should provide at least one of duration or request count"))
	}

	for i, proxy := range config.Proxies {
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
	if newConfig.Duration != nil {
		config.Duration = newConfig.Duration
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
	if config.Duration == nil {
		config.Duration = &types.Duration{Duration: DefaultDuration}
	}
	if config.Yes == nil {
		config.Yes = utils.ToPtr(DefaultYes)
	}
	config.Headers.SetIfNotExists("User-Agent", DefaultUserAgent)
}
