package main

import (
	"fmt"
	"log"
	config "rabbitMQ/l2/conf"

	"github.com/streadway/amqp"
)

func main() {

	conn, err := amqp.Dial(config.RMQADDR)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	forever := make(chan bool)

	for routine := 0; routine < config.CONSUMERCNT; routine++ {
		go func(routineNum int) {
			ch, err := conn.Channel()
			failOnError(err, "Failed to open a channel")
			defer ch.Close()

			q, err := ch.QueueDeclare(
				config.QUEUENAME,
				true, //durable
				false,
				false,
				false,
				nil,
			)

			failOnError(err, "Failed to declare a queue")

			err = ch.Qos(
				1,     // prefetch count
				0,     // prefetch size
				false, // global
			)

			failOnError(err, "Failed to set QoS")

			msgs, err := ch.Consume(
				q.Name,
				"MsgWorkConsumer",
				false, //Auto Ack
				false,
				false,
				false,
				nil,
			)

			if err != nil {
				log.Fatal(err)
			}
			//其中Auto ack可以设置为true。如果设为true则消费者一接收到就从queue中去除了，如果消费者处理消息中发生意外该消息就丢失了。
			//如果Auto ack设为false。consumer在处理完消息后，调用msg.Ack(false)后消息才从queue中去除。即便当前消费者处理该消息发生意外，只要没有执行msg.Ack(false)那该消息就仍然在queue中，不会丢失
			for msg := range msgs {
				log.Printf("In %d consume a message: %s\n", routineNum, msg.Body)
				log.Printf("Done")
				msg.Ack(false) //Ack
			}

		}(routine)
	}

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}
