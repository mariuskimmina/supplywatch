package warehouse

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

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

// Each warehouse is a subscriber and a publisher
// when a warehouse runs out of stock of something it will send a request to the queue for some other
// warehouse to send that stuff to it - each warehouse subscribes to these requests for supply
func (w *warehouse) SetupMessageQueue(storageChan, sendChan chan string, warehouseNames string) {
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
    for i, name := range whNames{
        if strings.Contains(name, host) {
            whNames = append(whNames[:i], whNames[i+1:]...)
        } else {
            dataQueuesToSub = append(dataQueuesToSub, name + "Data")
        }
    }

    w.logger.Infof("Other warehouses: %s", warehouseNames)
    w.logger.Infof("Other warehouses: %s", whNames)
    w.logger.Infof("Other warehouses data Queues: %s", dataQueuesToSub)
	var queueName string
	var queueName2 string
	var otherQueues []string
	if strings.Contains(host, "warehouse1") {
		queueName = warehouse1Queue
        queueName2 = warehouse1Data
	} else {
		otherQueues = append(otherQueues, warehouse1Queue)
	}
	if strings.Contains(host, "warehouse2") {
		queueName = warehouse2Queue
        queueName2 = warehouse2Data
	} else {
		otherQueues = append(otherQueues, warehouse2Queue)
	}
	// if strings.Contains(host, "warehouse3") {
	// queueName = warehouse3Queue
	// } else {
	// otherQueues = append(otherQueues, warehouse3Queue)
	// }

	_, err = channel.QueueDeclare(
		queueName,
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

	_, err = channel.QueueDeclare(
		queueName2,
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

    // wait so that all queues are delcared before we do stuff with them
    // TODO: find a better solution
    time.Sleep(10 *time.Second)

	// Subscribe
	go func() {
		w.publishMessages(channel, storageChan, otherQueues)
		wg.Done()
	}()

    // Here we publish data for the monitor
	go func() {
		w.publishDataMonitor(channel, queueName2)
		wg.Done()
	}()

	// Publish data to the other warehouse
	go func() {
		w.subscribeMessages(channel, sendChan, queueName)
		wg.Done()
	}()

    for _, queue := range dataQueuesToSub {
        go func(queueName string) {
            w.subscribeToOtherWarehouseData(channel, queueName)
            wg.Done()
        }(queue)
    }

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
