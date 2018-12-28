package apks

import (
	"errors"
	"fmt"
	"github.com/ybonjour/atr/aapt"
	"github.com/ybonjour/atr/files"
	"strings"
)

type Apk struct {
	Path        string
	PackageName string
}

type Apks interface {
	GetApk(path string) (*Apk, error)
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

func (apks apksImpl) GetApk(path string) (*Apk, error) {
	if !strings.HasSuffix(path, ".apk") {
		return nil, errors.New(fmt.Sprint("APK '%v' has no .apk ending.", path))
	}

	if !apks.files.CanAccess(path) {
		return nil, errors.New(fmt.Sprintf("Can not access APK '%v'.", path))
	}

	packageName, packageNameError := apks.aapt.PackageName(path)
	if packageNameError != nil {
		return nil, packageNameError
	}

	apk := Apk{
		Path:        path,
		PackageName: packageName,
	}

	return &apk, nil
}
