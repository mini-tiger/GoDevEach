package service

import (
	"github.com/streadway/amqp"
	"log"
)

func ConsumeExchageQueue(conn *amqp.Connection, exchange, routingKey, queueName string) {
	// 1. 建立RabbitMQ连接
	//conn, err := amqp.Dial("amqp://user:password@172.16.71.17:5672/")
	//failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 2. 创建channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 3. 声明exchange,routing key,queue name
	//exchange := "test_topic_exchange"
	//routingKey := "user.#"
	//routingKey := "user.*"
	//queueName := "test_topic_queue"

	// 4. 声明（创建）一个交换机
	//name:交换器的名称。
	//kind:也叫作type，表示交换器的类型。有四种常用类型：direct、fanout、topic、headers。
	//durable:是否持久化，true表示是。持久化表示会把交换器的配置存盘，当RMQ Server重启后，会自动加载交换器。
	//autoDelete:是否自动删除，true表示是。至少有一条绑定才可以触发自动删除，当所有绑定都与交换器解绑后，会自动删除此交换器。
	//internal:是否为内部，true表示是。客户端无法直接发送msg到内部交换器，只有交换器可以发送msg到内部交换器。
	//noWait:是否非阻塞，true表示是。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。（不推荐使用）
	//args:直接写nil，没研究过，不解释。
	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// 5. 声明（创建）一个队列
	//name：队列名称
	//durable：是否持久化，true为是。持久化会把队列存盘，服务器重启后，不会丢失队列以及队列内的信息。（注：1、不丢失是相对的，如果宕机时有消息没来得及存盘，还是会丢失的。2、存盘影响性能。）
	//autoDelete：是否自动删除，true为是。至少有一个消费者连接到队列时才可以触发。当所有消费者都断开时，队列会自动删除。
	//exclusive：是否设置排他，true为是。如果设置为排他，则队列仅对首次声明他的连接可见，并在连接断开时自动删除。（注意，这里说的是连接不是信道，相同连接不同信道是可见的）。
	//nowait：是否非阻塞，true表示是。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。（不推荐使用）
	//args：直接写nil，没研究过，不解释。
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 6. 队列绑定
	//name：队列名称
	//key：BandingKey，表示要绑定的键。
	//exchange：交换器名称
	//nowait：是否非阻塞，true表示是。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。（不推荐使用）
	//args：直接写nil，没研究过，不解释。

	err = ch.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	// 7. RMQ Server主动把消息推给消费者
	//queue:队列名称。
	//consumer:消费者标签，用于区分不同的消费者。
	//autoAck:是否自动回复ACK，true为是，回复ACK表示告诉服务器我收到消息了。建议为false，手动回复，这样可控性强。
	//exclusive:设置是否排他，排他表示当前队列只能给一个消费者使用。
	//noLocal:如果为true，表示生产者和消费者不能是同一个connect。
	//nowait：是否非阻塞，true表示是。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。（不推荐使用）
	//args：直接写nil，没研究过，不解释。

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
