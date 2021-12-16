package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

//errgroup 可以捕获和记录子协程的错误(只能记录最先出错的协程的错误)
//errgroup 可以控制协程并发顺序。确保子协程执行完成后再执行主协程
//errgroup 可以使用 context 实现协程撤销。或者超时撤销。子协程中使用 ctx.Done()来获取撤销信号

func main() {
	group, _ := errgroup.WithContext(context.Background())
	for i := 0; i < 5; i++ {
		index := i
		group.Go(func() error {
			fmt.Printf("start to execute the %d gorouting\n", index)
			time.Sleep(time.Duration(index) * time.Second)
			if index%2 == 0 {
				return fmt.Errorf("something has failed on grouting:%d", index)
			}
			fmt.Printf("gorouting:%d end\n", index)
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		fmt.Println(err) //只返回 第一个错误
	}
}
