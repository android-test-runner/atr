package test

import (
	"fmt"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
)

type Installer interface {
	Reinstall(apk *apks.Apk, device devices.Device) error
}

type installerImpl struct {
	adb adb.Adb
}

func NewInstaller() Installer {
	return installerImpl{
		adb: adb.New(),
	}
}

func (installer installerImpl) Reinstall(apk *apks.Apk, device devices.Device) error {
	apkUninstallError := installer.adb.Uninstall(apk.PackageName, device.Serial)
	if apkUninstallError != nil {
		fmt.Println("Could not uninstall apk. Try to install it anyways.")
	}

	apkInstallError := installer.adb.Install(apk.Path, device.Serial)
	if apkInstallError != nil {
		return apkInstallError
	}

	return nil
}
