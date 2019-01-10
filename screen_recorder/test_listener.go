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
	writer         output.Writer
	screenRecorder map[devices.Device]ScreenRecorder
}

func NewTestListener(writer output.Writer) test_listener.TestListener {
	return &testListener{
		writer:         writer,
		screenRecorder: map[devices.Device]ScreenRecorder{},
	}
}

func (listener *testListener) BeforeTestSuite(device devices.Device) {
	listener.screenRecorder[device] = New(device)
}

func (listener *testListener) AfterTestSuite(device devices.Device) {}

func (listener *testListener) BeforeTest(test test.Test, device devices.Device) {
	errStartScreenRecording := listener.screenRecorder[device].StartRecording(test)
	if errStartScreenRecording != nil {
		fmt.Printf("Could not start screen recording: '%v'\n", errStartScreenRecording)
	}
}

func (listener *testListener) AfterTest(r result.Result, device devices.Device) []result.Extra {
	errStopScreenRecording := listener.screenRecorder[device].StopRecording(r.Test)
	if errStopScreenRecording != nil {
		fmt.Printf("Could not save screen recording: '%v'\n", errStopScreenRecording)
	}

	extras := []result.Extra{}
	if r.IsFailure() {
		filePath, errSave := listener.screenRecorder[device].SaveResult(r.Test, listener.writer)
		if errSave != nil {
			fmt.Printf("Could not save screen recording: '%v'\n", errSave)
		} else {
			extras = append(extras, result.Extra{Name: "Screen Recording", Value: filePath, Type: result.File})
		}
	}

	errRemove := listener.screenRecorder[device].RemoveRecording(r.Test)

	if errRemove != nil {
		fmt.Printf("Could not remove screen recording: '%v'\n", errRemove)
	}

	return extras
}
