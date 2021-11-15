package warehouse

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	log "github.com/mariuskimmina/supplywatch/pkg/log"
)

type warehouse struct {
	logger *log.Logger
}

func NewWarehouse(logger *log.Logger) *warehouse {
	return &warehouse{
		logger: logger,
	}
}

var (
	address = net.UDPAddr{
		Port: 4444,
		IP:   net.ParseIP("0.0.0.0"),
	}
	logger = &log.Logger{}
)

const (
	maxBufferSize = 1024
)

func (w *warehouse) Start() {
	//ctx := context.Background()
	w.logger.Info("Warehouse Starting")
	logger = w.logger
	listen, err := net.ListenUDP("udp", &address)
	if err != nil {
		return
	}
	defer listen.Close()
	go recvDataFromSensor(listen)
	ln, err := net.Listen("tcp", "0.0.0.0:8000")
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
		go handleConnection(c)
	}
}

type HTTPRequest struct {
	method  string
	path    string
	version string
	query   string
}

type logEntry struct {
	SensorType string `json:"sensor_type"`
	Message    string `json:"message"`
	IP         net.IP `json:"ip"`
	Port       int    `json:"port"`
}

type SensorMesage struct {
	SensorType string `json:"sensor_type"`
	Message    string `json:"message"`
}

func recvDataFromSensor(listen *net.UDPConn) {
	f, err := os.Create("/tmp/warehouselog")
	if err != nil {
		return
	}
	defer f.Close()
	for {
		p := make([]byte, maxBufferSize)
		_, remoteaddr, err := listen.ReadFromUDP(p)
		if err != nil {
			logger.Error("Error reading data from UDP: ", err)
			return
		}
		sensorCleanBytes := bytes.Trim(p, "\x00")
		var sensorMessage SensorMesage
		err = json.Unmarshal(sensorCleanBytes, &sensorMessage)
		if err != nil {
			logger.Error("Error unmarshaling sensor data: ", err)
			return
		}
        logger.Infof("Received %s", sensorMessage.Message)
		logentry := &logEntry{
			SensorType: sensorMessage.SensorType,
			Message:    sensorMessage.Message,
			IP:         remoteaddr.IP,
			Port:       remoteaddr.Port,
		}
		jsonLogEntry, err := json.Marshal(logentry)
		if err != nil {
			logger.Error("Error marshaling log entry to json: ", err)
			return
		}
		f.Write(jsonLogEntry)
		f.Write([]byte("\n"))
	}
}

func handleConnection(c net.Conn) {
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
		handleGetAllSensorData(&request, c)
	} else if request.path == "/sensordata" {
		handleGetOneSensorData(&request, c)
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

func handleGetAllSensorData(request *HTTPRequest, c net.Conn) {
	response := HTTPResponse{
		HTTPVersion: "HTTP/1.1",
		statuscode:  200,
		reason:      "OK \r\n\r\n",
		body:        "All Sensor Data",
	}
	c.Write([]byte(fmt.Sprintf("%v", response)))
}

func handleGetOneSensorData(request *HTTPRequest, c net.Conn) {
	c.Write([]byte("One Sensor Data"))
}

func handleGetSensorHistorie() {

}
