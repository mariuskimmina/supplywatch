package warehouse

import (
	"bufio"
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
)

const (
	maxBufferSize = 1024
)

func (w *warehouse) Start() {
	//ctx := context.Background()
	w.logger.Info("Warehouse Starting")
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

func recvDataFromSensor(listen *net.UDPConn) {
    f, err := os.Create("/tmp/log2")
    if err != nil {
        return
    }
    defer f.Close()
	for {
		p := make([]byte, maxBufferSize)
		_, remoteaddr, err := listen.ReadFromUDP(p)
		if err != nil {
			// logger.Error(err.Error())
			return
		}
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
        f.Write(p)
        f.Write([]byte("\n"))
	}
}

func handleConnection(c net.Conn) {
	//fmt.Printf("Serving %s\n", c.RemoteAddr().String())
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
    fmt.Printf("Received Request: %s, %s, %s", request.method, request.path, request.version)
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


func handleGetAllSensorData(request *HTTPRequest, c net.Conn) {
    response, err := NewHTTPResponse()
    if err != nil {
        c.Write([]byte(err.Error()))
    }
    c.Write([]byte(fmt.Sprintf("%v", response)))
}

func handleGetOneSensorData(request *HTTPRequest, c net.Conn) {
	c.Write([]byte("One Sensor Data"))
}

func handleGetSensorHistorie() {

}
