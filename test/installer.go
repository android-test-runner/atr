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

type installerImpl struct{}

func NewInstaller() Installer {
	return installerImpl{}
}

func (installer installerImpl) Reinstall(apk *apks.Apk, device devices.Device) error {
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
