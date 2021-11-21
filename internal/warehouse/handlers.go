package warehouse

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func (w *warehouse) handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
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
		w.handleGetAllSensorData(&request, c)
	} else if request.path == "/sensordata" {
		w.handleGetOneSensorData(&request, c)
	} else {
		c.Write([]byte(string(request.path)))
	}
	c.Close()
	//}
}


func (w *warehouse) handleGetAllSensorData(request *HTTPRequest, c net.Conn) {
    response, err := NewHTTPResponse()
    if err != nil {
        c.Write([]byte(err.Error()))
    }
    c.Write([]byte(fmt.Sprintf("%v", response)))
}

func (w *warehouse) handleGetOneSensorData(request *HTTPRequest, c net.Conn) {
	c.Write([]byte("One Sensor Data"))
}

func (w *warehouse) handleGetSensorHistorie() {

}