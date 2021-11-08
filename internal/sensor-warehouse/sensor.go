package sensorwarehouse

import (
	"encoding/json"
	"math/rand"
	"net"
	"time"

	"github.com/mariuskimmina/supplywatch/pkg/log"
)

type Sensor struct {
	logger *log.Logger
}

type Message struct {
	SensorType  string
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
		"Wasser",
		"Zucker",
	}
	conn, err := net.Dial("udp", "supplywatch_warehouse_1:4444")
	if err != nil {
		s.logger.Error("Failed to dial")
	}

	for {
		message := Message{
			SensorType:  sensorType,
			MessageBody: products[rand.Intn(4)],
		}
		jsonMessage, err := json.Marshal(message)

		if err != nil {
			s.logger.Error("Failed to convert message to json")
		}
		s.logger.Info("Sending Mehl")
		conn.Write([]byte(jsonMessage))
		time.Sleep(5 * time.Second)
	}
}
