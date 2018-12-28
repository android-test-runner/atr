package test

import (
	"fmt"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apk"
	"github.com/ybonjour/atr/device"
)

type TestConfig struct {
	Apk          *apk.Apk
	TestApk      *apk.Apk
	TestRunner   string
	Tests        []Test
	OutputFolder string
}

func ExecuteTests(testConfig TestConfig, devices []device.Device) error {
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

func executeTests(testConfig TestConfig, device device.Device) []TestResult {
	var results []TestResult
	for _, t := range testConfig.Tests {
		output, err := adb.ExecuteTest(testConfig.TestApk.PackageName, testConfig.TestRunner, FullName(t), device.Serial)
		results = append(results, ResultFromOutput(t, err, output))
	}

	return results
}

func reinstall(apk *apk.Apk, device device.Device) error {
	apkUninstallError := adb.Uninstall(apk.PackageName, device.Serial)
	if apkUninstallError != nil {
		fmt.Println("Could not uninstall apk. Try to install it anyways.")
	}

	apkInstallError := adb.Install(apk.Path, device.Serial)
	if apkInstallError != nil {
		return apkInstallError
	}

	return nil
}
