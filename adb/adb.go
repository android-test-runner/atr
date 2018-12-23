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

func ExecuteAllTests(packageName string, testRunner string) error {
	packageArgument := fmt.Sprintf("%v/%v", packageName, testRunner)
	return command.Execute(exec.Command("adb", "shell", "am", "instrument", "-w", packageArgument))
}
