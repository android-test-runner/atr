package main

import (
	"fmt"
	"github.com/android-test-runner/atr/aapt"
	"github.com/android-test-runner/atr/adb"
	"github.com/android-test-runner/atr/console"
	"gopkg.in/urfave/cli.v1"
	"strings"
)

var doctorCommand = cli.Command{
	Name:   "doctor",
	Usage:  "Helps with diagnosing atr problems",
	Action: doctorAction,
}

type check func() bool

func doctorAction(c *cli.Context) {
	fmt.Printf("Hi, I am the atr doctor and I am going to do some examinations:\n%v\n", console.AsciiDoctor())
	checks := []check{
		checkAdb,
		checkAapt,
		checkConnectedDevices,
	}
	passed := executeChecks(checks)
	if passed {
		fmt.Printf(console.Color("\nAll good. Have a nice day!\n", console.Green))
	} else {
		fmt.Printf(console.Color("\nI detected some problems that need to be fixed before you can run tests with atr.\n", console.Red))
	}
}

func executeChecks(checks []check) bool {
	overallResult := true
	for _, c := range checks {
		result := c()
		overallResult = overallResult && result
	}
	return overallResult
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

func checkConnectedDevices() bool {
	connectedDevices, err := adb.New().ConnectedDevices()
	if err != nil {
		printError("unable to get connected devices", err.Error())
		return false
	} else if len(connectedDevices) == 0 {
		printError("no connected devices",
			"- Make sure a device is plugged in or an emluator is started.\n"+
				"- Make sure the devices are reachable through adb: \n"+
				"  Execute 'adb devices' and check if the devices are listed with status `device`.")
		return false
	} else {
		printOk(fmt.Sprintf("%v connected devices", len(connectedDevices)))
		return true
	}
}

func checkAapt() bool {
	aaptVersion, err := aapt.New().Version()
	if err != nil {
		printError("aapt is not installed or not in PATH",
			"- Install the Android Asset Packaging Tool (aapt). See https://developer.android.com/studio/command-line/aapt2 \n"+
				"- Make sure aapt can be executed from the command line. You might need to add the path to aapt to the PATH environment variable")
		return false
	} else {
		printOk(fmt.Sprintf("aapt version %v installed", aaptVersion))
		return true
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
	indentSmall := "  "
	fmt.Printf("%v%v %v\n", indentSmall, indicator, message)
	if details == "" {
		return
	}
	// Ensure indentation is correct
	indentBig := "    "
	details = strings.Replace(details, "\n", fmt.Sprintf("\n%v", indentBig), -1)
	fmt.Printf("%v%v\n", indentBig, details)
}
