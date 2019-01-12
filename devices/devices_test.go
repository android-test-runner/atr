package devices

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/ybonjour/atr/mock_adb"
	"testing"
)

func TestMultipleConnectedDevices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	deviceSerial1 := "abcd"
	deviceSerial2 := "efgh"
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().ConnectedDevices().Return([]string{deviceSerial1, deviceSerial2}, nil)
	devices := devicesImpl{
		adb: adbMock,
	}

	connectedDevices, err := devices.ConnectedDevices([]string{deviceSerial1, deviceSerial2})

	if err != nil {
		t.Error(fmt.Sprintf("Did not expect an error but got '%v'", err))
	}
	expected := []Device{{Serial: deviceSerial1}, {Serial: deviceSerial2}}
	if !AreEqual(expected, connectedDevices) {
		t.Error(fmt.Sprintf("Expected devices '%v' but got '%v'.", expected, connectedDevices))
	}
}

func TestIgnoresDeviceIfItIsNotConnected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	deviceSerial := "abcd"

	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().ConnectedDevices().Return([]string{}, nil)
	devices := devicesImpl{
		adb: adbMock,
	}

	connectedDevices, err := devices.ConnectedDevices([]string{deviceSerial})

	if err != nil {
		t.Error(fmt.Sprintf("Did not expect an error but got '%v'", err))
	}
	if len(connectedDevices) != 0 {
		t.Error(fmt.Sprintf("Expected no devices but got '%v'.", connectedDevices))
	}
}

func TestIgnoresDeviceIfItIsNotIncluded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	deviceSerial := "abcd"
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().ConnectedDevices().Return([]string{deviceSerial}, nil)
	devices := devicesImpl{
		adb: adbMock,
	}

	connectedDevices, err := devices.ConnectedDevices([]string{})

	if err != nil {
		t.Error(fmt.Sprintf("Did not expect an error but got '%v'", err))
	}
	if len(connectedDevices) != 0 {
		t.Error(fmt.Sprintf("Expected no devices but got '%v'.", connectedDevices))
	}
}

func TestReturnsErrorIfConnectedDevicesFails(t *testing.T) {
	expectedError := errors.New("can not get connected devices")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().ConnectedDevices().Return(nil, expectedError)
	devices := devicesImpl{
		adb: adbMock,
	}

	_, err := devices.ConnectedDevices([]string{})

	if expectedError != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'.", expectedError, err))
	}
}

func TestParsesScreenDimension(t *testing.T) {
	dimension, err := ParseToScreenDimension("1024x768")

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}

	if dimension.Width != 1024 {
		t.Error(fmt.Sprintf("Expected screen dimension with width 1024 but got %v", dimension.Width))
	}

	if dimension.Height != 768 {
		t.Error(fmt.Sprintf("Expected screen dimension with height 768 but got %v", dimension.Height))
	}
}

func TestParseScreenDimensionReturnsErrorIfDimensionIsInvalid(t *testing.T) {
	_, err := ParseToScreenDimension("notaresultion")

	if err == nil {
		t.Error("Expected no error because no resultion provided, but did not get one.")
	}
}

func TestFormatsScreenDimension(t *testing.T) {
	dimension := ScreenDimension{Width: 1024, Height: 768}

	result := dimension.ToString()

	if result != "1024x768" {
		t.Error(fmt.Sprintf("Expected format '1024x768' but got '%v'", result))
	}
}

func TestParsesDevicesWithProvidedScreenDimension(t *testing.T) {
	deviceDefinition := "abcd@1024x768"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().GetScreenDimensions(gomock.Any()).Times(0)
	adbMock.EXPECT().ConnectedDevices().Return([]string{"abcd"}, nil)
	devices := devicesImpl{adb: adbMock}

	parsedDevices, err := devices.ParseConnectedDevices([]string{deviceDefinition})

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}
	expectedDevices := []Device{{Serial: "abcd", ScreenDimension: ScreenDimension{Width: 1024, Height: 768}}}
	if !AreEqual(expectedDevices, parsedDevices) {
		t.Error(fmt.Sprintf("Expected devices '%v' but got '%v'", expectedDevices, parsedDevices))
	}
}

func TestParsesDevicesWithNoScreenDimension(t *testing.T) {
	deviceDefinition := "abcd"
	screenWidth := 1024
	screenHeight := 768
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().GetScreenDimensions(gomock.Any()).Return(screenWidth, screenHeight, nil)
	adbMock.EXPECT().ConnectedDevices().Return([]string{"abcd"}, nil)
	devices := devicesImpl{adb: adbMock}

	parsedDevices, err := devices.ParseConnectedDevices([]string{deviceDefinition})

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}
	expectedDevices := []Device{{Serial: "abcd", ScreenDimension: ScreenDimension{Width: 1024, Height: 768}}}
	if !AreEqual(expectedDevices, parsedDevices) {
		t.Error(fmt.Sprintf("Expected devices '%v' but got '%v'", expectedDevices, parsedDevices))
	}
}

