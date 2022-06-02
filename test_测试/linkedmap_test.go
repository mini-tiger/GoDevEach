package test

import (
	lm "github.com/mini-tiger/tjtools/LinkedMap"
	"strconv"
	"testing"
)

func Benchmark_Linkmap(b *testing.B) {
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
		a.Get(key)
		//fmt.Printf("key:%s,value:%v\n", key, v)

	}

}
