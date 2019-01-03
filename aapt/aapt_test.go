package aapt

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/mock_aapt"
	"github.com/ybonjour/atr/mock_command"
	"testing"
)

func TestPackageName(t *testing.T) {
	unparsedPackageName := "unparsedPackageName"
	parsedPackageName := "parsedPackageName"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commandExecutorMock := mock_command.NewMockExecutor(ctrl)
	commandExecutorMock.EXPECT().ExecuteOutput(gomock.Any()).Return(unparsedPackageName, nil)
	outputParserMock := mock_aapt.NewMockoutputParser(ctrl)
	outputParserMock.EXPECT().ParsePackageName(gomock.Eq(unparsedPackageName)).Return(parsedPackageName, nil)
	aapt := aaptImpl{
		commandExecutor: commandExecutorMock,
		outputParser:    outputParserMock,
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commandExecutorMock := mock_command.NewMockExecutor(ctrl)
	commandExecutorMock.EXPECT().ExecuteOutput(gomock.Any()).Return("", expectedErr)
	aapt := aaptImpl{
		commandExecutor: commandExecutorMock,
	}

	_, err := aapt.PackageName("apkPath")

	if expectedErr != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'", expectedErr, err))
	}
}

func TestPackageNameReturnsParseError(t *testing.T) {
	expectedErr := errors.New("Parsing failed.")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commandExecutorMock := mock_command.NewMockExecutor(ctrl)
	commandExecutorMock.EXPECT().ExecuteOutput(gomock.Any()).Return("", nil)
	outputParserMock := mock_aapt.NewMockoutputParser(ctrl)
	outputParserMock.EXPECT().ParsePackageName(gomock.Any()).Return("", expectedErr)
	aapt := aaptImpl{
		commandExecutor: commandExecutorMock,
		outputParser:    outputParserMock,
	}

	_, err := aapt.PackageName("apkPath")

	if expectedErr != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'", expectedErr, err))
	}
}

func TestTestRunner(t *testing.T) {
	unparsedTestRunner := "unparsedTestRunner"
	expectedTestRunner := "parsedTestRunner"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commandExecutorMock := mock_command.NewMockExecutor(ctrl)
	commandExecutorMock.
		EXPECT().
		ExecuteOutput(gomock.Any()).
		Return(unparsedTestRunner, nil)
	outputParserMock := mock_aapt.NewMockoutputParser(ctrl)
	outputParserMock.
		EXPECT().
		ParseTestRunner(unparsedTestRunner).
		Return(expectedTestRunner, nil)
	aapt := aaptImpl{
		commandExecutor: commandExecutorMock,
		outputParser:    outputParserMock,
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commandExecutorMock := mock_command.NewMockExecutor(ctrl)
	commandExecutorMock.EXPECT().ExecuteOutput(gomock.Any()).Return("", expectedErr)
	aapt := aaptImpl{
		commandExecutor: commandExecutorMock,
	}

	_, err := aapt.TestRunner("apkPath")

	if expectedErr != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'", expectedErr, err))
	}
}

func TestTestRunnerReturnsParseError(t *testing.T) {
	expectedErr := errors.New("Parsing failed.")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commandExecutorMock := mock_command.NewMockExecutor(ctrl)
	commandExecutorMock.EXPECT().ExecuteOutput(gomock.Any()).Return("", nil)
	outputParserMock := mock_aapt.NewMockoutputParser(ctrl)
	outputParserMock.EXPECT().ParseTestRunner(gomock.Any()).Return("", expectedErr)
	aapt := aaptImpl{
		commandExecutor: commandExecutorMock,
		outputParser:    outputParserMock,
	}

	_, err := aapt.TestRunner("apkPath")

	if expectedErr != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'", expectedErr, err))
	}
}
