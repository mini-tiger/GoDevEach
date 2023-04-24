package main

import (
	"fmt"
	"github.com/asaskevich/govalidator"
)

func main() {
	testNumber := "42"

	if govalidator.InRangeInt(testNumber, 0, 100) {
		fmt.Printf("%v 在 0 和 100 之间\n", testNumber)
	} else {
		fmt.Printf("%v 不在 0 和 100 之间\n", testNumber)
	}
	fmt.Println(govalidator.IsIPv4("172.22.50.2411"))
}
