package sensor

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/mariuskimmina/supplywatch/pkg/config"
	"github.com/mariuskimmina/supplywatch/pkg/log"
)

type Sensor struct {
	logger *log.Logger
	config *config.Config
}

type Message struct {
	SensorID    string `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	MessageBody string    `json:"message"`
}

func NewSensor(logger *log.Logger, config *config.Config) *Sensor {
	return &Sensor{
		logger: logger,
		config: config,
	}
}

var (
	products = []string{
		"butter",
		"sugar",
		"eggs",
		"baking powder",
		"cheese",
		"lemons",
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
)

func (s *Sensor) Start() {
	SeedRandom()
	n := rand.Int() % len(SensorType)
	sensorType := SensorType[n]
	sensorID, err := os.Hostname()
	if err != nil {
		s.logger.Fatal("Failed to create ID for sensor")
	}
	warehouses := []string{
		"warehouse_1:4444",
		"warehouse_2:4445",
	}
	n = rand.Int() % len(warehouses)
	warehouseAdr := warehouses[n]
	conn, err := net.Dial("udp", warehouseAdr)
	if err != nil {
		s.logger.Error("Failed to dial the warehouse")
	}

	var packetCounter = 0
	for {
		n = rand.Int() % len(products)
		product := products[n]
		message := Message{
			SensorID:    sensorID,
			SensorType:  sensorType,
			MessageBody: product,
		}
		jsonMessage, err := json.Marshal(message)

		if err != nil {
			s.logger.Error("Failed to convert message to json")
		}
		s.logger.Infof("Sending %s", product)
		conn.Write([]byte(jsonMessage))

		// If NumOfPackets is not 0 we stop sending once the NumOfPackets has been reached
		if s.config.SensorWarehouse.NumOfPackets != 0 {
			packetCounter += 1
			if packetCounter == s.config.SensorWarehouse.NumOfPackets {
				break
			}
		}

		time.Sleep(time.Duration(s.config.SensorWarehouse.Delay) * time.Millisecond)
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
