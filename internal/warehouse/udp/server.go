package udp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/mariuskimmina/supplywatch/internal/domain"
)

type udpServer struct {
    conn *net.UDPConn
}

const (
	maxBufferSize = 1024
    logFileDir = "/var/supplywatch/udpserver/"
    logFile = "sensorlog"
)

func NewUDPServer() (*udpServer, error) {
	var udpConn *net.UDPConn
	udpPort, err := strconv.Atoi(os.Getenv("SW_UDPPORT"))
	listenIP := os.Getenv("SW_LISTEN_IP")
	if err != nil {
		return nil, err
	}
	address := &net.UDPAddr{
		Port: udpPort,
		IP:   net.ParseIP(listenIP),
	}
	udpConn, err = net.ListenUDP("udp", address)
	if err != nil {
		return nil, err
	}

    return &udpServer{
        conn: udpConn,
    }, nil
}


func (s udpServer) Listen(c chan *domain.ReceivedProduct) error {
    fmt.Println("Hello From UDP Server")
    err := os.MkdirAll(logFileDir, 0644)
    if err != nil {
        return fmt.Errorf("Error creating log file directory: %v", err)
    }
    f, err := os.OpenFile(logFileDir + logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("Error opening log file: %v", err)
    }
    defer f.Close()
	for {
		p := make([]byte, maxBufferSize)
		_, _, err := s.conn.ReadFromUDP(p)
		if err != nil {
            return fmt.Errorf("Error reading data from UDP: %v", err)
		}
		sensorCleanBytes := bytes.Trim(p, "\x00")
		var sensorMessage SensorMessage
		err = json.Unmarshal(sensorCleanBytes, &sensorMessage)
		if err != nil {
            return fmt.Errorf("Error unmarshaling sensor data: %v", err)
		}
		//w.logger.Infof("Received %s, Incoming: %v", sensorMessage.Message, sensorMessage.Incoming)
		logentry := &domain.SensorLog{
			SensorType: sensorMessage.SensorType,
			SensorID:   sensorMessage.SensorID,
			Message:    sensorMessage.Message,
			Incoming:   sensorMessage.Incoming,
		}
        logjson, err := json.Marshal(logentry)
		if err != nil {
            return fmt.Errorf("Error marshaling log data: %v", err)
		}
        f.WriteString(string(logjson) + "\n")
        fmt.Println("Sending Product to channel!")
        c <- &domain.ReceivedProduct{
            ProductName: sensorMessage.Message,
            Incoming: sensorMessage.Incoming,
            Amount: 1,
        }
	}
}

type SensorMessage struct {
	SensorID   string `json:"sensor_id"`
	SensorType string `json:"sensor_type"`
	Message    string `json:"message"`
	Incoming   bool   `json:"incoming"`
}
