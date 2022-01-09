package warehouse

//
// import (
// 	"bufio"
// 	"encoding/json"
// 	"fmt"
// 	"net"
// 	"os"
// 	"strings"
//
// 	"gorm.io/gorm"
//
// 	"github.com/mariuskimmina/supplywatch/internal/tcp/http"
// )
//
// func (w *warehouse) tcpListen(tcpConn net.Listener, db *gorm.DB) {
// 	for {
// 		c, err := tcpConn.Accept()
// 		if err != nil {
// 			w.logger.Error(err.Error())
// 			return
// 		}
// 		go w.handleConnection(c, db)
// 	}
// 	//}
// }
//
// func (w *warehouse) handleConnection(c net.Conn, db *gorm.DB) {
// 	netData, err := bufio.NewReader(c).ReadString('\n')
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	requestData := strings.Split(netData, " ")
// 	queryString := strings.Split(requestData[1], "?")
// 	httpVersion := strings.TrimSuffix(requestData[2], "\n")
//
// 	request, err := http.NewRequest(
// 		http.WithMethod(requestData[0]),
// 		http.WithPath(queryString[0]),
// 		http.WithVersion(httpVersion),
// 	)
// 	if len(queryString) > 1 {
// 		request.Query = queryString[1]
// 	}
// 	fmt.Printf("Received Request: %s, %s, %s \n", request.Method, request.Path, request.Version)
// 	if request.Path == "/allsensordata" {
// 		w.handleGetAllSensorData(request, c)
// 	} else if request.Path == "/sensordata" {
// 		w.handleGetOneSensorData(request, c)
// 	} else if request.Path == "/sensorhistory" {
// 		w.handleGetSensorHistory(request, c)
// 	} else if request.Path == "/" {
// 		w.handleWarehouseRequest(request, c, db)
// 	} else {
// 		w.handleRessourceNotFound(request, c)
// 	}
// 	c.Close()
// }
//
// func (w *warehouse) handleWarehouseRequest(request *http.Request, c net.Conn, db *gorm.DB) {
// 	response, err := http.NewResponse()
// 	if err != nil {
// 		c.Write([]byte(err.Error()))
// 		return
// 	}
// 	response.SetHeader("Content-Type", "application/json")
// 	response.SetHeader("Server", "Supplywatch")
// 	//products, err := GetAllProductsAsBytes()
// 	if err != nil {
// 		c.Write([]byte(err.Error()))
// 		w.logger.Fatal("Failed to get all bytes")
// 	}
// 	var products []Product
// 	db.Find(&products)
// 	response.SetBody([]byte(fmt.Sprintf("%v", products)))
// 	jsonProducts, err := json.MarshalIndent(products, " ", "")
// 	if err != nil {
// 		c.Write([]byte(err.Error()))
// 		return
// 	}
// 	//byteResponse, _ := ResponseToBytes(jsonProducts)
// 	c.Write(jsonProducts)
// }
//
// func (w *warehouse) handleRessourceNotFound(request *http.Request, c net.Conn) {
// 	response, err := http.NewResponse(
// 		http.WithStatusCode(404),
// 	)
// 	if err != nil {
// 		c.Write([]byte(err.Error()))
// 	}
// 	response.SetReason("Not Found")
// 	response.SetHeader("Server", "Supplywatch")
// 	response.SetBody([]byte("404 Not Found"))
// 	fmt.Println(response)
// 	byteResponse, _ := http.ResponseToBytes(response)
// 	c.Write(byteResponse)
// }
//
// // handleGetAllSensorData handles requests to /allsensordata
// // we read the log file and return all entrys to the user
// func (w *warehouse) handleGetAllSensorData(request *http.Request, c net.Conn) {
// 	response, err := http.NewResponse()
// 	if err != nil {
// 		c.Write([]byte(err.Error()))
// 	}
// 	response.SetHeader("Access-Control-Allow-Origin", "*")
// 	response.SetHeader("Content-Type", "application/json")
// 	response.SetHeader("Server", "Supplywatch")
//
// 	allLogData, err := w.GetAllSensorLogs(w.config.LogFileDir)
// 	if err != nil {
// 		w.logger.Error(err)
// 		w.logger.Fatal("Failed to read all logs")
// 	}
// 	response.SetBody(allLogData)
// 	byteResponse, _ := http.ResponseToBytes(response)
// 	c.Write(byteResponse)
// }
//
// func (w *warehouse) handleGetOneSensorData(request *http.Request, c net.Conn) {
// 	response, err := http.NewResponse()
// 	if err != nil {
// 		c.Write([]byte(err.Error()))
// 	}
// 	response.SetHeader("Access-Control-Allow-Origin", "*")
// 	response.SetHeader("Content-Type", "application/json")
// 	response.SetHeader("Server", "Supplywatch")
// 	queryValue := strings.Split(request.Query, "=")
//
// 	sensorData, err := GetOneSensorLogs(w.config.LogFileDir, queryValue[1])
// 	if err != nil {
// 		w.logger.Error(err)
// 		w.logger.Fatal("Failed to read all logs")
// 	}
// 	response.SetBody(sensorData)
// 	byteResponse, _ := http.ResponseToBytes(response)
// 	c.Write(byteResponse)
// }
//
// // handleGetSensorHistory takes a query parameter `date` and returns all
// // all logs from that day
// func (w *warehouse) handleGetSensorHistory(request *http.Request, c net.Conn) {
// 	response, err := http.NewResponse()
// 	if err != nil {
// 		c.Write([]byte(err.Error()))
// 	}
// 	response.SetHeader("Access-Control-Allow-Origin", "*")
// 	response.SetHeader("Server", "Supplywatch")
// 	queryValue := strings.Split(request.Query, "=")
// 	if queryValue[0] != "date" {
// 		response.SetHeader("Content-Type", "text/plain")
// 		response.SetBody([]byte("Unkown query parameter: " + queryValue[0]))
// 		byteResponse, _ := http.ResponseToBytes(response)
// 		c.Write(byteResponse)
// 		return
// 	}
// 	hostname, err := os.Hostname()
// 	if err != nil {
// 		w.logger.Fatal("Failed to access hostname")
// 	}
// 	logfileName := w.config.LogFileDir + hostname + "-" + queryValue[1]
// 	sensorData, err := ReadLogsFromDate(logfileName)
// 	if err != nil {
// 		response.SetHeader("Content-Type", "text/plain")
// 		response.SetBody([]byte("No data was found for this date"))
// 		byteResponse, _ := http.ResponseToBytes(response)
// 		c.Write(byteResponse)
// 		return
// 	}
// 	response.SetHeader("Content-Type", "application/json")
// 	response.SetBody(sensorData)
// 	byteResponse, _ := http.ResponseToBytes(response)
// 	c.Write(byteResponse)
// }
