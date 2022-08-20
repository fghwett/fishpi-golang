package logger

type Logger interface {
	Logf(format string, a ...interface{})
	Log(msg string)
}
