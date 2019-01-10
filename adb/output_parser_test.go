package adb

import (
	"fmt"
	"testing"
)

func TestParsesConnectedDeviceSerial(t *testing.T) {
	expectedSerial := "abcd"
	out := fmt.Sprintf("%v\tdevice", expectedSerial)

	deviceSerials := newOutputParser().ParseConnectedDeviceSerials(out)

	verifySerials([]string{expectedSerial}, deviceSerials, t)
}

func TestParsesMultipleConnectedDeviceSerial(t *testing.T) {
	expectedSerials := []string{"abcd", "efgh"}
	out := fmt.Sprintf("%v\tdevice\n%v\tdevice", expectedSerials[0], expectedSerials[1])

	deviceSerials := newOutputParser().ParseConnectedDeviceSerials(out)

	verifySerials(expectedSerials, deviceSerials, t)
}

func TestIgnoresUnconnectedDevices(t *testing.T) {
	out := "abcd\tunauthorized"

	deviceSerials := newOutputParser().ParseConnectedDeviceSerials(out)

	if !AreEqual([]string{}, deviceSerials) {
		t.Error("Did not ignore unconnected device.")
	}
}

func TestIgnoresNonDeviceOutput(t *testing.T) {
	out := "Some other output"

	deviceSerials := newOutputParser().ParseConnectedDeviceSerials(out)

	if !AreEqual([]string{}, deviceSerials) {
		t.Error("Did not ignore different output.")
	}
}

func TestParsesVersion(t *testing.T) {
	expectedVersion := "1.0.40"
	out := fmt.Sprintf("Android Debug Bridge version %v", expectedVersion)

	version, err := newOutputParser().ParseVersion(out)

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}

	if version != expectedVersion {
		t.Error(fmt.Sprintf("Expected version '%v' but got '%v'", expectedVersion, version))
	}
}

func TestParsesVersionFromMultipleLines(t *testing.T) {
	expectedVersion := "1.0.40"
	out := fmt.Sprintf("some other line\nAndroid Debug Bridge version %v\nsome other line", expectedVersion)

	version, err := newOutputParser().ParseVersion(out)

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}

	if version != expectedVersion {
		t.Error(fmt.Sprintf("Expected version '%v' but got '%v'", expectedVersion, version))
	}
}

func TestReturnsErrorWhenNoVersionProvided(t *testing.T) {
	out := "some output without a version"

	_, err := newOutputParser().ParseVersion(out)

	if err == nil {
		t.Error("Expected an error because no version present but did not get any")
	}
}

func TestParsesDimensions(t *testing.T) {
	expectedWidth := 1440
	expectedHeight := 2880
	out := fmt.Sprintf("Physical size: %vx%v", expectedWidth, expectedHeight)

	width, height, err := newOutputParser().ParseDimensions(out)

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}

	if width != expectedWidth {
		t.Error(fmt.Sprintf("Expected width to be '%v' but it was '%v'", expectedWidth, width))
	}

	if height != expectedHeight {
		t.Error(fmt.Sprintf("Expected height to be '%v' but it was '%v'", expectedHeight, height))
	}
}

func TestParsesDimensionsFromMultipleLines(t *testing.T) {
	expectedWidth := 1440
	expectedHeight := 2880
	out := fmt.Sprintf("other line\nPhysical size: %vx%v\n other line", expectedWidth, expectedHeight)

	width, height, err := newOutputParser().ParseDimensions(out)

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}

	if width != expectedWidth {
		t.Error(fmt.Sprintf("Expected width to be '%v' but it was '%v'", expectedWidth, width))
	}

	if height != expectedHeight {
		t.Error(fmt.Sprintf("Expected height to be '%v' but it was '%v'", expectedHeight, height))
	}
}

func TestReturnsErrorWhenNoDimensionFound(t *testing.T) {
	out := "not a dimension"

	_, _, err := newOutputParser().ParseDimensions(out)

	if err == nil {
		t.Error("Expected an error because no dimension found but did not get any")
	}
}

func verifySerials(expected, actual []string, t *testing.T) {
	if !AreEqual(expected, actual) {
		t.Error(fmt.Sprintf("Got serials %v instead of %v.", actual, expected))
	}
}

func AreEqual(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}
