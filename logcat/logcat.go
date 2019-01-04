package logcat

import (
	"errors"
	"fmt"
	"github.com/ybonjour/atr/adb"
	"github.com/ybonjour/atr/devices"
	"github.com/ybonjour/atr/files"
	"github.com/ybonjour/atr/output"
	"github.com/ybonjour/atr/test"
)

type Logcat interface {
	StartRecording(test test.Test) error
	StopRecording(test test.Test) error
	SaveRecording(test test.Test, writer output.Writer) error
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
	return logcat.Adb.ClearLogcat(logcat.Device.Serial)
}

func (logcat *logcatImpl) StopRecording(test test.Test) error {
	if logcat.Test != test {
		return errors.New(fmt.Sprintf("never started recording for test '%v'", test))
	}

	out, err := logcat.Adb.GetLogcat(logcat.Device.Serial)
	logcat.Output = out
	return err
}

func (logcat *logcatImpl) SaveRecording(test test.Test, writer output.Writer) error {
	if logcat.Test != test {
		return errors.New(fmt.Sprintf("never started recording for test '%v'", test))
	}

	f := files.File{
		Name:    fmt.Sprintf("%v.log", test.FullName()),
		Content: logcat.Output,
	}

	return writer.WriteFile(f, logcat.Device)
}
