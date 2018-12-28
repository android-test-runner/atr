package aapt

import (
	"github.com/ybonjour/atr/command"
	"os/exec"
)

type Aapt interface {
	PackageName(apkPath string) (string, error)
	TestRunner(apkPath string) (string, error)
}

type aaptImpl struct {
	outputParser outputParser
}

func New() Aapt {
	return aaptImpl{
		outputParser: newOutputParser(),
	}
}

func (aapt aaptImpl) PackageName(apkPath string) (string, error) {
	out, err := command.ExecuteOutput(exec.Command("aapt", "dump", "badging", apkPath))
	if err != nil {
		return "", err
	}

	return aapt.outputParser.ParsePackageName(out)
}

func (aapt aaptImpl) TestRunner(apkPath string) (string, error) {
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

	return aapt.outputParser.ParseTestRunner(out)
}
