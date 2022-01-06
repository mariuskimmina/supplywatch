package warehouse

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mariuskimmina/supplywatch/pkg/backoff"
)

var (
	dbHost     = os.Getenv("SW_DATABASE")
	dbUser     = "supplywatch"
	dbName     = "supplywatch"
	dbPassword = "test"
)

const (
	dbPort = 5432
)

func (w *warehouse) initDB() {
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%d", dbHost, dbUser, dbName, dbPassword, dbPort)
	var attempt int
	for {
		time.Sleep(backoff.Default.Duration(attempt))
		db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
		if err != nil {
			w.logger.Error(err)
			w.logger.Error("Failed to connect to Database, going to retry")
			attempt++
			continue
		}
		w.logger.Info("Successfully connected to Database")
		db.AutoMigrate(&Product{})
		w.DB = db
		break
	}
}

