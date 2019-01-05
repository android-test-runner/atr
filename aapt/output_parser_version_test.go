package aapt

import (
	"fmt"
	"testing"
)

func TestParsesVersion(t *testing.T) {
	expectedVersion := "v0.2-4355572"
	out := fmt.Sprintf("Android Asset Packaging Tool, %v", expectedVersion)

	version, err := newOutputParser().ParseVersion(out)

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}

	if expectedVersion != version {
		t.Error(fmt.Sprintf("Expected version '%v' but got '%v'", expectedVersion, version))
	}
}

func TestParsesVersionWithOtherLines(t *testing.T) {
	expectedVersion := "v0.2-4355572"
	out := fmt.Sprintf("other line\nAndroid Asset Packaging Tool, %v\n other line", expectedVersion)

	version, err := newOutputParser().ParseVersion(out)

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}

	if expectedVersion != version {
		t.Error(fmt.Sprintf("Expected version '%v' but got '%v'", expectedVersion, version))
	}
}

func TestReturnsErrorIfNoVersionFound(t *testing.T) {
	out := "some output without a version"

	_, err := newOutputParser().ParseVersion(out)

	if err == nil {
		t.Error("Expected an error but did not get one")
	}
}
