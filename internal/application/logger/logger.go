package logger

type Logger interface {
	Log(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type Log struct {
	Message string
	Code string
}