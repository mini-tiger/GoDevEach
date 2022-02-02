package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"strings"
	"time"
)

func main() {
	servers := []string{"nats://127.0.0.1:4222"}

	nc, err := nats.Connect(strings.Join(servers, ","))

	if err != nil {
		log.Fatalln(err)
	}
	defer nc.Close()

	// 用来等待消息到来
	//wg := sync.WaitGroup{}
	//wg.Add(1)

	// 异步订阅
	for {
		_, err := nc.QueueSubscribe("demo.test1", "queue", func(msg *nats.Msg) {
			fmt.Println(string(msg.Data))
		})
		if err != nil {
			log.Fatal(err)
		}
		//msg, err := queue.NextMsg(50*time.Millisecond)
		//if err != nil {
		//	log.Println(err)
		//	continue
		//}
		//log.Printf("Msg: %s", msg.Data)
		time.Sleep(50 * time.Millisecond)
	}

	//wg.Wait()
	//select {}

}
