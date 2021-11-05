package main

import (
	"log"
	"time"

	amqp "github.com/streadway/amqp"
	rabbitmq "github.com/wagslane/go-rabbitmq"
)

func main() {
	//publisher, err, _ := rabbitmq.NewPublisher(
	//	"amqp://user:password@172.16.71.17:5672/", amqp.Config{},
	//	rabbitmq.WithPublisherOptionsLogging,
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	publisher, returns, err := rabbitmq.NewPublisher(
		"amqp://user:password@172.16.71.17:5672/",
		// can pass nothing for no logging
		//rabbitmq.WithPublisherOptionsLogging,
		amqp.Config{}, rabbitmq.WithPublisherOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = publisher.Publish(
		[]byte("hello, world"),
		[]string{"routing_key"},
		// leave blank for defaults
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange("taojun"),
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for r := range returns {
			log.Printf("message returned from server: %s", string(r.Body))
		}
	}()
	time.Sleep(10 * time.Second)
}
