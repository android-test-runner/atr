package junit_xml

import (
	"github.com/android-test-runner/atr/apks"
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/files"
	"github.com/android-test-runner/atr/mock_junit_xml"
	"github.com/android-test-runner/atr/mock_logging"
	"github.com/android-test-runner/atr/mock_output"
	"github.com/android-test-runner/atr/result"
	"github.com/android-test-runner/atr/test"
	"github.com/golang/mock/gomock"
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
	loggerMock := mock_logging.NewMockLogger(ctrl)
	allowLogging(loggerMock)
	listener := testListener{
		device:    device,
		formatter: formatterMock,
		writer:    writerMock,
		apk:       apk,
		results:   []result.Result{},
		logger:    loggerMock,
	}
	listener.BeforeTestSuite()
	listener.AfterTest(testResult1)
	listener.AfterTest(testResult2)

	listener.AfterTestSuite()

	ctrl.Finish()
}

func allowLogging(loggerMock *mock_logging.MockLogger) {
	loggerMock.EXPECT().Debug(gomock.Any()).AnyTimes()
	loggerMock.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
}
