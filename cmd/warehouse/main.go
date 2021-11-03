package main

import (
	"github.com/mariuskimmina/supplywatch/internal/warehouse"
	log "github.com/mariuskimmina/supplywatch/pkg/log"
)

func main() {
	logger := log.NewLogger()
	warehouse := warehouse.NewWarehouse(logger)
	warehouse.Start()
}
