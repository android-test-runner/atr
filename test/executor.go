package test

import (
	"fmt"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
)

type TestConfig struct {
	Apk          *apks.Apk
	TestApk      *apks.Apk
	TestRunner   string
	Tests        []Test
	OutputFolder string
}

func ExecuteTests(testConfig TestConfig, devices []devices.Device) error {
	for _, d := range devices {
		apkInstallError := reinstall(testConfig.Apk, d)
		if apkInstallError != nil {
			return apkInstallError
		}
		testApkInstallError := reinstall(testConfig.TestApk, d)
		if testApkInstallError != nil {
			return testApkInstallError
		}
		testResults := executeTests(testConfig, d)

		fmt.Printf("Results %v\n", testResults)
	}
	return nil
}

func executeTests(testConfig TestConfig, device devices.Device) []TestResult {
	var results []TestResult
	for _, t := range testConfig.Tests {
		output, err := adb.New().ExecuteTest(testConfig.TestApk.PackageName, testConfig.TestRunner, FullName(t), device.Serial)
		results = append(results, ResultFromOutput(t, err, output))
	}

	return results
}

func reinstall(apk *apks.Apk, device devices.Device) error {
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
