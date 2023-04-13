package tools

import (
	"collect_web/conf"
	"github.com/go-resty/resty/v2"
	"time"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/6
 * @Desc: sendhttp.go
**/

var HttpSendClient = resty.New().SetTimeout(conf.DefaultTimeOut * time.Second)

func HttpGetRes(url string, body interface{}, method string, result interface{}) (err error, resp *resty.Response) {
	resp, err = HttpSendClient.R().
		//EnableTrace().
		SetBody(body).
		SetHeader("Content-Type", "application/json"). // request header
		ForceContentType("application/json").          // 强制返回格式
		SetResult(result).                             // or SetResult(AuthSuccess{}).
		//SetError(&AuthError{}).       // or SetError(AuthError{}).

		Execute(method, url)
	return
}
