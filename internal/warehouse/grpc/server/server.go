package gserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/mariuskimmina/supplywatch/internal/domain"
	"github.com/mariuskimmina/supplywatch/internal/pb"
)

const (
	receivLogFileDir = "/var/supplywatch/grpc/server/"
	receivLogFile    = "receivLog"
	receivLog        = receivLogFileDir + receivLogFile
)

type gserver struct {
	pb.UnimplementedProductServiceServer
	Conn             net.Listener
	InOutProductChan chan *domain.InOutProduct
}

func New(inOutProductChan chan *domain.InOutProduct) (*gserver, error) {
	err := os.MkdirAll(receivLogFileDir, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error creating log file directory: %v", err)
	}
	var tcpConn net.Listener
	tcpPort := os.Getenv("SW_GRPCPORT")
	listenIP := os.Getenv("SW_LISTEN_IP")
	tcpListenIP := listenIP + ":" + tcpPort
	tcpConn, err = net.Listen("tcp", tcpListenIP)
	if err != nil {
		return nil, err
	}
	return &gserver{
		Conn:             tcpConn,
		InOutProductChan: inOutProductChan,
	}, nil
}

func (s *gserver) GetAllProducts(ctx context.Context, req *pb.GetAllProductsRequest) (*pb.GetAllProductsResponse, error) {
	//w.logger.Info("Received GRPC Request, GetAllProducts")
	//var allProducts []Product
	//w.DB.Find(&allProducts)

	//sendProducts := []*pb.Product{}
	//for _, product := range allProducts {
	//sendProduct := pb.Product{
	//Name: product.Name,
	//Id:   product.ID.String(),
	//}
	//sendProducts = append(sendProducts, &sendProduct)
	//}
	//
	//allProductsJson, err := json.MarshalIndent(sendProducts, "", "  ")
	//if err != nil {
	//return nil, err
	//}
	//err = ioutil.WriteFile(receivLog, allProductsJson, 0644)
	//if err != nil {
	//return nil, err
	//}
	//
	//return &pb.GetAllProductsResponse{
	//Products: sendProducts,
	//}, nil
	return nil, nil
}

func (s *gserver) ReceivProducts(ctx context.Context, req *pb.ReceivProductsRequest) (*pb.ReceivProductsResponse, error) {
	f, err := os.OpenFile(receivLogFileDir+receivLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error opening log file: %v", err)
	}
	defer f.Close()
	// first we log the Request
	// this way we can compare the shippingReceivLog with the shippingSendLog of the other warehouse
	// if they don't match something went wrong
	logentry := &domain.ShippingLog{
		ShippingProductID:   req.Product.Id,
		ShippingProductName: req.Product.Name,
	}
	if err != nil {
		return nil, err
	}
	logjson, err := json.Marshal(logentry)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling log data: %v", err)
	}
	f.WriteString(string(logjson) + ",\n")

	//if err != nil {
	//return nil, err
	//}
	//w.logger.Info("Received GRPC Request, ReceivProducts")
	//product := &Product{}
	//w.DB.First(product, "name = ?", req.Product.Name)
	//oldQuantity := product.Quantity
	//w.DB.Model(&Product{}).Where("name = ?", req.Product.Name).Update("quantity", req.Amount+int32(oldQuantity))
	fmt.Println("Received GRPC Request, ReceivProducts")
	p := &domain.InOutProduct{
		ProductName: req.Product.Name,
		Incoming:    true,
		Amount:      1,
		Reason:      "Another Warehouse",
	}
	// inform the storage channel about the product that has been shipped to us
	s.InOutProductChan <- p

	return &pb.ReceivProductsResponse{
		Success: true,
	}, nil
}
