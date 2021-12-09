package warehouse

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
    gorm.Model
    ID       uuid.UUID
    Name     string
    Quantity           int
    lastReceived    string
    lastDispatched  string
}

var Products []*Product


func (w *warehouse) IncrementorCreateProduct(name string) (error) {
    w.logger.Infof("Incrementing quantity of %s", name)
    productExists := false

    for _, product := range Products {
        if product.Name == name {
            product.Increment()
            w.logger.Infof("Updating database quantity of %s to %d", name, product.Quantity)
            w.DB.Model(&Product{}).Where("name = ?", product.Name).Update("quantity", product.Quantity)
            productExists = true
            break
        }
    }

    if !productExists {
        w.logger.Fatalf("Could not find a product with name: %s", name)
    }
    return nil
}

func (w *warehouse) DecrementProduct(name string) (error) {
    w.logger.Infof("Decrementing quantity of %s", name)
    productExists := false

    for _, product := range Products {
        if product.Name == name {
            product.Decrement()
            w.logger.Infof("Updating database quantity of %s to %d", name, product.Quantity)
            w.DB.Model(&Product{}).Where("name = ?", product.Name).Update("quantity", product.Quantity)
            productExists = true
            break
        }
    }

    if !productExists {
        w.logger.Fatalf("Could not find a product with name: %s", name)
    }
    return nil
}

func GetAllProductsAsBytes() ([]byte, error) {
	JsonBytes, err := json.MarshalIndent(&Products, "", "  ")
    if err != nil {
        return nil, err
    }
    return JsonBytes, nil
}

func NewProduct(name string) (*Product, error) {
    id := uuid.New()
    return &Product{
        ID: id,
        Name: name,
        Quantity: 1,
        lastReceived: time.Now().Format("01-02-2006"),
    }, nil
}

func (p *Product) Increment() {
    p.Quantity += 1
    p.lastReceived = time.Now().Format("01-02-2006")
}

func (p *Product) Decrement() {
    p.Quantity -= 1
    p.lastDispatched = time.Now().Format("01-02-2006")
}


func SaveProductsState() error {
    fmt.Println("Saving Products")
    jsonProducts, err := json.MarshalIndent(Products, "", "  ")
    if err != nil {
        return err
    }
    hostname, err := os.Hostname()
    if err != nil {
        return err
    }
    productsFileName := "/var/supplywatch/log/" + hostname + "-products.json"
    err = ioutil.WriteFile(productsFileName, jsonProducts, 0644)
    if err != nil {
        return err
    }
    return nil
}

func LoadProductsState() error {
    fmt.Println("Loading Products")
    hostname, err := os.Hostname()
    if err != nil {
        return err
    }
    productsFileName := "/var/supplywatch/log/" + hostname + "-products.json"
    if _, err := os.Stat(productsFileName); errors.Is(err, os.ErrNotExist) {
        return nil
    }
    productsFile, err := os.Open(productsFileName)
    if err != nil {
        return err
    }
    defer productsFile.Close()
    jsonProducts, err := ioutil.ReadAll(productsFile)
    if err != nil {
        return err
    }
    json.Unmarshal(jsonProducts, &Products)
    return nil
}

