package main

/**
 * @Author: Tao Jun
 * @Since: 2022/6/29
 * @Desc: main.go
**/
import (
	"fmt"

	"github.com/valyala/fasthttp"
)

//相应请求的函数 RequestCtx 传递数据
func testhandle(ctx *fasthttp.RequestCtx) {

	fmt.Fprintf(ctx, "hello world")
}

func main() {

	if err := fasthttp.ListenAndServe(":8081", testhandle); err != nil {

		fmt.Println("start fasthttp fail", err.Error())
	}
}
