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
	var connRabbit *amqp.Connection
	var err error

	var attempt int
	for {
		time.Sleep(backoff.Default.Duration(attempt))
		connRabbit, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
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
        queueName := "warehouse" + fmt.Sprint(i) + "DataExchange"
        go func() {
            s.SubtoWarehouseData(channel, queueName)
            wg.Done()
        }()
    }
    wg.Wait()
}

func (s *monitor) SubtoWarehouseData(channel *amqp.Channel, exchangeName string) {
    s.logger.Infof("Subbing to data exchange !! %s !!", exchangeName)
    err := os.MkdirAll(logFileDir, 0644)
	if err != nil {
        fmt.Println("doof")
	}
	f, err := os.OpenFile(logFileDir+logFilePrefix+exchangeName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
        fmt.Println("doof")
	}
	defer f.Close()

    q, err := channel.QueueDeclare(
        "",
        false,
        false,
        true,
        false,
        nil,
    )
    if err != nil {
        s.logger.Error("Failed to declare queue on exchange")
    }

    err = channel.QueueBind(
        q.Name,
        "",
        exchangeName,
        false,
        nil,
    )
    if err != nil {
        s.logger.Error("Failed to bind to queue to subscribe")
    }

    msgs, err := channel.Consume(
        q.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        s.logger.Error("Failed to register as a consumer here")
    }

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        s.logger.Info("Waiting for messages here")
        for d := range msgs {
            s.logger.Infof("Received Update from %s", exchangeName)
            //s.logger.Info(string(d.Body))
            var allProducts []domain.Producttype
            err := json.Unmarshal(d.Body, &allProducts)
            //fmt.Println(allProducts)
            if err != nil {
                //TODO: move the file writing somewhere else
                fmt.Println("doof")
            }
            allProductsJson, err := json.MarshalIndent(allProducts, " ", "")
            if err != nil {
                //TODO: move the file writing somewhere else
                fmt.Println("doof")
            }
            f.Truncate(0)
            f.Seek(0, 0)
            f.Write(allProductsJson)
            //productsFileName := logFileDir + logFilePrefix + queueName
            //err = ioutil.WriteFile(productsFileName, jsonProducts, 0644)
            if err != nil {
                //TODO: move the file writing somewhere else
                fmt.Println("doof")
            }
        }
        wg.Done()
        s.logger.Info("MONITOR IS DONE LISTENING")
    }()
    wg.Wait()
}
