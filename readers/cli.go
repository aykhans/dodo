package readers

import (
	"fmt"

	"github.com/aykhans/dodo/config"
	"github.com/aykhans/dodo/custom_errors"
	"github.com/aykhans/dodo/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func CLIConfigReader() (*config.CLIConfig, error) {
	var (
		returnNil    = false
		cliConfig    = &config.CLIConfig{}
		dodosCount   int
		requestCount int
		timeout      int
		rootCmd      = &cobra.Command{
			Use: "dodo [flags]",
			Example: `  Simple usage only with URL:
    dodo -u https://example.com

  Simple usage with config file:
    dodo -c /path/to/config/file/config.json

  Usage with all flags:
    dodo -c /path/to/config/file/config.json -u https://example.com -m POST -d 10 -r 1000 -t 2000`,
			Short: `
██████████      ███████    ██████████      ███████   
░░███░░░░███   ███░░░░░███ ░░███░░░░███   ███░░░░░███ 
 ░███   ░░███ ███     ░░███ ░███   ░░███ ███     ░░███
 ░███    ░███░███      ░███ ░███    ░███░███      ░███
 ░███    ░███░███      ░███ ░███    ░███░███      ░███
 ░███    ███ ░░███     ███  ░███    ███ ░░███     ███ 
 ██████████   ░░░███████░   ██████████   ░░░███████░  
░░░░░░░░░░      ░░░░░░░    ░░░░░░░░░░      ░░░░░░░    
`,
			Run:           func(cmd *cobra.Command, args []string) {},
			SilenceErrors: true,
			SilenceUsage:  true,
			Version:       config.VERSION,
		}
	)

	rootCmd.Flags().StringVarP(&cliConfig.ConfigFile, "config-file", "c", "", "Path to the config file")
	rootCmd.Flags().StringVarP(&cliConfig.Method, "method", "m", "", fmt.Sprintf("HTTP Method (default %s)", config.DefaultMethod))
	rootCmd.Flags().StringVarP(&cliConfig.URL, "url", "u", "", "URL for stress testing")
	rootCmd.Flags().IntVarP(&dodosCount, "dodos-count", "d", config.DefaultDodosCount, "Number of dodos(threads)")
	rootCmd.Flags().IntVarP(&requestCount, "request-count", "r", config.DefaultRequestCount, "Number of total requests")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", config.DefaultTimeout, "Timeout for each request in milliseconds")
	if err := rootCmd.Execute(); err != nil {
		utils.PrintErr(err)
		rootCmd.Println(rootCmd.UsageString())
		return nil, customerrors.CobraErrorFormater(err)
	}
	rootCmd.Flags().Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "help":
			returnNil = true
		case "version":
			returnNil = true
		case "dodos-count":
			cliConfig.DodosCount = dodosCount
		case "request-count":
			cliConfig.RequestCount = requestCount
		case "timeout":
			cliConfig.Timeout = timeout
		}
	})
	if returnNil {
		return nil, nil
	}
	return cliConfig, nil
}

func CLIYesOrNoReader(message string) bool {
	var answer string
	fmt.Printf("%s [y/N]: ", message)
	if _, err := fmt.Scanln(&answer); err != nil {
		return false
	}
	return answer == "y" || answer == "Y"
}
