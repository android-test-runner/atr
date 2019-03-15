package console

import (
	"fmt"
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/logging"
	"github.com/android-test-runner/atr/result"
	"github.com/android-test-runner/atr/test"
	"github.com/android-test-runner/atr/test_listener"
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
