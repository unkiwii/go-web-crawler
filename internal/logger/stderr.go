package logger

import (
	"fmt"
	"os"
	"time"
)

func NewStderr() Logger {
	return stderrlogger{}
}

type stderrlogger struct{}

func (l stderrlogger) Errorf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR %s :: %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, v...))
}

func (l stderrlogger) Infof(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, " INFO %s :: %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, v...))
}
