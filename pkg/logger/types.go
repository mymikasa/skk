package logger

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type LoggerV1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
}

type Field struct {
	Key string
	Val any
}

type LoggerV2 interface {
	//  args 必须是偶数，并且是以 key1,value1,key2,value2 的形式传递
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
