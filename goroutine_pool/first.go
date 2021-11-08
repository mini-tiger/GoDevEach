package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Jeffail/tunny"
)

func main() {
	pool := tunny.NewFunc(3, func(i interface{}) interface{} {
		log.Println(i.(int))
		time.Sleep(time.Second)
		return nil
	})

	defer pool.Close()
	fmt.Println(pool.GetSize(),pool.QueueLength())
	for i := 0; i < 10; i++ {
		go pool.Process(i)
	}

	time.Sleep(time.Second * 4)
}