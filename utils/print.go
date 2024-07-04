package utils

import (
	"fmt"
	"os"
)

var Colors = struct {
	reset   string
	Red     string
	Green   string
	Yellow  string
	Orange  string
	Blue    string
	Magenta string
	Cyan    string
	Gray    string
	White   string
}{
	reset:   "\033[0m",
	Red:     "\033[31m",
	Green:   "\033[32m",
	Yellow:  "\033[33m",
	Orange:  "\033[38;5;208m",
	Blue:    "\033[34m",
	Magenta: "\033[35m",
	Cyan:    "\033[36m",
	Gray:    "\033[37m",
	White:   "\033[97m",
}

func Colored(color string, a ...any) string {
	return color + fmt.Sprint(a...) + Colors.reset
}

func PrintfC(color string, format string, a ...any) {
	fmt.Printf(Colored(color, format), a...)
}

func PrintlnC(color string, a ...any) {
	fmt.Println(Colored(color, a...))
}

func PrintErr(err error) {
	PrintlnC(Colors.Red, err.Error())
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
