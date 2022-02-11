package main

import (
	"os"

	"github.com/mariuskimmina/supplywatch/internal/monitor"
	"github.com/mariuskimmina/supplywatch/pkg/config"
	"github.com/mariuskimmina/supplywatch/pkg/log"
)

func main() {
	logger := log.NewLogger()
	host, err := os.Hostname()
	if err != nil {
		panic("Failed to get hostname")
	}
	logger.Infof("Warehouse starting, hostname: %s", host)
	config, err := config.LoadSupplywatchConfig("./configurations")
	if err != nil {
		logger.Error(err)
		logger.Fatal("Failed to load warehouse configuration")
	}
	monitor := monitor.NewMonitor(logger, config)
	monitor.RunAndServe()
}
