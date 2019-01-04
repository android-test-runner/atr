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
	commandExecutor command.Executor
	outputParser    outputParser
}

func New() Aapt {
	return aaptImpl{
		commandExecutor: command.NewExecutor(),
		outputParser:    newOutputParser(),
	}
}

func (aapt aaptImpl) PackageName(apkPath string) (string, error) {
	result := aapt.commandExecutor.Execute(exec.Command("aapt", "dump", "badging", apkPath))
	if result.Error != nil {
		return "", result.Error
	}

	return aapt.outputParser.ParsePackageName(result.StdOut)
}

func (aapt aaptImpl) TestRunner(apkPath string) (string, error) {
	arguments := []string{
		"dunmp",
		"xmltree",
		apkPath,
		"AndroidManifest.xml",
	}

	out, err := aapt.commandExecutor.ExecuteOutput(exec.Command("aapt", arguments...))
	if err != nil {
		return "", err
	}

	return aapt.outputParser.ParseTestRunner(out)
}
