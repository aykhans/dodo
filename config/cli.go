package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aykhans/dodo/types"
	"github.com/aykhans/dodo/utils"
)

const cliUsageText = `Usage:
  dodo [flags]

Examples:

Simple usage only with URL:
  dodo -u https://example.com

Usage with config file:
  dodo -f /path/to/config/file/config.json

Usage with all flags:
  dodo -f /path/to/config/file/config.json \
    -u https://example.com -m POST \
    -d 10 -r 1000 -t 3s \
    -b "body1" -body "body2" \
    -H "header1:value1" -header "header2:value2" \
    -p "param1=value1" -param "param2=value2" \
    -c "cookie1=value1" -cookie "cookie2=value2" \
    -x "http://proxy.example.com:8080" -proxy "socks5://proxy2.example.com:8080" \
    -y

Flags:
  -h, -help                   help for dodo
  -v, -version                version for dodo
  -y, -yes          bool      Answer yes to all questions (default %v)
  -f, -config-file  string    Path to the local config file or http(s) URL of the config file
  -d, -dodos        uint      Number of dodos(threads) (default %d)
  -r, -requests     uint      Number of total requests (default %d)
  -t, -timeout      Duration  Timeout for each request (e.g. 400ms, 15s, 1m10s) (default %v)
  -u, -url          string    URL for stress testing
  -m, -method       string    HTTP Method for the request (default %s)
  -b, -body         [string]  Body for the request (e.g. "body text")
  -p, -param        [string]  Parameter for the request (e.g. "key1=value1")
  -H, -header       [string]  Header for the request (e.g. "key1: value1")
  -c, -cookie       [string]  Cookie for the request (e.g. "key1=value1")
  -x, -proxy        [string]  Proxy for the request (e.g. "http://proxy.example.com:8080")`

func (config *Config) ReadCLI() (types.ConfigFile, error) {
	flag.Usage = func() {
		fmt.Printf(
			cliUsageText+"\n",
			DefaultYes,
			DefaultDodosCount,
			DefaultRequestCount,
			DefaultTimeout,
			DefaultMethod,
		)
	}

	var (
		version      = false
		configFile   = ""
		yes          = false
		method       = ""
		url          types.RequestURL
		dodosCount   = uint(0)
		requestCount = uint(0)
		timeout      time.Duration
	)

	{
		flag.BoolVar(&version, "version", false, "Prints the version of the program")
		flag.BoolVar(&version, "v", false, "Prints the version of the program")

		flag.StringVar(&configFile, "config-file", "", "Path to the configuration file")
		flag.StringVar(&configFile, "f", "", "Path to the configuration file")

		flag.BoolVar(&yes, "yes", false, "Answer yes to all questions")
		flag.BoolVar(&yes, "y", false, "Answer yes to all questions")

		flag.StringVar(&method, "method", "", "HTTP Method")
		flag.StringVar(&method, "m", "", "HTTP Method")

		flag.Var(&url, "url", "URL to send the request")
		flag.Var(&url, "u", "URL to send the request")

		flag.UintVar(&dodosCount, "dodos", 0, "Number of dodos(threads)")
		flag.UintVar(&dodosCount, "d", 0, "Number of dodos(threads)")

		flag.UintVar(&requestCount, "requests", 0, "Number of total requests")
		flag.UintVar(&requestCount, "r", 0, "Number of total requests")

		flag.DurationVar(&timeout, "timeout", 0, "Timeout for each request (e.g. 400ms, 15s, 1m10s)")
		flag.DurationVar(&timeout, "t", 0, "Timeout for each request (e.g. 400ms, 15s, 1m10s)")

		flag.Var(&config.Params, "param", "URL parameter to send with the request")
		flag.Var(&config.Params, "p", "URL parameter to send with the request")

		flag.Var(&config.Headers, "header", "Header to send with the request")
		flag.Var(&config.Headers, "H", "Header to send with the request")

		flag.Var(&config.Cookies, "cookie", "Cookie to send with the request")
		flag.Var(&config.Cookies, "c", "Cookie to send with the request")

		flag.Var(&config.Body, "body", "Body to send with the request")
		flag.Var(&config.Body, "b", "Body to send with the request")

		flag.Var(&config.Proxies, "proxy", "Proxy to use for the request")
		flag.Var(&config.Proxies, "x", "Proxy to use for the request")
	}

	flag.Parse()

	if len(os.Args) <= 1 {
		flag.CommandLine.Usage()
		os.Exit(0)
	}

	if args := flag.Args(); len(args) > 0 {
		return types.ConfigFile(configFile), fmt.Errorf("unexpected arguments: %v", strings.Join(args, ", "))
	}

	if version {
		fmt.Printf("dodo version %s\n", VERSION)
		os.Exit(0)
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "method", "m":
			config.Method = utils.ToPtr(method)
		case "url", "u":
			config.URL = utils.ToPtr(url)
		case "dodos", "d":
			config.DodosCount = utils.ToPtr(dodosCount)
		case "requests", "r":
			config.RequestCount = utils.ToPtr(requestCount)
		case "timeout", "t":
			config.Timeout = &types.Timeout{Duration: timeout}
		case "yes", "y":
			config.Yes = utils.ToPtr(yes)
		}
	})

	return types.ConfigFile(configFile), nil
}

// CLIYesOrNoReader reads a yes or no answer from the command line.
// It prompts the user with the given message and default value,
// and returns true if the user answers "y" or "Y", and false otherwise.
// If there is an error while reading the input, it returns false.
// If the user simply presses enter without providing any input,
// it returns the default value specified by the `dft` parameter.
func CLIYesOrNoReader(message string, dft bool) bool {
	var answer string
	defaultMessage := "Y/n"
	if !dft {
		defaultMessage = "y/N"
	}
	fmt.Printf("%s [%s]: ", message, defaultMessage)
	if _, err := fmt.Scanln(&answer); err != nil {
		if err.Error() == "unexpected newline" {
			return dft
		}
		return false
	}
	if answer == "" {
		return dft
	}
	return answer == "y" || answer == "Y"
}
