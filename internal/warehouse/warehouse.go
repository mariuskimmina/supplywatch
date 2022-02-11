package warehouse

import (
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/mariuskimmina/supplywatch/internal/domain"
	"github.com/mariuskimmina/supplywatch/internal/pb"
	gclient "github.com/mariuskimmina/supplywatch/internal/warehouse/grpc/client"
	gserver "github.com/mariuskimmina/supplywatch/internal/warehouse/grpc/server"
	"github.com/mariuskimmina/supplywatch/internal/warehouse/udp"

	//"github.com/mariuskimmina/supplywatch/internal/tcp"
	"github.com/mariuskimmina/supplywatch/pkg/backoff"
	"github.com/mariuskimmina/supplywatch/pkg/config"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type warehouse struct {
	logger   Logger
	config   *config.WarehouseConfig
	swConfig *config.SupplywatchConfig
	DB       *gorm.DB
	pb.UnimplementedProductServiceServer
}

// Logger is a generic interface that can be implemented by any logging engine
// this allows for dependency injection which results in easier testing
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

func NewWarehouse(logger Logger, config *config.WarehouseConfig, swConfig *config.SupplywatchConfig, db *gorm.DB) *warehouse {
	return &warehouse{
		logger:   logger,
		config:   config,
		swConfig: swConfig,
		DB:       db,
	}
}

const (
	maxBufferSize = 1024
)

var (
	todayTimeStamp = time.Now().Format("01-02-2006")
)

var (
	hostname        = os.Getenv("SW_OTHER_WAREHOUSE_HOST")
	port            = os.Getenv("SW_OTHER_WAREHOUSE_PORT")
	address         = hostname + ":" + port
	warehouses      = os.Getenv("SW_WAREHOUSES") //list of all warehouses (hostnames)
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

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// Start starts the warehouse server
// The warehouse listens on a UPD Port to reiceive data from sensors
// and it also listens on a TCP Port to handle HTTP requests
func (w *warehouse) Start() {
	// this channel is for publishing messages once the capacity of an item reaches zero
	storageChan := make(chan string)
	sendChan := make(chan string)
	inOutProductChan := make(chan *domain.InOutProduct)

	w.DB.AutoMigrate(Product{})
	//w.DB.AutoMigrate(&Product{})
	sqlDB, err := w.DB.DB()
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to connect to Database")
	}
	defer sqlDB.Close()
	w.logger.Info("Successfully Connected to Database")

	// create all products with quanitity five
	err = w.setupProducts()
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to setup Product Database")
	}
	w.logger.Info("Successfully setup Products on Database")

	var wg sync.WaitGroup
	wg.Add(4)

	w.logger.Info("Setting up UDP Server")
	udpServer, err := udp.NewUDPServer()
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to create UPD Server")
	}
	go func() {
		w.logger.Info("Starting UDP Server")
		udpServer.Listen(inOutProductChan)
		wg.Done()
	}()

	go func() {
		for {
			newProduct := <-inOutProductChan
			var inOut string
			if newProduct.Incoming {
				inOut = "coming in from"
			} else {
				inOut = "leaving because"
			}
			w.logger.Infof("Product %s is %s %s ", newProduct.ProductName, inOut, newProduct.Reason)
			w.HandleProduct(newProduct, storageChan)
		}
	}()

	// RabbitMQ
	go func() {
		w.SetupPublishing(storageChan, sendChan, w.swConfig.Warehouses)
		wg.Done()
	}()

	//go func() {
	//w.SetupConsuming(storageChan, sendChan, w.swConfig.Warehouses)
	//wg.Done()
	//}()

	// GRPC Server part
	gserver, err := gserver.New(inOutProductChan)
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to setup TCP Listener")
	}
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, gserver)
	go func() {
		w.logger.Info("GRPC Server Starts")
		if err := grpcServer.Serve(gserver.Conn); err != nil {
			w.logger.Fatal("Failed to setup GRPC Listener")
		}
		wg.Done()
		w.logger.Info("GRPC Server Ends")
	}()

	// GRPC Client starts here
	var conn *grpc.ClientConn
	var attempt int
	//var err error
	for {
		time.Sleep(backoff.Default.Duration(attempt))
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			attempt++
			w.logger.Infof("Failed to Connect via GRPC, trying again in %d seconds\n", backoff.Default.Duration(attempt))
			continue
		}
		break
	}
	defer conn.Close()
	gc, err := gclient.New(conn)
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to setup GRPC Client")
	}
	go func() {
		w.logger.Info("GRPC Client Starts")
		//w.grpcClient(sendChan)
		gc.Start(sendChan, inOutProductChan)
		wg.Done()
		w.logger.Info("GRPC Client Ends")
	}()
	wg.Wait()
}

func setupTCPConn() (net.Listener, error) {
	var tcpConn net.Listener
	tcpPort := os.Getenv("SW_TCPPORT")
	listenIP := os.Getenv("SW_LISTEN_IP")
	tcpListenIP := listenIP + ":" + tcpPort
	tcpConn, err := net.Listen("tcp", tcpListenIP)
	if err != nil {
		return tcpConn, err
	}
	return tcpConn, nil
}

func setupTCPConnGRPC() (net.Listener, error) {
	var tcpConn net.Listener
	tcpPort := os.Getenv("SW_GRPCPORT")
	listenIP := os.Getenv("SW_LISTEN_IP")
	tcpListenIP := listenIP + ":" + tcpPort
	tcpConn, err := net.Listen("tcp", tcpListenIP)
	if err != nil {
		return tcpConn, err
	}
	return tcpConn, nil
}

func setupUDPConn() (*net.UDPConn, error) {
	var udpConn *net.UDPConn
	udpPort, err := strconv.Atoi(os.Getenv("SW_UDPPORT"))
	listenIP := os.Getenv("SW_LISTEN_IP")
	if err != nil {
		return udpConn, err
	}
	address := &net.UDPAddr{
		Port: udpPort,
		IP:   net.ParseIP(listenIP),
	}
	udpConn, err = net.ListenUDP("udp", address)
	if err != nil {
		return udpConn, err
	}
	return udpConn, nil
}

// SensorMesage represents the data we hope to receive from a sensor
type SensorMessage struct {
	SensorID   string `json:"sensor_id"`
	SensorType string `json:"sensor_type"`
	Message    string `json:"message"`
	Incoming   bool   `json:"incoming"`
}
