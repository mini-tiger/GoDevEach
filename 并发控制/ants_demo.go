package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"log"
	"sync"
	"time"
)

//xxx https://gitee.com/mirrors/ants

func myFunc(i interface{}) {
	n := i.(*Entry)
	//atomic.AddInt32(&sum, n)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("run with %+v\n", *n)
}

type Entry struct {
	int
	Massage string
}

func main() {
	runTimes := 1000
	var wg sync.WaitGroup
	// 创建一个容量为10的goroutine池
	//p,_:=ants.NewPool(10)

	// xxx 性能更好
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		myFunc(i)
		wg.Done()
	},
		ants.WithMaxBlockingTasks(2), // 设置等待队列的最大长度。超过这个长度，提交任务直接返回错误
		ants.WithNonblocking(false),  // 设置其为非阻塞。非阻塞的ants池中，在所有 goroutine 都在处理任务时，提交新任务会直接返回错误
	)

	defer p.Release() // xxx 使用完必须释放

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		err := p.Invoke(&Entry{i, "Hello"}) //xxx 提交任务  不是执行
		if err != nil {
			log.Fatalln(err)
		}
	}
	wg.Wait()

	//p.Reboot()
	fmt.Printf("running goroutines: %d\n", p.Running()) //执行
	//fmt.Printf("finish all tasks, result is %d\n", sum)
}