func TestParsesDevicesIgnoreDevicesWithSpecifiedScreenDimensionThatAreNotConnected(t *testing.T) {
	deviceDefinition := "abcd@1024x768"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().ConnectedDevices().Return([]string{"efgh"}, nil)
	devices := devicesImpl{adb: adbMock}

	parsedDevices, err := devices.ParseConnectedDevices([]string{deviceDefinition})

	verifyParsedDevicesEmpty(parsedDevices, err, t)
}

func TestParsesDevicesIgnoreDevicesWithoutSpecifiedScreenDimensionThatAreNotConnected(t *testing.T) {
	deviceDefinition := "abcd"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().ConnectedDevices().Return([]string{"efgh"}, nil)
	adbMock.EXPECT().GetScreenDimensions("abcd").Return(1024, 768, nil)
	devices := devicesImpl{adb: adbMock}

	parsedDevices, err := devices.ParseConnectedDevices([]string{deviceDefinition})

	verifyParsedDevicesEmpty(parsedDevices, err, t)
}

func TestParsesDevicesIgnoresUnparsableDevices(t *testing.T) {
	deviceDefinition := "abcd@1072x768@unknown"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().GetScreenDimensions(gomock.Any()).Times(0)
	adbMock.EXPECT().ConnectedDevices().Return([]string{"abcd"}, nil)
	devices := devicesImpl{adb: adbMock}

	parsedDevices, err := devices.ParseConnectedDevices([]string{deviceDefinition})

	verifyParsedDevicesEmpty(parsedDevices, err, t)
}

func TestParsesDevicesIgnoresDevicesWithUnparsableScreenDimension(t *testing.T) {
	deviceDefinition := "abcd@unknown"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().GetScreenDimensions(gomock.Any()).Times(0)
	adbMock.EXPECT().ConnectedDevices().Return([]string{"abcd"}, nil)
	devices := devicesImpl{adb: adbMock}

	parsedDevices, err := devices.ParseConnectedDevices([]string{deviceDefinition})

	verifyParsedDevicesEmpty(parsedDevices, err, t)
}

func TestParsesDevicesIgnoresDevicesWhereDefaultScreenDimensionCanNotBeRetrieved(t *testing.T) {
	deviceDefinition := "abcd"
	err := errors.New("can not get dimension")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().GetScreenDimensions(gomock.Any()).Return(0, 0, err)
	adbMock.EXPECT().ConnectedDevices().Return([]string{"abcd"}, nil)
	devices := devicesImpl{adb: adbMock}

	parsedDevices, err := devices.ParseConnectedDevices([]string{deviceDefinition})

	verifyParsedDevicesEmpty(parsedDevices, err, t)
}

func TestParseDevicesReturnsErrorIfConnectedDevicesCanNotBeRetrieved(t *testing.T) {
	deviceDefinition := "abcd"
	expectedErr := errors.New("can not get connected devices")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().ConnectedDevices().Return(nil, expectedErr)
	devices := devicesImpl{adb: adbMock}

	_, err := devices.ParseConnectedDevices([]string{deviceDefinition})

	if expectedErr != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'", expectedErr, err))
	}
}

func TestAllConnectedDevices(t *testing.T) {
	connectedDeviceSerial := "abcd"
	screenDimension := ScreenDimension{Width: 1024, Height: 768}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().ConnectedDevices().Return([]string{connectedDeviceSerial}, nil)
	adbMock.EXPECT().GetScreenDimensions(connectedDeviceSerial).Return(screenDimension.Width, screenDimension.Height, nil)
	devices := devicesImpl{adb: adbMock}

	connectedDevices, err := devices.AllConnectedDevices()

	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}

	expectedDevices := []Device{{Serial: connectedDeviceSerial, ScreenDimension: screenDimension}}
	if !AreEqual(expectedDevices, connectedDevices) {
		t.Error(fmt.Sprintf("Expected connected devices to be '%v' but it was '%v'", expectedDevices, connectedDevices))
	}

}

func TestAllConnectedDevicesReturnsErrorOnFailure(t *testing.T) {
	expectedError := errors.New("could not get connected devices")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adbMock := mock_adb.NewMockAdb(ctrl)
	adbMock.EXPECT().ConnectedDevices().Return(nil, expectedError)
	devices := devicesImpl{adb: adbMock}

	_, err := devices.AllConnectedDevices()

	if expectedError != err {
		t.Error(fmt.Sprintf("Expected error '%v' but got '%v'", expectedError, err))
	}
}

func verifyParsedDevicesEmpty(parsedDevices []Device, err error, t *testing.T) {
	if err != nil {
		t.Error(fmt.Sprintf("Expected no error but got '%v'", err))
	}
	if len(parsedDevices) != 0 {
		t.Error(fmt.Sprintf("Expected no devices but got '%v'", parsedDevices))
	}
}

func AreEqual(slice1, slice2 []Device) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}
