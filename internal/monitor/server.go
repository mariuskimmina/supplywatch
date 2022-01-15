package monitor

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/mariuskimmina/supplywatch/pkg/config"
)

type monitor struct {
	logger Logger
    config config.SupplywatchConfig
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


func NewMonitor(logger Logger, config config.SupplywatchConfig) *monitor {
	return &monitor{
		logger: logger,
        config: config,
	}
}

func (s *monitor) RunAndServe() {
    time.Sleep(30 * time.Second)
	var wg sync.WaitGroup
	wg.Add(1)
	// RabbitMQ
	go func() {
		s.SetupMessageQueue(s.config.NumOfWarehouses)
		wg.Done()
	}()
    http.HandleFunc("/hello", hello)
    http.ListenAndServe(":9000", nil)
    wg.Wait()
}

func hello(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "hello\n")
}
