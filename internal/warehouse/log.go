package warehouse

// Disclaimer: how we handle our log files here is terrible inefficent
// we constantly overwrite data with the same data or read more than we would need to
// the focus of this application is to build a distributed system and not to be the most
// performant it could be

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
)

// LogEntry represents a new entry in the log file
type LogEntry struct {
	SensorID   string `json:"sensor_id"`
	SensorType string `json:"sensor_type"`
	Message    string `json:"message"`
	Incoming   bool   `json:"incoming"`
	IP         net.IP `json:"ip"`
	Port       int    `json:"port"`
}

// LogFile is a wrapper around os.File which represents our log file
type LogFile struct {
	*os.File
	Logs []LogEntry
}

func (l *LogFile) addLog(log LogEntry) {
	l.Logs = append(l.Logs, log)
}

// NewLogFile Creates a new log file if one does not exist for today
// if a log file for today already exist it will open it at the end of the NewLogFile
// thus new data will be appended to an existing log file
func (w *warehouse) NewLogFile(path string) *LogFile {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			w.logger.Error(err)
			w.logger.Fatal("Failed to create LogFile at os.MkdirAll")
		}
		file, err := os.Create(path)
		if err != nil {
			w.logger.Error(err)
			w.logger.Fatal("Failed to create LogFile at os.Create")
		}
		logFile := &LogFile{
			File: file,
		}

		jsonStruct, err := json.MarshalIndent(&logFile, "", "  ")
		if err != nil {
			w.logger.Error(err)
			w.logger.Fatal("Failed to create LogFile at MarshalIndent")
		}
		logFile.WriteString(string(jsonStruct))
		return logFile
	} else {
		file, err := os.OpenFile(path, os.O_WRONLY, 0644)
		if err != nil {
			w.logger.Error(err)
			w.logger.Fatal("Failed to create LogFile at os.OpenFile")
		}
		logFile := &LogFile{
			File: file,
		}
		filecontent, err := ioutil.ReadFile(path)
		if err != nil {
			w.logger.Error(err)
			w.logger.Fatal("Failed to create LogFile at ioutil.ReadFile")
		}
		err = json.Unmarshal(filecontent, &logFile)
		if err != nil {
			w.logger.Error(err)
			w.logger.Fatal("Failed to create LogFile at Unmarshal")
		}

		// this should not be done
		jsonStruct, err := json.MarshalIndent(&logFile, "", "  ")
		if err != nil {
			w.logger.Error(err)
			w.logger.Fatal("Failed to create LogFile at MarshalIndent (2)")
		}
		logFile.WriteString(string(jsonStruct))

		return logFile
	}
}

const (
	sensorLogFile = "/var/supplywatch/udpserver/sensorlog"
)

func (w *warehouse) AllSensorLogs() (logs []byte, err error) {
	sensorLog, err := os.ReadFile(sensorLogFile)
	if err != nil {
		return nil, err
	}
	return sensorLog, nil
}

// GetAllSensorLogs goes over all log files in the LogFileDir defined in config.yml
// and puts them together
func (w *warehouse) GetAllSensorLogs(path string) (logs []byte, err error) {
	allLogFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	var allLogs []LogEntry
	for _, f := range allLogFiles {
		var logfile LogFile
		filename := filepath.Join(path, f.Name())
		jsonLogFile, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer jsonLogFile.Close()
		logs, _ = ioutil.ReadAll(jsonLogFile)
		err = json.Unmarshal(logs, &logfile)
		for _, log := range logfile.Logs {
			allLogs = append(allLogs, log)
		}
	}
	JsonBytes, err := json.MarshalIndent(&allLogs, "", "  ")
	if err != nil {
		return nil, err
	}
	return JsonBytes, nil
}

// GetAllSensorLogs goes over all log files in the LogFileDir defined in config.yml
// and puts all logs together that match the specified sensorID
func GetOneSensorLogs(path string, sensorID string) (logs []byte, err error) {
	allLogFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	var allLogs []LogEntry
	for _, f := range allLogFiles {
		var logfile LogFile
		filename := filepath.Join(path, f.Name())
		jsonLogFile, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer jsonLogFile.Close()
		logs, _ = ioutil.ReadAll(jsonLogFile)
		err = json.Unmarshal(logs, &logfile)
		for _, log := range logfile.Logs {
			if log.SensorID == sensorID {
				allLogs = append(allLogs, log)
			}
		}
	}
	JsonBytes, err := json.MarshalIndent(&allLogs, "", "  ")
	if err != nil {
		return nil, err
	}
	return JsonBytes, nil
}

//func ReadOneSensorLogs(file string, SensorID string) ([]byte, error) {
//jsonLogFile, err := os.Open(file)
//if err != nil {
//return nil, err
//}
//jsonLogs, _ := ioutil.ReadAll(jsonLogFile)
//var logfile LogFile
//var filteredLogs []LogEntry
//err = json.Unmarshal(jsonLogs, &logfile)
//if err != nil {
//panic("failed to unmarshal logfile")
//}
//for _, log := range logfile.Logs {
//if log.SensorID.String() == SensorID {
//filteredLogs = append(filteredLogs, log)
//}
//}

//JsonBytes, err := json.MarshalIndent(&filteredLogs, "", "  ")
//if err != nil {
//return nil, err
//}
//return JsonBytes, nil
//}

func ReadLogsFromDate(file string) ([]byte, error) {
	jsonLogFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	jsonLogs, _ := ioutil.ReadAll(jsonLogFile)
	var logfile LogFile
	err = json.Unmarshal(jsonLogs, &logfile)
	if err != nil {
		panic("failed to unmarshal logfile")
	}
	JsonBytes, err := json.MarshalIndent(&logfile.Logs, "", "  ")
	if err != nil {
		return nil, err
	}
	return JsonBytes, nil
}

func (l *LogFile) WriteToFile() error {
	jsonStruct, err := json.MarshalIndent(&l, "", "  ")
	if err != nil {
		return err
	}
	_, err = l.WriteAt(jsonStruct, 0)
	if err != nil {
		log.Fatal("WriteAt error")
	}
	return nil
}

func (l *LogFile) UpdateLogFile() error {
	return nil
}
