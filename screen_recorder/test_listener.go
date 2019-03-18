package screen_recorder

import (
	"fmt"
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/logging"
	"github.com/android-test-runner/atr/output"
	"github.com/android-test-runner/atr/result"
	"github.com/android-test-runner/atr/test"
	"github.com/android-test-runner/atr/test_listener"
)

type testListener struct {
	device         devices.Device
	writer         output.Writer
	screenRecorder ScreenRecorder
	logger         logging.Logger
}

func NewTestListener(device devices.Device, writer output.Writer) test_listener.TestListener {
	return &testListener{
		device:         device,
		writer:         writer,
		screenRecorder: New(device),
		logger:         logging.NewForDevice(device),
	}
}

func (listener *testListener) BeforeTestSuite() {}

func (listener *testListener) AfterTestSuite() {}

func (listener *testListener) BeforeTest(test test.Test) {
	listener.logger.Debug(logging.TestPrefix("Start screen recording", test))
	err := listener.screenRecorder.StartRecording(test)
	if err != nil {
		listener.logger.Error(logging.TestPrefix("Could not start screen recording", test), err)
	} else {
		listener.logger.Debug(logging.TestPrefix("Successfully started screen recording", test))
	}
}

func (listener *testListener) AfterTest(r result.Result) []result.Extra {
	listener.logger.Debug(logging.TestPrefix("Stop screen recording", r.Test))
	errStopScreenRecording := listener.screenRecorder.StopRecording(r.Test)
	if errStopScreenRecording != nil {
		listener.logger.Error(logging.TestPrefix("Could not stop screen recording", r.Test), errStopScreenRecording)
	}

	extras := []result.Extra{}
	if r.IsFailure() {
		listener.logger.Debug(logging.TestPrefix("Save screen recording to file", r.Test))
		filePath, errSave := listener.screenRecorder.SaveResult(r.Test, listener.writer)
		if errSave != nil {
			listener.logger.Error(logging.TestPrefix("Could not save screen recording to file", r.Test), errSave)
		} else {
			listener.logger.Debug(logging.TestPrefix(fmt.Sprintf("Successfully saved screen recording to file %v", filePath), r.Test))
			extras = append(extras, result.Extra{Name: "Screen Recording", Value: filePath, Type: result.Video})
		}
	}

	listener.logger.Debug(logging.TestPrefix("Removes screen recording on device", r.Test))
	errRemove := listener.screenRecorder.RemoveRecording(r.Test)

	if errRemove != nil {
		listener.logger.Error(logging.TestPrefix("Could not remove screen recording from device", r.Test), errRemove)
	}

	return extras
}
