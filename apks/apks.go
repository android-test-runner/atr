package apks

import (
	"errors"
	"fmt"
	"github.com/android-test-runner/atr/aapt"
	"github.com/android-test-runner/atr/files"
	"strings"
)

type Apk struct {
	Path        string
	PackageName string
}

type Apks interface {
	GetApk(path string) (Apk, error)
}

type apksImpl struct {
	aapt  aapt.Aapt
	files files.Files
}

func New() Apks {
	return apksImpl{
		aapt:  aapt.New(),
		files: files.New(),
	}
}

func (apks apksImpl) GetApk(path string) (Apk, error) {
	if !strings.HasSuffix(path, ".apk") {
		return Apk{}, errors.New(fmt.Sprintf("apk '%v' has no .apk ending", path))
	}

	if !apks.files.CanAccess(path) {
		return Apk{}, errors.New(fmt.Sprintf("can not access APK '%v'", path))
	}

	packageName, packageNameError := apks.aapt.PackageName(path)
	if packageNameError != nil {
		return Apk{}, packageNameError
	}

	apk := Apk{
		Path:        path,
		PackageName: packageName,
	}

	return apk, nil
}
