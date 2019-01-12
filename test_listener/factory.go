package test_listener

import (
	"github.com/ybonjour/atr/devices"
)

type Factory interface {
	ForDevice(device devices.Device) []TestListener
}
