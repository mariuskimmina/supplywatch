package sensorwarehouse

import (
	"encoding/json"
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/mariuskimmina/supplywatch/pkg/log"
)

type Sensor struct {
	logger *log.Logger
}

type Message struct {
    SensorID    uuid.UUID `json:"sensor_id"`
	SensorType  string `json:"sensor_type"`
	MessageBody string `json:"message"`
}

func NewSensor(logger *log.Logger) *Sensor {
	return &Sensor{
		logger: logger,
	}
}

var(
	products = []string{
		"Mehl",
		"Backpulver",
		"Wasser",
		"Zucker",
    }
    SensorType = []string{
        "BarcodeReader",
        "RFID-Reader",
    }
)

func (s *Sensor) Start() {
    sensorType := SensorType[rand.Intn(2)]
    sensorID, err := uuid.NewUUID()
	if err != nil {
		s.logger.Fatal("Failed to create ID for sensor")
	}
	conn, err := net.Dial("udp", "warehouse:4444")
	if err != nil {
		s.logger.Error("Failed to dial the warehouse")
	}

	for {
		product := products[rand.Intn(4)]
		message := Message{
            SensorID: sensorID,
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
