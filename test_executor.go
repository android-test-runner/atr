package main

import (
	"fmt"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apk"
	"github.com/ybonjour/atr/device"
)

type TestConfig struct {
	Apk        *apk.Apk
	TestApk    *apk.Apk
	TestRunner string
	Tests      []string
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
		testError := executeTests(testConfig, d)
		if testError != nil {
			return testError
		}
	}
	return nil
}

func executeTests(testConfig TestConfig, device device.Device) error {
	for _, test := range testConfig.Tests {
		testError := adb.ExecuteTest(testConfig.TestApk.PackageName, testConfig.TestRunner, test, device.Serial)
		if testError != nil {
			return testError
		}
	}

	return nil
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
