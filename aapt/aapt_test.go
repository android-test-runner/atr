package aapt

import (
	"errors"
	"fmt"
	"os/exec"
	"testing"
)

type commandExecutorMock struct {
	executeError  error
	executeOutput string
}

func (commandExecutor commandExecutorMock) Execute(cmd *exec.Cmd) error {
	return commandExecutor.executeError
}

func (commandExecutor commandExecutorMock) ExecuteOutput(cmd *exec.Cmd) (string, error) {
	return commandExecutor.executeOutput, commandExecutor.executeError
}

type outputParserMock struct {
	parsedPackageName      string
	parsedPackageNameError error
	parsedTestRunner       string
	parsedTestRunnerError  error
}

func (outputParser outputParserMock) ParsePackageName(out string) (string, error) {
	return outputParser.parsedPackageName, outputParser.parsedPackageNameError
}

func (outputParser outputParserMock) ParseTestRunner(out string) (string, error) {
	return outputParser.parsedTestRunner, outputParser.parsedTestRunnerError
}

func TestPackageName(t *testing.T) {
	parsedPackageName := "parsedPackageName"
	aapt := aaptImpl{
		commandExecutor: commandExecutorMock{executeError: nil, executeOutput: "unparsedPackageName"},
		outputParser:    outputParserMock{parsedPackageName: parsedPackageName, parsedPackageNameError: nil},
	}

	result, err := aapt.PackageName("apkPath")

	if err != nil {
		t.Error(fmt.Sprintf("Got an error but non was expected"))
	}
	if result != parsedPackageName {
		t.Error(fmt.Sprintf("Got packagename %v instead of %v", result, parsedPackageName))
	}
}

func TestPackageNameReturnsCommandError(t *testing.T) {
	expectedErr := errors.New("Command execution failed.")
	aapt := aaptImpl{
		commandExecutor: commandExecutorMock{executeError: expectedErr},
		outputParser:    outputParserMock{parsedPackageName: "", parsedPackageNameError: nil},
	}

	_, err := aapt.PackageName("apkPath")

	if expectedErr != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'", expectedErr, err))
	}
}

func TestPackageNameReturnsParseError(t *testing.T) {
	expectedErr := errors.New("Parsing failed.")
	aapt := aaptImpl{
		commandExecutor: commandExecutorMock{executeError: nil},
		outputParser:    outputParserMock{parsedPackageName: "", parsedPackageNameError: expectedErr},
	}

	_, err := aapt.PackageName("apkPath")

	if expectedErr != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'", expectedErr, err))
	}
}

func TestTestRunner(t *testing.T) {
	expectedTestRunner := "parsedTestRunner"
	aapt := aaptImpl{
		commandExecutor: commandExecutorMock{executeError: nil, executeOutput: "unparsed"},
		outputParser:    outputParserMock{parsedTestRunner: expectedTestRunner, parsedTestRunnerError: nil},
	}

	testRunner, err := aapt.TestRunner("apkPath")

	if err != nil {
		t.Error(fmt.Sprintf("Got an error but non was expected"))
	}
	if expectedTestRunner != testRunner {
		t.Error(fmt.Sprintf("Got test runner '%v' instead of '%v'", testRunner, expectedTestRunner))
	}
}

func TestTestRunnerReturnsCommandError(t *testing.T) {
	expectedErr := errors.New("Command execution failed.")
	aapt := aaptImpl{
		commandExecutor: commandExecutorMock{executeError: expectedErr},
		outputParser:    outputParserMock{parsedTestRunner: "", parsedTestRunnerError: nil},
	}

	_, err := aapt.TestRunner("apkPath")

	if expectedErr != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'", expectedErr, err))
	}
}

func TestTestRunnerReturnsParseError(t *testing.T) {
	expectedErr := errors.New("Parsing failed.")
	aapt := aaptImpl{
		commandExecutor: commandExecutorMock{executeError: nil},
		outputParser:    outputParserMock{parsedTestRunner: "", parsedTestRunnerError: expectedErr},
	}

	_, err := aapt.TestRunner("apkPath")

	if expectedErr != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'", expectedErr, err))
	}
}
