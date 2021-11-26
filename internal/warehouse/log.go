package warehouse

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/google/uuid"
)

// LogEntry represents a new entry in the log file
type LogEntry struct {
	SensorID   uuid.UUID `json:"sensor_id"`
	SensorType string    `json:"sensor_type"`
	Message    string    `json:"message"`
	IP         net.IP    `json:"ip"`
	Port       int       `json:"port"`
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
func NewLogFile(path string) *LogFile {
    if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
        file, err := os.Create(path)
        if err != nil {
            panic("Failed to create LogFile")
        }
        logFile := &LogFile{
           File: file,
        }

        jsonStruct, err := json.MarshalIndent(&logFile, "", "  ")
        if err != nil {
            panic("Failed to create LogFile")
        }
        logFile.WriteString(string(jsonStruct))
        return logFile
    } else {
        file, err := os.OpenFile(path, os.O_WRONLY, 0644) 
        if err != nil {
            panic("Failed to open LogFile")
        }
        logFile := &LogFile{
            File: file,
        }
        filecontent, err := ioutil.ReadFile(path)
        if err != nil {
            panic("Failed to read LogFile")
        }
        err = json.Unmarshal(filecontent, &logFile)
        if err != nil {
            panic("Failed to unmarshal LogFile")
        }

        // this should not be done
        jsonStruct, err := json.MarshalIndent(&logFile, "", "  ")
        if err != nil {
            panic("Failed to create LogFile")
        }
        logFile.WriteString(string(jsonStruct))

        return logFile
    }
}

func ReadAllLogs(file string) (logs []byte, err error) {
	jsonLogFile, err := os.Open(file)
	if err != nil {
		return
	}
	defer jsonLogFile.Close()
	logs, _ = ioutil.ReadAll(jsonLogFile)
	return logs, nil
}

func ReadOneSensorLogs(file string, SensorID string) ([]byte, error) {
	jsonLogFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	jsonLogs, _ := ioutil.ReadAll(jsonLogFile)
	var logfile LogFile
	var filteredLogs []LogEntry
	err = json.Unmarshal(jsonLogs, &logfile)
	if err != nil {
		panic("failed to unmarshal logfile")
	}
	for _, log := range logfile.Logs {
		if log.SensorID.String() == SensorID {
			filteredLogs = append(filteredLogs, log)
		}
	}

	JsonBytes, err := json.MarshalIndent(&filteredLogs, "", "  ")
	if err != nil {
		return nil, err
	}
	return JsonBytes, nil
}

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
    b, err := l.WriteAt(jsonStruct, 0)
	if err != nil {
        fmt.Println(err)
        fmt.Println(b)
        log.Fatal("writeat error")
	}
	return nil
}

func (l *LogFile) UpdateLogFile() error {
	return nil
}
