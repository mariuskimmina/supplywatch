package main

import "github.com/mariuskimmina/supplywatch/internal/warehouse"


func main() {
    warehouse := warehouse.NewWarehouse()
    warehouse.Start()
}
