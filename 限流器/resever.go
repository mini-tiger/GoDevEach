package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  resever
 * @Version: 1.0.0
 * @Date: 2021/4/23 上午11:22
 */

func main() {
	l := rate.NewLimiter(1, 3)
	for {
		r := l.ReserveN(time.Now(), 3) // Reserve/ReserveN 当没有可用或足够的事件时，返回 Reservation，和要等待多久才能获得足够的事件。
		fmt.Printf("需要等待%vs\n",r.Delay().Seconds())
		time.Sleep(r.Delay())
		fmt.Println(time.Now().Format("2016-01-02 15:04:05.000"))
	}
}