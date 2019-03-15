package logcat

import (
	"fmt"
	"github.com/android-test-runner/atr/mock_logcat"
	"github.com/android-test-runner/atr/mock_logging"
	"github.com/android-test-runner/atr/mock_output"
	"github.com/android-test-runner/atr/result"
	"github.com/android-test-runner/atr/test"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestBeforeTestStartsLogcatRecording(t *testing.T) {
	targetTest := test.Test{}
	ctrl := gomock.NewController(t)
	logcatMock := mock_logcat.NewMockLogcat(ctrl)
	logcatMock.EXPECT().StartRecording(targetTest).Return(nil)
	loggerMock := mock_logging.NewMockLogger(ctrl)
	allowLogging(loggerMock)
	listener := testListener{
		logcat: logcatMock,
		logger: loggerMock,
	}

	listener.BeforeTest(targetTest)

	ctrl.Finish()
}

func TestAfterTestStopsLogcatRecordingAndSavesForFailedResult(t *testing.T) {
	targetTest := test.Test{}
	testResult := result.Result{Test: targetTest, Status: result.Failed}
	ctrl := gomock.NewController(t)
	writer := mock_output.NewMockWriter(ctrl)
	logcatMock := mock_logcat.NewMockLogcat(ctrl)
	logcatMock.EXPECT().StopRecording(targetTest).Return(nil)
	logcatMock.EXPECT().SaveRecording(targetTest, writer).Return("", nil)
	loggerMock := mock_logging.NewMockLogger(ctrl)
	allowLogging(loggerMock)
	listener := testListener{
		writer: writer,
		logcat: logcatMock,
		logger: loggerMock,
	}

	listener.AfterTest(testResult)

	ctrl.Finish()
}

func TestAfterTestStopsLogcatRecordingForPassedResult(t *testing.T) {
	targetTest := test.Test{}
	testResult := result.Result{Test: targetTest, Status: result.Passed}
	ctrl := gomock.NewController(t)
	writer := mock_output.NewMockWriter(ctrl)
	logcatMock := mock_logcat.NewMockLogcat(ctrl)
	logcatMock.EXPECT().StopRecording(targetTest).Return(nil)
	logcatMock.EXPECT().SaveRecording(targetTest, writer).Return("", nil).Times(0)
	loggerMock := mock_logging.NewMockLogger(ctrl)
	allowLogging(loggerMock)
	listener := testListener{
		writer: writer,
		logcat: logcatMock,
		logger: loggerMock,
	}
	listener.AfterTest(testResult)

	ctrl.Finish()
}

func TestAfterTestRetunrnsFileAsExtra(t *testing.T) {
	targetTest := test.Test{}
	testResult := result.Result{Test: targetTest, Status: result.Failed}
	pathToLogcatFile := "path/to/logcat"
	ctrl := gomock.NewController(t)
	writer := mock_output.NewMockWriter(ctrl)
	logcatMock := mock_logcat.NewMockLogcat(ctrl)
	logcatMock.EXPECT().StopRecording(targetTest).Return(nil)
	logcatMock.EXPECT().SaveRecording(targetTest, writer).Return(pathToLogcatFile, nil)
	loggerMock := mock_logging.NewMockLogger(ctrl)
	allowLogging(loggerMock)
	listener := testListener{
		writer: writer,
		logcat: logcatMock,
		logger: loggerMock,
	}
	extras := listener.AfterTest(testResult)

	if len(extras) != 1 {
		t.Error(fmt.Sprintf("Expected 1 extra but got %v", len(extras)))
	}
	expectedExtra := result.Extra{Name: "Logcat", Value: pathToLogcatFile, Type: result.File}
	if expectedExtra != extras[0] {
		t.Error(fmt.Sprintf("Expected extra '%v' but got '%v'", expectedExtra, extras[0]))
	}
}

func allowLogging(loggerMock *mock_logging.MockLogger) {
	loggerMock.EXPECT().Debug(gomock.Any()).AnyTimes()
	loggerMock.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
}
