package test

import (
	"fmt"
	"sync"
	"testing"
)

/**
 * @Author: Tao Jun
 * @Since: 2022/6/27
 * @Desc: chanl_test.go
**/

func Benchmark_chan1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		chan1()
	}
}

func Benchmark_chan2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		chan2()
	}
}

func chan1() {
	var sync1 sync.WaitGroup
	var chan11 chan int = make(chan int, 0)
	go func() {
		for {
			select {
			case i, ok := <-chan11:
				if !ok {
					break
				}
				fmt.Sprint(i)
				sync1.Done()
			}
		}
	}()

	for i := 0; i < 1000; i++ {
		sync1.Add(1)
		chan11 <- i
	}
	//time.Sleep(1000 * time.Microsecond)

	sync1.Wait()
	//close(chan11)

}

func chan2() {
	var sync2 sync.WaitGroup
	var chan12 chan int = make(chan int, 0)
	go func() {
		for i := range chan12 {
			fmt.Sprint(i)
			sync2.Done()
		}
	}()

	for i := 0; i < 1000; i++ {
		sync2.Add(1)

		chan12 <- i
	}
	//time.Sleep(1000 * time.Microsecond)

	sync2.Wait()
	close(chan12)
}
