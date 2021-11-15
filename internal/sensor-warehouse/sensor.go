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
	SensorType  string `json:"sensor_type"`
	MessageBody string `json:"message"`
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
	conn, err := net.Dial("udp", "warehouse:4444")
	if err != nil {
		s.logger.Error("Failed to dial")
	}

	for {
		product := products[rand.Intn(4)]
		message := Message{
			SensorType:  sensorType,
			MessageBody: product,
		}
		jsonMessage, err := json.Marshal(message)

		if err != nil {
			s.logger.Error("Failed to convert message to json")
		}
		s.logger.Infof("Sending %s", product)
		conn.Write([]byte(jsonMessage))
		time.Sleep(5 * time.Second)
	}
}
