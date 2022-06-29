package main

import (
	"fmt"
)

func main() {
	var chan_exit chan int = make(chan int, 0)
	var chan1 chan int = make(chan int, 1)
	go func() {
		for {
			select {
			case i, ok := <-chan1:
				if !ok {
					//sync1.Done()
					fmt.Println("exit")
					chan_exit <- 0
					return
				}
				fmt.Println(i, ok)
			}
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			//sync1.Add(1)
			chan1 <- i
		}
		close(chan1)
	}()
	//sync1.Add(1)
	//sync1.Wait()
	<-chan_exit
}
