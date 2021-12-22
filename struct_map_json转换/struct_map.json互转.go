package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/goinggo/mapstructure"
	"math"
	"reflect"
)

type Host struct {
	IP   string `json:"ip_1"`
	Name string `json:"name"`
}

func JsonToStructDemo() {

	b := []byte(`{"ip_1": "192.168.11.22", "name": "SKY"}`)

	m := Host{}

	err := json.Unmarshal(b, &m)
	if err != nil {

		fmt.Println("Umarshal failed:", err)
		return
	}

	fmt.Println("m:", m)
	fmt.Println("m.IP:", m.IP)
	fmt.Println("m.Name:", m.Name)
}

func JsonToMapDemo() {
	jsonStr := `
        {
                "name": "jqw",
                "age": 18
        }
        `
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}
	fmt.Println(mapResult)

}

func MapToJsonDemo1() {
	mapInstances := []map[string]interface{}{}
	instance_1 := map[string]interface{}{"name": "John", "age": 10}
	instance_2 := map[string]interface{}{"name": "Alex", "age": 12}
	mapInstances = append(mapInstances, instance_1, instance_2)

	jsonStr, err := json.Marshal(mapInstances)

	if err != nil {
		fmt.Println("MapToJsonDemo err: ", err)
	}
	fmt.Println(string(jsonStr))
}

func MapToJsonDemo2() {
	b, _ := json.Marshal(map[string]int{"test": 1, "try": 2})
	fmt.Println(string(b))
}

type People struct {
	Name    string  `json:"name_title"`
	Age     int     `json:"age_size"`
	Float64 float64 `json:"float64"`
}

func MapToStructDemo() {
	mapInstance := make(map[string]interface{})
	mapInstance["Name"] = "jqw"
	mapInstance["Age"] = 18
	mapInstance["Float64"] = math.MaxFloat64

	var people People
	err := mapstructure.Decode(mapInstance, &people)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", people)
}

func StructToMapDemo(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}

	//
	m := structs.Map(obj)
	fmt.Printf("%#v\n", m)
	return data
}

func main() {
	// xxx jsonstr to struct , 注意json里面的key和struct里面的key要一致，struct中的key的首字母必须大写，而json中大小写都可以。
	JsonToStructDemo()

	// xxx jsonstr to map
	JsonToMapDemo()

	// xxx map to json   [{"age":10,"name":"John"},{"age":12,"name":"Alex"}]
	MapToJsonDemo1()

	// xxx map to json   {"test":1,"try":2}
	MapToJsonDemo2()

	// xxx map to struct
	MapToStructDemo()

	// xxx struct to map

	data := StructToMapDemo(People{"aa", 123, math.MaxFloat64 - 1})
	fmt.Println(data)
}
