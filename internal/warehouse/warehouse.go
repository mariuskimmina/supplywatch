package warehouse

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	log "github.com/mariuskimmina/supplywatch/pkg/log"
)

type warehouse struct {
    logger  *log.Logger
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
    w.logger.Info("test")
	listen, err := net.ListenUDP("udp", &address)
	if err != nil {
		return
	}
	defer listen.Close()
	go recvDataFromSensor(listen)
	ln, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
        // logger.Error(err.Error())
		return
	}
	defer ln.Close()
	for {
        // logger.Debug("Warehouse running")
		c, err := ln.Accept()
		if err != nil {
            // logger.Error(err.Error())
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
	for {
		p := make([]byte, maxBufferSize)
		_, remoteaddr, err := listen.ReadFromUDP(p)
		if err != nil {
			// logger.Error(err.Error())
			return
		}
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
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
