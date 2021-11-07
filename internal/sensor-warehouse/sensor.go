package sensorwarehouse

import (
	"net"
	"time"

	"github.com/mariuskimmina/supplywatch/pkg/log"
)

type Sensor struct {
	logger *log.Logger
}

type Message struct {
	SensorType 	string
	MessageBody string
}

func NewSensor(logger *log.Logger) *Sensor {
	return &Sensor{
		logger: logger,
	}
}

const (
	sensorType = "Warehouse-sensor"
)

func (s *Sensor) Start() {
	products := []string{
		"Mehl",
		"Backpulver",
	}
	conn, err := net.Dial("udp", "supplywatch_warehouse_1:4444")
	message := Message{
		SensorType: sensorType,
		MessageBody: products[0],
	}
	if err != nil {
		s.logger.Error("Failed to dial")
	}
	for {
		s.logger.Info("Sending Mehl")
		conn.Write([]byte(message)
		time.Sleep(5 * time.Second)
	}
}
