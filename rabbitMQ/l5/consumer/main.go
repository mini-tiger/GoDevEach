package main

import (
	"fmt"
	"log"
	config "rabbitMQ/l5/conf"

	"github.com/streadway/amqp"
)

/*


./consumer "#" info.payment.* *.log debug.payment.#

*/

func main() {

	conn, err := amqp.Dial(config.RMQADDR)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	forever := make(chan bool)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		config.EXCHANGENAME, //exchange name
		"topic",             //exchange kind
		true,                //durable
		false,               //autodelete
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	//if len(os.Args) < 2 {
	//	log.Println(len(os.Args))
	//	log.Println(`"Arguments error(Example: ./consumer "#" info.payment.* *.log debug.payment.#"`)
	//	return
	//}

	topics := []string{"debug.#", "#.info"} // 通配符，符号“#”匹配一个或多个词，符号“*”仅匹配一个词
	topicsCnt := len(topics)

	for routing := 0; routing < topicsCnt; routing++ {
		go func(routingNum int) {

			q, err := ch.QueueDeclare(
				"",
				false, //durable
				false, //delete when unused
				true,  //exclusive
				false, //no-wait
				nil,   //arguments
			)

			failOnError(err, "Failed to declare a queue")

			err = ch.QueueBind(
				q.Name,
				topics[routingNum],
				config.EXCHANGENAME,
				false,
				nil,
			)
			failOnError(err, "Failed to bind exchange")

			msgs, err := ch.Consume(
				q.Name,
				"",
				true, //Auto Ack
				false,
				false,
				false,
				nil,
			)

			failOnError(err, "Failed to register a consumer")

			for msg := range msgs {
				log.Printf("In %s consume a message: %s\n", topics[routingNum], msg.Body)
			}

		}(routing)
	}

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}
