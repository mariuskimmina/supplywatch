package warehouse

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mariuskimmina/supplywatch/internal/pb"
	"google.golang.org/grpc"
)

var (
	hostname        = os.Getenv("SW_OTHER_WAREHOUSE_HOST")
	port            = os.Getenv("SW_OTHER_WAREHOUSE_PORT")
	address         = hostname + ":" + port
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

func (w *warehouse) grpcClient(sendChan chan string) {
	var conn *grpc.ClientConn
	var err error
	for {
		w.logger.Infof("Try GRPC Dialing %s", address)
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			w.logger.Error("Could not connect to the other warehouse, trying again in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		}
		w.logger.Info("Connected to the oter warehouse successfully")
		break
	}
	defer conn.Close()
	c := pb.NewProductServiceClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), 6 * time.Second)
	ctx := context.Background()
	//defer cancel()

	for {
		w.logger.Info("Ready to send a product - waiting for trigger")
		sendingProduct := <-sendChan
		w.logger.Infof("Received a message [ %s ] shipping product if possible", sendingProduct)
		s := strings.Split(sendingProduct, ":")
		product := s[0]
		warehouse := s[1]
		//choose a random product to ship to the other warehouse
		// SeedRandom()
		// n := rand.Int() % len(allProductNames)
		productName := product

		// remove the product from this warehouse
		// then send it to the other warehouse via grpc
		for _, product := range Products {
			if product.Name == productName {
				oldquantity := product.Quantity
				if oldquantity <= 0 {
					w.logger.Infof("Unable to send %s to %s because there are 0 left in stock", productName, warehouse)
				}
				w.logger.Infof("Sending %s to %s", productName, warehouse)
				product.Decrement()
				w.logger.Infof("Removing %s from this warehouse, quantity drops from %d to %d", productName, oldquantity, product.Quantity)
				w.DB.Model(&Product{}).Where("name = ?", product.Name).Update("quantity", product.Quantity)
				sendingProduct := &pb.Product{
					Name: product.Name,
					Id:   product.ID.String(),
				}
				sendProdcuts := &pb.ReceivProductsRequest{
					Product: sendingProduct,
					Amount:  1,
				}
				allOutgoingRequests = append(allOutgoingRequests, sendProdcuts)
				allJsonReq, err := json.MarshalIndent(allOutgoingRequests, "", "  ")
				if err != nil {
					w.logger.Fatal("Failed to marhal Request")
				}
				err = ioutil.WriteFile(shippingSendLog, allJsonReq, 0644)
				if err != nil {
					w.logger.Fatal("Failed to write Log")
				}
				for {
					_, err := c.ReceivProducts(ctx, sendProdcuts)
					if err != nil {
						w.logger.Error("Failed to send Products, trying again in 5 seconds")
						time.Sleep(5 * time.Second)
						continue
					}
					w.logger.Info("Send product successfully")
					break
				}
				break
			}
		}
		req := &pb.GetAllProductsRequest{}
		allProducts, err := c.GetAllProducts(ctx, req)
		if err != nil {
			w.logger.Error(err)
			w.logger.Fatal("Failed to get all Products")
		}
		allProductsJson, err := json.MarshalIndent(allProducts.Products, "", "  ")
		if err != nil {
			w.logger.Fatal("Failed to marhal Request")
		}
		err = ioutil.WriteFile(allProductsReceivLog, allProductsJson, 0644)
		if err != nil {
			w.logger.Fatal("Failed to write Log")
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
