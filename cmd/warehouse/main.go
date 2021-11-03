package main

import (
	"github.com/mariuskimmina/supplywatch/internal/warehouse"
	log "github.com/mariuskimmina/supplywatch/pkg/log"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := log.NewLogger()
	logger.SetLevel(logrus.InfoLevel)
	logger.Info("Starting Warehouse")
	warehouse := warehouse.NewWarehouse(logger)
	warehouse.Start()
}
