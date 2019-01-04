package logcat

import (
	"fmt"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
)

type testListener struct {
	writer output.Writer
	logcat Logcat
}

func NewTestListener(writer output.Writer) test_listener.TestListener {
	return &testListener{
		writer: writer,
	}
}

func (listener *testListener) BeforeTestSuite(device devices.Device) {
	listener.logcat = New(device)
}

func (listener *testListener) AfterTestSuite() {}

func (listener *testListener) BeforeTest(test test.Test) {
	errStartLogcat := listener.logcat.StartRecording(test)
	if errStartLogcat != nil {
		fmt.Printf("Could not clear logcat: '%v'\n", errStartLogcat)
	}
}

func (listener *testListener) AfterTest(result result.Result) {
	errStopLogcat := listener.logcat.StopRecording(result.Test)
	if errStopLogcat != nil {
		fmt.Printf("Could not get logcat: '%v'\n", errStopLogcat)
	}

	if result.IsFailure() {
		errSave := listener.logcat.SaveRecording(result.Test, listener.writer)
		if errSave != nil {
			fmt.Printf("Could not save logcat: '%v'\n", errSave)
		}
	}
}
