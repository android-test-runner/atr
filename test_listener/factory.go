package test_listener

import (
	"github.com/android-test-runner/atr/devices"
)

type Factory interface {
	ForDevice(device devices.Device) []TestListener
}
