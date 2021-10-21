package warehouse

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"github.com/mariuskimmina/supplywatch/pkg/logger"
)

type warehouse struct {

}

func NewWarehouse() *warehouse {
    return &warehouse{}
}

var (
    address = "0.0.0.0:1234"
)

const (
    maxBufferSize = 1024
)

func (w *warehouse) Start() {
    //ctx := context.Background()
    logger.Debug("Starting warehouse")
    //listen, err := net.ListenPacket("udp", address)
    //if err != nil {
        //return
    //}
    //defer listen.Close()
    ln, err := net.Listen("tcp", "127.0.0.1:8000")
    if err != nil {
        logger.Error(err.Error())
        return
    }
    defer ln.Close()
    for {
        logger.Debug("Warehouse running")
        c, err := ln.Accept()
        if err != nil {
            logger.Error(err.Error())
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
        method:     requestData[0],
        path:       queryString[0],
        version:    httpVersion,
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

func handleGetAllSensorData(request *HTTPRequest, c net.Conn) {
    c.Write([]byte("All Sensor Data"))
}

func handleGetOneSensorData(request *HTTPRequest, c net.Conn) {
    c.Write([]byte("One Sensor Data"))
}

func handleGetSensorHistorie() {

}
