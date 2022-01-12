package gserver

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/mariuskimmina/supplywatch/internal/domain"
	"github.com/mariuskimmina/supplywatch/internal/pb"
)

const (
    receivLogFileDir = "/var/supplywatch/grpc/client/"
    receivLogFile = "sendlog"
    receivLog = receivLogFileDir + receivLogFile
)

type gserver struct {
	pb.UnimplementedProductServiceServer
    Conn net.Listener
    InOutProductChan chan *domain.InOutProduct
}

func New(inOutProductChan chan *domain.InOutProduct) (*gserver, error) {
	var tcpConn net.Listener
	tcpPort := os.Getenv("SW_GRPCPORT")
	listenIP := os.Getenv("SW_LISTEN_IP")
	tcpListenIP := listenIP + ":" + tcpPort
	tcpConn, err := net.Listen("tcp", tcpListenIP)
	if err != nil {
		return nil, err
	}
    return &gserver{
        Conn: tcpConn,
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
	// first we log the Request
	// this way we can compare the shippingReceivLog with the shippingSendLog of the other warehouse
	// if they don't match something went wrong
	//allRequests = append(allRequests, req)
	//allJsonReq, err := json.MarshalIndent(allRequests, "", "  ")
	//if err != nil {
		//return nil, err
	//}
	//err = ioutil.WriteFile(shippingReceivLog, allJsonReq, 0644)
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
        Incoming: true,
        Amount: 1,
        Reason: "Another Warehouse",
    }
    // inform the storage channel about the product that has been shipped to us
    s.InOutProductChan <- p

	return &pb.ReceivProductsResponse{
		Success: true,
	}, nil
}
