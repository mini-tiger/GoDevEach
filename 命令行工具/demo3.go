package main

import (
	"fmt"
	flag "github.com/spf13/pflag" //替换原生的flag，并兼容
)

// ./demo3 -b=false  要使用等号
var flagvar1 int
var flagvar2 bool

func init() {
	// 不包含短参数
	flag.IntVar(&flagvar1, "varname1", 1, "help message for flagname")
	// 包含短参数
	flag.BoolVarP(&flagvar2, "boolname1", "b", true, "help message")
}

func main() {
	// 不包含短参数
	var ip1 *int = flag.Int("flagname1", 1, "help message for flagname")
	// 包含短参数
	var ip2 = flag.IntP("flagname2", "f", 2, "help message")

	flag.Parse()

	fmt.Println("ip1 has value ", *ip1)
	fmt.Println("ip2 has value ", *ip2)
	fmt.Println("flagvar1 has value ", flagvar1)
	fmt.Println("flagvar2 has value ", flagvar2)
}
