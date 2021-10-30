package sensorwarehouse

import (
	"net"
	"time"

	"github.com/mariuskimmina/supplywatch/pkg/logger"
)

type Sensor struct {
}

func NewSensor() *Sensor {
    return &Sensor{}
}

func (s *Sensor) Start() {
    logger.Debug("Starting Sensor")
    products := []string{
        "Mehl",
        "Backpulver",
    }
    conn, err := net.Dial("udp", "supplywatch_warehouse_1:4444")
    if err != nil {
        logger.Error("Failed to dial")
    }
    for {
        logger.Info("Sending Mehl")
        conn.Write([]byte(products[0]))
        time.Sleep(5 * time.Second)
    }
}
