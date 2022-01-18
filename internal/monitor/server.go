package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/mariuskimmina/supplywatch/internal/domain"
	"github.com/mariuskimmina/supplywatch/pkg/config"
)

type monitor struct {
	logger Logger
    config config.SupplywatchConfig
}

// Logger is a generic interface that can be implemented by any logging engine
// this allows for dependency injection which results in easier testing
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}


func NewMonitor(logger Logger, config config.SupplywatchConfig) *monitor {
	return &monitor{
		logger: logger,
        config: config,
	}
}

func (s *monitor) RunAndServe() {
    //time.Sleep(30 * time.Second)
	var wg sync.WaitGroup
	wg.Add(1)
	// RabbitMQ
	go func() {
		s.SetupMessageQueue(s.config.NumOfWarehouses)
		wg.Done()
	}()
    for i := 1; i <= s.config.NumOfWarehouses; i++ {
        path := "/w" + strconv.Itoa(i)
        http.HandleFunc(path, warehousedata)
    }
    http.HandleFunc("/", s.overview)
    http.ListenAndServe(":9000", nil)
    wg.Wait()
}

func (s *monitor) overview(w http.ResponseWriter, req *http.Request) {
    NumOfWarehouses := strconv.Itoa(s.config.NumOfWarehouses)
    fmt.Fprintf(w, "Welcome to the Monitor!\n\nNumber of Warehouses: %s", NumOfWarehouses)
}

func warehousedata(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    path := req.URL.Path
    number := path[len(path) - 1:]
	productsFileName := "/var/supplywatch/monitor/products-warehouse" + number + "Data"
	if _, err := os.Stat(productsFileName); errors.Is(err, os.ErrNotExist) {
        fmt.Print(w, err)
	}
	productsFile, err := os.Open(productsFileName)
	if err != nil {
        fmt.Fprintf(w, err.Error())
	}
	defer productsFile.Close()
	jsonProducts, err := ioutil.ReadAll(productsFile)
	if err != nil {
        fmt.Fprintf(w, err.Error())
	}
    var products []domain.Producttype
    err = json.Unmarshal(jsonProducts, &products)
	if err != nil {
        fmt.Fprintf(w, err.Error())
	}
    allProductsJson, err := json.MarshalIndent(&products, " ", "")
	if err != nil {
        fmt.Fprintf(w, err.Error())
	}
    fmt.Fprintf(w, string(allProductsJson))
}
