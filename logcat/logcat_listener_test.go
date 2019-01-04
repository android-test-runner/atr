package logcat

import (
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/mock_logcat"
	"github.com/ybonjour/atr/mock_output"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"testing"
)

func TestBeforeTestStartsLogcatRecording(t *testing.T) {
	targetTest := test.Test{}
	ctrl := gomock.NewController(t)
	logcatMock := mock_logcat.NewMockLogcat(ctrl)
	logcatMock.EXPECT().StartRecording(targetTest).Return(nil)

	listener := logcatListener{
		logcat: logcatMock,
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
	logcatMock.EXPECT().SaveRecording(targetTest, writer).Return(nil)
	listener := logcatListener{
		writer: writer,
		logcat: logcatMock,
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
	logcatMock.EXPECT().SaveRecording(targetTest, writer).Return(nil).Times(0)
	listener := logcatListener{
		writer: writer,
		logcat: logcatMock,
	}
	listener.AfterTest(testResult)

	ctrl.Finish()
}
