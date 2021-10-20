package sensorwarehouse

import (
    "github.com/mariuskimmina/supplywatch/pkg/logger"
)

type Sensor struct {

}

func NewSensor() *Sensor {
    return &Sensor{}
}

func (s *Sensor) Start() {
    logger.Debug("Starting Sensor")
}
