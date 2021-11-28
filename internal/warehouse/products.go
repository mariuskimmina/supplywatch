package warehouse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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


func SaveProductsState() {
    fmt.Println("Saving Products")
    jsonProducts, err := json.MarshalIndent(Products, "", "  ")
    if err != nil {
        log.Fatal("Failed to marhal products - cannot save products")
    }
    err = ioutil.WriteFile("/var/supplywatch/log/products.json", jsonProducts, 0644)
    if err != nil {
        log.Fatal("Failed to write products to file - cannot save products")
    }
}

func LoadProductsState() {

}

