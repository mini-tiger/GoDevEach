package main

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

func main() {
	// Create a Resty Client
	client := resty.New()

	// Retries are configured per client
	client.
		//SetLogger(nil). // xxx 修改 默认日志
		SetDebug(true).
		// xxx timeout
		SetTimeout(1 * time.Minute).
		// xxx Setting a Proxy URL and Port
		//SetProxy("http://proxyserver:8888")

		// xxx Want to remove proxy setting
		//	client.RemoveProxy().
		// Set retry count to non zero to enable retries
		SetRetryCount(3).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(5 * time.Second).
		// MaxWaitTime can be overridden as well.
		// Default is 2 seconds.
		SetRetryMaxWaitTime(20 * time.Second).
		// SetRetryAfter sets callback to calculate wait time between retries.
		// Default (nil) implies exponential backoff with jitter
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		}).
		OnError(func(req *resty.Request, err error) { // xxx 请求错误hook
			if v, ok := err.(*resty.ResponseError); ok {
				// v.Response contains the last response from the server
				// v.Err contains the original error
				fmt.Printf("111111 %+v\n", v)
			}
			fmt.Printf("22222 req:%+v err:%+v\n", req, err)
			// Log the error, increment a metric, etc...
		})

	type respStruct struct {
		UserName string `json:"username"`
		PassWord string `json:"password"`
		id       int    `json:"id"`
	}

	//respbody:=make(map[string]interface{},0)
	respbody := &respStruct{}
	resp, err := client.R().SetBody(map[string]interface{}{"username": "testuser", "password": "testpass"}).
		SetHeader("Content-Type", "application/json"). // request header
		ForceContentType("application/json").          // 强制返回格式
		SetResult(respbody).                           // or SetResult(AuthSuccess{}).
		//SetError(&AuthError{}).       // or SetError(AuthError{}).
		Post("https://jsonplaceholder.typicode1.com/posts")
	if err != nil {
		log.Printf("err !!! :%v\n", err)
	} else {
		fmt.Println(string(resp.Body()))
	}

	fmt.Println(respbody)
}
