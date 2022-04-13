package main

import (
	"fmt"
)

type Addable interface {
	int | int64 | string
}

func add[T Addable](a, b T) T {
	return a + b
}
func main() {
	fmt.Println(add(1, 2))
	// FIXME
	fmt.Println(add("foo", "bar"))
}
