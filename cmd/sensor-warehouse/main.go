package main

import (
	"github.com/mariuskimmina/supplywatch/internal/sensor-warehouse"
	"github.com/mariuskimmina/supplywatch/pkg/config"
	"github.com/mariuskimmina/supplywatch/pkg/log"
)

func main() {
	logger := log.NewLogger()
	config, err := config.LoadConfig(".")
	if err != nil {
		logger.Fatalf("Failed to load warehouse configuration: %v", err)
	}
	sensor := sensorwarehouse.NewSensor(logger, &config)

	sensor.Start()
}
