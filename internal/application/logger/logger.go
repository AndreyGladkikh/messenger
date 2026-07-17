package logger

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// type Log struct {
// 	Message string
// 	Code    string
// }
