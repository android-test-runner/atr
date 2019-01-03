package screen_recorder

import (
	"fmt"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
)

type screenRecorderListener struct {
	screenRecorderFactory Factory
	screenRecorder        ScreenRecorder
}

func NewScreenRecorderListener(writer output.Writer) test_listener.TestListener {
	return &screenRecorderListener{
		screenRecorderFactory: NewFactory(writer),
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
	errStopScreenRecording := listener.screenRecorder.StopRecording(result.Test, result.ShallSaveResult())
	if errStopScreenRecording != nil {
		fmt.Printf("Could not save screen recording: '%v'\n", errStopScreenRecording)
	}
}
