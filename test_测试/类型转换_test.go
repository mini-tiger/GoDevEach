package test

import (
	"reflect"
	"strconv"
	"testing"
)

func Benchmark_Type(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			switchType1(i)
		} else {
			switchType1(strconv.Itoa(i))
		}
	}
}

func Benchmark_type1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			switchType(i)
		} else {
			switchType(strconv.Itoa(i))
		}
	}
}

func switchType1(v interface{}) {
	vv := reflect.TypeOf(v)
	switch vv.Kind() {
	case reflect.Int:

		return
	case reflect.Map:
		return
	case reflect.String:
		return
	default:
		return
	}
}
func switchType(v interface{}) {
	switch v.(type) {
	case int:

		return
	case map[string]interface{}:
		return
	case string:
		return
	default:
		return
	}
}
