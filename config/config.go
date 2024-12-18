package config

import (
	"fmt"
	"net/url"
	"os"
	"time"

	. "github.com/aykhans/dodo/types"
	"github.com/aykhans/dodo/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	VERSION                 string = "0.5.4"
	DefaultUserAgent        string = "Dodo/" + VERSION
	ProxyCheckURL           string = "https://www.google.com"
	DefaultMethod           string = "GET"
	DefaultTimeout          uint32 = 10000 // Milliseconds (10 seconds)
	DefaultDodosCount       uint   = 1
	DefaultRequestCount     uint   = 1
	MaxDodosCountForProxies uint   = 20 // Max dodos count for proxy check
)

type RequestConfig struct {
	Method       string
	URL          *url.URL
	Timeout      time.Duration
	DodosCount   uint
	RequestCount uint
	Params       map[string][]string
	Headers      map[string][]string
	Cookies      map[string][]string
	Proxies      []Proxy
	Body         []string
	Yes          bool
	NoProxyCheck bool
}

func (config *RequestConfig) Print() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, WidthMax: 50},
	})

	newHeaders := make(map[string][]string)
	newHeaders["User-Agent"] = []string{DefaultUserAgent}
	for k, v := range config.Headers {
		newHeaders[k] = v
	}

	t.AppendHeader(table.Row{"Request Configuration"})
	t.AppendRow(table.Row{"Method", config.Method})
	t.AppendSeparator()
	t.AppendRow(table.Row{"URL", config.URL})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Timeout", fmt.Sprintf("%dms", config.Timeout/time.Millisecond)})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Dodos", config.DodosCount})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Requests", config.RequestCount})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Params", utils.MarshalJSON(config.Params, 3)})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Headers", utils.MarshalJSON(newHeaders, 3)})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Cookies", utils.MarshalJSON(config.Cookies, 3)})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Proxies Count", len(config.Proxies)})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Proxy Check", !config.NoProxyCheck})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Body", utils.MarshalJSON(config.Body, 3)})

	t.Render()
}

func (config *RequestConfig) GetValidDodosCountForRequests() uint {
	return min(config.DodosCount, config.RequestCount)
}

func (config *RequestConfig) GetValidDodosCountForProxies() uint {
	return min(config.DodosCount, uint(len(config.Proxies)), MaxDodosCountForProxies)
}

func (config *RequestConfig) GetMaxConns(minConns uint) uint {
	maxConns := max(
		minConns, uint(config.GetValidDodosCountForRequests()),
	)
	return ((maxConns * 50 / 100) + maxConns)
}

type Config struct {
	Method       string       `json:"method" validate:"http_method"` // custom validations: http_method
	URL          string       `json:"url" validate:"http_url,required"`
	Timeout      uint32       `json:"timeout" validate:"gte=1,lte=100000"`
	DodosCount   uint         `json:"dodos_count" validate:"gte=1"`
	RequestCount uint         `json:"request_count" validation_name:"request-count" validate:"gte=1"`
	NoProxyCheck Option[bool] `json:"no_proxy_check"`
}

func NewConfig(
	method string,
	timeout uint32,
	dodosCount uint,
	requestCount uint,
	noProxyCheck Option[bool],
) *Config {
	if noProxyCheck == nil {
		noProxyCheck = NewNoneOption[bool]()
	}

	return &Config{
		Method:       method,
		Timeout:      timeout,
		DodosCount:   dodosCount,
		RequestCount: requestCount,
		NoProxyCheck: noProxyCheck,
	}
}

func (config *Config) MergeConfigs(newConfig *Config) {
	if newConfig.Method != "" {
		config.Method = newConfig.Method
	}
	if newConfig.URL != "" {
		config.URL = newConfig.URL
	}
	if newConfig.Timeout != 0 {
		config.Timeout = newConfig.Timeout
	}
	if newConfig.DodosCount != 0 {
		config.DodosCount = newConfig.DodosCount
	}
	if newConfig.RequestCount != 0 {
		config.RequestCount = newConfig.RequestCount
	}
	if !newConfig.NoProxyCheck.IsNone() {
		config.NoProxyCheck = newConfig.NoProxyCheck
	}
}

func (config *Config) SetDefaults() {
	if config.Method == "" {
		config.Method = DefaultMethod
	}
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}
	if config.DodosCount == 0 {
		config.DodosCount = DefaultDodosCount
	}
	if config.RequestCount == 0 {
		config.RequestCount = DefaultRequestCount
	}
	if config.NoProxyCheck.IsNone() {
		config.NoProxyCheck = NewOption(false)
	}
}

type Proxy struct {
	URL      string `json:"url" validate:"required,proxy_url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type JSONConfig struct {
	*Config
	Params  map[string][]string `json:"params"`
	Headers map[string][]string `json:"headers"`
	Cookies map[string][]string `json:"cookies"`
	Proxies []Proxy             `json:"proxies" validate:"dive"`
	Body    []string            `json:"body"`
}

func NewJSONConfig(
	config *Config,
	params map[string][]string,
	headers map[string][]string,
	cookies map[string][]string,
	proxies []Proxy,
	body []string,
) *JSONConfig {
	return &JSONConfig{
		config, params, headers, cookies, proxies, body,
	}
}

func (config *JSONConfig) MergeConfigs(newConfig *JSONConfig) {
	config.Config.MergeConfigs(newConfig.Config)
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

type CLIConfig struct {
	*Config
	Yes        Option[bool]   `json:"yes" validate:"omitempty"`
	ConfigFile string `validation_name:"config-file" validate:"omitempty,filepath"`
}

func NewCLIConfig(
	config *Config,
	yes Option[bool],
	configFile string,
) *CLIConfig {
	return &CLIConfig{
		config, yes, configFile,
	}
}

func (config *CLIConfig) MergeConfigs(newConfig *CLIConfig) {
	config.Config.MergeConfigs(newConfig.Config)
	if newConfig.ConfigFile != "" {
		config.ConfigFile = newConfig.ConfigFile
	}
	if !newConfig.Yes.IsNone() {
		config.Yes = newConfig.Yes
	}
}
