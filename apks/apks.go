package apks

import (
	"errors"
	"github.com/ybonjour/atr/aapt"
	"os"
	"strings"
)

type Apk struct {
	Path        string
	PackageName string
}

func GetApk(path string) (*Apk, error) {
	if !strings.HasSuffix(path, ".apk") {
		return nil, errors.New("APK has no .apk ending")
	}
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	packageName, packageNameError := aapt.New().PackageName(path)
	if packageNameError != nil {
		return nil, packageNameError
	}

	apk := Apk{
		Path:        path,
		PackageName: packageName,
	}

	return &apk, nil
}
