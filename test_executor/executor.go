package test_executor

import (
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
)

type Config struct {
	Apk          apks.Apk
	TestApk      apks.Apk
	TestRunner   string
	Tests        []test.Test
	OutputFolder string
}

type Executor interface {
	Execute(config Config, devices []devices.Device) (map[devices.Device][]result.Result, error)
}

type executorImpl struct {
	installer    Installer
	resultParser result.ResultParser
	adb          adb.Adb
}

func NewExecutor() Executor {
	return executorImpl{
		installer:    NewInstaller(),
		resultParser: result.NewResultParser(),
		adb:          adb.New(),
	}
}

func (executor executorImpl) Execute(config Config, targetDevices []devices.Device) (map[devices.Device][]result.Result, error) {
	resultsByDevice := map[devices.Device][]result.Result{}
	for _, targetDevice := range targetDevices {
		apkInstallError := executor.installer.Reinstall(config.Apk, targetDevice)
		if apkInstallError != nil {
			return nil, apkInstallError
		}
		testApkInstallError := executor.installer.Reinstall(config.TestApk, targetDevice)
		if testApkInstallError != nil {
			return nil, testApkInstallError
		}
		testResults := executor.executeTests(config, targetDevice)

		resultsByDevice[targetDevice] = testResults
	}
	return resultsByDevice, nil
}

func (executor executorImpl) executeTests(testConfig Config, device devices.Device) []result.Result {
	var results []result.Result
	for _, t := range testConfig.Tests {
		output, err := executor.adb.ExecuteTest(testConfig.TestApk.PackageName, testConfig.TestRunner, t.FullName(), device.Serial)
		results = append(results, executor.resultParser.ParseFromOutput(t, err, output))
	}

	return results
}
