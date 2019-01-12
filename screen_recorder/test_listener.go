package screen_recorder

import (
	"fmt"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
)

type testListener struct {
	device         devices.Device
	writer         output.Writer
	screenRecorder ScreenRecorder
}

func NewTestListener(device devices.Device, writer output.Writer) test_listener.TestListener {
	return &testListener{
		device:         device,
		writer:         writer,
		screenRecorder: New(device),
	}
}

func (listener *testListener) BeforeTestSuite() {}

func (listener *testListener) AfterTestSuite() {}

func (listener *testListener) BeforeTest(test test.Test) {
	errStartScreenRecording := listener.screenRecorder.StartRecording(test)
	if errStartScreenRecording != nil {
		fmt.Printf("Could not start screen recording: '%v'\n", errStartScreenRecording)
	}
}

func (listener *testListener) AfterTest(r result.Result) []result.Extra {
	errStopScreenRecording := listener.screenRecorder.StopRecording(r.Test)
	if errStopScreenRecording != nil {
		fmt.Printf("Could not save screen recording: '%v'\n", errStopScreenRecording)
	}

	extras := []result.Extra{}
	if r.IsFailure() {
		filePath, errSave := listener.screenRecorder.SaveResult(r.Test, listener.writer)
		if errSave != nil {
			fmt.Printf("Could not save screen recording: '%v'\n", errSave)
		} else {
			extras = append(extras, result.Extra{Name: "Screen Recording", Value: filePath, Type: result.File})
		}
	}

	errRemove := listener.screenRecorder.RemoveRecording(r.Test)

	if errRemove != nil {
		fmt.Printf("Could not remove screen recording: '%v'\n", errRemove)
	}

	return extras
}
