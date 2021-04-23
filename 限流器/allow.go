package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  allow
 * @Version: 1.0.0
 * @Date: 2021/4/23 上午10:41
 */

func main() {
	limit := rate.Every(2000 * time.Millisecond) //xxx 每2 秒 产出一个
	l := rate.NewLimiter(limit, 3)               // 第一个参数为每秒产生多少次token，第二个参数是其缓存最大可存多少个事件


	/*
	前3个 马上打印（缓存3个），后面每2秒打印
	*/
	for {
		if l.AllowN(time.Now(), 1) { // AllowN 方法表示，截止到某一时刻，目前桶中数目是否至少为 n 个，满足则返回 true，同时从桶中消费 n 个 token。 反之返回不消费 Token，false。
			fmt.Println(time.Now().Format("2016-01-02 15:04:05.000"))
		} else {

			time.Sleep(1 * time.Second)
		}
	}
}
