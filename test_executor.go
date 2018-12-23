package main

import (
	"fmt"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apk"
)

type TestConfig struct {
	Apk        *apk.Apk
	TestApk    *apk.Apk
	TestRunner string
	Tests      []string
}

func ExecuteTests(testConfig TestConfig) error {
	apkInstallError := reinstall(testConfig.Apk)
	if apkInstallError != nil {
		return apkInstallError
	}
	testApkInstallError := reinstall(testConfig.TestApk)
	if testApkInstallError != nil {
		return testApkInstallError
	}
	return executeTests(testConfig)
}

func executeTests(testConfig TestConfig) error {
	for _, test := range testConfig.Tests {
		testError := adb.Execute(testConfig.TestApk.PackageName, testConfig.TestRunner, test)
		if testError != nil {
			return testError
		}
	}
	return nil
}

func reinstall(apk *apk.Apk) error {
	apkUninstallError := adb.Uninstall(apk.PackageName)
	if apkUninstallError != nil {
		fmt.Println("Could not uninstall apk. Try to install it anyways.")
	}

	apkInstallError := adb.Install(apk.Path)
	if apkInstallError != nil {
		return apkInstallError
	}

	return nil
}
