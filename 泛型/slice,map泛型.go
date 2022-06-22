package main

/**
 * @Author: Tao Jun
 * @Since: 2022/6/21
 * @Desc: slice,map泛型.go
**/
import (
	"fmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func main() {
	m := make(map[string]interface{}, 0)
	m["1"] = 1
	m["2"] = "2"
	fmt.Println(maps.Keys(m))

	a := []string{"1", "2"}
	fmt.Println(slices.Contains(a, "1"))
}
