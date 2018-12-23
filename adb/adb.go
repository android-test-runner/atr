package adb

import (
	"fmt"
	"github.com/ybonjour/atr/command"
	"os/exec"
)

func Install(apkPath string) error {
	return command.Execute(exec.Command("adb", "install", apkPath))
}

func Uninstall(name string) error {
	return command.Execute(exec.Command("adb", "uninstall", name))
}

func Execute(packageName string, testRunner string, test string) error {
	arguments := []string{
		"shell",
		"am",
		"instrument",
		"-w",
		fmt.Sprintf("-e class %v", test),
		fmt.Sprintf("%v/%v", packageName, testRunner),
	}
	return command.Execute(exec.Command("adb", arguments...))
}
