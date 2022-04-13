package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/nats-io/nats.go"
)

// todo 启动服务端 ./nats-server -js
// docker run -d --name nats-main -p 4222:4222 -p 6222:6222 -p 8222:8222 nats -js
// https://www.jianshu.com/p/27a49b9d4306

const (
	streamName     = "ORDERS"
	streamSubjects = "ORDERS.*"
	subjectName    = "ORDERS.created"
)

type Order struct {
	OrderID    int
	CustomerID string
	Status     string
}

func main() {
	// 连接NATS
	nc, _ := nats.Connect(nats.DefaultURL)
	// 创建JetStreamContext

	js, err := nc.JetStream()

	checkErr(err)
	// 创建stream流
	err = createStream(js)
	checkErr(err)
	// 通过发布消息创建订单
	err = createOrder(js)
	checkErr(err)
	//subHandle(js)
	// Simple Pull Consumer
	sub, err := js.PullSubscribe(streamSubjects, "MONITOR")
	msgs, err := sub.Fetch(10)
	//fmt.Println(msgs)
	for _, msg := range msgs {
		msgHandle(msg)
	}
}

func msgHandle(msg *nats.Msg) {
	fmt.Println("data: ", string(msg.Data))
}

/*
	nats.DeliverNew()这个选项。如果不声明，则默认为nats.DeliverAll()；
	除了这两个参数，还有一个nats.DeliverLast()参数，这分别对应了3种开始订阅时的方式：
	默认方式nats.DeliverAll()是会读取有效生命周期内的所有消息，甚至包含已被处理的消息；
	nats.DeliverLast()是会包含消息队列中的最后一条消息，即使被处理过的消息；
	nats.DeliverNew()则只处理订阅之后的新消息
*/
func subHandle(js nats.JetStreamContext) {
	// Create durable consumer monitor
	js.Subscribe(streamSubjects, func(msg *nats.Msg) {
		msg.Ack()
		var order Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("monitor service subscribes from subject:%s\n", msg.Subject)
		log.Printf("OrderID:%d, CustomerID: %s, Status:%s\n", order.OrderID, order.CustomerID, order.Status)
	}, nats.Durable("monitor"), nats.ManualAck())
}

// createOrder 以 "ORDERS.created"主题发布事件流
func createOrder(js nats.JetStreamContext) error {
	var order Order
	for i := 1; i <= 10; i++ {
		order = Order{
			OrderID:    i,
			CustomerID: "Cust-" + strconv.Itoa(i),
			Status:     "created",
		}
		orderJSON, _ := json.Marshal(order)
		_, err := js.Publish(subjectName, orderJSON)
		if err != nil {
			return err
		}
		log.Printf("Order with OrderID:%d has been published\n", i)
	}
	return nil
}

// createStream 使用JetStreamContext创建流
func createStream(js nats.JetStreamContext) error {
	// Check if the ORDERS stream already exists; if not, create it.
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Println(err)
	}
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", streamName, streamSubjects)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamSubjects},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
