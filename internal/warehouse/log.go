package warehouse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type LogFile struct {
    *os.File
    Logs []LogEntry
}

func (l *LogFile) addLog(log LogEntry) {
    l.Logs = append(l.Logs, log)
}

func NewLogFile(path string) (*LogFile) {
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
}

func ReadAllLogs() (logs []byte, err error) {
    jsonLogFile, err := os.Open("/tmp/warehouselog")
    if err != nil {
        return
    }
    defer jsonLogFile.Close()
    logs, _ = ioutil.ReadAll(jsonLogFile)
    return logs, nil
}

func ReadOneSensorLogs(SensorID string) ([]byte, error) {
    jsonLogFile, err := os.Open("/tmp/warehouselog")
    if err != nil {
        return nil, err
    }
    jsonLogs, _ := ioutil.ReadAll(jsonLogFile)
    var logfile LogFile
    var filteredLogs []LogEntry
    err = json.Unmarshal(jsonLogs, &logfile)
    if err != nil {
        panic("FAILED")
    }
    for _, log := range logfile.Logs {
        fmt.Println("comparing")
        fmt.Println(log.SensorID.String())
        fmt.Println(SensorID)
        if log.SensorID.String() == SensorID {
            fmt.Println("match")
            filteredLogs = append(filteredLogs, log)
        }
    }

    JsonBytes, err := json.MarshalIndent(&filteredLogs, "", "  ")
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
    l.WriteAt(jsonStruct, 0)
    return nil
}

func (l *LogFile) UpdateLogFile() error {
    return nil
}
