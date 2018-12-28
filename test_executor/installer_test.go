package test_executor

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/mock_adb"
	"testing"
)

func TestReInstallUninstallsAndInstallsApk(t *testing.T) {
	apk := apks.Apk{
		Path:        "path",
		PackageName: "packageName",
	}
	device := devices.Device{
		Serial: "abcde",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().Uninstall(apk.PackageName, device.Serial).Return(nil)
	adbMock.EXPECT().Install(apk.Path, device.Serial).Return(nil)
	installer := installerImpl{
		adb: adbMock,
	}

	err := installer.Reinstall(apk, device)

	if err != nil {
		t.Error(fmt.Sprintf("Did not expect an error but got '%v'", err))
	}
}

func TestReInstallContinuesIfUninstallFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().Uninstall(gomock.Any(), gomock.Any()).Return(errors.New("uninstall failed"))
	adbMock.EXPECT().Install(gomock.Any(), gomock.Any()).Return(nil)
	installer := installerImpl{
		adb: adbMock,
	}

	err := installer.Reinstall(apks.Apk{}, devices.Device{})

	if err != nil {
		t.Error(fmt.Sprintf("Expected to ignore uninstall error but got '%v'", err))
	}
}

func TestReInstallFailsIfInstallFails(t *testing.T) {
	installError := errors.New("install failed")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().Uninstall(gomock.Any(), gomock.Any()).Return(nil)
	adbMock.EXPECT().Install(gomock.Any(), gomock.Any()).Return(installError)
	installer := installerImpl{
		adb: adbMock,
	}

	err := installer.Reinstall(apks.Apk{}, devices.Device{})

	if err != installError {
		t.Error(fmt.Sprintf("Install error '%v' expected but got '%v", installError, err))
	}
}