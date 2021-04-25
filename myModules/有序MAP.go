package main

import (
	"fmt"
	//lm "gitee.com/taojun319/tjtools/LinkedMap"
	lm "github.com/mini-tiger/tjtools/LinkedMap"
	"strconv"
)

func main() {
	a := lm.NewLinkedMap()

	for i := 0; i < 2000; i++ {
		a.Put(strconv.Itoa(i), i)
	}

	for i := 0; i < 100; i++ {
		a.Remove(strconv.Itoa(i))
	}
	//fmt.Println(a.Max())
	//fmt.Printf("%+v\n", a.MData)
	//fmt.Printf("%+v\n", a.MLink)
	for _, key := range a.SortLinkMap() {
		//fmt.Println(key)
		if v, e := a.Get(key); e {
			fmt.Printf("key:%s,value:%v\n", key, v)
		}
	}
}
