package test_executor

import (
	"errors"
	"fmt"
	"github.com/android-test-runner/atr/apks"
	"github.com/android-test-runner/atr/command"
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/mock_adb"
	"github.com/golang/mock/gomock"
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
	adbMock.EXPECT().Uninstall(apk.PackageName, device.Serial).Return(executionResultOk())
	adbMock.EXPECT().Install(apk.Path, device.Serial).Return(executionResultOk())
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
	uninstallError := errors.New("uninstall failed")
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().Uninstall(gomock.Any(), gomock.Any()).Return(executionResultError(uninstallError))
	adbMock.EXPECT().Install(gomock.Any(), gomock.Any()).Return(executionResultOk())
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
	adbMock.EXPECT().Uninstall(gomock.Any(), gomock.Any()).Return(executionResultOk())
	adbMock.EXPECT().Install(gomock.Any(), gomock.Any()).Return(executionResultError(installError))
	installer := installerImpl{
		adb: adbMock,
	}

	err := installer.Reinstall(apks.Apk{}, devices.Device{})

	if err != installError {
		t.Error(fmt.Sprintf("Install error '%v' expected but got '%v", installError, err))
	}
}

func executionResultOk() command.ExecutionResult {
	return command.ExecutionResult{Error: nil}
}

func executionResultError(err error) command.ExecutionResult {
	return command.ExecutionResult{Error: err}
}
