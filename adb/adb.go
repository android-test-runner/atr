package adb

import (
	"github.com/ybonjour/atr/command"
	"os/exec"
)

func Install(apkPath string) error {
	return command.Execute(exec.Command("adb", "install", apkPath))
}

func Uninstall(name string) error {
	return command.Execute(exec.Command("adb", "uninstall", name))
}
