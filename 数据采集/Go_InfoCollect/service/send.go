package service

import (
	"collect_web/conf"
	"collect_web/log"
	"fmt"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/10
 * @Desc: send.go
**/
func SendHttpRes(body interface{}) error {
	err, resp := HttpGetRes(conf.SendHttpUrl, body)
	if err != nil {
		//log.Glog.Error(fmt.Sprintf("SendHttpRes Err:%s", err))
		return err
	}
	log.Glog.Debug(fmt.Sprintf("SendHttpRes Resp:%v", resp))

	return nil
}
