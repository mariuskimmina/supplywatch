package main

import (
	"github.com/mariuskimmina/supplywatch/internal/warehouse"
	log "github.com/mariuskimmina/supplywatch/pkg/log"
	"github.com/mariuskimmina/supplywatch/pkg/config"
)

func main() {
	logger := log.NewLogger()
    config, err := config.LoadConfig(".")
    if err != nil {
        logger.Error(err)
        logger.Fatal("Failed to load warehouse configuration")
    }
    warehouse := warehouse.NewWarehouse(logger, &config)
	warehouse.Start()
}
