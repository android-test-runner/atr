package main

import (
	"fmt"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/console"
	"gopkg.in/urfave/cli.v1"
	"strings"
)

var doctorCommand = cli.Command{
	Name:   "doctor",
	Usage:  "Helps with diagnosing atr problems",
	Action: doctorAction,
}

func doctorAction(c *cli.Context) {
	fmt.Printf("Results:\n")
	adbVersion, adbError := adb.New().Version()
	if adbError != nil {
		printError("adb is not installed or not reachable",
			"- Please install the Android Debug Bridge (ADB). See https://developer.android.com/studio/command-line/adb for more information.\n"+
				"- Make sure adb can be executed from the command line. You might need to add the path to adb to the PATH directory.")
	} else {
		printOk(fmt.Sprintf("adb version %v installed", adbVersion))
	}
}

func printOk(message string) {
	checkMark := console.Color("\u2713", console.Green)
	printResult(checkMark, message, "")
}

func printError(message string, details string) {
	errorSign := console.Color("\u2718", console.Red)
	printResult(errorSign, message, details)
}

func printResult(indicator string, message string, details string) {
	fmt.Printf("\t%v %v\n", indicator, message)
	if details == "" {
		return
	}
	// Ensure indentation is correct
	details = strings.Replace(details, "\n", "\n\t  ", -1)
	fmt.Printf("\t  %v\n", details)
}
