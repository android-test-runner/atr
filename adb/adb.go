package adb

import (
	"fmt"
	"github.com/ybonjour/atr/command"
	"os/exec"
)

func Install(apkPath string, deviceSerial string) error {
	return command.Execute(exec.Command("adb", "-s", deviceSerial, "install", apkPath))
}

func Uninstall(packageName string, deviceSerial string) error {
	return command.Execute(exec.Command("adb", "-s", deviceSerial, "uninstall", packageName))
}

func ExecuteTest(packageName string, testRunner string, test string, deviceSerial string) error {
	arguments := []string{
		"-s",
		deviceSerial,
		"shell",
		"am",
		"instrument",
		"-w",
		fmt.Sprintf("-e class %v", test),
		fmt.Sprintf("%v/%v", packageName, testRunner),
	}
	return command.Execute(exec.Command("adb", arguments...))
}
