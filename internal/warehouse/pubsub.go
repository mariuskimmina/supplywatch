package warehouse

import (
	"strings"
	"sync"
	"time"

	"github.com/mariuskimmina/supplywatch/pkg/backoff"
	"github.com/streadway/amqp"
)

// Each warehouse is a subscriber and a publisher
// when a warehouse runs out of stock of something it will send a request to the queue for some other
// warehouse to send that stuff to it - each warehouse subscribes to these requests for supply
func (w *warehouse) SetupMessageQueue(storageChan, sendChan chan string) {
    var wg sync.WaitGroup
    wg.Add(2)

    var connRabbit *amqp.Connection
    var err error

    var attempt int
    for {
        time.Sleep(backoff.Default.Duration(attempt))
        connRabbit, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
        if err != nil {
            w.logger.Error(err)
            w.logger.Error("Failed to connect to RabbitMQ")
            attempt++
            continue
        }
        break
    }
    w.logger.Info("Successfully Connected to RabbitMQ")
    defer connRabbit.Close()
    channel, err := connRabbit.Channel()
    if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to setup a channel to RabbitMQ")
    }
    defer channel.Close()
    _, err = channel.QueueDeclare(
        "RequestProducts",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed declare a Queue")
    }

    // Subscribe
    go func(){
        w.publishMessages(channel, storageChan)
        wg.Done()
    }()

    // Publish
    go func(){
        w.subscribeMessages(channel, sendChan)
        wg.Done()
    }()

    wg.Wait()
}

// publishMessages waits for a message on the storageChan
// when a product quantity drops to zero a message will be send to the message Queue
// ideally another warehouse receives this message and sends some of that product
func (w *warehouse) publishMessages(c *amqp.Channel, storageChan chan string) {
    for {
        w.logger.Info("Waiting for a product to drop to zero")

        zeroProduct := <- storageChan
        if !strings.Contains(zeroProduct, ":") {
            w.logger.Error("Cannot publish Message %s because the format is invalid", zeroProduct)
            continue
        }

        w.logger.Info("--------------------Sending to Queue!-------------------")
        err := c.Publish(
            "",
            "RequestProducts",
            false,
            false,
            amqp.Publishing{
                ContentType: "text/plain",
                Body: []byte(zeroProduct),
            },
        )
        if err != nil {
            w.logger.Error(err)
            w.logger.Fatal("Failed publish a Testmessage")
        }
    }
}

//subscribeMessages subscribes to the RequestProducts queue and whenever a request for products comes in
//it trys to initalize the transfer of this product
func (w *warehouse) subscribeMessages(c *amqp.Channel, sendChan chan string) {
    var wg sync.WaitGroup
    wg.Add(1)
    w.logger.Info("Subscribe to RequestProducts queue!")
    msgs, err := c.Consume(
        "RequestProducts",
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed to subscribe to a Testmessage")
    }
    //forever := make(chan bool)

    //listen to incoming messages
    go func(){
        for d := range msgs{
            w.logger.Info("--------------------Receiving from Queue!-------------------")
            w.logger.Infof("Received Message: %s from queue", string(d.Body))
            //incoming messages should have the following format: product:hostname
            if !strings.Contains(string(d.Body), ":") {
                w.logger.Error("Received Message from queue with invalid format, this message will be ignored")
                continue
            }
            w.logger.Info("Initalizing transport of products")
            sendChan <- string(d.Body)
        }
        wg.Done()
    }()

    wg.Wait()
}
