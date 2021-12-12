package warehouse

import (
	"context"
	"encoding/json"
	"io/ioutil"

	pb "github.com/mariuskimmina/supplywatch/internal/warehouse/grpc"
)

//type ProductGrpcServer struct {
//pb.UnimplementedProductServiceServer
//}

var (
    allRequests []*pb.ReceivProductsRequest
)

const (
    shippingReceivLog = "/var/supplywatch/log/shipping_receiv_log"
    shippingSendLog = "/var/supplywatch/log/shipping_send_log"
)

func (w *warehouse) GetAllProducts(ctx context.Context, req *pb.GetAllProductsRequest) (*pb.GetAllProductsResponse, error) {
    w.logger.Info("Received GRPC Request, GetAllProducts")
    var products []*pb.Product
    return &pb.GetAllProductsResponse{
        Products: products,
    }, nil
}

func (w *warehouse) ReceivProducts(ctx context.Context, req *pb.ReceivProductsRequest) (*pb.ReceivProductsResponse, error) {
    // first we log the Request
    // this way we can compare the shippingReceivLog with the shippingSendLog of the other warehouse
    // if they don't match something went wrong
    allRequests = append(allRequests, req)
    allJsonReq, err := json.MarshalIndent(allRequests, "", "  ")
    if err != nil {
        return nil, err
    }
    err = ioutil.WriteFile(shippingReceivLog, allJsonReq, 0644)
    if err != nil {
        return nil, err
    }
    w.logger.Info("Received GRPC Request, ReceivProducts")
    product := &Product{}
    w.DB.First(product, "name = ?", req.Product.Name)
    oldQuantity := product.Quantity
    w.logger.Infof("Updating database quantity of %s from %d to %d", req.Product.Name, oldQuantity, req.Amount + int32(oldQuantity))
    w.DB.Model(&Product{}).Where("name = ?", req.Product.Name).Update("quantity", req.Amount + int32(oldQuantity))
    return &pb.ReceivProductsResponse{
        Success: true,
    }, nil
}