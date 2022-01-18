package warehouse

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mariuskimmina/supplywatch/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Product struct {
	gorm.Model
	ID             uuid.UUID
	Name           string
	Quantity       int
	lastReceived   string
	lastDispatched string
}

var Products []*Product

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
		//"walnuts",
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
		//"eff79513-5425-4b9c-8ff5-b597f7a7f67e",
	}
	for index := range products {
        w.logger.Infof("Setting up Product %d out of %d \n", index, len(products) - 1)
		id, err := uuid.Parse(ids[index])
		if err != nil {
            w.logger.Error("Error Setting up Products")
			return err
		}

		newProduct := &Product{Name: products[index], ID: id, Quantity: 3}
        w.logger.Info(newProduct.ID.String())
        result := w.DB.Clauses(clause.OnConflict{
            Columns: []clause.Column{{Name: "id"}},
            DoUpdates: clause.AssignmentColumns([]string{"quantity"}),
        }).Create(&newProduct)
        //result := w.DB.FirstOrCreate(&newProduct, Product{ID: id})
        w.logger.Info(result.Error)
        //result := w.DB.Create(&newProduct)
        //if err != nil {
            //w.logger.Error("Error during database operation")
            //w.logger.Error(err)
        //}
		//w.DB.Create(&newProduct)
		Products = append(Products, newProduct)
        w.logger.Infof("Done Setting up Product %d out of %d \n", index, len(products) - 1)

	}
    w.logger.Info("Done setting up Products")
	return nil
}

func (w *warehouse) HandleProduct(product *domain.InOutProduct, storageChan chan<- string) error {
	if product.Incoming {
		w.IncrementorCreateProduct(product.ProductName)
	} else {
		w.DecrementProduct(product.ProductName, storageChan)
	}
	return nil
}

func (w *warehouse) AllProducts() (domain.Products, error) {
	var products []*domain.Product
	w.DB.Find(&products)
	return products, nil
}

func (w *warehouse) SensorLogByID() error {
	return nil
}

func (w *warehouse) ProductByID() error {
	return nil
}

func (w *warehouse) IncrementorCreateProduct(name string) error {
	//w.logger.Infof("Incrementing quantity of %s", name)
	productExists := false

	for _, product := range Products {
		if product.Name == name {
			oldquantity := product.Quantity
			product.Increment()
			w.logger.Infof("Icrementing database quantity of %s from %d to %d", name, oldquantity, product.Quantity)
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

func (w *warehouse) DecrementProduct(name string, storageChan chan<- string) error {
	//w.logger.Infof("Decrementing quantity of %s", name)
	productExists := false

	for _, product := range Products {
		if product.Name == name {
			oldquantity := product.Quantity
			product.Decrement()
			w.logger.Infof("Decrementing database quantity of %s from %d to %d", name, oldquantity, product.Quantity)
			w.DB.Model(&Product{}).Where("name = ?", product.Name).Update("quantity", product.Quantity)
			productExists = true
			if product.Quantity == 0 {
				// if the product quantity is zero, we send a request to the message queue so that another warehouse is hopefully
				// going to send this product, the message send has to have the format productname:productid:hostname
				pname := product.Name
				pid := product.ID.String()
				host, err := os.Hostname()
				if err != nil {
					w.logger.Fatal("Failed to get hostname")
				}
				storageChan <- pname + ":" + pid + ":" + host
			}
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
		ID:           id,
		Name:         name,
		Quantity:     1,
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
	//fmt.Println("Saving Products")
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
	//fmt.Println("Loading Products")
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
