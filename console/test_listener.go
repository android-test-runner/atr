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

func NewTestListener() test_listener.TestListener {
	return &testListener{}
}

func (listener *testListener) BeforeTestSuite(device devices.Device) {
	listener.device = device
}

func (listener *testListener) AfterTestSuite() {}

func (listener *testListener) BeforeTest(test test.Test) {}

func (listener *testListener) AfterTest(result result.Result) {
	var resultOutput string
	noColor := "\033[0m"
	if result.IsFailure() {
		red := "\033[0;31m"
		resultOutput = fmt.Sprintf("%vFAILED%v", red, noColor)
	} else {
		green := "\033[0;32m"
		resultOutput = fmt.Sprintf("%vPASSED%v", green, noColor)
	}
	fmt.Printf(
		"'%v' on '%v': %v\n",
		result.Test.FullName(),
		listener.device.Serial,
		resultOutput)
}
