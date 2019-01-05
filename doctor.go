package main

import (
	"fmt"
	"github.com/ybonjour/atr/aapt"
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
	checkAdbResult := checkAdb()
	checkAapt()
	if checkAdbResult {
		checkConnectedDevices()
	}
}

func checkAdb() bool {
	adbVersion, err := adb.New().Version()
	if err != nil {
		printError("adb is not installed or not in PATH",
			"- Install the Android Debug Bridge (ADB). See https://developer.android.com/studio/command-line/adb for more information.\n"+
				"- Make sure adb can be executed from the command line. You might need to add the path to adb to the PATH environment variable.")
		return false
	} else {
		printOk(fmt.Sprintf("adb version %v installed", adbVersion))
		return true
	}
}

func checkConnectedDevices() {
	connectedDevices, err := adb.New().ConnectedDevices()
	if err != nil {
		printError("unable to get connected devices", err.Error())
	} else if len(connectedDevices) == 0 {
		printError("no connected devices",
			"- Make sure a device is plugged in or an emluator is started.\n"+
				"- Make sure the devices are reachable through adb: \n"+
				"  Execute 'adb devices' and check if the devices are listed with status `device`.")
	} else {
		printOk(fmt.Sprintf("%v connected devices", len(connectedDevices)))
	}
}

func checkAapt() {
	aaptVersion, err := aapt.New().Version()
	if err != nil {
		printError("aapt is not installed or not in PATH",
			"- Install the Android Asset Packaging Tool (aapt). See https://developer.android.com/studio/command-line/aapt2 \n"+
				"- Make sure aapt can be executed from the command line. You might need to add the path to aapt to the PATH environment variable")
	} else {
		printOk(fmt.Sprintf("aapt version %v installed", aaptVersion))
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
