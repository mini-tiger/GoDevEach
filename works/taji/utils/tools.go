package utils

import (
	"reflect"
	"unsafe"
)

/**
 * @Author: Tao Jun
 * @Description: utils
 * @File:  tools
 * @Version: 1.0.0
 * @Date: 2021/6/30 下午3:48
 */

func StringToBytes(s string) []byte {
	str := (*reflect.StringHeader)(unsafe.Pointer(&s))
	by := reflect.SliceHeader{
		Data: str.Data,
		Len:  str.Len,
		Cap:  str.Len,
	}
	//在把by从sliceheader转为[]byte类型
	return *(*[]byte)(unsafe.Pointer(&by))
}
