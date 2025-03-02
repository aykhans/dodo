package utils

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func PrintErr(err error) {
	color.New(color.FgRed).Fprintln(os.Stderr, err.Error())
}

func PrintErrAndExit(err error) {
	if err != nil {
		PrintErr(err)
		os.Exit(1)
	}
}

func PrintAndExit(message string) {
	fmt.Println(message)
	os.Exit(0)
}
