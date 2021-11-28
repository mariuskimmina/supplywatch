package warehouse

import (
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mariuskimmina/supplywatch/pkg/config"
	log "github.com/mariuskimmina/supplywatch/pkg/log"
)

type warehouse struct {
	logger *log.Logger
	config *config.Config
}

// Create a new warehouse object
// TODO: the arguments here should probably be interfaces, I think..
// this way, I think I'm doing depency injection wrong here...
func NewWarehouse(logger *log.Logger, config *config.Config) *warehouse {
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

// Start starts the warehouse server
// The warehouse listens on a UPD Port to reiceive data from sensors
// and it also listens on a TCP Port to handle HTTP requests
func (w *warehouse) Start() {
    var wg sync.WaitGroup
    wg.Add(2)
    udpConn, err := setupUDPConn()
	if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to setup UPD Listener")
	}
	defer udpConn.Close()
    go func() {
        w.udpListen(udpConn)
        wg.Done()
    }()

    tcpConn, err := setupTCPConn()
	if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to setup TCP Listener")
	}
	defer tcpConn.Close()
    go func() {
        w.tcpListen(tcpConn)
        wg.Done()
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
	SensorID   uuid.UUID `json:"sensor_id"`
	SensorType string    `json:"sensor_type"`
	Message    string    `json:"message"`
}



