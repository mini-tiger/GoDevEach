package main

import (
	"github.com/robfig/cron"
	"log"
)

func main() {
	log.Println("Starting...")

	c := cron.New()                   // 新建一个定时任务对象
	c.AddFunc("* * * * * *", func() { //每秒
		log.Println("hello world")
	}) // 给对象增加定时任务
	c.Start()
	select {}
}
