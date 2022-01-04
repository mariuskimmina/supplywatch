package warehouse

import (
	"bytes"
	"encoding/json"
	"net"
	"os"
)

// recvDataFromSensor handles incoming UPD Packets
func (w *warehouse) udpListen(listen *net.UDPConn, storageChan chan string) {
    err := LoadProductsState()
    if err != nil {
        w.logger.Fatal("Failed load products")
    }
	hostname, err := os.Hostname()
	if err != nil {
		w.logger.Fatal("Failed to access hostname")
	}
	logfileName := w.config.Warehouse.LogFileDir + hostname + "-" + todayTimeStamp
	logfile := w.NewLogFile(logfileName)
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
		var sensorMessage SensorMessage
		err = json.Unmarshal(sensorCleanBytes, &sensorMessage)
		if err != nil {
			w.logger.Error("Error unmarshaling sensor data: ", err)
			return
		}
        w.logger.Infof("Received %s, Incoming: %v", sensorMessage.Message, sensorMessage.Incoming)
		logentry := &LogEntry{
			SensorType: sensorMessage.SensorType,
			SensorID:   sensorMessage.SensorID,
			Message:    sensorMessage.Message,
			Incoming:    sensorMessage.Incoming,
			IP:         remoteaddr.IP,
			Port:       remoteaddr.Port,
		}
        if sensorMessage.Incoming {
            w.IncrementorCreateProduct(sensorMessage.Message)
        } else {
            w.DecrementProduct(sensorMessage.Message, storageChan)
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
        SaveProductsState()
	}
}
