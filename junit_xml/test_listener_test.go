package junit_xml

import (
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/apks"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/files"
	"github.com/ybonjour/atr/mock_junit_xml"
	"github.com/ybonjour/atr/mock_output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"testing"
)

func TestCollectsAndWritesResults(t *testing.T) {
	apk := apks.Apk{}
	device := devices.Device{}
	testResult1 := result.Result{Test: test.Test{Class: "class", Method: "method1"}}
	testResult2 := result.Result{Test: test.Test{Class: "class", Method: "method2"}}
	xmlFile := files.File{}
	ctrl := gomock.NewController(t)
	formatterMock := mock_junit_xml.NewMockFormatter(ctrl)
	formatterMock.EXPECT().Format([]result.Result{testResult1, testResult2}, apk).Return(xmlFile, nil)
	writerMock := mock_output.NewMockWriter(ctrl)
	writerMock.EXPECT().WriteFile(xmlFile, device)
	listener := testListener{
		formatter: formatterMock,
		writer:    writerMock,
		apk:       apk,
	}
	listener.BeforeTestSuite(device)
	listener.AfterTest(testResult1)
	listener.AfterTest(testResult2)

	listener.AfterTestSuite()

	ctrl.Finish()
}
