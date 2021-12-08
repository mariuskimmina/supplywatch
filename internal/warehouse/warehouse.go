package warehouse

import (
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	pb "github.com/mariuskimmina/supplywatch/internal/warehouse/grpc"
	"github.com/mariuskimmina/supplywatch/pkg/config"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type warehouse struct {
	logger Logger
	config *config.Config
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

// Create a new warehouse object
// TODO: config.Config should also be replaced by a generic interface
func NewWarehouse(logger Logger, config *config.Config) *warehouse {
	return &warehouse{
		logger: logger,
		config: config,
	}
}

const (
	maxBufferSize = 1024
)

var (
	todayTimeStamp = time.Now().Format("01-02-2006")
)

var db *gorm.DB


// Start starts the warehouse server
// The warehouse listens on a UPD Port to reiceive data from sensors
// and it also listens on a TCP Port to handle HTTP requests
func (w *warehouse) Start() {
    db, err := initDB()
	if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to connect to Database")
	} else {
        w.logger.Info("Successfully connected to Database")
    }
    sqlDB, err := db.DB()
	if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to connect to Database")
    }
    defer sqlDB.Close()

    db.AutoMigrate(&Product{})

    // create all products with quanitity zero
    err = setupProducts(db)
	if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to setup Product Database")
    }


    var wg sync.WaitGroup
    wg.Add(2)
    udpConn, err := setupUDPConn()
	if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to setup UPD Listener")
	}
	defer udpConn.Close()
    go func() {
        w.udpListen(udpConn, db)
        wg.Done()
    }()

    tcpConn, err := setupTCPConn()
	if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to setup TCP Listener")
	}
	defer tcpConn.Close()
    go func() {
        w.tcpListen(tcpConn, db)
        wg.Done()
    }()

    tcpConnGrpc, err := setupTCPConnGRPC()
	if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to setup TCP Listener")
	}
    defer tcpConnGrpc.Close()
    grpcServer := grpc.NewServer()
    pb.RegisterProductServiceServer(grpcServer, &ProductGrpcServer{})
    go func() {
        if err := grpcServer.Serve(tcpConnGrpc); err != nil {
            w.logger.Fatal("Failed to setup GRPC Listener")
        }
        wg.Done()
    }()
    Connect()
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
type SensorMesage struct {
	SensorID   string `json:"sensor_id"`
	SensorType string    `json:"sensor_type"`
	Message    string    `json:"message"`
}



