package warehouse

import (
	"bytes"
	"encoding/json"
	"net"
	"os"
	"strconv"

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

// Start starts the warehouse server
// The warehouse listens on a UPD Port to reiceive data from sensors
// and it also listens on a TCP Port to handle HTTP requests
func (w *warehouse) Start() {
	address := &net.UDPAddr{
		Port: w.config.Warehouse.UDPPort,
		IP:   net.ParseIP(w.config.Warehouse.ListenIP),
	}
	listen, err := net.ListenUDP("udp", address)
	if err != nil {
		return
	}
	defer listen.Close()
	go w.recvDataFromSensor(listen)
	tcpPort := strconv.Itoa(w.config.Warehouse.TCPPort)
	tcpListenIP := w.config.Warehouse.ListenIP + ":" + tcpPort
	ln, err := net.Listen("tcp", tcpListenIP)
	if err != nil {
		w.logger.Error(err.Error())
		return
	}
	defer ln.Close()
	for {
		c, err := ln.Accept()
		if err != nil {
			w.logger.Error(err.Error())
			return
		}
		go w.handleConnection(c)
	}
}

// LogEntry represents a new entry in the log file
type logEntry struct {
	SensorID   uuid.UUID `json:"sensor_id"`
	SensorType string    `json:"sensor_type"`
	Message    string    `json:"message"`
	IP         net.IP    `json:"ip"`
	Port       int       `json:"port"`
}

// SensorMesage represents the data we hope to receive from a sensor
type SensorMesage struct {
	SensorID   uuid.UUID `json:"sensor_id"`
	SensorType string    `json:"sensor_type"`
	Message    string    `json:"message"`
}

// recvDataFromSensor handles incoming UPD Packets
func (w *warehouse) recvDataFromSensor(listen *net.UDPConn) {
	f, err := os.Create("/tmp/warehouselog")
	if err != nil {
		return
	}
	defer f.Close()
	for {
		p := make([]byte, maxBufferSize)
		_, remoteaddr, err := listen.ReadFromUDP(p)
		if err != nil {
			w.logger.Error("Error reading data from UDP: ", err)
			return
		}
		sensorCleanBytes := bytes.Trim(p, "\x00")
		var sensorMessage SensorMesage
		err = json.Unmarshal(sensorCleanBytes, &sensorMessage)
		if err != nil {
			w.logger.Error("Error unmarshaling sensor data: ", err)
			return
		}
		w.logger.Infof("Received %s", sensorMessage.Message)
		logentry := &logEntry{
			SensorType: sensorMessage.SensorType,
			SensorID:   sensorMessage.SensorID,
			Message:    sensorMessage.Message,
			IP:         remoteaddr.IP,
			Port:       remoteaddr.Port,
		}
		jsonLogEntry, err := json.Marshal(logentry)
		if err != nil {
			w.logger.Error("Error marshaling log entry to json: ", err)
			return
		}
		f.Write(jsonLogEntry)
		f.Write([]byte("\n"))
	}
}
