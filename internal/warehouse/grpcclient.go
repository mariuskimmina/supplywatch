package warehouse

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/mariuskimmina/supplywatch/internal/warehouse/grpc"
	"google.golang.org/grpc"
)

var (
    hostname = os.Getenv("SW_OTHER_WAREHOUSE_HOST")
    port = os.Getenv("SW_OTHER_WAREHOUSE_PORT")
    address = hostname + ":" + port
)

func Connect() {
    fmt.Println("Connecting via GRPC")
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatal("Could not connect to the other warehouse")
    }
    defer conn.Close()
    c := pb.NewProductServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    req := &pb.GetAllProductsRequest{}

    products, err := c.GetAllProducts(ctx, req)
    fmt.Println(products)
}
