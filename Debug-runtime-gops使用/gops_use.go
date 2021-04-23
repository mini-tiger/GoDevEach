package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/google/gops/agent"
)

func main() {
	if err := agent.Listen(agent.Options{
		Addr: "0.0.0.0:8848",
		// ConfigDir:       "/home/centos/gopsconfig", // 最好使用默认
		ShutdownCleanup: true}); err != nil {
		log.Fatal(err)
	}



	//_ = make([]int, 1000, 2000)
	//runtime.GC()


	for i:=0;i<10;i++{
		// 测试代码
		A := make([]int, 1000, 1000)
		runtime.GC()
		time.Sleep(time.Minute)
		fmt.Sprintf("%+v",A)
	}

	time.Sleep(time.Hour)
}
