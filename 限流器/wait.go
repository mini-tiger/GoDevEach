package main

import (
	"context"
	"fmt"
	"time"
	"golang.org/x/time/rate"
)

/**
 * @Author: Tao Jun
 * @Description: 限流器
 * @File:  first
 * @Version: 1.0.0
 * @Date: 2021/4/23 上午10:25
 */
// xxx https://pkg.go.dev/golang.org/x/time/rate
// xxx https://www.cyhone.com/articles/usage-of-golang-rate/

func main() {
	limit :=  rate.Every(500 * time.Millisecond) //xxx 每0.5 秒 产出一个
	l := rate.NewLimiter(limit, 10)  // 第一个参数为每秒产生多少次token，第二个参数是其缓存最大可存多少个事件


	c, _ := context.WithCancel(context.TODO())

	for {
		//l.Wait(c) //xxx 当没有可用或足够的事件时，将阻塞等待 推荐实际程序中使用这个方法
		l.WaitN(c,1) // 如果此时桶内 Token 数组不足 (小于 N)，那么 Wait 方法将会阻塞一段时间，直至 Token 满足条件。如果充足则直接返回。

		fmt.Printf("可用:%v,time:%v\n",l.Burst(),time.Now().Format("2016-01-02 15:04:05.000"))
	}
}
