package tcp

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/mariuskimmina/supplywatch/internal/tcp/http"
)

func (s *tcpServer) handleWarehouseRequest(request *http.Request, c net.Conn) {
	response, err := http.NewResponse()
	if err != nil {
		c.Write([]byte(err.Error()))
		return
	}
	response.SetHeader("Access-Control-Allow-Origin", "*")
	response.SetHeader("Content-Type", "application/json")
	response.SetHeader("Server", "Supplywatch")
	//products, err := GetAllProductsAsBytes()
	if err != nil {
		c.Write([]byte(err.Error()))
	}
    products, err := s.wh.AllProducts()
	if err != nil {
		c.Write([]byte(err.Error()))
		return
	}
	//db.Find(&products)
	response.SetBody([]byte(fmt.Sprintf("%v", products)))
	jsonProducts, err := json.Marshal(products)
	if err != nil {
		c.Write([]byte(err.Error()))
		return
	}
	//byteResponse, _ := ResponseToBytes(jsonProducts)
	c.Write(jsonProducts)
}

func (s *tcpServer) handleRessourceNotFound(request *http.Request, c net.Conn) {
	response, err := http.NewResponse(
		http.WithStatusCode(404),
	)
	if err != nil {
		c.Write([]byte(err.Error()))
	}
	response.SetReason("Not Found")
	response.SetHeader("Server", "Supplywatch")
	response.SetBody([]byte("404 Not Found"))
	fmt.Println(response)
	byteResponse, _ := http.ResponseToBytes(response)
	c.Write(byteResponse)
}

// handleGetAllSensorData handles requests to /allsensordata
// we read the log file and return all entrys to the user
func (s *tcpServer) handleGetAllSensorData(request *http.Request, c net.Conn) {
	response, err := http.NewResponse()
	if err != nil {
		c.Write([]byte(err.Error()))
	}
	response.SetHeader("Access-Control-Allow-Origin", "*")
	response.SetHeader("Content-Type", "application/json")
	response.SetHeader("Server", "Supplywatch")


	allLogData, err := s.wh.AllSensorLogs()
	if err != nil {
		c.Write([]byte(err.Error()))
	}

    // json does not allow a , after the last entry in a list, thus we remove it here
    if len(allLogData) > 0 {
        allLogData = allLogData[:len(allLogData)-2]
    }

    logDataList := string("[\n" + string(allLogData) + "\n]")
	response.SetBody([]byte(logDataList))
	byteResponse, _ := http.ResponseToBytes(response)
	c.Write(byteResponse)
}

func (s *tcpServer) handleGetOneSensorData(request *http.Request, c net.Conn) {
	response, err := http.NewResponse()
	if err != nil {
		c.Write([]byte(err.Error()))
	}
	response.SetHeader("Access-Control-Allow-Origin", "*")
	response.SetHeader("Content-Type", "application/json")
	response.SetHeader("Server", "Supplywatch")
	//queryValue := strings.Split(request.Query, "=")

	//sensorData, err := GetOneSensorLogs(w.config.LogFileDir, queryValue[1])
	//if err != nil {
        // TODO: better handling
        //return
	//}
	//response.SetBody(sensorData)
	byteResponse, _ := http.ResponseToBytes(response)
	c.Write(byteResponse)
}

// handleGetSensorHistory takes a query parameter `date` and returns all
// all logs from that day
func (s *tcpServer) handleGetSensorHistory(request *http.Request, c net.Conn) {
	response, err := http.NewResponse()
	if err != nil {
		c.Write([]byte(err.Error()))
	}
	response.SetHeader("Access-Control-Allow-Origin", "*")
	response.SetHeader("Server", "Supplywatch")
	queryValue := strings.Split(request.Query, "=")
	if queryValue[0] != "date" {
		response.SetHeader("Content-Type", "text/plain")
		response.SetBody([]byte("Unkown query parameter: " + queryValue[0]))
		byteResponse, _ := http.ResponseToBytes(response)
		c.Write(byteResponse)
		return
	}
	//hostname, err := os.Hostname()
	//if err != nil {
		//w.logger.Fatal("Failed to access hostname")
	//}
	//logfileName := w.config.LogFileDir + hostname + "-" + queryValue[1]
	//sensorData, err := ReadLogsFromDate(logfileName)
	if err != nil {
		response.SetHeader("Content-Type", "text/plain")
		response.SetBody([]byte("No data was found for this date"))
		byteResponse, _ := http.ResponseToBytes(response)
		c.Write(byteResponse)
		return
	}
	response.SetHeader("Content-Type", "application/json")
	//response.SetBody(sensorData)
	byteResponse, _ := http.ResponseToBytes(response)
	c.Write(byteResponse)
}
