package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"time"
)

var Address = []string{"172.22.159.22:31082", "172.22.159.22:31083", "172.22.159.22:31084", "172.22.159.22:31085",
	"172.22.159.22:31086", "172.22.159.22:31087", "172.22.159.22:31088", "172.22.159.22:31089", "172.22.159.22:31090"}

func main() {
	syncProducer(Address)
	//asyncProducer1(Address)
}

//同步消息模式
func syncProducer(address []string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}
	defer p.Close()
	topic := "test-topic"
	srcValue := "sync: this is a message. index=%d"
	for i := 0; i < 100000; i++ {
		value := fmt.Sprintf(srcValue, i)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(value),
		}
		part, offset, err := p.SendMessage(msg)
		if err != nil {
			log.Printf("send message(%s) err=%s \n", value, err)
		} else {
			fmt.Fprintf(os.Stdout, value+"发送成功，partition=%d, offset=%d \n", part, offset)
		}
		time.Sleep(2 * time.Second)
	}
}
