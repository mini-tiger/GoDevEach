package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
)

type taskFunc func()

func taskFuncWrapper(i *Entry, wg *sync.WaitGroup) taskFunc {
	return func() {
		n := i
		//atomic.AddInt32(&sum, n)
		fmt.Printf("run with %+v\n", *n)
		wg.Done()
	}
}

type Entry struct {
	int
	Massage string
}

func main() {
	runTimes := 1000
	var wg sync.WaitGroup
	// 创建一个容量为10的goroutine池
	p, _ := ants.NewPool(10)

	defer p.Release() // xxx 使用完必须释放

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		p.Submit(taskFuncWrapper(&Entry{
			int:     i,
			Massage: "Hello",
		}, &wg))
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	//fmt.Printf("finish all tasks, result is %d\n", sum)
}
