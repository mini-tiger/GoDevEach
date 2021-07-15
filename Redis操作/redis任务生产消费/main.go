package main

import "redis任务生产消费/funcs"

func main() {
	//funcs.Demo1()

	funcs.Demo2()
	select {}
}
