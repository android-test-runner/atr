package logcat

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/mock_logcat"
	"github.com/ybonjour/atr/mock_output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"testing"
)

func TestBeforeTestStartsLogcatRecording(t *testing.T) {
	targetTest := test.Test{}
	device := devices.Device{Serial: "abcd"}
	ctrl := gomock.NewController(t)
	logcatMock := mock_logcat.NewMockLogcat(ctrl)
	logcatMock.EXPECT().StartRecording(targetTest).Return(nil)

	listener := testListener{
		logcat: map[devices.Device]Logcat{device: logcatMock},
	}
	listener.BeforeTest(targetTest, device)

	ctrl.Finish()
}

func TestAfterTestStopsLogcatRecordingAndSavesForFailedResult(t *testing.T) {
	targetTest := test.Test{}
	testResult := result.Result{Test: targetTest, Status: result.Failed}
	device := devices.Device{Serial: "abcd"}
	ctrl := gomock.NewController(t)
	writer := mock_output.NewMockWriter(ctrl)
	logcatMock := mock_logcat.NewMockLogcat(ctrl)
	logcatMock.EXPECT().StopRecording(targetTest).Return(nil)
	logcatMock.EXPECT().SaveRecording(targetTest, writer).Return("", nil)
	listener := testListener{
		writer: writer,
		logcat: map[devices.Device]Logcat{device: logcatMock},
	}
	listener.AfterTest(testResult, device)

	ctrl.Finish()
}

func TestAfterTestStopsLogcatRecordingForPassedResult(t *testing.T) {
	targetTest := test.Test{}
	testResult := result.Result{Test: targetTest, Status: result.Passed}
	device := devices.Device{Serial: "abcd"}
	ctrl := gomock.NewController(t)
	writer := mock_output.NewMockWriter(ctrl)
	logcatMock := mock_logcat.NewMockLogcat(ctrl)
	logcatMock.EXPECT().StopRecording(targetTest).Return(nil)
	logcatMock.EXPECT().SaveRecording(targetTest, writer).Return("", nil).Times(0)
	listener := testListener{
		writer: writer,
		logcat: map[devices.Device]Logcat{device: logcatMock},
	}
	listener.AfterTest(testResult, device)

	ctrl.Finish()
}

func TestAfterTestRetunrnsFileAsExtra(t *testing.T) {
	targetTest := test.Test{}
	testResult := result.Result{Test: targetTest, Status: result.Failed}
	device := devices.Device{Serial: "abcd"}
	pathToLogcatFile := "path/to/logcat"
	ctrl := gomock.NewController(t)
	writer := mock_output.NewMockWriter(ctrl)
	logcatMock := mock_logcat.NewMockLogcat(ctrl)
	logcatMock.EXPECT().StopRecording(targetTest).Return(nil)
	logcatMock.EXPECT().SaveRecording(targetTest, writer).Return(pathToLogcatFile, nil)
	listener := testListener{
		writer: writer,
		logcat: map[devices.Device]Logcat{device: logcatMock},
	}
	extras := listener.AfterTest(testResult, device)

	if len(extras) != 1 {
		t.Error(fmt.Sprintf("Expected 1 extra but got %v", len(extras)))
	}
	expectedExtra := result.Extra{Name: "Logcat", Value: pathToLogcatFile, Type: result.File}
	if expectedExtra != extras[0] {
		t.Error(fmt.Sprintf("Expected extra '%v' but got '%v'", expectedExtra, extras[0]))
	}
}
