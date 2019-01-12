package logcat

import (
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/logging"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
)

type testListener struct {
	device devices.Device
	writer output.Writer
	logcat Logcat
	logger logging.Logger
}

func NewTestListener(device devices.Device, writer output.Writer) test_listener.TestListener {
	return &testListener{
		device: device,
		writer: writer,
		logcat: New(device),
		logger: logging.NewForDevice(device),
	}
}

func (listener *testListener) BeforeTestSuite() {}

func (listener *testListener) AfterTestSuite() {}

func (listener *testListener) BeforeTest(test test.Test) {
	errStartLogcat := listener.logcat.StartRecording(test)
	if errStartLogcat != nil {
		listener.logger.Error("Could not clear logcat", errStartLogcat)
	}
}

func (listener *testListener) AfterTest(r result.Result) []result.Extra {
	errStopLogcat := listener.logcat.StopRecording(r.Test)
	if errStopLogcat != nil {
		listener.logger.Error("Coud not get logcat", errStopLogcat)
	}

	extras := []result.Extra{}
	if r.IsFailure() {
		pathToFile, errSave := listener.logcat.SaveRecording(r.Test, listener.writer)
		if errSave != nil {
			listener.logger.Error("Could not save logcat", errSave)
		} else {
			extras = append(extras, result.Extra{Name: "Logcat", Value: pathToFile, Type: result.File})
		}
	}
	return extras
}
