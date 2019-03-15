package logging

import (
	"fmt"
	"github.com/android-test-runner/atr/devices"
)

type deviceLogger struct {
	device   devices.Device
	delegate Logger
}

func NewForDevice(device devices.Device) Logger {
	return deviceLogger{
		device:   device,
		delegate: NewLogger(),
	}
}

func (logger deviceLogger) Debug(message string) {
	logger.delegate.Debug(logger.addDevice(message))
}

func (logger deviceLogger) Info(message string) {
	logger.delegate.Info(logger.addDevice(message))
}

func (logger deviceLogger) Warn(message string) {
	logger.delegate.Warn(logger.addDevice(message))
}

func (logger deviceLogger) Error(message string, err error) {
	logger.delegate.Error(logger.addDevice(message), err)
}

func (logger deviceLogger) addDevice(message string) string {
	return fmt.Sprintf("%v: %v", logger.device.Serial, message)
}
