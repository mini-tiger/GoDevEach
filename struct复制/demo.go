package main

import (
	"fmt"
)

type Dog struct {
	age  int
	name string
}

func main() {
	roger := Dog{5, "Roger"}
	mydog := roger
	fmt.Printf("roger addr %p\n", &roger)
	fmt.Printf("mydog addr %p\n", &mydog)
	fmt.Println("Roger and mydog are equal structs?", roger == mydog)
	mydog.name = "piggie"
	fmt.Println("Roger and mydog are equal structs?", roger == mydog)
	fmt.Println(roger)
	// ___________________
	r1 := &Dog{4, "111"}
	m1 := new(Dog)
	*m1 = *r1
	fmt.Println("r1", r1)
	m1.age = 5
	fmt.Println("r1", r1)
	fmt.Println("m1", m1)
	fmt.Printf("roger addr %p\n", r1)
	fmt.Printf("mydog addr %p\n", m1)
}
