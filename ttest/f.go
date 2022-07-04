package main

import "fmt"

func GenSliceAppend() []int {
	s := make([]int, 0)
	for i := 0; i < 1023; i++ {
		s = append(s, i)
	}

	return s
}
func main() {
	append1()
	copy1()
}
func append1() {
	s := make([]int, 0, 100)
	fmt.Printf(" GC took %v\n", s)
	for j := 0; j < 10; j++ {
		s = append(s, []int{j, j + 1, j + 2, j + 3, j + 4}...)

	}
	fmt.Printf(" GC took %v,%d\n", len(s), cap(s))
}
func copy1() {
	s := make([]int, 3, 100)
	ss := []int{1, 2, 3}

	copy(ss, s)
	fmt.Printf("%v\n", s)
	fmt.Printf("%v\n", ss)
}
