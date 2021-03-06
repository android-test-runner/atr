package logcat

import (
	"errors"
	"fmt"
	"github.com/android-test-runner/atr/adb"
	"github.com/android-test-runner/atr/devices"
	"github.com/android-test-runner/atr/files"
	"github.com/android-test-runner/atr/output"
	"github.com/android-test-runner/atr/test"
)

type Logcat interface {
	StartRecording(test test.Test) error
	StopRecording(test test.Test) error
	SaveRecording(test test.Test, writer output.Writer) (string, error)
}

type logcatImpl struct {
	Device devices.Device
	Adb    adb.Adb
	Test   test.Test
	Output string
}

func New(device devices.Device) Logcat {
	return &logcatImpl{
		Device: device,
		Adb:    adb.New(),
	}
}

func (logcat *logcatImpl) StartRecording(test test.Test) error {
	logcat.Test = test
	result := logcat.Adb.ClearLogcat(logcat.Device.Serial)
	return result.Error
}

func (logcat *logcatImpl) StopRecording(test test.Test) error {
	err := logcat.ensureTest(test)
	if err != nil {
		return err
	}

	out, err := logcat.Adb.GetLogcat(logcat.Device.Serial)
	logcat.Output = out
	return err
}

func (logcat *logcatImpl) SaveRecording(test test.Test, writer output.Writer) (string, error) {
	err := logcat.ensureTest(test)
	if err != nil {
		return "", err
	}

	f := files.File{
		Name:    fmt.Sprintf("%v.log", test.FullName()),
		Content: logcat.Output,
	}

	return writer.WriteFile(f, logcat.Device)
}

func (logcat *logcatImpl) ensureTest(test test.Test) error {
	if logcat.Test != test {
		return errors.New(fmt.Sprintf("never started recording for test %v on device %v", test.FullName(), logcat.Device.Serial))
	} else {
		return nil
	}
}
