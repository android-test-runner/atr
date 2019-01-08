package test_listener

import (
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
)

type TestListener interface {
	BeforeTestSuite(device devices.Device)
	AfterTestSuite()
	BeforeTest(test test.Test)
	AfterTest(r result.Result) []result.Extra
}
