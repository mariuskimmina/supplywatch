package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mariuskimmina/supplywatch/internal/tcp"
	"github.com/mariuskimmina/supplywatch/internal/warehouse"
	"github.com/mariuskimmina/supplywatch/pkg/backoff"
	"github.com/mariuskimmina/supplywatch/pkg/config"
	"github.com/mariuskimmina/supplywatch/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	logger log.Logger
	dbHost string
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	host, err := os.Hostname()
	logger := log.NewLogger()
	logger.Infof("Warehouse starting, hostname: %s", host)
	config, err := config.LoadWarehouseConfig("./configurations")
	if err != nil {
		logger.Error(err)
		logger.Fatal("Failed to load warehouse configuration")
	}

	// to find out which warehouse should connect to which database at runtime we use the last char of the Hostname
	// we do this so that we do not have to have a seperate configuration file for each warehouse instance
	if strings.Contains(host, "warehouse1") {
		dbHost = "database1"
        logger.Infof("Trying to Connect to: %s", dbHost)
	}
	if strings.Contains(host, "warehouse2") {
		dbHost = "database2"
        logger.Infof("Trying to Connect to: %s", dbHost)
	}
	if dbHost == "" {
		logger.Fatal("Failed to determine database hostname")
	}

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%d", dbHost, config.DBUser, config.DBDatabase, config.DBPassword, config.DBPort)
	logger.Infof(dbURI)
	db := dbConnect(dbURI)
	warehouse := warehouse.NewWarehouse(logger, &config, db)
    go func() {
        logger.Info("Starting Warehouse")
        warehouse.Start()
        wg.Done()
    }()
    tcpServer, err := tcp.NewTCPServer(warehouse)
    go func() {
        logger.Info("Starting TCP Server")
        tcpServer.Listen()
        wg.Done()
    }()
    wg.Wait()
}

func dbConnect(dbURI string) *gorm.DB {
	var attempt int
	var db *gorm.DB
	var err error
	for {
		time.Sleep(backoff.Default.Duration(attempt))
		db, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to connect to Database, going to retry")
			attempt++
			continue
		}
		fmt.Println("Successfully connected to Database")
		break
	}
	return db
}
