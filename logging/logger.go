package logging

import "fmt"

type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string, err error)
}

type loggerImpl struct{}

func NewLogger() Logger {
	return loggerImpl{}
}

func (logger loggerImpl) Debug(message string) {
	logger.log(message, Debug)
}

func (logger loggerImpl) Info(message string) {
	logger.log(message, Info)
}

func (logger loggerImpl) Warn(message string) {
	logger.log(message, Warn)
}

func (logger loggerImpl) Error(message string, err error) {
	logger.log(fmt.Sprintf("%v: %v", message, err), Error)
}

func (logger loggerImpl) log(message string, level LogLevel) {
	if level >= globalLogLevel {
		fmt.Printf("%v\n", message)
	}
}
