package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	config "rabbitMQ/l5/conf"
)

func main() {

	conn, err := amqp.Dial(config.RMQADDR)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		config.EXCHANGENAME, //exchange name
		"topic",             //exchange kind,routeing key 通配符
		true,                //durable
		false,               //autodelete
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	//if len(os.Args) < 3 {
	//	fmt.Println("Arguments error(ex:producer topic msg1 msg2 msg3")
	//	return
	//}

	routingKey1 := "debug.mail.key1"
	routingKey2 := "mail.info"

	msgs := []string{"a", "b", "c", "d", "e"}

	msgNum := len(msgs)
	for i := 0; i < 65535; i++ {
		for cnt := 0; cnt < msgNum; cnt++ {
			msgBody := msgs[cnt]
			err = ch.Publish(
				config.EXCHANGENAME, //exchange
				routingKey1,         //routing key
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(msgBody),
				})

			log.Printf("routingkey %s, [x] Sent %s", routingKey1, msgBody)
			err = ch.Publish(
				config.EXCHANGENAME, //exchange
				routingKey2,         //routing key
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(msgBody),
				})

			log.Printf("routingkey %s, [x] Sent %s", routingKey2, msgBody)
		}
	}
	failOnError(err, "Failed to publish a message")

}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}
