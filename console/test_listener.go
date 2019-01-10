package console

import (
	"fmt"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
)

type testListener struct {
}

func NewTestListener() test_listener.TestListener {
	return &testListener{}
}

func (listener *testListener) BeforeTestSuite(device devices.Device) {}

func (listener *testListener) AfterTestSuite(device devices.Device) {}

func (listener *testListener) BeforeTest(test test.Test, device devices.Device) {}

func (listener *testListener) AfterTest(r result.Result, device devices.Device) []result.Extra {
	var resultOutput string
	if r.IsFailure() {
		resultOutput = Color("FAILED", Red)
	} else {
		resultOutput = Color("PASSED", Green)
	}
	fmt.Printf(
		"'%v' on '%v': %v\n",
		r.Test.FullName(),
		device.Serial,
		resultOutput)

	return []result.Extra{}
}
