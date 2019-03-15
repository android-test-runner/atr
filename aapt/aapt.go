package aapt

import (
	"github.com/android-test-runner/atr/command"
	"os/exec"
)

type Aapt interface {
	Version() (string, error)
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

func (aapt aaptImpl) Version() (string, error) {
	result := aapt.commandExecutor.Execute(exec.Command("aapt", "v"))
	if result.Error != nil {
		return "", result.Error
	}

	version, parseError := aapt.outputParser.ParseVersion(result.StdOut)
	if parseError != nil {
		return "Unknown Version", nil
	} else {
		return version, nil
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

	result := aapt.commandExecutor.Execute(exec.Command("aapt", arguments...))
	if result.Error != nil {
		return "", result.Error
	}

	return aapt.outputParser.ParseTestRunner(result.StdOut)
}
