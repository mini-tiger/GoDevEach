package service

import (
	"collect_web/tools"
	"github.com/go-resty/resty/v2"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/6
 * @Desc: sendhttp.go
**/

func HttpGetRes(url string, body interface{}) (err error, resp *resty.Response) {
	resp, err = tools.HttpSendCleint.R().
		//EnableTrace().
		SetBody(body).
		SetHeader("Content-Type", "application/json"). // request header
		ForceContentType("application/json").          // 强制返回格式
		//SetResult(respbody).                           // or SetResult(AuthSuccess{}).
		//SetError(&AuthError{}).       // or SetError(AuthError{}).
		Post(url)
	return
}
