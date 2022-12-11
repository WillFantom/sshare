package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	errorHeading   = color.New(color.Bold, color.FgWhite, color.BgRed).SprintFunc()
	warnHeading    = color.New(color.Bold, color.FgWhite, color.BgYellow).SprintFunc()
	successHeading = color.New(color.Bold, color.FgHiWhite, color.BgHiGreen).SprintFunc()
	infoHeading    = color.New(color.Bold, color.FgHiWhite, color.BgHiCyan).SprintFunc()
)

func Errorln(message string, fatal bool) {
	fmt.Printf("%s %s\n", errorHeading(" ERROR "), message)
	if fatal {
		os.Exit(1)
	}
}

func Warnln(message string) {
	fmt.Printf("%s %s\n", warnHeading(" WARN "), message)
}

func Successln(message string) {
	fmt.Printf("%s %s\n", successHeading(" SUCCESS "), message)
}

func Infoln(message string) {
	fmt.Printf("%s %s\n", infoHeading(" INFO "), message)
}
