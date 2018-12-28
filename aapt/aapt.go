package aapt

import (
	"github.com/ybonjour/atr/command"
	"os/exec"
)

func PackageName(apkPath string) (string, error) {
	out, err := command.ExecuteOutput(exec.Command("aapt", "dump", "badging", apkPath))
	if err != nil {
		return "", err
	}

	return ParsePackageName(out)
}

func TestRunner(apkPath string) (string, error) {
	arguments := []string{
		"dunmp",
		"xmltree",
		apkPath,
		"AndroidManifest.xml",
	}

	out, err := command.ExecuteOutput(exec.Command("aapt", arguments...))
	if err != nil {
		return "", err
	}

	return ParseTestRunner(out)
}
