package logcat

import (
	"fmt"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
)

type logcatListener struct {
	writer output.Writer
	logcat Logcat
}

func NewLogcatListener(writer output.Writer) test_listener.TestListener {
	return &logcatListener{
		writer: writer,
	}
}

func (listener *logcatListener) BeforeTestSuite(device devices.Device) {
	listener.logcat = NewLogcat(device)
}

func (listener *logcatListener) BeforeTest(test test.Test) {
	errStartLogcat := listener.logcat.StartRecording(test)
	if errStartLogcat != nil {
		fmt.Printf("Could not clear logcat: '%v'\n", errStartLogcat)
	}
}

func (listener *logcatListener) AfterTest(result result.Result) {
	errStopLogcat := listener.logcat.StopRecording(result.Test)
	if errStopLogcat != nil {
		fmt.Printf("Could not get logcat: '%v'\n", errStopLogcat)
	}

	if result.ShallSaveResult() {
		errSave := listener.logcat.SaveRecording(result.Test, listener.writer)
		if errSave != nil {
			fmt.Printf("Could not save logcat: '%v'\n", errSave)
		}
	}
}
