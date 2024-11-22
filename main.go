package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aykhans/dodo/config"
	customerrors "github.com/aykhans/dodo/custom_errors"
	"github.com/aykhans/dodo/readers"
	"github.com/aykhans/dodo/requests"
	"github.com/aykhans/dodo/utils"
	"github.com/aykhans/dodo/validation"
	goValidator "github.com/go-playground/validator/v10"
)

func main() {
	validator := validation.NewValidator()
	conf := config.Config{}
	jsonConf := config.JSONConfig{}

	cliConf, err := readers.CLIConfigReader()
	if err != nil || cliConf == nil {
		os.Exit(0)
	}
	if err := validator.StructPartial(cliConf, "ConfigFile"); err != nil {
		utils.PrintErrAndExit(
			customerrors.ValidationErrorsFormater(
				err.(goValidator.ValidationErrors),
			),
		)
	}
	if cliConf.ConfigFile != "" {
		jsonConfNew, err := readers.JSONConfigReader(cliConf.ConfigFile)
		if err != nil {
			utils.PrintErrAndExit(err)
		}
		if err := validator.StructFiltered(
			jsonConfNew,
			func(ns []byte) bool {
				return strings.LastIndex(string(ns), "Proxies") == -1
			}); err != nil {
			utils.PrintErrAndExit(
				customerrors.ValidationErrorsFormater(
					err.(goValidator.ValidationErrors),
				),
			)
		}
		jsonConf = *jsonConfNew
		conf.MergeConfigs(&jsonConf.Config)
	}

	conf.MergeConfigs(&cliConf.Config)
	conf.SetDefaults()
	if err := validator.Struct(conf); err != nil {
		utils.PrintErrAndExit(
			customerrors.ValidationErrorsFormater(
				err.(goValidator.ValidationErrors),
			),
		)
	}

	parsedURL, err := url.Parse(conf.URL)
	if err != nil {
		utils.PrintErrAndExit(err)
	}
	requestConf := &config.RequestConfig{
		Method:       conf.Method,
		URL:          parsedURL,
		Timeout:      time.Duration(conf.Timeout) * time.Millisecond,
		DodosCount:   conf.DodosCount,
		RequestCount: conf.RequestCount,
		Params:       jsonConf.Params,
		Headers:      jsonConf.Headers,
		Cookies:      jsonConf.Cookies,
		Proxies:      jsonConf.Proxies,
		Body:         jsonConf.Body,
		Yes:          cliConf.Yes,
		NoProxyCheck: *conf.NoProxyCheck.ValueOrPanic(),
	}
	requestConf.Print()
	if !cliConf.Yes {
		response := readers.CLIYesOrNoReader("Do you want to continue?", true)
		if !response {
			utils.PrintAndExit("Exiting...")
		}
		fmt.Println()
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	responses, err := requests.Run(ctx, requestConf)
	if err != nil {
		if customerrors.Is(err, customerrors.ErrInterrupt) {
			utils.PrintlnC(utils.Colors.Yellow, err.Error())
			return
		} else if customerrors.Is(err, customerrors.ErrNoInternet) {
			utils.PrintAndExit("No internet connection")
			return
		}
		panic(err)
	}

	responses.Print()
}
