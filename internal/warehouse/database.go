package warehouse

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func initDB()(*gorm.DB, error){
    dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%d", dbHost, dbUser, dbName, dbPassword, dbPort)
    db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
    return db, err
}

func setupProducts(db *gorm.DB) error {
    products := []string{
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
    for _, product := range products {
        newProduct := &Product{Name: product, ID: uuid.New(), Quantity: 0}
        db.FirstOrCreate(&Products, &newProduct)
        Products = append(Products, newProduct)
    }
    return nil
}
