package warehouse

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/mariuskimmina/supplywatch/pkg/backoff"
)


var (
    dbHost = os.Getenv("SW_DATABASE")
    dbUser = "supplywatch"
    dbName = "supplywatch"
    dbPassword = "test"
)

const (
    dbPort = 5432
)

func (w *warehouse) initDB()(){
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

func (w *warehouse) setupProducts() error {
    products := []string{
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
    ids := []string{
        "24476e37-558d-4068-8a76-fd0af990465f",
        "470de33a-c08f-4867-b316-cecce280a37a",
        "344f8e1d-bb13-4a85-86af-3a201193ef8b",
        "4ad03013-21f8-4d19-a4d7-eccdff210b16",
        "27a40683-392c-496f-8022-9ee4ff119aaa",
        "2ac00095-2a7d-416c-a250-8c776e1476a7",
        "bacd7beb-061f-444c-82bf-18454d0cf0ca",
        "d9eb4c17-1a87-44a4-98f3-dcb6923b485f",
        "39cc2513-66be-462d-99d6-969951ac1f93",
        "eff79513-5425-4b9c-8ff5-b597f7a7f67e",
    }
    for index, _ := range products {
        id, err := uuid.Parse(ids[index])
        if err != nil {
            return err
        }
        newProduct := &Product{Name: products[index], ID: id, Quantity: 5}
        w.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&newProduct)
        //w.DB.Create(&newProduct)
        Products = append(Products, newProduct)
    }
    return nil
}
