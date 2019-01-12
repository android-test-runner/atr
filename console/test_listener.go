package console

import (
	"fmt"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
)

type testListener struct {
	device devices.Device
}

func NewTestListener(device devices.Device) test_listener.TestListener {
	return &testListener{device: device}
}

func (listener *testListener) BeforeTestSuite() {}

func (listener *testListener) AfterTestSuite() {}

func (listener *testListener) BeforeTest(test test.Test) {}

func (listener *testListener) AfterTest(r result.Result) []result.Extra {
	var resultOutput string
	if r.IsFailure() {
		resultOutput = Color("FAILED", Red)
	} else {
		resultOutput = Color("PASSED", Green)
	}
	fmt.Printf(
		"'%v' on '%v': %v\n",
		r.Test.FullName(),
		listener.device.Serial,
		resultOutput)

	return []result.Extra{}
}
