package warehouse

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mariuskimmina/supplywatch/internal/domain"
	"github.com/mariuskimmina/supplywatch/pkg/backoff"
	"github.com/streadway/amqp"
)

const (
	warehouse1Queue = "warehouse1Queue"
	warehouse2Queue = "warehouse2Queue"
	warehouse3Queue = "warehouse3Queue"

	warehouse1Data = "warehouse1Data"
	warehouse2Data = "warehouse2Data"
	warehouse3Data = "warehouse3Data"
)

const (
	logFileDir    = "/var/supplywatch/warehouse/"
	logFilePrefix       = "products-"
)

// Each warehouse is a subscriber and a publisher
// when a warehouse runs out of stock of something it will send a request to the queue for some other
// warehouse to send that stuff to it - each warehouse subscribes to these requests for supply
func (w *warehouse) SetupPublishing(storageChan, sendChan chan string, warehouseNames string) {
	var wg sync.WaitGroup
	wg.Add(4)

	var connRabbit *amqp.Connection
	var err error

	var attempt int
	for {
		time.Sleep(backoff.Default.Duration(attempt))
		connRabbit, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			//w.logger.Error(err)
			attempt++
			w.logger.Infof("Failed to connect to RabbitMQ, trying again in %f seconds", backoff.Default.Duration(attempt).Seconds())
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
	host, err := os.Hostname()
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("failed to get hostname")
	}

	// Each warehouse creates a queue for itself, based on its hostname
	// otherQueues contains the names of all queues but the one of the current warehouse
	// we will publish to otherQueues and subscribe to our own

    // get list of all warehouses from config
    // remove this warehouse from the list
    // subscribe to the "data" queues of all other warehouses
    // when a product drops to zero, we publish a message to the product queues of all other warehouses

    whNames := strings.Split(warehouseNames, ",")
    var dataQueuesToSub []string
    var productOrderQueuesToSub []string
    for i, name := range whNames{
        if strings.Contains(name, host) {
            whNames = append(whNames[:i], whNames[i+1:]...)
        } else {
            dataQueuesToSub = append(dataQueuesToSub, name + "DataExchange")
            productOrderQueuesToSub = append(productOrderQueuesToSub, name + "OrderProductsExchange")

        }
    }

    w.logger.Infof("Other warehouses: %s", warehouseNames)
    w.logger.Infof("Other warehouses: %s", whNames)
    w.logger.Infof("Other warehouses data Queues: %s", dataQueuesToSub)


    // This exchange is used to publish data to other warehouses and the monitor
    dataExchangeName := host + "DataExchange"
	err = channel.ExchangeDeclare(
        dataExchangeName, // name
		"fanout", //type
		true, // durable
		false, // auto-deleted
		false, // internal
		false, //no-vait
		nil, //arguments
	)


    // This exchange is used to order Products from all the other warehouses
    // Once a product drops to zero we publish a message to this exchange to request
    // this product from all the other warehosues
    productExchangeName := host + "OrderProductsExchange"
	err = channel.ExchangeDeclare(
        productExchangeName, // name
		"fanout", //type
		true, // durable
		false, // auto-deleted
		false, // internal
		false, //no-vait
		nil, //arguments
	)

	// Subscribe
	go func() {
		w.publishProductRequests(channel, storageChan, productExchangeName)
		wg.Done()
	}()

    // Here we publish data for the monitor
	go func() {
		w.publishData(channel, dataExchangeName)
		wg.Done()
	}()

	go func() {
		w.SubscribeToAllWarehouseData(channel, dataQueuesToSub)
		wg.Done()
	}()

	go func() {
		w.SubscribeToAllProductOrders(channel, sendChan, productOrderQueuesToSub)
		wg.Done()
	}()

	wg.Wait()
}

func (w *warehouse) SubscribeToAllProductOrders(channel *amqp.Channel, sendChan chan string, exchangeNames []string) {
    w.logger.Info("Subscribing to the Product Orders of all other warehouses")
	var wg sync.WaitGroup
	wg.Add(len(exchangeNames))
    for _, name := range exchangeNames {
        go func(name string) {
            w.SubtoWarehouseOrder(channel, sendChan, name)
            wg.Done()
        }(name)
    }
    wg.Wait()

}

