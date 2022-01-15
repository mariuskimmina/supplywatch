package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/mariuskimmina/supplywatch/pkg/backoff"
	"github.com/streadway/amqp"
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
	channel, err := connRabbit.Channel()
    if err != nil {
        s.logger.Error(err)
        s.logger.Error("Failed to setup Channel with RabbitMQ")
    }
	s.logger.Info("Successfully setup Channel with RabbitMQ")
	defer connRabbit.Close()
    defer channel.Close()
    //for {
        //time.Sleep(10 *time.Second)
        //s.logger.Info("KEEPING RABBITMQ GOING")
    //}

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

	//listen to incoming messages
	go func() {
		for d := range msgs {
			s.logger.Info("--------------------Receiving from Queue!-------------------")
			s.logger.Infof("Received Message: %s from queue", string(d.Body))
			//incoming messages should have the following format: product:hostname
			//if !strings.Contains(string(d.Body), ":") {
				//s.logger.Error("Received Message from queue with invalid format, this message will be ignored")
				//continue
			//}
			//s.logger.Info("Initalizing transport of products")
			//sendChan <- string(d.Body)
		}
		wg.Done()
	}()

	wg.Wait()
}
