package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

// https://github.com/robfig/cron
//https://godoc.org/github.com/robfig/cron

type TestJob struct {
}

func (this TestJob) Run() {
	fmt.Println("testJob1...")
}

type Test2Job struct {
}

func (this Test2Job) Run() {
	fmt.Println("testJob2...", time.Now().Format(time.RFC3339))

}

type Test3Job struct {
}

func (this Test3Job) Run() {
	time.Sleep(time.Second * 3)
	fmt.Println("testJob3...", time.Now().Format(time.RFC3339))

}

//启动多个任务
func main() {

	ttttt1()
	ttttt2()
	ttttt3()
	select {}
}
func ttttt3() {

	// xxx 可以配置如果当前任务正在进行，那么跳过
	//c := cron.New(cron.WithChain(cron.SkipIfStillRunning(logger)))  // https://blog.cugxuan.cn/2020/06/04/Go/golang-cron-v3/

	// 官方也提供了旧版本的秒级的定义，这个注意你需要传入的 cron 表达式不再是标准 cron 表达式
	c := cron.New(cron.WithSeconds())

	spec := "*/10 * * * * *" // 每10s
	c.AddJob(spec, Test3Job{})
	c.Start()
}

func ttttt2() {
	defaultSecond := cron.Every(time.Duration(time.Second * 2)) // 每两秒执行
	c := cron.New()
	_ = c.Schedule(defaultSecond, Test2Job{})
	c.Start()
}

func ttttt1() {
	c := cron.New()
	i := 0

	/*
		─分鐘（0 - 59）
		# │  ┌──小時（0 - 23）
		# │  │  ┌──日（1 - 31）
		# │  │  │  ┌─月（1 - 12）
		# │  │  │  │  ┌─星期（0 - 6，表示从周日到周六）
		# │  │  │  │  │
	*/

	spec := "*/1 * * * * " // 每一分钟
	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})

	//c.AddFunc()
	c.AddJob(spec, TestJob{})
	c.Start()
}
