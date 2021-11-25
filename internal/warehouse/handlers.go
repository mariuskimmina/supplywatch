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
	} else if request.path == "/sensorhistory" {
		w.handleGetSensorHistory(&request, c)
	} else {
		c.Write([]byte(string(request.path)))
	}
	c.Close()
	//}
}

// handleGetAllSensorData handles requests to /allsensordata 
// we read the log file and return all entrys to the user
func (w *warehouse) handleGetAllSensorData(request *HTTPRequest, c net.Conn) {
	response, err := NewHTTPResponse()
	if err != nil {
		c.Write([]byte(err.Error()))
	}
	response.SetHeader("Access-Control-Allow-Origin", "*")
	response.SetHeader("Content-Type", "application/json")
	response.SetHeader("Server", "Supplywatch")

    allLogData, err := ReadAllLogs()
	if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to read all logs")
	}

    w.logger.Infof("Read the logs: %s", string(allLogData))

	response.SetBody(allLogData)
	byteResponse, _ := ResponseToBytes(response)
	c.Write(byteResponse)
}

func (w *warehouse) handleGetOneSensorData(request *HTTPRequest, c net.Conn) {
	response, err := NewHTTPResponse()
	if err != nil {
		c.Write([]byte(err.Error()))
	}
	response.SetHeader("Access-Control-Allow-Origin", "*")
	response.SetHeader("Content-Type", "application/json")
	response.SetHeader("Server", "Supplywatch")
	queryValue := strings.Split(request.query, "=")
    sensorData, err := ReadOneSensorLogs(queryValue[1])
	if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to read all logs")
	}
	response.SetBody(sensorData)
	byteResponse, _ := ResponseToBytes(response)
    c.Write(byteResponse)
}

func (w *warehouse) handleGetSensorHistory(request *HTTPRequest, c net.Conn) {

}


