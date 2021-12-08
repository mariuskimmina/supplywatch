package warehouse

import (
	"context"
	"log"

	pb "github.com/mariuskimmina/supplywatch/internal/warehouse/grpc"
)

type ProductGrpcServer struct {
    pb.UnimplementedProductServiceServer
}

func (ps *ProductGrpcServer) GetAllProducts(ctx context.Context, req *pb.GetAllProductsRequest) (*pb.GetAllProductsResponse, error) {
    log.Println("GetAllProducts")
    var products []*pb.Product
    return &pb.GetAllProductsResponse{
        Products: products,
    }, nil
}
