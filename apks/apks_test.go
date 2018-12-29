package apks

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/mock_aapt"
	"github.com/ybonjour/atr/mock_files"
	"testing"
)

func TestGetApk(t *testing.T) {
	path := "path/apk.apk"
	packageName := "packageName"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	aaptMock := mock_aapt.NewMockAapt(ctrl)
	aaptMock.EXPECT().PackageName(path).Return(packageName, nil)
	filesMock := mock_files.NewMockFiles(ctrl)
	filesMock.EXPECT().CanAccess(path).Return(true)
	apks := apksImpl{
		aapt:  aaptMock,
		files: filesMock,
	}

	apk, err := apks.GetApk(path)

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'.", err))
	}
	expectedApk := Apk{Path: path, PackageName: packageName}
	if expectedApk != apk {
		t.Error(fmt.Sprintf("Expected apk '%v' but got '%v'", expectedApk, apk))
	}
}

func TestGetsErrorIfPathIsNotApk(t *testing.T) {
	path := "path/somefile.txt"
	apks := apksImpl{}

	_, err := apks.GetApk(path)

	if err == nil {
		t.Error("Expected error because file is not an apk, but didn't get an error.")
	}
}

func TestGetsErrorIfApkDoesNotExist(t *testing.T) {
	path := "non-existing-path/apk.apk"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filesMock := mock_files.NewMockFiles(ctrl)
	filesMock.EXPECT().CanAccess(path).Return(false)
	apks := apksImpl{
		files: filesMock,
	}

	_, err := apks.GetApk(path)

	if err == nil {
		t.Error("Expected error because file does not exist, but didn't get an error.")
	}
}

func TestGetsErrorIfPackageNameCanNotBeRetrieved(t *testing.T) {
	path := "path/apk.apk"
	expectedError := errors.New("could not get packagename")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	aaptMock := mock_aapt.NewMockAapt(ctrl)
	aaptMock.EXPECT().PackageName(path).Return("", expectedError)
	filesMock := mock_files.NewMockFiles(ctrl)
	filesMock.EXPECT().CanAccess(path).Return(true)
	apks := apksImpl{
		aapt:  aaptMock,
		files: filesMock,
	}

	_, err := apks.GetApk(path)

	if expectedError != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'.", expectedError, err))
	}
}
