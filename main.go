package main

import (
	"os"
	"time"

	"github.com/aykhans/dodo/config"
	"github.com/aykhans/dodo/custom_errors"
	"github.com/aykhans/dodo/readers"
	"github.com/aykhans/dodo/requests"
	"github.com/aykhans/dodo/utils"
	"github.com/aykhans/dodo/validation"
	goValidator "github.com/go-playground/validator/v10"
	"github.com/jedib0t/go-pretty/v6/table"
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
		if err := validator.StructPartial(jsonConfNew, "Proxies"); err != nil {
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

	dodoConf := &config.DodoConfig{
		Method:       conf.Method,
		URL:          conf.URL,
		Timeout:      time.Duration(conf.Timeout) * time.Millisecond,
		DodosCount:   conf.DodosCount,
		RequestCount: conf.RequestCount,
		Params:       jsonConf.Params,
		Headers:      jsonConf.Headers,
		Cookies:      jsonConf.Cookies,
		Proxies:      jsonConf.Proxies,
		Body:         jsonConf.Body,
	}

	dodoConf.Print()
	responses, err := requests.Run(dodoConf)
	if err != nil {
		utils.PrintErrAndExit(err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{
		"Response",
		"Count",
		"Min Time",
		"Max Time",
		"Average Time",
	})
	for _, mergedResponse := range responses.MergeDodoResponses() {
		t.AppendRow(table.Row{
			mergedResponse.Response,
			mergedResponse.Count,
			mergedResponse.MinTime,
			mergedResponse.MaxTime,
			mergedResponse.AvgTime,
		})
		t.AppendSeparator()
	}
	t.AppendFooter(table.Row{
		"Total",
		responses.Len(),
		responses.MinTime(),
		responses.MaxTime(),
		responses.AvgTime(),
	})
	t.Render()
}
