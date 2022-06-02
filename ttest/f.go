package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var jsonstr string = `
{
        "auths": {
                "harbor.dev.21vianet.com": {
                        "auth": "dGFvLmp1bjpUYW9qdW4yMDg="
                }
        }
}

`

func main() {
	jsonm := make(map[string]interface{}, 0)
	err := json.Unmarshal([]byte(jsonstr), &jsonm)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%+v\n", jsonm)
	//for key, value := range jsonm {
	//	fmt.Println(key, value)
	//}
	value, ok := jsonm["auths"]
	if ok {
		//fmt.Println(value)
		//mm := value.(map[string]interface{})
		mm := rt(value)
		if mm != nil {
			fmt.Println(mm)
		}
	}
	//fmt.Println(rt(map[string]interface{}{"a": 1}))
}

func rt(v interface{}) map[string]interface{} {
	vv := reflect.TypeOf(v)
	switch vv.Kind() {
	case reflect.Map:
		return v.(map[string]interface{})
	}
	return nil
}
