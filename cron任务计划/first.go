package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
)

//https://godoc.org/github.com/robfig/cron
type TestJob struct {
}

func (this TestJob) Run() {
	fmt.Println("testJob1...")
}

type Test2Job struct {
}

func (this Test2Job) Run() {
	fmt.Println("testJob2...")

}

//启动多个任务
func main() {
	i := 0
	c := cron.New()

	//defaultSecond:=cron.Every(t)

	//AddFunc
	spec := "0 */1 * * * " //每小时 只在 整分钟0秒执行

	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})

	//AddJob方法
	c.AddJob(spec, TestJob{})
	c.AddJob(spec, Test2Job{})

	//启动计划任务
	c.Start()

	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()

	select {}
}
