package main

import (
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apk"
)

type TestConfig struct {
	Apk     *apk.Apk
	TestApk *apk.Apk
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

	return nil
}

func reinstall(apk *apk.Apk) error {
	apkUninstallError := adb.Uninstall(apk.PackageName)
	if apkUninstallError != nil {
		return apkUninstallError
	}

	apkInstallError := adb.Install(apk.Path)
	if apkInstallError != nil {
		return apkInstallError
	}

	return nil
}
