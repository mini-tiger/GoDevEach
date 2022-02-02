package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	servers := []string{"nats://127.0.0.1:4222"}

	nc, err := nats.Connect(strings.Join(servers, ","))

	if err != nil {
		log.Fatalln(err)
	}
	defer nc.Close()
	for i := 0; ; i++ {
		str := fmt.Sprintf("Hello nats num: %d", i)
		log.Println(str)
		err = nc.Publish("demo.test", []byte(str))
		if err != nil {
			log.Fatalln("Publish message err: ", err)
		}
		nc.Flush()
		//log.Println("Publish message success")
		time.Sleep(500 * time.Millisecond)
	}

}
