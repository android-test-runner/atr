package screen_recorder

import (
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/mock_output"
	"github.com/ybonjour/atr/mock_screen_recorder"
	"github.com/ybonjour/atr/result"
	"github.com/ybonjour/atr/test"
	"testing"
)

func TestBeforeTestStartsScreenRecording(t *testing.T) {
	targetTest := test.Test{}
	device := devices.Device{Serial: "abcd"}
	ctrl := gomock.NewController(t)
	screenRecorderMock := mock_screen_recorder.NewMockScreenRecorder(ctrl)
	screenRecorderMock.EXPECT().StartRecording(targetTest).Return(nil)

	listener := testListener{
		screenRecorder: map[devices.Device]ScreenRecorder{device: screenRecorderMock},
	}
	listener.BeforeTest(targetTest, device)

	ctrl.Finish()
}

func TestAfterTestStopsScreenRecordingAndSavesForFailedResult(t *testing.T) {
	targetTest := test.Test{}
	testResult := result.Result{Test: targetTest, Status: result.Failed}
	device := devices.Device{Serial: "abcd"}
	ctrl := gomock.NewController(t)
	writer := mock_output.NewMockWriter(ctrl)
	screenRecorderMock := mock_screen_recorder.NewMockScreenRecorder(ctrl)
	screenRecorderMock.EXPECT().StopRecording(targetTest).Return(nil)
	screenRecorderMock.EXPECT().SaveResult(targetTest, writer).Return(nil)
	screenRecorderMock.EXPECT().RemoveRecording(targetTest)
	listener := testListener{
		writer:         writer,
		screenRecorder: map[devices.Device]ScreenRecorder{device: screenRecorderMock},
	}
	listener.AfterTest(testResult, device)

	ctrl.Finish()
}

func TestAfterTestStopsScreenRecordingForPassedResult(t *testing.T) {
	targetTest := test.Test{}
	testResult := result.Result{Test: targetTest, Status: result.Passed}
	device := devices.Device{Serial: "abcd"}
	ctrl := gomock.NewController(t)
	writer := mock_output.NewMockWriter(ctrl)
	screenRecorderMock := mock_screen_recorder.NewMockScreenRecorder(ctrl)
	screenRecorderMock.EXPECT().StopRecording(targetTest).Return(nil)
	screenRecorderMock.EXPECT().SaveResult(targetTest, writer).Return(nil).Times(0)
	screenRecorderMock.EXPECT().RemoveRecording(targetTest)
	listener := testListener{
		writer:         writer,
		screenRecorder: map[devices.Device]ScreenRecorder{device: screenRecorderMock},
	}
	listener.AfterTest(testResult, device)

	ctrl.Finish()
}
