package aapt

import (
	"fmt"
	"testing"
)

func TestParsesPackageName(t *testing.T) {
	expectedPackageName := "ch.yvu.atr"
	out := fmt.Sprintf("package: name='%v'", expectedPackageName)

	packageName, err := ParsePackageName(out)

	verifyPackageName(expectedPackageName, packageName, err, t)
}

func TestParsesPackageNameMultiLines(t *testing.T) {
	expectedPackageName := "ch.yvu.atr"
	out := fmt.Sprintf("some output\n....\npackage: name='%v'", expectedPackageName)

	packageName, err := ParsePackageName(out)

	verifyPackageName(expectedPackageName, packageName, err, t)
}

func TestParsesPackageNameWithOuterInformationOnSameLine(t *testing.T) {
	expectedPackageName := "ch.yvu.atr"
	out := fmt.Sprintf("package: name='%v' some other information", expectedPackageName)

	packageName, err := ParsePackageName(out)

	verifyPackageName(expectedPackageName, packageName, err, t)
}

func TestDoesNotParsePackageNameWithInvalidFormat(t *testing.T) {
	out := "invalid format"

	_, err := ParsePackageName(out)

	if err == nil {
		t.Error("Expected 'package name not found error' but didn't get one.")
	}
}

func verifyPackageName(expected string, actual string, err error, t *testing.T) {
	if err != nil {
		t.Error(fmt.Sprintf("Got unexpected error: %v", err))
	}
	if expected != actual {
		t.Error(fmt.Sprintf("Got package name %v isntead of %v", actual, expected))
	}
}
