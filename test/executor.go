package test

import (
	"fmt"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
)

type Config struct {
	Apk          *apks.Apk
	TestApk      *apks.Apk
	TestRunner   string
	Tests        []Test
	OutputFolder string
}

type Executor interface {
	ExecuteTests(config Config, devices []devices.Device) error
}

type executorImpl struct {
	installer Installer
}

func NewExecutor() Executor {
	return executorImpl{
		installer: NewInstaller(),
	}
}

func (executor executorImpl) ExecuteTests(config Config, devices []devices.Device) error {
	for _, d := range devices {
		apkInstallError := executor.installer.Reinstall(config.Apk, d)
		if apkInstallError != nil {
			return apkInstallError
		}
		testApkInstallError := executor.installer.Reinstall(config.TestApk, d)
		if testApkInstallError != nil {
			return testApkInstallError
		}
		testResults := executor.executeTests(config, d)

		fmt.Printf("Results %v\n", testResults)
	}
	return nil
}

func (executorImpl) executeTests(testConfig Config, device devices.Device) []TestResult {
	var results []TestResult
	for _, t := range testConfig.Tests {
		output, err := adb.New().ExecuteTest(testConfig.TestApk.PackageName, testConfig.TestRunner, FullName(t), device.Serial)
		results = append(results, ResultFromOutput(t, err, output))
	}

	return results
}
