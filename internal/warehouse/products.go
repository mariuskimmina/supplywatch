package warehouse

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/uuid"
)

type Product struct {
    ID       uuid.UUID
    Name     string
    Count           int
    lastReceived    string
    lastDispatched  string
}

var Products []*Product


func GetorCreateProduct(name string) (*Product, error) {
    var product *Product
    productExists := false

    for _, product := range Products {
        if product.Name == name {
            product.Increment()
            productExists = true
            break
        }
    }

    if !productExists {
        product, err := NewProduct(name)
        if err != nil {
            return product, err
        }
        Products = append(Products, product)
    }

    return product, nil
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
        Count: 1,
        lastReceived: time.Now().Format("01-02-2006"),
    }, nil
}

func (p *Product) Increment() {
    p.Count += 1
    p.lastReceived = time.Now().Format("01-02-2006")
}

func (p *Product) Decrement() {
    p.Count -= 1
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

