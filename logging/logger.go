package logging

import "fmt"

type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string, err error)
}

type loggerImpl struct{}

func (logger loggerImpl) Debug(message string) {
	logger.print(message)
}

func (logger loggerImpl) Info(message string) {
	logger.print(message)
}

func (logger loggerImpl) Warn(message string) {
	logger.print(message)
}

func (logger loggerImpl) Error(message string, err error) {
	logger.print(fmt.Sprintf("%v: %v", message, err))
}

func (logger loggerImpl) print(message string) {
	fmt.Printf("%v\n", message)
}
