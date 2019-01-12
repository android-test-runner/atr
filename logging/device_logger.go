package logging

import (
	"fmt"
	"github.com/ybonjour/atr/devices"
)

type deviceLogger struct {
	Device devices.Device
}

func NewForDevice(device devices.Device) Logger {
	return deviceLogger{
		Device: device,
	}
}

func (logger deviceLogger) Debug(message string) {
	logger.print(message)
}

func (logger deviceLogger) Info(message string) {
	logger.print(message)
}

func (logger deviceLogger) Warn(message string) {
	logger.print(message)
}

func (logger deviceLogger) Error(message string, err error) {
	logger.print(fmt.Sprintf("%v: %v", message, err))
}

func (logger deviceLogger) print(message string) {
	fmt.Printf("%v: %v\n", logger.Device.Serial, message)
}
