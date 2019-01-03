package test_executor

import (
	"fmt"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/screen_recorder"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
)

type screenRecorderListener struct {
	screenRecorderFactory screen_recorder.Factory
	screenRecorder        screen_recorder.ScreenRecorder
}

func NewScreenRecorderListener(writer output.Writer) test_listener.TestListener {
	return &screenRecorderListener{
		screenRecorderFactory: screen_recorder.NewFactory(writer),
	}
}

func (listener *screenRecorderListener) BeforeTestSuite(device devices.Device) {
	listener.screenRecorder = listener.screenRecorderFactory.ForDevice(device)
}

func (listener *screenRecorderListener) BeforeTest(test test.Test) {
	errStartScreenRecording := listener.screenRecorder.StartRecording(test)
	if errStartScreenRecording != nil {
		fmt.Printf("Could not start screen recording: '%v'\n", errStartScreenRecording)
	}
}

func (listener *screenRecorderListener) AfterTest(result result.Result) {
	errStopScreenRecording := listener.screenRecorder.StopRecording(result.Test)
	if errStopScreenRecording != nil {
		fmt.Printf("Could not save screen recording: '%v'\n", errStopScreenRecording)
	}
}
