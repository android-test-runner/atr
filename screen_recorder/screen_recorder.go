package screen_recorder

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/test"
	"math/rand"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type ScreenRecorder interface {
	StartRecording(test test.Test) error
	StopRecording(test test.Test) error
	SaveResult(test test.Test, writer output.Writer) (string, error)
	RemoveRecording(test test.Test) error
}

type screenRecorderImpl struct {
	Device   devices.Device
	Adb      adb.Adb
	Test     test.Test
	pid      int
	filePath string
}

func New(device devices.Device) ScreenRecorder {
	return &screenRecorderImpl{
		Device: device,
		Adb:    adb.New(),
	}
}

func (screenRecorder *screenRecorderImpl) StartRecording(test test.Test) error {
	errTest := screenRecorder.testScreenRecording(screenRecorder.Device)
	if errTest != nil {
		return errTest
	}

	filepPath := fmt.Sprintf("/sdcard/%v.mp4", test.FullName())
	device := screenRecorder.Device
	pid, err := screenRecorder.recordScreenInBackground(device, filepPath)
	if err != nil {
		return err
	}
	screenRecorder.pid = pid
	screenRecorder.filePath = filepPath
	screenRecorder.Test = test

	return nil
}

func (screenRecorder *screenRecorderImpl) testScreenRecording(device devices.Device) error {
	testDurationInSeconds := 1
	randomFileName := fmt.Sprintf("/sdcard/screen-record-test-%v.mp4", randomText())
	var result error
	testResult := screenRecorder.Adb.RecordScreen(device.Serial, randomFileName, device.ScreenDimension.ToString(), testDurationInSeconds)
	if testResult.Error != nil {
		result = errors.New(fmt.Sprintf("can not record video on device '%v': %v\n", device.Serial, testResult.Error) +
			"Some devices might not be able to record at their native display resolution. If you encounter problems with screen recording, try using a lower screen resolution. (https://developer.android.com/studio/command-line/adb#screenrecord)\n" +
			"You can control the resolution at which atr records the screen in the device definition that you pass in the --device flag. Just pass --<deviceSerial>@<width>x<height> (e.g. --device emulator-5554@720x1440)")
	}

	resultRemove := screenRecorder.Adb.RemoveFile(device.Serial, randomFileName)
	if resultRemove.Error != nil {
		fmt.Printf("Could not remove screen recording testfile '%v' on device '%v': %v\n", randomFileName, device.Serial, resultRemove.Error)
	}

	return result
}

func (screenRecorder *screenRecorderImpl) recordScreenInBackground(device devices.Device, filepPath string) (int, error) {
	return screenRecorder.Adb.RecordScreenInBackground(device.Serial, filepPath, device.ScreenDimension.ToString())
}

func (screenRecorder screenRecorderImpl) StopRecording(test test.Test) error {
	if screenRecorder.Test != test {
		return errors.New(fmt.Sprintf("never started recording for test '%v'", test))
	}

	return interruptProcess(screenRecorder.pid)
}

func (screenRecorder *screenRecorderImpl) RemoveRecording(test test.Test) error {
	if screenRecorder.Test != test {
		return errors.New(fmt.Sprintf("never started recording for test '%v'", test))
	}
	result := screenRecorder.Adb.RemoveFile(screenRecorder.Device.Serial, screenRecorder.filePath)
	return result.Error
}

func (screenRecorder *screenRecorderImpl) SaveResult(test test.Test, writer output.Writer) (string, error) {
	if screenRecorder.Test != test {
		return "", errors.New(fmt.Sprintf("never started recording for test '%v'", test))
	}

	deviceDirectory, directoryErr := writer.MakeDeviceDirectory(screenRecorder.Device)
	if directoryErr != nil {
		return "", directoryErr
	}

	localFile := filepath.Join(writer.ToAbsolute(deviceDirectory), fmt.Sprintf("%v.mp4", test.FullName()))

	// Give screen recorder some time to properly write the video file
	time.Sleep(2 * time.Second)

	r := screenRecorder.Adb.PullFile(screenRecorder.Device.Serial, screenRecorder.filePath, localFile)

	return deviceDirectory, r.Error
}

func interruptProcess(pid int) error {
	process, findError := os.FindProcess(pid)
	if findError != nil {
		return findError
	}

	interruptError := process.Signal(syscall.SIGINT)
	if interruptError != nil {
		return interruptError
	}

	return nil
}

func randomText() string {
	randomId, err := uuid.NewV4()
	if err != nil {
		return fmt.Sprintf("static-%v", rand.Intn(100))
	}

	return randomId.String()
}
