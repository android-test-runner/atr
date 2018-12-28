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

type executorImpl struct{}

func NewExecutor() Executor {
	return executorImpl{}
}

func (executor executorImpl) ExecuteTests(config Config, devices []devices.Device) error {
	for _, d := range devices {
		apkInstallError := executor.reinstall(config.Apk, d)
		if apkInstallError != nil {
			return apkInstallError
		}
		testApkInstallError := executor.reinstall(config.TestApk, d)
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

func (executor executorImpl) reinstall(apk *apks.Apk, device devices.Device) error {
	apkUninstallError := adb.New().Uninstall(apk.PackageName, device.Serial)
	if apkUninstallError != nil {
		fmt.Println("Could not uninstall apk. Try to install it anyways.")
	}

	apkInstallError := adb.New().Install(apk.Path, device.Serial)
	if apkInstallError != nil {
		return apkInstallError
	}

	return nil
}
