package main

import (
	"context"
	"fmt"
)

func makeChan[T chan E, E any](ctx context.Context, arr []E) T {
	ch := make(T)
	go func() {
		defer close(ch)
		for _, v := range arr {
			select {
			case <-ctx.Done():
				return
			case ch <- v:
			}

		}
	}()
	return ch
}

type able interface {
	int | int64 | string
}

func makeChan1[T chan e, e able](ctx context.Context, arr []e) T {
	ch := make(T)
	go func() {
		defer close(ch)
		for _, v := range arr {
			select {
			case <-ctx.Done():
				return
			case ch <- v:
			}

		}
	}()
	return ch
}
func main() {
	// xxx 有限制
	ch1 := makeChan1(context.Background(), []int{1, 2, 3})

	for v := range ch1 {
		fmt.Println(v)
	}
	// xxx 无限制
	makeChan(context.Background(), []int{1, 2, 3})
}
