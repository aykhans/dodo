package utils

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/text"
)

func PrintErr(err error) {
	fmt.Fprintln(os.Stderr, text.FgRed.Sprint(err.Error()))
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
