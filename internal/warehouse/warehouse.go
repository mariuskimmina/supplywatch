package warehouse

import (
	"bytes"
	"encoding/json"
	"net"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/mariuskimmina/supplywatch/pkg/config"
	log "github.com/mariuskimmina/supplywatch/pkg/log"
)

type warehouse struct {
	logger *log.Logger
	config *config.Config
}

// Create a new warehouse object
// TODO: the arguments here should probably be interfaces, I think..
// this way, I think I'm doing depency injection wrong here...
func NewWarehouse(logger *log.Logger, config *config.Config) *warehouse {
	return &warehouse{
		logger: logger,
		config: config,
	}
}

const (
	maxBufferSize = 1024
)

// Start starts the warehouse server
// The warehouse listens on a UPD Port to reiceive data from sensors
// and it also listens on a TCP Port to handle HTTP requests
func (w *warehouse) Start() {
	address := &net.UDPAddr{
		Port: w.config.Warehouse.UDPPort,
		IP:   net.ParseIP(w.config.Warehouse.ListenIP),
	}
	listen, err := net.ListenUDP("udp", address)
	if err != nil {
		return
	}
	defer listen.Close()
	go w.recvDataFromSensor(listen)
	tcpPort := strconv.Itoa(w.config.Warehouse.TCPPort)
	tcpListenIP := w.config.Warehouse.ListenIP + ":" + tcpPort
	ln, err := net.Listen("tcp", tcpListenIP)
	if err != nil {
		w.logger.Error(err.Error())
		return
	}
	defer ln.Close()
	for {
		c, err := ln.Accept()
		if err != nil {
			w.logger.Error(err.Error())
			return
		}
		go w.handleConnection(c)
	}
}


// SensorMesage represents the data we hope to receive from a sensor
type SensorMesage struct {
	SensorID   uuid.UUID `json:"sensor_id"`
	SensorType string    `json:"sensor_type"`
	Message    string    `json:"message"`
}

// recvDataFromSensor handles incoming UPD Packets
func (w *warehouse) recvDataFromSensor(listen *net.UDPConn) {
	logfile := NewLogFile("/tmp/warehouselog")
    defer logfile.Close()
	logcount, err := os.Create("/tmp/logcount")
    defer logcount.Close()
	sensorCounter := []*SensorMessageCounter{}
	if err != nil {
		return
	}
	for {
		p := make([]byte, maxBufferSize)
		_, remoteaddr, err := listen.ReadFromUDP(p)
		if err != nil {
			w.logger.Error("Error reading data from UDP: ", err)
			return
		}
		sensorCleanBytes := bytes.Trim(p, "\x00")
		var sensorMessage SensorMesage
		err = json.Unmarshal(sensorCleanBytes, &sensorMessage)
		if err != nil {
			w.logger.Error("Error unmarshaling sensor data: ", err)
			return
		}
		w.logger.Infof("Received %s", sensorMessage.Message)
		logentry := &LogEntry{
			SensorType: sensorMessage.SensorType,
			SensorID:   sensorMessage.SensorID,
			Message:    sensorMessage.Message,
			IP:         remoteaddr.IP,
			Port:       remoteaddr.Port,
		}

		// to keep track of how many messages we have received form each sensor
		// check if we know any sensor yet, if not create a new one
		// else check if we have seen this sensor before
		// if yes, we increase it's counter
		// if not, we create a new counter for it
		var found bool
		if len(sensorCounter) == 0 {
			w.logger.Debug("Sensor added to list of sensors")
			newSensorCounter := NewSensorMessageCounter(logentry.SensorID)
			sensorCounter = append(sensorCounter, newSensorCounter)
		} else {
			for _, counter := range sensorCounter {
				if counter.SensorID == logentry.SensorID {
					found = true
					counter.increment()
					w.logger.Debug("Increased Counter")
					break
				} else {
					found = false
				}
			}
			if !found {
				w.logger.Debug("Sensor added to list of sensors")
				newSensorCounter := NewSensorMessageCounter(logentry.SensorID)
				sensorCounter = append(sensorCounter, newSensorCounter)
			}
		}

        logfile.addLog(*logentry)
        err = logfile.WriteToFile()
        if err != nil {
            w.logger.Fatalf("Failed to write to logfile: %v", err)
        }


		var jsonLogCount []byte
		for _, counter := range sensorCounter {
			jsonLogCountEntry, err := json.Marshal(counter)
			jsonLogCount = append(jsonLogCount, jsonLogCountEntry...)
			if err != nil {
				w.logger.Error("Error marshaling log counter to json: ", err)
				return
			}
			jsonLogCount = append(jsonLogCount, []byte("\n")...)
		}

		// We write to the start of the file meaning everytime we receive a packet we update the
		// /tmp/logcount file with the new counter - this way the file always contains only 1 line for
		// each Sensor with updated values
		logcount.WriteAt(jsonLogCount, 0)
	}
}
