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
	listener.logger.Debug(logging.TestPrefix("Start logcat recording", test))
	errStartLogcat := listener.logcat.StartRecording(test)
	if errStartLogcat != nil {
		listener.logger.Error(
			logging.TestPrefix("Could not start logcat recording", test),
			errStartLogcat)
	} else {
		listener.logger.Debug(logging.TestPrefix("Successfully started logcat recording", test))
	}

}

func (listener *testListener) AfterTest(r result.Result) []result.Extra {
	listener.logger.Debug(logging.TestPrefix("Stop logcat recording", r.Test))
	errStopLogcat := listener.logcat.StopRecording(r.Test)
	if errStopLogcat != nil {
		listener.logger.Error(logging.TestPrefix("Could not stop logcat recording", r.Test), errStopLogcat)
	} else {
		listener.logger.Debug(logging.TestPrefix("Successfully stopped logcat recording", r.Test))
	}

	extras := []result.Extra{}
	if r.IsFailure() {
		listener.logger.Debug(logging.TestPrefix("Save logcat to file", r.Test))
		pathToFile, errSave := listener.logcat.SaveRecording(r.Test, listener.writer)
		if errSave != nil {
			listener.logger.Error(logging.TestPrefix("Could not save logcat to file", r.Test), errSave)
		} else {
			listener.logger.Debug(logging.TestPrefix("Successfully saved logcat to file", r.Test))
			extras = append(extras, result.Extra{Name: "Logcat", Value: pathToFile, Type: result.File})
		}
	}
	return extras
}
