package main

import "fmt"

func switchType1(st any) {
	switch st.(type) {
	case int:
		fmt.Println("this is int")
		break
	case string:
		fmt.Println("this is string")
		break
	default:
		fmt.Println("other type")
	}
}

func switchType2(st interface{}) {
	switch st.(type) {
	case int:
		fmt.Println("this is int")
		break
	case string:
		fmt.Println("this is string")
		break
	default:
		fmt.Println("other type")
	}
}

func main() {
	switchType1("abc")
	switchType2("abc")
}
