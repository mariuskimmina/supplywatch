package monitor

import (
	"fmt"
	"net/http"
)

type monitor struct {
	logger Logger
}

// Logger is a generic interface that can be implemented by any logging engine
// this allows for dependency injection which results in easier testing
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}


func NewMonitor(logger Logger) *monitor {
	return &monitor{
		logger: logger,
	}
}

func (s *monitor) RunAndServe() {
    http.HandleFunc("/hello", hello)
    http.ListenAndServe(":9000", nil)
}

func hello(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "hello\n")
}
