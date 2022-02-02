package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"strings"
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
	if _, err := nc.Subscribe("demo.test", func(m *nats.Msg) {
		fmt.Println(string(m.Data))
		//wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	//wg.Wait()
	select {}

}
