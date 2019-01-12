package logging

var globalLogLevel = Info

func SetLogLevel(level LogLevel) {
	globalLogLevel = level
}

type LogLevel int

const (
	Debug LogLevel = iota
	Info  LogLevel = iota
	Warn  LogLevel = iota
	Error LogLevel = iota
)
