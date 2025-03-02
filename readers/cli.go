package readers

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aykhans/dodo/config"
	. "github.com/aykhans/dodo/types"
	"github.com/fatih/color"
)

const usageText = `Usage:
  dodo [flags]

Examples:

Simple usage only with URL:
  dodo -u https://example.com

Simple usage with config file:
  dodo -c /path/to/config/file/config.json

Usage with all flags:
  dodo -c /path/to/config/file/config.json -u https://example.com -m POST -d 10 -r 1000 -t 2000 --no-proxy-check -y

Flags:
  -h, --help                 help for dodo
  -v, --version              version for dodo
  -c, --config-file string   Path to the config file
  -d, --dodos uint           Number of dodos(threads) (default %d)
  -m, --method string        HTTP Method (default %s)
  -r, --request uint         Number of total requests (default %d)
  -t, --timeout uint32       Timeout for each request in milliseconds (default %d)
  -u, --url string           URL for stress testing
      --no-proxy-check bool  Do not check for proxies (default false)
  -y, --yes bool             Answer yes to all questions (default false)`

func CLIConfigReader() (*config.CLIConfig, error) {
	flag.Usage = func() {
		fmt.Printf(
			usageText+"\n",
			config.DefaultDodosCount,
			config.DefaultMethod,
			config.DefaultRequestCount,
			config.DefaultTimeout,
		)
	}

	var (
		cliConfig          = config.NewCLIConfig(config.NewConfig("", 0, 0, 0, nil), NewOption(false), "")
		configFile         = ""
		yes                = false
		method             = ""
		url                = ""
		dodosCount    uint = 0
		requestsCount uint = 0
		timeout       uint = 0
		noProxyCheck  bool = false
	)
	{
		flag.Bool("version", false, "Prints the version of the program")
		flag.Bool("v", false, "Prints the version of the program")

		flag.StringVar(&configFile, "config-file", "", "Path to the configuration file")
		flag.StringVar(&configFile, "c", "", "Path to the configuration file")

		flag.BoolVar(&yes, "yes", false, "Answer yes to all questions")
		flag.BoolVar(&yes, "y", false, "Answer yes to all questions")

		flag.StringVar(&method, "method", "", "HTTP Method")
		flag.StringVar(&method, "m", "", "HTTP Method")

		flag.StringVar(&url, "url", "", "URL to send the request")
		flag.StringVar(&url, "u", "", "URL to send the request")

		flag.UintVar(&dodosCount, "dodos", 0, "Number of dodos(threads)")
		flag.UintVar(&dodosCount, "d", 0, "Number of dodos(threads)")

		flag.UintVar(&requestsCount, "requests", 0, "Number of total requests")
		flag.UintVar(&requestsCount, "r", 0, "Number of total requests")

		flag.UintVar(&timeout, "timeout", 0, "Timeout for each request in milliseconds")
		flag.UintVar(&timeout, "t", 0, "Timeout for each request in milliseconds")

		flag.BoolVar(&noProxyCheck, "no-proxy-check", false, "Do not check for active proxies")
	}

	flag.Parse()

	args := flag.Args()
	if len(args) > 0 {
		return nil, fmt.Errorf("unexpected arguments: %v", strings.Join(args, ", "))
	}

	returnNil := false
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "version", "v":
			fmt.Printf("dodo version %s\n", config.VERSION)
			returnNil = true
		case "config-file", "c":
			cliConfig.ConfigFile = configFile
		case "yes", "y":
			cliConfig.Yes.SetValue(yes)
		case "method", "m":
			cliConfig.Method = method
		case "url", "u":
			cliConfig.URL = url
		case "dodos", "d":
			cliConfig.DodosCount = dodosCount
		case "requests", "r":
			cliConfig.RequestCount = requestsCount
		case "timeout", "t":
			var maxUint32 uint = 4294967295
			if timeout > maxUint32 {
				color.Yellow("timeout value is too large, setting to %d", maxUint32)
				timeout = maxUint32
			}
			cliConfig.Timeout = uint32(timeout)
		case "no-proxy-check":
			cliConfig.NoProxyCheck.SetValue(noProxyCheck)
		}
	})

	if returnNil {
		return nil, nil
	}
	return cliConfig, nil
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
