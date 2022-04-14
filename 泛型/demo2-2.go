package main

import "fmt"

type Custom[T any] [1]T

func (c *Custom[T]) change(elem T) {
	(*c)[0] = elem
}
func (c *Custom[T]) print() {
	fmt.Println(*c)
}
func main() {
	cc := Custom[string]{}
	cc.change("abc")
	cc.print()
}
