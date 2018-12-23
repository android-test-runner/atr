package apk

import "github.com/ybonjour/atr/aapt"

type Apk struct {
	Path        string
	PackageName string
}

func GetApk(path string) (*Apk, error) {
	packageName, packageNameError := aapt.PackageName(path)
	if packageNameError != nil {
		return nil, packageNameError
	}

	apk := Apk{
		Path:        path,
		PackageName: packageName,
	}

	return &apk, nil
}
