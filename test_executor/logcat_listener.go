package test_executor

import (
	"fmt"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/logcat"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
)

type logcatListener struct {
	logcatFactory logcat.Factory
	logcat        logcat.Logcat
}

func NewLogcatListener(writer output.Writer) test_listener.TestListener {
	return &logcatListener{
		logcatFactory: logcat.NewFactory(writer),
	}
}

func (listener *logcatListener) BeforeTestSuite(device devices.Device) {
	listener.logcat = listener.logcatFactory.ForDevice(device)
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
		fmt.Printf("Could not save logcat: '%v'\n", errStopLogcat)
	}
}
