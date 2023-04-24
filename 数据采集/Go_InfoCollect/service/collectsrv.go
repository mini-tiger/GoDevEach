package service

import (
	"collect_web/conf"
	"collect_web/log"
	"fmt"
	"net"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/21
 * @Desc: collectsrv.go
**/

var CollectScan CollectSendInter
var CollectFlag bool = false

func CollectSrv() {
	conf.GetEnv()
	_ = conf.ReadYAML()

	// init zap log
	log.InitLog("CollectAgent", "agent.log")
	log.InitWrapLog()

	CollectScan = NewCollectSend()

	if ip := conf.GetServerAddr(); ip != "" {
		log.Glog.Info(fmt.Sprintf("!!! Read Success , Current ServerAddr:%s ServerPort:%v ", ip, conf.GetServerPort()))
		go CollectScan.ScanIPHttpSend()
		CollectScan.GetHttpAuthSuccess() <- net.ParseIP(ip)
	} else {
		log.Glog.Info(fmt.Sprintf("!!! Read Fail ,Current ServerPort:%v ", conf.GetServerPort()))
		go CollectScan.ScanIPHttpSend()
		go CollectScan.ScanIPHttpAuth()
		go CollectScan.LoopScanIpTcp()
	}
}
