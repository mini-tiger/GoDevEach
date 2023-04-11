package service

import (
	"collect_web/collect"
	"collect_web/conf"
	"collect_web/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/4
 * @Desc: cronentry.go
**/

type CronEntry struct {
	ctx context.Context
	lx  collect.GetMetricInter
}

func (c *CronEntry) Run() {

	data := c.lx.GetMetrics(conf.GlobalCtx)

	// log
	b, err := json.Marshal(data)
	if err != nil {
		log.Glog.Error(fmt.Sprintf("jsonMarshal Err:%v", err))
		return
	}
	log.Glog.Info(fmt.Sprintf("collect SN:%s", c.lx.FormatData()["SN"]))
	log.Glog.Debug(fmt.Sprintf("cron collect data:%v", string(b)))
	log.Glog.Error(fmt.Sprintf("cron collect errs:%v", c.lx.GetErrors()))
	//
	if err = SendHttpRes(c.lx.FormatData()); err != nil {
		log.Glog.Error(fmt.Sprintf("cron sendHttp :%v", err))
	}

}

func CronStart(ctx context.Context) {

	cl := new(log.Wraplog)
	cl.ZapLog = log.GlogSkip1
	lx := new(LinuxMetric)
	lx.RegMetrics()
	var lxface collect.GetMetricInter = lx

	defaultSecond := cron.Every(time.Duration(time.Minute * 1)) //
	c := cron.New(cron.WithLogger(cl))
	_ = c.Schedule(defaultSecond, &CronEntry{ctx: ctx, lx: lxface})
	c.Start()

}
