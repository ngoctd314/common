package ghttp

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// LoggerFunc is a bridge between Logger and any third party logger
type LoggerFunc func(msg string, args ...any)

func (f LoggerFunc) Info(msg string, args ...interface{})  { f(msg, args...) }
func (f LoggerFunc) Warn(msg string, args ...interface{})  { f(msg, args...) }
func (f LoggerFunc) Error(msg string, args ...interface{}) { f(msg, args...) }
