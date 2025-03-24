package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aykhans/dodo/config"
	"github.com/aykhans/dodo/requests"
	"github.com/aykhans/dodo/types"
	"github.com/aykhans/dodo/utils"
	"github.com/jedib0t/go-pretty/v6/text"
)

func main() {
	conf := config.NewConfig()
	configFile, err := conf.ReadCLI()
	if err != nil {
		utils.PrintErrAndExit(err)
	}

	if configFile.String() != "" {
		tempConf := config.NewConfig()
		if err := tempConf.ReadFile(configFile); err != nil {
			utils.PrintErrAndExit(err)
		}
		tempConf.MergeConfig(conf)
		conf = tempConf
	}
	conf.SetDefaults()

	if errs := conf.Validate(); len(errs) > 0 {
		utils.PrintErrAndExit(errors.Join(errs...))
	}

	requestConf := config.NewRequestConfig(conf)
	requestConf.Print()

	if !requestConf.Yes {
		response := config.CLIYesOrNoReader("Do you want to continue?", false)
		if !response {
			utils.PrintAndExit("Exiting...\n")
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	go listenForTermination(func() { cancel() })

	if requestConf.Duration > 0 {
		time.AfterFunc(requestConf.Duration, func() { cancel() })
	}

	responses, err := requests.Run(ctx, requestConf)
	if err != nil {
		if err == types.ErrInterrupt {
			fmt.Println(text.FgYellow.Sprint(err.Error()))
			return
		}
		utils.PrintErrAndExit(err)
	}

	responses.Print()
}

func listenForTermination(do func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	do()
}
