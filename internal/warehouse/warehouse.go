package warehouse

import (
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/mariuskimmina/supplywatch/internal/domain"
	pb "github.com/mariuskimmina/supplywatch/internal/warehouse/grpc"
	"github.com/mariuskimmina/supplywatch/internal/warehouse/udp"
	"github.com/mariuskimmina/supplywatch/pkg/config"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type warehouse struct {
	logger Logger
	config *config.WarehouseConfig
	DB     *gorm.DB
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

func NewWarehouse(logger Logger, config *config.WarehouseConfig, db *gorm.DB) *warehouse {
	return &warehouse{
		logger: logger,
		config: config,
		DB:     db,
	}
}

const (
	maxBufferSize = 1024
)

var (
	todayTimeStamp = time.Now().Format("01-02-2006")
)

// Start starts the warehouse server
// The warehouse listens on a UPD Port to reiceive data from sensors
// and it also listens on a TCP Port to handle HTTP requests
func (w *warehouse) Start() {
	// this channel is for publishing messages once the capacity of an item reaches zero
	storageChan := make(chan string)
	sendChan := make(chan string)
    receivProdcutChan := make(chan *domain.ReceivedProduct)

    w.DB.AutoMigrate(&Product{})
	sqlDB, err := w.DB.DB()
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to connect to Database")
	}
	defer sqlDB.Close()

	// create all products with quanitity zero
	err = w.setupProducts()
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to setup Product Database")
	}

	var wg sync.WaitGroup
	wg.Add(4)

    udpServer, err := udp.NewUDPServer()
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to create UPD Server")
	}
    go func() {
        w.logger.Info("Starting UDP Server")
        udpServer.Listen(receivProdcutChan)
		wg.Done()
    }()

    go func() {
        for {
            newProduct := <- receivProdcutChan
            w.logger.Infof("Received new product from udp server: %s", newProduct.ProductName)
            w.HandleProduct(newProduct, storageChan)
        }
    }()


	//udpConn, err := setupUDPConn()
	//if err != nil {
		//w.logger.Error(err)
		//w.logger.Fatal("Failed to setup UPD Listener")
	//}
	//defer udpConn.Close()
	//go func() {
		//w.udpListen(udpConn, storageChan)
		//wg.Done()
	//}()

	tcpConn, err := setupTCPConn()
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to setup TCP Listener")
	}
	defer tcpConn.Close()
	go func() {
		w.tcpListen(tcpConn, w.DB)
		wg.Done()
	}()

	// RabbitMQ
	go func() {
		w.SetupMessageQueue(storageChan, sendChan)
		wg.Done()
	}()

	// GRPC part
	tcpConnGrpc, err := setupTCPConnGRPC()
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to setup TCP Listener")
	}
	defer tcpConnGrpc.Close()
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, w)
	go func() {
		w.logger.Info("GRPC Server Starts")
		if err := grpcServer.Serve(tcpConnGrpc); err != nil {
			w.logger.Fatal("Failed to setup GRPC Listener")
		}
		wg.Done()
		w.logger.Info("GRPC Server Ends")
	}()
	// wait a bit before we start the client
	time.Sleep(8 * time.Second)
	// GRPC Client starts here
	go func() {
		w.logger.Info("GRPC Client Starts")
		w.grpcClient(sendChan)
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
