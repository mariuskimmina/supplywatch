package domain

import "net"

type ReceivedProduct struct {
    ProductName string
    Incoming    bool
    Amount      int
}

type SensorLog struct {
	SensorID   string `json:"sensor_id"`
	SensorType string `json:"sensor_type"`
	Message    string `json:"message"`
	Incoming   bool   `json:"incoming"`
	IP         net.IP `json:"ip"`
	Port       int    `json:"port"`
}

