package domain

import (
	"context"

	//"github.com/google/uuid"
	//"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/mariuskimmina/supplywatch/internal/pb"
	"gorm.io/gorm"
)

type InOutProduct struct {
	ProductName string
	Incoming    bool
	Amount      int
	Reason      string
}

type Producttype struct {
	gorm.Model
	ID             uuid.UUID
	Name           string
	Quantity       int
	lastReceived   string
	lastDispatched string
}

type Product interface {
	Increment()
	Decrement()
	//gorm.Model
	//ID             uuid.UUID
	//Name           string
	//Quantity       int
	//lastReceived   string
	//lastDispatched string
}

type Products []*Product

type SensorLog struct {
	SensorID   string `json:"sensor_id"`
	SensorType string `json:"sensor_type"`
	Message    string `json:"message"`
	Incoming   bool   `json:"incoming"`
}

type ShippingLog struct {
	ShippingProductID   string `json:"shipping_product_id"`
	ShippingProductName string `json:"shipping_product_name"`
}

type Warehouse interface {
	AllProducts() (Products, error)
	ProductByID() error
	AllSensorLogs() ([]byte, error)
	GetAllProducts(ctx context.Context, req *pb.GetAllProductsRequest) (response *pb.GetAllProductsResponse, err error)
	ReceivProducts(ctx context.Context, req *pb.ReceivProductsRequest) (*pb.ReceivProductsResponse, error)
}
