package main

import (
	"fmt"
	"github.com/mohae/deepcopy"
)

/**
 * @Author: Tao Jun
 * @Since: 2022/7/8
 * @Desc: main.go
**/
func main() {
	m := make(map[string]interface{}, 0)
	m["1"] = 0
	cpy := deepcopy.Copy(m)
	fmt.Println(cpy)
	m["2"] = 1
	fmt.Println(cpy)
}
