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
		dodosCount   uint
		requestCount uint
		timeout      uint32
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
	rootCmd.Flags().BoolVarP(&cliConfig.Yes, "yes", "y", false, "Answer yes to all questions")
	rootCmd.Flags().StringVarP(&cliConfig.Method, "method", "m", "", fmt.Sprintf("HTTP Method (default %s)", config.DefaultMethod))
	rootCmd.Flags().StringVarP(&cliConfig.URL, "url", "u", "", "URL for stress testing")
	rootCmd.Flags().UintVarP(&dodosCount, "dodos-count", "d", config.DefaultDodosCount, "Number of dodos(threads)")
	rootCmd.Flags().UintVarP(&requestCount, "request-count", "r", config.DefaultRequestCount, "Number of total requests")
	rootCmd.Flags().Uint32VarP(&timeout, "timeout", "t", config.DefaultTimeout, "Timeout for each request in milliseconds")
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
