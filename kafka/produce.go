package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

var wg sync.WaitGroup

func main() {
	config := sarama.NewConfig()
	//设置
	//ack应答机制
	config.Producer.RequiredAcks = sarama.WaitForAll

	//发送分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	//回复确认
	config.Producer.Return.Successes = true

	//连接kafka
	client, err := sarama.NewSyncProducer([]string{"192.168.43.111:9092"}, config)
	if err != nil {
		fmt.Println("producer closed,err:", err)
	}
	defer client.Close()

	//发送消息
	for i := 0; ; i++ {

		//构造一个消息
		msg := &sarama.ProducerMessage{}
		msg.Topic = "test-topic"
		msg.Value = sarama.StringEncoder(fmt.Sprintf("test:weatherStation device [%d]", i))

		pid, offset, err := client.SendMessage(msg)
		if err != nil {
			fmt.Println("send msg failed,err:", err)
			return
		}
		fmt.Printf("pid:%v offset:%v msg:%v \n ", pid, offset, msg.Value)
		time.Sleep(1 * time.Second)
	}

}
