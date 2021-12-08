package warehouse

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"gorm.io/gorm"
)

func (w *warehouse) tcpListen(tcpConn net.Listener, db *gorm.DB) {
	for {
		c, err := tcpConn.Accept()
		if err != nil {
			w.logger.Error(err.Error())
			return
		}
		go w.handleConnection(c, db)
	}
	//}
}

func (w *warehouse) handleConnection(c net.Conn, db *gorm.DB) {
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
	fmt.Printf("Received Request: %s, %s, %s \n", request.method, request.path, request.version)
	if request.path == "/allsensordata" {
		w.handleGetAllSensorData(&request, c)
	} else if request.path == "/sensordata" {
		w.handleGetOneSensorData(&request, c)
	} else if request.path == "/sensorhistory" {
		w.handleGetSensorHistory(&request, c)
	} else if request.path == "/" {
		w.handleWarehouseRequest(&request, c, db)
	} else {
		w.handleRessourceNotFound(&request, c)
	}
	c.Close()
}

func (w *warehouse) handleWarehouseRequest(request *HTTPRequest, c net.Conn, db *gorm.DB) {
	response, err := NewHTTPResponse()
	if err != nil {
		c.Write([]byte(err.Error()))
	}
	response.SetHeader("Content-Type", "application/json")
	response.SetHeader("Server", "Supplywatch")
    //products, err := GetAllProductsAsBytes()
    if err != nil {
		c.Write([]byte(err.Error()))
        w.logger.Fatal("Failed to get all bytes")
    }
    result := db.Find(&Products)
    response.SetBody([]byte(fmt.Sprintf("%v", result)))
	byteResponse, _ := ResponseToBytes(response)
	c.Write(byteResponse)
}

func (w *warehouse) handleRessourceNotFound(request *HTTPRequest, c net.Conn) {
	response, err := NewHTTPResponse()
	if err != nil {
		c.Write([]byte(err.Error()))
	}
	response.SetStatusCode(404)
	response.SetReason("Not Found")
	response.SetHeader("Server", "Supplywatch")
	response.SetBody([]byte("404 Not Found"))
	fmt.Println(response)
	byteResponse, _ := ResponseToBytes(response)
	c.Write(byteResponse)
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

	allLogData, err := GetAllSensorLogs(w.config.Warehouse.LogFileDir)
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to read all logs")
	}
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

	sensorData, err := GetOneSensorLogs(w.config.Warehouse.LogFileDir, queryValue[1])
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to read all logs")
	}
	response.SetBody(sensorData)
	byteResponse, _ := ResponseToBytes(response)
	c.Write(byteResponse)
}

// handleGetSensorHistory takes a query parameter `date` and returns all
// all logs from that day
func (w *warehouse) handleGetSensorHistory(request *HTTPRequest, c net.Conn) {
	response, err := NewHTTPResponse()
	if err != nil {
		c.Write([]byte(err.Error()))
	}
	response.SetHeader("Access-Control-Allow-Origin", "*")
	response.SetHeader("Server", "Supplywatch")
	queryValue := strings.Split(request.query, "=")
	if queryValue[0] != "date" {
		response.SetHeader("Content-Type", "text/plain")
		response.SetBody([]byte("Unkown query parameter: " + queryValue[0]))
		byteResponse, _ := ResponseToBytes(response)
		c.Write(byteResponse)
		return
	}
	hostname, err := os.Hostname()
	if err != nil {
		w.logger.Fatal("Failed to access hostname")
	}
	logfileName := w.config.Warehouse.LogFileDir + hostname + "-" + queryValue[1]
	sensorData, err := ReadLogsFromDate(logfileName)
	if err != nil {
		response.SetHeader("Content-Type", "text/plain")
		response.SetBody([]byte("No data was found for this date"))
		byteResponse, _ := ResponseToBytes(response)
		c.Write(byteResponse)
		return
	}
	response.SetHeader("Content-Type", "application/json")
	response.SetBody(sensorData)
	byteResponse, _ := ResponseToBytes(response)
	c.Write(byteResponse)
}
