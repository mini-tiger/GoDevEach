package main

import (
	"collect_web/cmd"
	"collect_web/conf"
	"collect_web/log"
	"collect_web/service"
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	// 不加参数 往下执行
	cmd.Execute()

	//_ = conf.ReadYAML()
	//var collectScan service.CollectSendInter = service.NewCollectSend()
	//
	//if ip := conf.GetServerAddr(); ip != "" {
	//	log.Glog.Info(fmt.Sprintf("!!! Read Success , Current ServerAddr:%s ServerPort:%v ", ip, conf.GetServerPort()))
	//	go collectScan.ScanIPHttpSend()
	//	collectScan.GetHttpAuthSuccess() <- net.ParseIP(ip)
	//} else {
	//	log.Glog.Info(fmt.Sprintf("!!! Read Fail ,Current ServerPort:%v ", conf.GetServerPort()))
	//	go collectScan.ScanIPHttpSend()
	//	go collectScan.ScanIPHttpAuth()
	//	go collectScan.LoopScanIpTcp()
	//}

	if service.CollectFlag {
		if b := <-service.CollectScan.GetHttpSendSuccess(); b {
			log.Glog.Info(fmt.Sprintf("!!!! Success Push Collect Data To %s:%s", conf.GetServerAddr(), conf.GetServerPort()))
		} else {
			log.Glog.Error(fmt.Sprintf("!!!! Fail Push Collect Data To %s:%s", conf.GetServerAddr(), conf.GetServerPort()))
		}
	}

	//ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	//go func() {
	//	// 启动后发送数据
	//	lx := new(service.LinuxMetric)
	//	lx.RegMetrics()
	//	var lxface collect.GetMetricInter = lx
	//	data := lxface.GetMetrics(conf.GlobalCtx)
	//	b, err := json.Marshal(data)
	//	if err != nil {
	//		log.Glog.Error(fmt.Sprintf("jsonMarshal Err:%v", err))
	//		return
	//	}
	//	log.Glog.Info(fmt.Sprintf("collect SN:%s", lxface.FormatData()["SN"]))
	//	log.Glog.Debug(fmt.Sprintf("collect data:%v", string(b)))
	//	log.Glog.Error(fmt.Sprintf("collect err:%v", lxface.GetErrors()))
	//	if err = service.SendHttpRes(lxface.FormatData()); err != nil {
	//		log.Glog.Error(fmt.Sprintf("cron sendHttp :%v", err))
	//	}
	//}()

	//go service.CronStart(conf.GlobalCtx)
	////fmt.Println(runtime.Caller(0))
	//service.StartWeb()
	//
	//exitChan := make(chan os.Signal, 1)
	//signal.Notify(exitChan, syscall.SIGINT, syscall.SIGHUP) // https://colobu.com/2015/10/09/Linux-Signals/  kill -2  kill -1
	//
	//select {
	//// 等待退出信号
	//case <-exitChan:
	//	log.Glog.Info("get exit signal")
	//	service.HttpSrv.Shutdown(context.TODO())
	//	//ctx, _ := context.WithTimeout(valueCtx1, 1*time.Second)
	//	conf.GlobalCtx.Done()
	//	return
	//}

}
