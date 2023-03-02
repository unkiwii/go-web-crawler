package logger

type Logger interface {
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
}

func Discard() Logger {
	return defaultDiscardLogger
}

var defaultDiscardLogger = discard{}

type discard struct{}

func (discard) Errorf(format string, v ...interface{}) {}
func (discard) Infof(format string, v ...interface{})  {}
