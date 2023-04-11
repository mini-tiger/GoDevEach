package main

import (
	"collect_web/collect"
	"collect_web/conf"
	"collect_web/log"
	"collect_web/service"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//os.Setenv("SPF_ID", "13") // 通常是在应用程序之外完成的

	// init zap log
	log.InitLog("CollectAgent", "agent.log")
	log.InitWrapLog()

	//ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	go func() {
		// 启动后发送数据
		lx := new(service.LinuxMetric)
		lx.RegMetrics()
		var lxface collect.GetMetricInter = lx
		data := lxface.GetMetrics(conf.GlobalCtx)
		b, err := json.Marshal(data)
		if err != nil {
			log.Glog.Error(fmt.Sprintf("jsonMarshal Err:%v", err))
			return
		}
		log.Glog.Info(fmt.Sprintf("collect SN:%s", lxface.FormatData()["SN"]))
		log.Glog.Debug(fmt.Sprintf("collect data:%v", string(b)))
		log.Glog.Error(fmt.Sprintf("collect err:%v", lxface.GetErrors()))
		if err = service.SendHttpRes(lxface.FormatData()); err != nil {
			log.Glog.Error(fmt.Sprintf("cron sendHttp :%v", err))
		}
	}()

	go service.CronStart(conf.GlobalCtx)
	//fmt.Println(runtime.Caller(0))
	service.StartWeb()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGHUP) // https://colobu.com/2015/10/09/Linux-Signals/  kill -2  kill -1

	select {
	// 等待退出信号
	case <-exitChan:
		log.Glog.Info("get exit signal")
		service.HttpSrv.Shutdown(context.TODO())
		//ctx, _ := context.WithTimeout(valueCtx1, 1*time.Second)
		conf.GlobalCtx.Done()
		return
	}

}
