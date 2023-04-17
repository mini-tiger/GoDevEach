package collect

import (
	"collect_web/log"
	"context"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/3/31
 * @Desc: iface.go
**/

type ErrorCollect map[string]interface{}

type GetInfoInter interface {
	GetInfo(wraplog log.WrapLogInter) (interface{}, ErrorCollect)
	GetName() string
}

//type GetHostInfoInter interface {
//	GetInfo() (interface{}, ErrorCollect)
//}

type GetMetricInter interface {
	GetMetrics(ctx context.Context) interface{}
	GetErrors() map[string]interface{}
	FormatData() map[string]interface{}
}
