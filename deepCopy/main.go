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

type GuestbookSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Guestbook. Edit guestbook_types.go to remove/update
	Foo           string  `json:"foo,omitempty"`
	ConfigMapData *CMSpec `json:"configmap,omitempty"`
}

type CMSpec struct {
	Key1 string `json:"key1,omitempty"`
	Key2 int    `json:"key2,omitempty"`
}

func main() {
	m := make(map[string]interface{}, 0)
	m["1"] = 0
	cpy := deepcopy.Copy(m)
	fmt.Println(cpy)
	m["2"] = 1
	fmt.Println(cpy)

	gbs := new(GuestbookSpec)
	gbs.Foo = "abc"
	fmt.Printf("%p\n", gbs)

	gbs.ConfigMapData = &CMSpec{
		Key1: "",
		Key2: 0,
	}
	fmt.Printf("%p\n", gbs)
	fmt.Printf("%p\n", gbs.ConfigMapData)

	gbs1 := deepcopy.Copy(gbs)
	fmt.Printf("%p\n", gbs1)
	gbs2 := gbs1.(*GuestbookSpec)
	fmt.Printf("%p\n", gbs2)
	gbs2.Foo = "1"
	gbs2.ConfigMapData = &CMSpec{
		Key1: "1",
		Key2: 2,
	}
	fmt.Printf("%+v\n", gbs2.ConfigMapData)
	fmt.Printf("%+v\n", gbs.ConfigMapData)
}
