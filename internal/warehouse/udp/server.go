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
    f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("Error opening log file: %v", err)
    }
    defer f.Close()
	for {
		p := make([]byte, maxBufferSize)
		_, remoteaddr, err := s.conn.ReadFromUDP(p)
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
			IP:         remoteaddr.IP,
			Port:       remoteaddr.Port,
		}
        logjson, err := json.MarshalIndent(logentry, " ", "")
		if err != nil {
            return fmt.Errorf("Error marshaling log data: %v", err)
		}
        f.Write(logjson)
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

		//if sensorMessage.Incoming {
			//w.IncrementorCreateProduct(sensorMessage.Message)
		//} else {
			//w.DecrementProduct(sensorMessage.Message, storageChan)
		//}

		// to keep track of how many messages we have received form each sensor
		// check if we know any sensor yet, if not create a new one
		// else check if we have seen this sensor before
		// if yes, we increase it's counter
		// if not, we create a new counter for it
		//var found bool
		//if len(sensorCounter) == 0 {
			//w.logger.Debug("Sensor added to list of sensors")
			//newSensorCounter := NewSensorMessageCounter(logentry.SensorID)
			//sensorCounter = append(sensorCounter, newSensorCounter)
		//} else {
			//for _, counter := range sensorCounter {
				//if counter.SensorID == logentry.SensorID {
					//found = true
					//counter.increment()
					//w.logger.Debug("Increased Counter")
					//break
				//} else {
					//found = false
				//}
			//}
			//if !found {
				//w.logger.Debug("Sensor added to list of sensors")
				//newSensorCounter := NewSensorMessageCounter(logentry.SensorID)
				//sensorCounter = append(sensorCounter, newSensorCounter)
			//}
		//}

		//logfile.addLog(*logentry)
		//err = logfile.WriteToFile()
		//if err != nil {
			//w.logger.Fatalf("Failed to write to logfile: %v", err)
		//}

		//var jsonLogCount []byte
		//for _, counter := range sensorCounter {
			//jsonLogCountEntry, err := json.Marshal(counter)
			//jsonLogCount = append(jsonLogCount, jsonLogCountEntry...)
			//if err != nil {
				//w.logger.Error("Error marshaling log counter to json: ", err)
				//return
			//}
			//jsonLogCount = append(jsonLogCount, []byte("\n")...)
		//}

		// We write to the start of the file meaning everytime we receive a packet we update the
		// /tmp/logcount file with the new counter - this way the file always contains only 1 line for
		// each Sensor with updated values
		//logcount.WriteAt(jsonLogCount, 0)
		//SaveProductsState()
