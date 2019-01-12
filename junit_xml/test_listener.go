package junit_xml

import (
	"fmt"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/logging"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
)

type testListener struct {
	device    devices.Device
	formatter Formatter
	apk       apks.Apk
	writer    output.Writer
	logger    logging.Logger
	results   []result.Result
}

func NewTestListener(device devices.Device, writer output.Writer, apk apks.Apk) test_listener.TestListener {
	return &testListener{
		device:    device,
		formatter: NewFormatter(),
		apk:       apk,
		writer:    writer,
		logger:    logging.NewForDevice(device),
		results:   []result.Result{},
	}
}

func (listener *testListener) BeforeTestSuite() {}

func (listener *testListener) BeforeTest(test test.Test) {}

func (listener *testListener) AfterTest(r result.Result) []result.Extra {
	listener.results = append(listener.results, r)
	return []result.Extra{}
}

func (listener *testListener) AfterTestSuite() {
	listener.logger.Debug(fmt.Sprintf("Save xml junit results for %v tests", len(listener.results)))
	file, errFormat := listener.formatter.Format(listener.results, listener.apk)
	if errFormat != nil {
		listener.logger.Error("Could not format xml junit results", errFormat)
		return
	}
	filePath, errWrite := listener.writer.WriteFile(file, listener.device)
	if errWrite != nil {
		listener.logger.Error("Could not write xml junit results to file", errWrite)
	} else {
		listener.logger.Debug(fmt.Sprintf("Successfully saved xml junit reports to file %v", filePath))
	}
}
