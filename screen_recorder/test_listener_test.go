package screen_recorder

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/mock_logging"
	"github.com/ybonjour/atr/mock_output"
	"github.com/ybonjour/atr/mock_screen_recorder"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"testing"
)

func TestBeforeTestStartsScreenRecording(t *testing.T) {
	targetTest := test.Test{}
	ctrl := gomock.NewController(t)
	screenRecorderMock := mock_screen_recorder.NewMockScreenRecorder(ctrl)
	screenRecorderMock.EXPECT().StartRecording(targetTest).Return(nil)
	loggerMock := mock_logging.NewMockLogger(ctrl)
	allowLogging(loggerMock)
	listener := testListener{screenRecorder: screenRecorderMock, logger: loggerMock}

	listener.BeforeTest(targetTest)

	ctrl.Finish()
}

func TestAfterTestStopsScreenRecordingAndSavesForFailedResult(t *testing.T) {
	targetTest := test.Test{}
	testResult := result.Result{Test: targetTest, Status: result.Failed}
	ctrl := gomock.NewController(t)
	writer := mock_output.NewMockWriter(ctrl)
	screenRecorderMock := mock_screen_recorder.NewMockScreenRecorder(ctrl)
	screenRecorderMock.EXPECT().StopRecording(targetTest).Return(nil)
	screenRecorderMock.EXPECT().SaveResult(targetTest, writer).Return("", nil)
	screenRecorderMock.EXPECT().RemoveRecording(targetTest)
	loggerMock := mock_logging.NewMockLogger(ctrl)
	allowLogging(loggerMock)
	listener := testListener{
		writer:         writer,
		screenRecorder: screenRecorderMock,
		logger:         loggerMock,
	}
	listener.AfterTest(testResult)

	ctrl.Finish()
}

func TestAfterTestStopsScreenRecordingForPassedResult(t *testing.T) {
	targetTest := test.Test{}
	testResult := result.Result{Test: targetTest, Status: result.Passed}
	ctrl := gomock.NewController(t)
	writer := mock_output.NewMockWriter(ctrl)
	screenRecorderMock := mock_screen_recorder.NewMockScreenRecorder(ctrl)
	screenRecorderMock.EXPECT().StopRecording(targetTest).Return(nil)
	screenRecorderMock.EXPECT().SaveResult(targetTest, writer).Return("", nil).Times(0)
	screenRecorderMock.EXPECT().RemoveRecording(targetTest)
	loggerMock := mock_logging.NewMockLogger(ctrl)
	allowLogging(loggerMock)
	listener := testListener{
		writer:         writer,
		screenRecorder: screenRecorderMock,
		logger:         loggerMock,
	}
	listener.AfterTest(testResult)

	ctrl.Finish()
}

func TestAfterTestReturnsScreenRecordingfileExtra(t *testing.T) {
	targetTest := test.Test{Class: "TestClass", Method: "testMethod"}
	testResult := result.Result{Test: targetTest, Status: result.Failed}
	filePath := "filepath"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	writer := mock_output.NewMockWriter(ctrl)
	screenRecorderMock := mock_screen_recorder.NewMockScreenRecorder(ctrl)
	screenRecorderMock.EXPECT().StopRecording(targetTest).Return(nil)
	screenRecorderMock.EXPECT().SaveResult(targetTest, writer).Return(filePath, nil)
	screenRecorderMock.EXPECT().RemoveRecording(targetTest)
	loggerMock := mock_logging.NewMockLogger(ctrl)
	allowLogging(loggerMock)
	listener := testListener{
		writer:         writer,
		screenRecorder: screenRecorderMock,
		logger:         loggerMock,
	}
	extras := listener.AfterTest(testResult)

	if len(extras) != 1 {
		t.Error(fmt.Sprintf("Expected 1 extra but got %v", len(extras)))
		return
	}
	expectedExtras := result.Extra{Name: "Screen Recording", Value: filePath, Type: result.File}
	if expectedExtras != extras[0] {
		t.Error(fmt.Sprintf("Expected extra '%v' but got '%v'", expectedExtras, extras[0]))
	}
}

func allowLogging(loggerMock *mock_logging.MockLogger) {
	loggerMock.EXPECT().Debug(gomock.Any()).AnyTimes()
	loggerMock.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
}
