package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReceivedProduct struct {
    ProductName string
    Incoming    bool
    Amount      int
}

type Product struct {
	gorm.Model
	ID             uuid.UUID
	Name           string
	Quantity       int
	lastReceived   string
	lastDispatched string
}

type Products []*Product

type SensorLog struct {
	SensorID   string `json:"sensor_id"`
	SensorType string `json:"sensor_type"`
	Message    string `json:"message"`
	Incoming   bool   `json:"incoming"`
}

type Warehouse interface {
    AllProducts() (Products, error)
    ProductByID() error
    AllSensorLogs() ([]byte, error)
}
