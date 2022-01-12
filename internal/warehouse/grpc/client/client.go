package gclient

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mariuskimmina/supplywatch/internal/domain"
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

const (
	sendLogFileDir = "/var/supplywatch/grpc/client/"
	sendLogFile    = "sendlog"
	sendLog        = sendLogFileDir + sendLogFile
)

type gclient struct {
	pb.ProductServiceClient
}

func New(conn *grpc.ClientConn) (*gclient, error) {
	//var conn *grpc.ClientConn
	//var attempt int
	//var err error
	//for {
	//time.Sleep(backoff.Default.Duration(attempt))
	//conn, err = grpc.Dial(address, grpc.WithInsecure())
	//if err != nil {
	//attempt++
	//fmt.Printf("Failed to Connect via GRPC, trying again in %d seconds\n", backoff.Default.Duration(attempt))
	//continue
	//}
	//break
	//}
	//defer conn.Close()
	c := pb.NewProductServiceClient(conn)
	return &gclient{
		c,
	}, nil
}

func (c *gclient) Start(sendChan chan string, inOutProductChan chan *domain.InOutProduct) error {
	fmt.Println("Hello From gRPC Client")
	err := os.MkdirAll(sendLogFileDir, 0644)
	if err != nil {
		return fmt.Errorf("Error creating log file directory: %v", err)
	}
	f, err := os.OpenFile(sendLogFileDir+sendLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Error opening log file: %v", err)
	}
	defer f.Close()
	for {
		fmt.Println("gRPC Client ready - waiting for trigger")
		//w.logger.Info("Ready to send a product - waiting for trigger")
		sendingProduct := <-sendChan
		fmt.Printf("Received a message [ %s ] shipping product if possible \n", sendingProduct)
		s := strings.Split(sendingProduct, ":")
		pname := s[0]
		pid := s[1]
		//warehouse := s[2]
		p := &domain.InOutProduct{
			Amount:      1,
			ProductName: pname,
			Incoming:    false,
			Reason:      "Shipping to another warehouse",
		}
		// tell the warehouse to remove the product we are going to send
		inOutProductChan <- p

		sendProduct := &pb.Product{
			Name: pname,
			Id:   pid,
		}
		sendProdcuts := &pb.ReceivProductsRequest{
			Product: sendProduct,
			Amount:  1,
		}

		logentry := &domain.ShippingLog{
			ShippingProductID:   pid,
			ShippingProductName: pname,
		}
		logjson, err := json.Marshal(logentry)
		if err != nil {
			return fmt.Errorf("Error marshaling log data: %v", err)
		}
		f.WriteString(string(logjson) + ",\n")
		if err != nil {
			return fmt.Errorf("Error writing log file: %v", err)
		}
		ctx := context.Background()
		for {
			_, err := c.ReceivProducts(ctx, sendProdcuts)
			if err != nil {
				fmt.Println("Failed to send Products, trying again in 5 seconds")
				time.Sleep(5 * time.Second)
				continue
			}
			fmt.Println("Send product successfully")
			break
		}

		//req := &pb.GetAllProductsRequest{}
		//allProducts, err := c.GetAllProducts(ctx, req)
		//if err != nil {
		//w.logger.Error(err)
		//w.logger.Fatal("Failed to get all Products")
		//}
		//allProductsJson, err := json.MarshalIndent(allProducts.Products, "", "  ")
		//if err != nil {
		//w.logger.Fatal("Failed to marhal Request")
		//}
		//err = ioutil.WriteFile(allProductsReceivLog, allProductsJson, 0644)
		//if err != nil {
		//w.logger.Fatal("Failed to write Log")
		//}
		//time.Sleep(5 * time.Second)
	}
}
