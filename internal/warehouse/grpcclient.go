package warehouse

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	pb "github.com/mariuskimmina/supplywatch/internal/warehouse/grpc"
	"google.golang.org/grpc"
)

var (
    hostname = os.Getenv("SW_OTHER_WAREHOUSE_HOST")
    port = os.Getenv("SW_OTHER_WAREHOUSE_PORT")
    address = hostname + ":" + port
	allProductNames = []string{
		"butter",
		"sugar",
		"eggs",
		"baking powder",
		"cheese",
		"cinnamon",
		"oil",
		"carrots",
		"raisins",
		"walnuts",
	}
)

var (
    allOutgoingRequests []*pb.ReceivProductsRequest
)

func (w *warehouse) grpcClient() {
    w.logger.Infof("GRPC Dialing %s", address)
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    w.logger.Info("After Dial")
    if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to connect to the other warehouse")
    }
    w.logger.Info("Connected to the oter warehouse successfully")
    defer conn.Close()
    c := pb.NewProductServiceClient(conn)

    //ctx, cancel := context.WithTimeout(context.Background(), 6 * time.Second)
    ctx := context.Background()
    //defer cancel()

    req := &pb.GetAllProductsRequest{}

    products, err := c.GetAllProducts(ctx, req)
    w.logger.Info(products)
    for {
        w.logger.Info("Sending a product")
        //choose a random product to ship to the other warehouse
        SeedRandom()
        n := rand.Int() % len(allProductNames)
        productName := allProductNames[n]

        // remove the product from this warehouse
        // then send it to the other warehouse via grpc
        for _, product := range Products {
            if product.Name == productName {
                oldquantity := product.Quantity
                product.Decrement()
                w.logger.Infof("Removing %s from this warehouse, quantity drops from %d to %d", productName, oldquantity, product.Quantity)
                w.DB.Model(&Product{}).Where("name = ?", product.Name).Update("quantity", product.Quantity)
                sendingProduct := &pb.Product{
                    Name: product.Name,
                    Id: product.ID.String(),
                }
                sendProdcuts := &pb.ReceivProductsRequest{
                    Product: sendingProduct,
                    Amount: 1,
                }
                allOutgoingRequests = append(allRequests, sendProdcuts)
                allJsonReq, err := json.MarshalIndent(allRequests, "", "  ")
                if err != nil {
                    w.logger.Fatal("Failed to marhal Request")
                }
                err = ioutil.WriteFile(shippingSendLog, allJsonReq, 0644)
                if err != nil {
                    w.logger.Fatal("Failed to write Log")
                }
                resp, err := c.ReceivProducts(ctx, sendProdcuts)
                if err != nil {
                    w.logger.Error(err)
                    w.logger.Fatal("Failed to send Products or something")
                }
                fmt.Println(resp)
                break
            }
        }
		time.Sleep(5 * time.Second)
    }
}

func SeedRandom() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}
