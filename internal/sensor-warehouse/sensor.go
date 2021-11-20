package sensorwarehouse

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mariuskimmina/supplywatch/pkg/config"
	"github.com/mariuskimmina/supplywatch/pkg/log"
)

type Sensor struct {
	logger *log.Logger
    config *config.Config

}

type Message struct {
    SensorID    uuid.UUID `json:"sensor_id"`
	SensorType  string `json:"sensor_type"`
	MessageBody string `json:"message"`
}

func NewSensor(logger *log.Logger, config *config.Config) *Sensor {
	return &Sensor{
		logger: logger,
        config: config,
	}
}

var(
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
    sensorID, err := uuid.NewUUID()
	if err != nil {
		s.logger.Fatal("Failed to create ID for sensor")
	}
    warehousePort :=strconv.Itoa(s.config.SensorWarehouse.UDPPort)
    warehouseAdr := "warehouse" + ":" +  warehousePort
	conn, err := net.Dial("udp", warehouseAdr)
	if err != nil {
		s.logger.Error("Failed to dial the warehouse")
	}

	for {
        n = rand.Int() % len(products)
		product := products[n]
		message := Message{
            SensorID: sensorID,
			SensorType:  sensorType,
			MessageBody: product,
		}
		jsonMessage, err := json.Marshal(message)

		if err != nil {
			s.logger.Error("Failed to convert message to json")
		}
		s.logger.Infof("Sending %s", product)
		conn.Write([]byte(jsonMessage))
		time.Sleep(5 * time.Second)
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
