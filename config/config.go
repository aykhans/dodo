package config

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	VERSION                 = "0.3.0"
	DefaultUserAgent        = "Dodo/" + VERSION
	ProxyCheckURL           = "https://www.google.com"
	DefaultMethod           = "GET"
	DefaultTimeout          = 10000 // Milliseconds (10 seconds)
	DefaultDodosCount       = 1
	DefaultRequestCount     = 1000
	MaxDodosCountForProxies = 20 // Max dodos count for proxy check
)

type IConfig interface {
	MergeConfigs(newConfig IConfig) IConfig
}

type RequestConfig struct {
	Method       string
	URL          *url.URL
	Timeout      time.Duration
	DodosCount   int
	RequestCount int
	Params       map[string]string
	Headers      map[string]string
	Cookies      map[string]string
	Proxies      []Proxy
	Body         string
}

func (config *RequestConfig) Print() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendRow(table.Row{
		"Method", "URL", "Timeout", "Dodos",
		"Request Count", "Params Count",
		"Headers Count", "Cookies Count",
		"Proxies Count", "Body"})
	t.AppendSeparator()
	t.AppendRow(table.Row{
		config.Method, config.URL,
		fmt.Sprintf("%dms", config.Timeout/time.Millisecond),
		config.DodosCount, config.RequestCount,
		len(config.Params), len(config.Headers),
		len(config.Cookies), len(config.Proxies), config.Body})
	t.Render()
}

func (config *RequestConfig) GetValidDodosCountForRequests() int {
	return min(config.DodosCount, config.RequestCount)
}

func (config *RequestConfig) GetValidDodosCountForProxies() int {
	return min(config.DodosCount, len(config.Proxies), MaxDodosCountForProxies)
}

type Config struct {
	Method       string `json:"method" validate:"http_method"` // custom validations: http_method
	URL          string `json:"url" validate:"http_url,required"`
	Timeout      int    `json:"timeout" validate:"gte=1,lte=100000"`
	DodosCount   int    `json:"dodos_count" validate:"gte=1"`
	RequestCount int    `json:"request_count" validation_name:"request-count" validate:"gte=1"`
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
}

type Proxy struct {
	URL      string `json:"url" validate:"required,proxy_url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type JSONConfig struct {
	Config
	Params  map[string]string `json:"params"`
	Headers map[string]string `json:"headers"`
	Cookies map[string]string `json:"cookies"`
	Proxies []Proxy           `json:"proxies" validate:"dive"`
	Body    string            `json:"body"`
}

func (config *JSONConfig) MergeConfigs(newConfig *JSONConfig) {
	config.Config.MergeConfigs(&newConfig.Config)
	if len(newConfig.Params) != 0 {
		config.Params = newConfig.Params
	}
	if len(newConfig.Headers) != 0 {
		config.Headers = newConfig.Headers
	}
	if len(newConfig.Cookies) != 0 {
		config.Cookies = newConfig.Cookies
	}
	if newConfig.Body != "" {
		config.Body = newConfig.Body
	}
	if len(newConfig.Proxies) != 0 {
		config.Proxies = newConfig.Proxies
	}
}

type CLIConfig struct {
	Config
	ConfigFile string `validation_name:"config-file" validate:"omitempty,filepath"`
}

func (config *CLIConfig) MergeConfigs(newConfig *CLIConfig) {
	config.Config.MergeConfigs(&newConfig.Config)
	if newConfig.ConfigFile != "" {
		config.ConfigFile = newConfig.ConfigFile
	}
}
