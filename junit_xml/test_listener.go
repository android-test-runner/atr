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
	formatter Formatter
	writer    output.Writer
	apk       apks.Apk
	results   map[devices.Device][]result.Result
}

func NewTestListener(writer output.Writer, apk apks.Apk) test_listener.TestListener {
	return &testListener{
		formatter: NewFormatter(),
		apk:       apk,
		writer:    writer,
		results:   map[devices.Device][]result.Result{},
	}
}

func (listener *testListener) BeforeTestSuite(device devices.Device) {}

func (listener *testListener) BeforeTest(test test.Test, device devices.Device) {}

func (listener *testListener) AfterTest(r result.Result, device devices.Device) []result.Extra {
	listener.results[device] = append(listener.results[device], r)
	return []result.Extra{}
}

func (listener *testListener) AfterTestSuite(device devices.Device) {
	file, errFormat := listener.formatter.Format(listener.results[device], listener.apk)
	if errFormat != nil {
		fmt.Printf("Could not format xml junit results for device '%v': '%v'", device, errFormat)
		return
	}
	filePath, errWrite := listener.writer.WriteFile(file, device)
	if errWrite != nil {
		fmt.Printf("Could not write xml junit report to file '%v': '%v'", filePath, errWrite)
	}
}
