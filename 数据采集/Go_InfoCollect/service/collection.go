package service

import (
	"collect_web/collect"
	"collect_web/conf"
	"collect_web/log"
	"collect_web/tools"
	"context"
	"fmt"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/3/31
 * @Desc: collection.go
**/

//type Linux struct {
//	Host         *host.InfoStat         `json:"host"`
//	Cpu          collect.Cpu            `json:"cpu"`
//	Disk         collect.Disk           `json:"disk"`
//	Memory       collect.Memory         `json:"memory"`
//	Network      collect.Network        `json:"network"`
//	NetInterface []collect.NetInterface `json:"netInterface"`
//	Process      []collect.Process      `json:"process"`
//	Service      collect.Service        `json:"service"`
//	Device       collect.Device         `json:"device"`
//	Log          collect.Log            `json:"log"`
//	FireWall     collect.FireWall       `json:"firewall"`
//	DataAmount   collect.DataAmount     `json:"dataAmount"`
//	Application  collect.Application    `json:"application"`
//	Errors       tools.MapStr
//}
//
//type HostHWInfo struct {
//	Cpu    interface{} `json:"cpu"`
//	Disk   interface{} `json:"disk"`
//	Memory interface{} `json:"memory"`
//	//Network      interface{}        `json:"network"`
//	NetInterface interface{} `json:"netInterface"`
//	OutBoundIP   interface{} `json:"outBoundIP"`
//	Baseboard    interface{} `json:"baseboard"`
//	Gpu          interface{} `json:"gpu"`
//}
//type HostInfo struct {
//	HostStat interface{} `json:"hostStat"`
//	Product  interface{} `json:"product"`
//	Bios     interface{} `json:"bios"`
//}

type LinuxMetric struct {
	CollectErrors tools.MapStr
	MetricsFn     []collect.GetInfoInter
	metricsData   tools.MapStr
}

func (l *LinuxMetric) GetErrors() map[string]interface{} {
	return l.CollectErrors.ToMapInterface()
}

func (l *LinuxMetric) CollectionWithCtx(ctx context.Context) interface{} {
	finchan := make(chan interface{}, 0)

	go func() {
		//time.Sleep(120 * time.Second)
		finchan <- l.getMetrics()
		//finchan <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		//fmt.Println(111111111111)
		//log.Glog.Sugar().Error("Collect Err:%s", ctx.Err())
		l.CollectErrors.Set("ctx timeout", ctx.Err())
		return nil
	case data := <-finchan:
		return data

	}

}

func (l *LinuxMetric) RegMetrics() {
	//var errors tools.MapStr = make(map[string]interface{})
	//l.CollectErrors = errors
	var cpu collect.GetInfoInter = collect.GetCpu()
	var mem collect.GetInfoInter = collect.GetMemory()
	var hostinfo collect.GetInfoInter = collect.GetHostInfo()
	var netinfo collect.GetInfoInter = collect.GetNetIfaces()
	var diskinfo collect.GetInfoInter = collect.GetDisk()
	var outip collect.GetInfoInter = collect.GetOutboundIP()
	var gpu collect.GetInfoInter = collect.GetGPU()

	l.MetricsFn = append(l.MetricsFn, cpu, mem, hostinfo, netinfo, diskinfo, outip, gpu)
	//l.MetricsFn = append(l.MetricsFn, hostinfo)
	l.CollectErrors = make(map[string]interface{})

}

func (l *LinuxMetric) GetMetrics(gctx context.Context) interface{} {
	ctx, _ := context.WithTimeout(gctx, conf.DefaultTimeOut)
	return l.CollectionWithCtx(ctx)

}

func (l *LinuxMetric) wrapFn(lg *log.Wraplog, inter collect.GetInfoInter) (interface{}, collect.ErrorCollect) {
	defer func() {
		if err := recover(); err != nil {
			l.CollectErrors.Set(inter.GetName(), fmt.Sprintf("panic %v", err))
		}
	}()

	return inter.GetInfo(lg)
}

func (l *LinuxMetric) getMetrics() interface{} {
	//var resData tools.MapStr = make(map[string]interface{})
	l.metricsData = make(map[string]interface{})
	l.CollectErrors.Reset()

	for _, fn := range l.MetricsFn {
		//data, err := fn.GetInfo(log.Wlog)
		data, err := l.wrapFn(log.Wlog, fn)

		l.metricsData.Set(fn.GetName(), data)
		//resData.Set(fn.GetName(), string(bytestr))
		//log.Glog.Debug(fmt.Sprintf("Metric : %s,Data: %v,err:%v", fn.GetName(), string(bytestr), err))

		l.CollectErrors.Merge(map[string]interface{}(err))
	}
	//l.Result = resData
	return l.metricsData
}

func (l *LinuxMetric) FormatData() map[string]interface{} {
	var sn string
	if val, b := l.metricsData.Get(collect.HostInfoStatStr); b {
		if m, ok := val.(*collect.HostInfoStat); ok {
			sn = m.SN
		}
	}

	return map[string]interface{}{"MetricData": l.metricsData, "Errors": l.GetErrors(), "SN": sn}
}
