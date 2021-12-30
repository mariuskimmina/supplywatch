package warehouse

import (
	"fmt"
	"sync"
	"time"

	"github.com/mariuskimmina/supplywatch/pkg/backoff"
	"github.com/streadway/amqp"
)

func (w *warehouse) SetupMessageQueue() {
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
    q, err := channel.QueueDeclare(
        "TestQueue",
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

    fmt.Println(q)

    // Subscribe
    go func(){
        w.publishMessages(channel)
        wg.Done()
    }()

    // Publish
    go func(){
        w.subscribeMessages(channel)
        wg.Done()
    }()

    wg.Wait()
}

func (w *warehouse) publishMessages(c *amqp.Channel) {
    for {
        w.logger.Info("Publish an important message!")
        err := c.Publish(
            "",
            "TestQueue",
            false,
            false,
            amqp.Publishing{
                ContentType: "text/plain",
                Body: []byte("Hello World"),
            },
        )
        if err != nil {
            w.logger.Error(err)
            w.logger.Fatal("Failed publish a Testmessage")
        }
        time.Sleep(10 * time.Second)
    }
}

func (w *warehouse) subscribeMessages(c *amqp.Channel) {
    var wg sync.WaitGroup
    wg.Add(1)
    w.logger.Info("Subscribe to an important message!")
    msgs, err := c.Consume(
        "TestQueue",
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
            w.logger.Infof("Received Message: %s\n", d.Body)
        }
        wg.Done()
    }()

    wg.Wait()
}
