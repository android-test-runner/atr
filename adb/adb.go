package adb

import (
	"fmt"
	"github.com/ybonjour/atr/command"
	"os/exec"
)

type Adb interface {
	ConnectedDevices() ([]string, error)
	Install(apkPath string, deviceSerial string) error
	Uninstall(packageName string, deviceSerial string) error
	ExecuteTest(packageName string, testRunner string, test string, deviceSerial string) (string, error)
}

type adbImpl struct {
	outputParser outputParser
}

func New() Adb {
	return adbImpl{
		outputParser: newOutputParser(),
	}
}

func (adb adbImpl) ConnectedDevices() ([]string, error) {
	out, err := command.ExecuteOutput(exec.Command("adb", "devices"))
	if err != nil {
		return nil, err
	}

	return adb.outputParser.ParseConnectedDeviceSerials(out), nil
}

func (adbImpl) Install(apkPath string, deviceSerial string) error {
	return command.Execute(exec.Command("adb", "-s", deviceSerial, "install", apkPath))
}

func (adbImpl) Uninstall(packageName string, deviceSerial string) error {
	return command.Execute(exec.Command("adb", "-s", deviceSerial, "uninstall", packageName))
}

func (adbImpl) ExecuteTest(packageName string, testRunner string, test string, deviceSerial string) (string, error) {
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
	return command.ExecuteOutput(exec.Command("adb", arguments...))
}
