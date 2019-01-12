package console

import (
	"fmt"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/logging"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
)

type testListener struct {
	device devices.Device
	logger logging.Logger
}

func NewTestListener(device devices.Device) test_listener.TestListener {
	return &testListener{device: device, logger: logging.NewForDevice(device)}
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

	resultMessage := fmt.Sprintf(
		"%v: %v\n",
		r.Test.FullName(),
		resultOutput)
	listener.logger.Info(resultMessage)

	return []result.Extra{}
}
