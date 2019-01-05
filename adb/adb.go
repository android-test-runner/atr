package adb

import (
	"fmt"
	"github.com/ybonjour/atr/command"
	"os/exec"
)

type Adb interface {
	Version() (string, error)
	ConnectedDevices() ([]string, error)
	Install(apkPath string, deviceSerial string) command.ExecutionResult
	Uninstall(packageName string, deviceSerial string) command.ExecutionResult
	ExecuteTest(packageName string, testRunner string, test string, deviceSerial string) (string, error)
	ClearLogcat(deviceSerial string) command.ExecutionResult
	GetLogcat(deviceSerial string) (string, error)
	RecordScreen(deviceSerial string, filePath string) (int, error)
	PullFile(deviceSerial string, filePathOnDevice string, filePathLocal string) command.ExecutionResult
	RemoveFile(deviceSerial string, filePathOnDevice string) command.ExecutionResult
}

type adbImpl struct {
	commandExecutor command.Executor
	outputParser    outputParser
}

func New() Adb {
	return adbImpl{
		commandExecutor: command.NewExecutor(),
		outputParser:    newOutputParser(),
	}
}

func (adb adbImpl) Version() (string, error) {
	result := adb.commandExecutor.Execute(exec.Command("adb", "--version"))
	if result.Error != nil {
		return "", result.Error
	}
	version, parseError := adb.outputParser.ParseVersion(result.StdOut)
	if parseError != nil {
		return "Unknown Version", nil
	} else {
		return version, nil
	}
}

func (adb adbImpl) ConnectedDevices() ([]string, error) {
	result := adb.commandExecutor.Execute(exec.Command("adb", "devices"))
	if result.Error != nil {
		return nil, result.Error
	}

	return adb.outputParser.ParseConnectedDeviceSerials(result.StdOut), nil
}

func (adb adbImpl) Install(apkPath string, deviceSerial string) command.ExecutionResult {
	return adb.commandExecutor.Execute(exec.Command("adb", "-s", deviceSerial, "install", apkPath))
}

func (adb adbImpl) Uninstall(packageName string, deviceSerial string) command.ExecutionResult {
	return adb.commandExecutor.Execute(exec.Command("adb", "-s", deviceSerial, "uninstall", packageName))
}

func (adb adbImpl) ClearLogcat(deviceSerial string) command.ExecutionResult {
	return adb.commandExecutor.Execute(exec.Command("adb", "-s", deviceSerial, "logcat", "-c"))
}

func (adb adbImpl) GetLogcat(deviceSerial string) (string, error) {
	result := adb.commandExecutor.Execute(exec.Command("adb", "-s", deviceSerial, "logcat", "-d"))
	return result.StdOut, result.Error
}

func (adb adbImpl) RecordScreen(deviceSerial string, filePath string) (int, error) {
	return adb.commandExecutor.ExecuteInBackground(exec.Command("adb", "-s", deviceSerial, "shell", "screenrecord", filePath))
}

func (adb adbImpl) PullFile(deviceSerial string, filePathOnDevice string, filePathLocal string) command.ExecutionResult {
	return adb.commandExecutor.Execute(exec.Command("adb", "-s", deviceSerial, "pull", filePathOnDevice, filePathLocal))
}

func (adb adbImpl) RemoveFile(deviceSerial string, filePathOnDevice string) command.ExecutionResult {
	return adb.commandExecutor.Execute(exec.Command("adb", "-s", deviceSerial, "shell", "rm", filePathOnDevice))
}

func (adb adbImpl) ExecuteTest(packageName string, testRunner string, test string, deviceSerial string) (string, error) {
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
	result := adb.commandExecutor.Execute(exec.Command("adb", arguments...))
	return result.StdOut, result.Error
}
