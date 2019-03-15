package test_executor

import (
	"fmt"
	"github.com/android-test-runner/atr/adb"
	"github.com/android-test-runner/atr/apks"
	"github.com/android-test-runner/atr/devices"
)

type Installer interface {
	Reinstall(apk apks.Apk, device devices.Device) error
}

type installerImpl struct {
	adb adb.Adb
}

func NewInstaller() Installer {
	return installerImpl{
		adb: adb.New(),
	}
}

func (installer installerImpl) Reinstall(apk apks.Apk, device devices.Device) error {
	resultUninstall := installer.adb.Uninstall(apk.PackageName, device.Serial)
	if resultUninstall.Error != nil {
		// Most likely the uninstall failed, because the package has never been installed.
		// That is why we ignore the error and try to install the package anyways.
		fmt.Printf("Could not uninstall apk on device '%v'. Try to install it anyways.\n", device.Serial)
	}

	resultInstall := installer.adb.Install(apk.Path, device.Serial)
	if resultInstall.Error != nil {
		return resultInstall.Error
	}

	return nil
}
