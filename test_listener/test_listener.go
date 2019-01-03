package test_listener

import (
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
)

type TestListener interface {
	BeforeTestSuite(device devices.Device)
	BeforeTest(test test.Test)
	AfterTest(result result.Result)
}