func (w *warehouse) SubtoWarehouseOrder(channel *amqp.Channel, sendChan chan string, exchangeName string) {
    w.logger.Infof("Subbing to data exchange !! %s !!", exchangeName)
    q, err := channel.QueueDeclare(
        "",
        false,
        false,
        true,
        false,
        nil,
    )
    if err != nil {
        w.logger.Error("Failed to declare queue on exchange")
    }

    err = channel.QueueBind(
        q.Name,
        "",
        exchangeName,
        false,
        nil,
    )
    if err != nil {
        w.logger.Error("Failed to bind to queue to subscribe")
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
        w.logger.Error("Failed to register as a consumer here")
    }

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        w.logger.Info("Waiting for messages here")
        for d := range msgs {
            w.logger.Infof("Recevied Order from another warehouse: %s", exchangeName)
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

func (w *warehouse) publishProductRequests(channel *amqp.Channel, storageChan chan string, exchangeName string) {
    w.logger.Info("Waiting for a product to drop to zero")

    zeroProduct := <-storageChan
    if !strings.Contains(zeroProduct, ":") {
        w.logger.Error("Cannot publish Message %s because the format is invalid", zeroProduct)
        //continue
    }
    w.logger.Infof("Publishing Request for product to exchange: %s", exchangeName)
    err := channel.Publish(
        exchangeName,
        "",
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(zeroProduct),
        },
    )
    if err != nil {
        w.logger.Error(err)
        w.logger.Fatal("Failed publish a Testmessage")
    }
}

func (w *warehouse) publishData(channel *amqp.Channel, exchangeName string) {
	for {
        time.Sleep(30 * time.Second)
		//w.logger.Info("Waiting for a product to drop to zero")

		//zeroProduct := <-storageChan
		//if !strings.Contains(zeroProduct, ":") {
			//w.logger.Error("Cannot publish Message %s because the format is invalid", zeroProduct)
			//continue
		//}

		w.logger.Info("Sending Info to Data Exchange!")
        var allProducts []Product
        w.DB.Find(&allProducts)
        productBytes, err := json.Marshal(allProducts)
        if err != nil {
            w.logger.Error(err)
            w.logger.Fatal("Failed to marshal info for monitor")
        }
        //w.logger.Info(allProducts)
        err = channel.Publish(
            exchangeName,
            "",
            false,
            false,
            amqp.Publishing{
                ContentType: "text/plain",
                Body:        []byte(productBytes),
            },
        )
        if err != nil {
            w.logger.Error(err)
            w.logger.Fatal("Failed publish a Testmessage")
        }
	}
}

func (w *warehouse) SubscribeToAllWarehouseData(channel *amqp.Channel, exchangeNames []string) {
    w.logger.Info("Subscribing to the Data of all other warehouses")
	var wg sync.WaitGroup
	wg.Add(len(exchangeNames))
    for _, name := range exchangeNames {
        go func(name string) {
            w.SubtoWarehouseData(channel, name)
            wg.Done()
        }(name)
    }
    wg.Wait()
}

func (w *warehouse) SubtoWarehouseData(channel *amqp.Channel, exchangeName string) {
    err := os.MkdirAll(logFileDir, 0644)
	if err != nil {
        fmt.Println("doof")
	}
	f, err := os.OpenFile(logFileDir+logFilePrefix+exchangeName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
        fmt.Println("doof")
	}
	defer f.Close()

    w.logger.Infof("Subbing to data exchange !! %s !!", exchangeName)
    q, err := channel.QueueDeclare(
        "",
        false,
        false,
        true,
        false,
        nil,
    )
    if err != nil {
        w.logger.Error("Failed to declare queue on exchange")
    }

    err = channel.QueueBind(
        q.Name,
        "",
        exchangeName,
        false,
        nil,
    )
    if err != nil {
        w.logger.Error("Failed to bind to queue to subscribe")
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
        w.logger.Error("Failed to register as a consumer here")
    }

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        w.logger.Info("Waiting for messages here")
        for d := range msgs {
            w.logger.Infof("Recevied Update from another warehouse: %s", exchangeName)
            var allProducts []domain.Producttype
            err := json.Unmarshal(d.Body, &allProducts)
            //fmt.Println(allProducts)
            if err != nil {
                w.logger.Error(err)
            }
            allProductsJson, err := json.MarshalIndent(allProducts, " ", "")
            if err != nil {
                w.logger.Error(err)
            }
            //w.logger.Info(string(allProductsJson))
            f.Truncate(0)
            f.Seek(0, 0)
            _, err = f.Write(allProductsJson)
            //productsFileName := logFileDir + logFilePrefix + queueName
            //err = ioutil.WriteFile(productsFileName, jsonProducts, 0644)
            if err != nil {
                w.logger.Error(err)
            }
        }
        wg.Done()
    }()
    wg.Wait()
}

