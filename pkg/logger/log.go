package logger

import (
        "fmt"
        "path"
        "runtime"

        "github.com/sirupsen/logrus"
)

type logger struct {

}

type Logger interface {

}

var (
    log *logrus.Logger
)

func init() {
    log = logrus.New()
    log.SetReportCaller(true)
    log.SetLevel(logrus.DebugLevel)
    log.Formatter = &logrus.TextFormatter{
        CallerPrettyfier: func(f *runtime.Frame) (string, string) {
            filename := path.Base(f.File)
            return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
        },
    }
}

// Debug ...
func Debug(format string, v ...interface{}) {
    log.Debugf(format, v...)
}

// Info ...
func Info(format string, v ...interface{}) {
    log.Infof(format, v...)
}

// Warn ...
func Warn(format string, v ...interface{}) {
    log.Warnf(format, v...)
}

// Error ...
func Error(format string, v ...interface{}) {
    log.Errorf(format, v...)
}
