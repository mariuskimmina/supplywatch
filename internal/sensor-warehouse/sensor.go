package sensorwarehouse

import (
	"net"
	"time"

	"github.com/mariuskimmina/supplywatch/pkg/log"
)

type Sensor struct {
	logger *log.Logger
}

func NewSensor(logger *log.Logger) *Sensor {
	return &Sensor{
		logger: logger,
	}
}

func (s *Sensor) Start() {
	products := []string{
		"Mehl",
		"Backpulver",
	}
	conn, err := net.Dial("udp", "supplywatch_warehouse_1:4444")
	if err != nil {
		s.logger.Error("Failed to dial")
	}
	for {
		s.logger.Info("Sending Mehl")
		conn.Write([]byte(products[0]))
		time.Sleep(5 * time.Second)
	}
}
