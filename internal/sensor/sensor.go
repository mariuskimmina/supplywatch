package sensor

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/mariuskimmina/supplywatch/internal/domain"
	"github.com/mariuskimmina/supplywatch/pkg/backoff"
	"github.com/mariuskimmina/supplywatch/pkg/config"
	"github.com/mariuskimmina/supplywatch/pkg/log"
)

type Sensor struct {
	logger *log.Logger
	config *config.SensorConfig
}

type Message struct {
	SensorID    string `json:"sensor_id"`
	SensorType  string `json:"sensor_type"`
	MessageBody string `json:"message"`
	Incoming    bool   `json:"incoming"`
}

func NewSensor(logger *log.Logger, config *config.SensorConfig) *Sensor {
	return &Sensor{
		logger: logger,
		config: config,
	}
}

const (
    logFileDir = "/var/supplywatch/udpserver/"
    logFile = "sensorlog"
)

var (
	products = []string{
		"butter",
		"sugar",
		"eggs",
		"baking powder",
		"cheese",
		"cinnamon",
		"oil",
		"carrots",
		"raisins",
		"walnuts",
	}
	SensorType = []string{
		"BarcodeReader",
		"RFID-Reader",
	}
	Incoming = []bool{
		true,
		false,
	}
)

func (s *Sensor) Start() {
	var err error
	var attempt int
	var conn net.Conn

    err = os.MkdirAll(logFileDir, 0644)
    if err != nil {
        s.logger.Fatalf("Error creating log file directory: %v", err)
    }
    f, err := os.OpenFile(logFileDir + logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        s.logger.Fatalf("Error opening log file: %v", err)
    }
    defer f.Close()

	SeedRandom()
	n := rand.Int() % len(SensorType)
	sensorType := SensorType[n]
	sensorID, err := os.Hostname()
	if err != nil {
		s.logger.Fatal("Failed to create ID for sensor")
	}
	warehouses := []string{
		"warehouse1:4444",
		"warehouse2:4444",
	}
	n = rand.Int() % len(warehouses)
	warehouseAdr := warehouses[n]

	for {
		time.Sleep(backoff.Default.Duration(attempt))
		conn, err = net.Dial("udp", warehouseAdr)
		if err != nil {
			s.logger.Info("Failed to dial the warehouse, going to retry")
			s.logger.Error(err)
			attempt++
			continue
		}
		break
	}

	var packetCounter = 0
	for {
		n = rand.Int() % len(products)
		incoming := IncomingOrOutgoing()
		product := products[n]
		message := Message{
			SensorID:    sensorID,
			SensorType:  sensorType,
			MessageBody: product,
			Incoming:    incoming,
		}
		jsonMessage, err := json.Marshal(message)

		if err != nil {
			s.logger.Error("Failed to convert message to json")
		}
		s.logger.Infof("Sending %s, Incoming: %v", product, incoming)
		conn.Write([]byte(jsonMessage))

		logentry := &domain.SensorLog{
			SensorType: message.SensorType,
			SensorID:   message.SensorID,
			Message:    message.MessageBody,
			Incoming:   message.Incoming,
		}
        logjson, err := json.Marshal(logentry)
		if err != nil {
            s.logger.Fatalf("Error marshaling log data: %v", err)
		}
        f.WriteString(string(logjson) + ",\n")

		// If NumOfPackets is not 0 we stop sending once the NumOfPackets has been reached
		if s.config.NumberOfPackets != 0 {
			packetCounter += 1
			if packetCounter == s.config.NumberOfPackets {
				break
			}
		}
		time.Sleep(time.Duration(s.config.Delay) * time.Millisecond)
	}
}

// SeedRandom makes sure that multiple sensors will send different random products
// using the simple time.Now() seeding did not work for this case as both containers start at the same time
func SeedRandom() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

// IncomingOrOutgoing creates a "random" boolean that is highly favoured to be true
// this way we get a lot more incoming products than outgoing products which leads to less negative numbers
func IncomingOrOutgoing() bool {
	return rand.Float32() > 0.5
}
