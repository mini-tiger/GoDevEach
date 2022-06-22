package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

type Addable interface {
	int | int64 | string
}

func add[T Addable](a, b T) T {
	return a + b
}
func add1[T constraints.Integer](a, b T) T {
	return a + b
}
func main() {
	fmt.Println(add1(1, 2))
	// FIXME
	fmt.Println(add("foo", "bar"))
}
