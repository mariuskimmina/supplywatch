package main

import (
	"github.com/mariuskimmina/supplywatch/internal/sensor-warehouse"
	"github.com/mariuskimmina/supplywatch/pkg/log"
)

func main() {
    logger := log.NewLogger()
	sensor := sensorwarehouse.NewSensor(logger)

	sensor.Start()
}
