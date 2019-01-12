package junit_xml

import (
	"fmt"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"github.com/ybonjour/atr/test_listener"
)

type testListener struct {
	device    devices.Device
	formatter Formatter
	writer    output.Writer
	apk       apks.Apk
	results   []result.Result
}

func NewTestListener(device devices.Device, writer output.Writer, apk apks.Apk) test_listener.TestListener {
	return &testListener{
		device:    device,
		formatter: NewFormatter(),
		apk:       apk,
		writer:    writer,
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
	file, errFormat := listener.formatter.Format(listener.results, listener.apk)
	if errFormat != nil {
		fmt.Printf("Could not format xml junit results for device '%v': '%v'", listener.device, errFormat)
		return
	}
	filePath, errWrite := listener.writer.WriteFile(file, listener.device)
	if errWrite != nil {
		fmt.Printf("Could not write xml junit report to file '%v': '%v'", filePath, errWrite)
	}
}
