package logger

type LogLevel int

const (
	LogLevelError LogLevel = iota + 1
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
)

var toString = map[LogLevel]string{
	LogLevelError:   "ERROR",
	LogLevelWarning: "WARNING",
	LogLevelInfo:    "INFO",
	LogLevelDebug:   "DEBUG",
}

type Logger interface {
	Error(format string, v ...interface{})
	// Warning(format string, v ...interface{})
	// Info(format string, v ...interface{})
	// Debug(format string, v ...interface{})
}
