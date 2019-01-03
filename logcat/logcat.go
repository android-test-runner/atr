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
	StopRecording(test test.Test, saveResult bool) error
}

type Factory interface {
	ForDevice(device devices.Device) Logcat
}

type factoryImpl struct {
	Writer output.Writer
}

func NewFactory(writer output.Writer) Factory {
	return factoryImpl{
		Writer: writer,
	}
}

type logcatImpl struct {
	Device devices.Device
	Adb    adb.Adb
	Writer output.Writer
	Test   test.Test
}

// One logcat instance per device to avoid problems with parallelism
func (factory factoryImpl) ForDevice(device devices.Device) Logcat {
	return &logcatImpl{
		Device: device,
		Adb:    adb.New(),
		Writer: factory.Writer,
	}
}

func (logcat *logcatImpl) StartRecording(test test.Test) error {
	logcat.Test = test
	return logcat.Adb.ClearLogcat(logcat.Device.Serial)
}

func (logcat *logcatImpl) StopRecording(test test.Test, saveResult bool) error {
	if logcat.Test != test {
		return errors.New(fmt.Sprintf("never started recording for test '%v'", test))
	}

	if !saveResult {
		return nil
	}

	logcatOutput, err := logcat.Adb.GetLogcat(logcat.Device.Serial)
	if err != nil {
		return err
	}

	f := files.File{
		Name:    fmt.Sprintf("%v.log", test.FullName()),
		Content: logcatOutput,
	}

	return logcat.Writer.WriteFile(f, logcat.Device)
}
