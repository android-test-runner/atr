package screen_recorder

import (
	"errors"
	"fmt"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/test"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type ScreenRecorder interface {
	StartRecording(test test.Test) error
	StopRecording(test test.Test) error
	SaveResult(test test.Test, writer output.Writer) error
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
	filepPath := fmt.Sprintf("/sdcard/%v.mp4", test.FullName())
	pid, err := screenRecorder.Adb.RecordScreen(screenRecorder.Device.Serial, filepPath)
	if err != nil {
		return err
	}
	screenRecorder.pid = pid
	screenRecorder.filePath = filepPath
	screenRecorder.Test = test

	return nil
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

func (screenRecorder *screenRecorderImpl) SaveResult(test test.Test, writer output.Writer) error {
	if screenRecorder.Test != test {
		return errors.New(fmt.Sprintf("never started recording for test '%v'", test))
	}

	deviceDirectory, directoryErr := writer.GetDeviceDirectory(screenRecorder.Device)
	if directoryErr != nil {
		return directoryErr
	}

	localFile := filepath.Join(deviceDirectory, fmt.Sprintf("%v.mp4", test.FullName()))

	// Give screen recorder some time to properly write the video file
	time.Sleep(2 * time.Second)

	r := screenRecorder.Adb.PullFile(screenRecorder.Device.Serial, screenRecorder.filePath, localFile)

	return r.Error
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
