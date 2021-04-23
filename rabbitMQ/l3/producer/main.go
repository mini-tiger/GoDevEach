package main

import (
	"fmt"
	"log"
	config "rabbitMQ/l3/conf"

	"github.com/streadway/amqp"
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
		"fanout",            //exchange kind
		true,                //durable
		false,               //autodelete
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	msgs := []string{"a", "b", "c", "d"}
	msgNum := len(msgs)

	for cnt := 0; cnt < msgNum; cnt++ {
		msgBody := msgs[cnt]
		err = ch.Publish( // 直接通过exchange ,queue 随机名称
			config.EXCHANGENAME, //exchange
			"",                  //routing key， fanout 类型exchange 不匹配routeing key
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msgBody),
			})

		log.Printf(" [x] Sent %s", msgBody)
	}
	failOnError(err, "Failed to publish a message")

}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}