// publishMessages waits for a message on the storageChan
// when a product quantity drops to zero a message will be send to the message Queue
// ideally another warehouse receives this message and sends some of that product
func (w *warehouse) publishMessages(c *amqp.Channel, storageChan chan string, otherQueues []string) {
	for {
		w.logger.Info("Waiting for a product to drop to zero")

		zeroProduct := <-storageChan
		if !strings.Contains(zeroProduct, ":") {
			w.logger.Error("Cannot publish Message %s because the format is invalid", zeroProduct)
			continue
		}
		w.logger.Info("--------------------Sending to Queue!-------------------")
		for _, queueName := range otherQueues {
			err := c.Publish(
				"",
				queueName,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
                    Body:        []byte(zeroProduct),
				},
			)
			if err != nil {
				w.logger.Error(err)
				w.logger.Fatal("Failed publish a Testmessage")
			}
		}
	}
}

func (w *warehouse) publishDataMonitor(c *amqp.Channel, queueName string) {
	for {
        time.Sleep(30 * time.Second)
		//w.logger.Info("Waiting for a product to drop to zero")

		//zeroProduct := <-storageChan
		//if !strings.Contains(zeroProduct, ":") {
			//w.logger.Error("Cannot publish Message %s because the format is invalid", zeroProduct)
			//continue
		//}

		w.logger.Info("Sending Info to Data Queue!")
        var allProducts []Product
        w.DB.Find(&allProducts)
        productBytes, err := json.Marshal(allProducts)
        if err != nil {
            w.logger.Error(err)
            w.logger.Fatal("Failed to marshal info for monitor")
        }
        //w.logger.Info(allProducts)
        err = c.Publish(
            "",
            queueName,
            false,
            false,
            amqp.Publishing{
                ContentType: "text/plain",
                Body:        []byte(productBytes),
            },
        )
        if err != nil {
            w.logger.Error(err)
            w.logger.Fatal("Failed publish a Testmessage")
        }
	}
}
func (w *warehouse) subscribeToOtherWarehouseData(c *amqp.Channel, queueName string) {
    w.logger.Infof("Subscribing to Queue: %s", queueName)
	var wg sync.WaitGroup
	wg.Add(1)
	var msgs <-chan amqp.Delivery
	var err error

	//err = os.MkdirAll(logFileDir, 0644)
	//if err != nil {
        //fmt.Println("doof")
	//}
	//f, err := os.OpenFile(logFileDir+logFilePrefix+queueName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	//if err != nil {
        //fmt.Println("doof")
	////}
	//defer f.Close()
	var attempt int
	for {
		time.Sleep(backoff.Default.Duration(attempt))
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
            w.logger.Infof("Failed to subscribe to queue %s, trying again", queueName)
            w.logger.Error(err)
            attempt++
            continue
		}
		break
	}
    w.logger.Infof("Successfully subscribed to queue %s \n", queueName)

	go func() {
        for d := range msgs {
            w.logger.Infof("Receiving Status Update from a Warehouse queue name: %s", queueName)
            w.logger.Infof(string(d.Body))
			//s.logger.Infof("Received Message: %s from queue", string(d.Body))
            //var allProducts []domain.Producttype
            //err := json.Unmarshal(d.Body, &allProducts)
            //fmt.Println(allProducts)
            //if err != nil {
                //TODO: move the file writing somewhere else
                //fmt.Println("doof")
            //}
            //allProductsJson, err := json.MarshalIndent(allProducts, " ", "")
            //if err != nil {
                ////TODO: move the file writing somewhere else
                //fmt.Println("doof")
            //}
            //f.Truncate(0)
            //f.Seek(0, 0)
            //f.Write(allProductsJson)
            //productsFileName := logFileDir + logFilePrefix + queueName
            //err = ioutil.WriteFile(productsFileName, jsonProducts, 0644)
            //if err != nil {
                ////TODO: move the file writing somewhere else
                //fmt.Println("doof")
            //}
			//sendChan <- string(d.Body)
		}
		wg.Done()
	}()

	wg.Wait()

}
//subscribeMessages subscribes to the RequestProducts queue and whenever a request for products comes in
//it trys to initalize the transfer of this product
func (w *warehouse) subscribeMessages(c *amqp.Channel, sendChan chan string, queueName string) {
	var wg sync.WaitGroup
	wg.Add(1)
	var msgs <-chan amqp.Delivery
	var err error
    var attempt int
	for {
		time.Sleep(backoff.Default.Duration(attempt))
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
            w.logger.Infof("Failed to subscribe to queue %s, trying again \n", queueName)
            attempt++
            continue
		}
		break
	}
    w.logger.Infof("Successfully subscribed to queue %s \n", queueName)
	//msgs, err = c.Consume(
		//queueName,
		//"",
		//true,
		//false,
		//false,
		//false,
		//nil,
	//)
	if err != nil {
		w.logger.Error(err)
		w.logger.Fatal("Failed to subscribe to a Testmessage")
	}
	w.logger.Infof("Subscribed to %s queue!", queueName)
	//forever := make(chan bool)

	//listen to incoming messages
	go func() {
		for d := range msgs {
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
