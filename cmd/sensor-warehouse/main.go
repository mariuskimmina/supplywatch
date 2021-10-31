package main

import (
	"github.com/mariuskimmina/supplywatch/internal/sensor-warehouse"
	"github.com/mariuskimmina/supplywatch/pkg/logger"
)

func main() {
	logger.Info("start")
	sensor := sensorwarehouse.NewSensor()

	sensor.Start()
}
