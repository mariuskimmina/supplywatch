package main

import (
	"github.com/mariuskimmina/supplywatch/internal/sensor"
	"github.com/mariuskimmina/supplywatch/pkg/config"
	"github.com/mariuskimmina/supplywatch/pkg/log"
)

func main() {
	logger := log.NewLogger()
	config, err := config.LoadSensorConfig("./configurations")
	if err != nil {
		logger.Error(err)
		logger.Fatal("Failed to load warehouse configuration")
	}
	sensor := sensor.NewSensor(logger, &config)

	sensor.Start()
}
