package warehouse

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

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

// HTTPRequest represents an incoming HTTP request
type HTTPRequest struct {
	method  string
	path    string
	version string
	query   string
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

// handleConnection handles incoming HTTP requests
func (w *warehouse) handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	//for {
	netData, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	requestData := strings.Split(netData, " ")
	queryString := strings.Split(requestData[1], "?")
	httpVersion := strings.TrimSuffix(requestData[2], "\n")
	request := HTTPRequest{
		method:  requestData[0],
		path:    queryString[0],
		version: httpVersion,
	}
	if len(queryString) > 1 {
		request.query = queryString[1]
	}
	fmt.Println(request.method)
	fmt.Println(request.path)
	fmt.Println(request.version)
	fmt.Println(request.query)
	if request.path == "/allsensordata" {
		w.handleGetAllSensorData(&request, c)
	} else if request.path == "/sensordata" {
		w.handleGetOneSensorData(&request, c)
	} else {
		c.Write([]byte(string(request.path)))
	}
	c.Close()
	//}
}

type HTTPResponse struct {
	HTTPVersion string // HTTP/1.1
	statuscode  int    // 200
	reason      string // ok
	body        string // content
}

func (w *warehouse) handleGetAllSensorData(request *HTTPRequest, c net.Conn) {
	response := HTTPResponse{
		HTTPVersion: "HTTP/1.1",
		statuscode:  200,
		reason:      "OK \r\n\r\n",
		body:        "All Sensor Data",
	}
	c.Write([]byte(fmt.Sprintf("%v", response)))
}

func (w *warehouse) handleGetOneSensorData(request *HTTPRequest, c net.Conn) {
	c.Write([]byte("One Sensor Data"))
}

func handleGetSensorHistorie() {

}
