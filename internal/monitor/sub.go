package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/mariuskimmina/supplywatch/internal/domain"
	"github.com/mariuskimmina/supplywatch/pkg/backoff"
	"github.com/streadway/amqp"
)

const (
	logFileDir    = "/var/supplywatch/monitor/"
	logFilePrefix       = "products-"
)

func (s *monitor) SetupMessageQueue(numOfWarehouses int) {
	//var wg sync.WaitGroup
	//wg.Add(2)
	var connRabbit *amqp.Connection
	var err error

	var attempt int
	for {
		time.Sleep(backoff.Default.Duration(attempt))
		connRabbit, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			//s.logger.Error(err)
			attempt++
			s.logger.Infof("Failed to connect to RabbitMQ, trying again in %f seconds",
                backoff.Default.Duration(attempt).Seconds())
			continue
		}
		break
	}
	s.logger.Info("Successfully Connected to RabbitMQ")
	defer connRabbit.Close()
	channel, err := connRabbit.Channel()
    if err != nil {
        s.logger.Error(err)
        s.logger.Error("Failed to setup Channel with RabbitMQ")
    }
	s.logger.Info("Successfully setup Channel with RabbitMQ")
    defer channel.Close()

    var wg sync.WaitGroup
    wg.Add(numOfWarehouses)
    s.logger.Info("Number Of Warehouses: ", numOfWarehouses)
    for i := 1; i <= numOfWarehouses; i++ {
        queueName := "warehouse" + fmt.Sprint(i) + "Data"
        go func() {
            s.subscribeToWarehouseData(channel, queueName)
            wg.Done()
        }()
    }
    wg.Wait()
}

func (s *monitor) subscribeToWarehouseData(c *amqp.Channel, queueName string) {
    s.logger.Infof("Subscribing to Queue: %s", queueName)
	var wg sync.WaitGroup
	wg.Add(1)
	var msgs <-chan amqp.Delivery
	var err error

	err = os.MkdirAll(logFileDir, 0644)
	if err != nil {
        fmt.Println("doof")
	}
	f, err := os.OpenFile(logFileDir+logFilePrefix+queueName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
        fmt.Println("doof")
	}
	defer f.Close()

	msgs, err = c.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.logger.Error(err)
		s.logger.Fatal("Failed to subscribe to a Testmessage")
	}
	s.logger.Infof("Subscribed to %s queue!", queueName)
	//forever := make(chan bool)

    //decoder := json.NewDecoder()
	//listen to incoming messages
	go func() {
		for d := range msgs {
			s.logger.Info("Receiving Status Update from a Warehouse")
			s.logger.Infof("Received Message: %s from queue", string(d.Body))
            var allProducts []domain.Producttype
            err := json.Unmarshal(d.Body, &allProducts)
            fmt.Println(allProducts)
            if err != nil {
                //TODO: move the file writing somewhere else
                fmt.Println("doof")
            }
            allProductsJson, err := json.MarshalIndent(allProducts, " ", "")
            if err != nil {
                //TODO: move the file writing somewhere else
                fmt.Println("doof")
            }
            f.Write(allProductsJson)
            //productsFileName := logFileDir + logFilePrefix + queueName
            //err = ioutil.WriteFile(productsFileName, jsonProducts, 0644)
            if err != nil {
                //TODO: move the file writing somewhere else
                fmt.Println("doof")
            }
			//sendChan <- string(d.Body)
		}
		wg.Done()
	}()

	wg.Wait()
}
