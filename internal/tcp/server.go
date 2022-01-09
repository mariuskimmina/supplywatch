package tcp

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/mariuskimmina/supplywatch/internal/domain"
	"github.com/mariuskimmina/supplywatch/internal/tcp/http"
)

type tcpServer struct {
    conn net.Listener
    wh domain.Warehouse
}

func NewTCPServer(wh domain.Warehouse) (*tcpServer, error) {
	var tcpConn net.Listener
	tcpPort := os.Getenv("SW_TCPPORT")
	listenIP := os.Getenv("SW_LISTEN_IP")
	tcpListenIP := listenIP + ":" + tcpPort
	tcpConn, err := net.Listen("tcp", tcpListenIP)
	if err != nil {
		return nil, err
	}
	return &tcpServer{
        conn: tcpConn,
        wh: wh,
    }, nil
}

func (s *tcpServer) Listen() error {
    fmt.Println("Hello from TCP Server")
	for {
		c, err := s.conn.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(c)
	}
}

func (s *tcpServer) handleConnection(c net.Conn) {
	netData, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	requestData := strings.Split(netData, " ")
	queryString := strings.Split(requestData[1], "?")
	httpVersion := strings.TrimSuffix(requestData[2], "\n")

	request, err := http.NewRequest(
		http.WithMethod(requestData[0]),
		http.WithPath(queryString[0]),
		http.WithVersion(httpVersion),
	)
	if len(queryString) > 1 {
		request.Query = queryString[1]
	}
	fmt.Printf("Received Request: %s, %s, %s \n", request.Method, request.Path, request.Version)
	if request.Path == "/allsensordata" {
		s.handleGetAllSensorData(request, c)
	} else if request.Path == "/sensordata" {
		s.handleGetOneSensorData(request, c)
	} else if request.Path == "/sensorhistory" {
		s.handleGetSensorHistory(request, c)
	} else if request.Path == "/" {
		s.handleWarehouseRequest(request, c)
	} else {
		s.handleRessourceNotFound(request, c)
	}
	c.Close()
}
