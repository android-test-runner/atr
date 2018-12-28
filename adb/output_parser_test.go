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
