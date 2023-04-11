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

var HttpSendCleint = resty.New().SetTimeout(conf.DefaultTimeOut * time.Second)
