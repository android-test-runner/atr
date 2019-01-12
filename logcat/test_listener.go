package logcat

import (
	"fmt"
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
	listener.logger.Debug(testPrefix("Starts logcat recording", test))
	errStartLogcat := listener.logcat.StartRecording(test)
	if errStartLogcat != nil {
		listener.logger.Error(
			testPrefix("Could not start logcat recording", test),
			errStartLogcat)
	} else {
		listener.logger.Debug(testPrefix("Successfully started logcat recording", test))
	}

}

func (listener *testListener) AfterTest(r result.Result) []result.Extra {
	listener.logger.Debug(testPrefix("Stop logcat recording", r.Test))
	errStopLogcat := listener.logcat.StopRecording(r.Test)
	if errStopLogcat != nil {
		listener.logger.Error(testPrefix("Coud not stop logcatrecording", r.Test), errStopLogcat)
	} else {
		listener.logger.Debug(testPrefix("Successfully stopped logcat recording", r.Test))
	}

	extras := []result.Extra{}
	if r.IsFailure() {
		listener.logger.Debug(testPrefix("Save logcat to file", r.Test))
		pathToFile, errSave := listener.logcat.SaveRecording(r.Test, listener.writer)
		if errSave != nil {
			listener.logger.Error(testPrefix("Could not save logcat to file", r.Test), errSave)
		} else {
			listener.logger.Debug(testPrefix("Successfully saved logcat to file", r.Test))
			extras = append(extras, result.Extra{Name: "Logcat", Value: pathToFile, Type: result.File})
		}
	}
	return extras
}

func testPrefix(message string, test test.Test) string {
	return fmt.Sprintf("%v: %v", test.FullName(), message)
}
