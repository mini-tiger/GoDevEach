package service

import (
	"collect_web/collect"
	"collect_web/conf"
	"collect_web/log"
	"collect_web/tools"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/panjf2000/ants/v2"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/13
 * @Desc: scaner.go
**/

const maxPoolNum = 100

type taskFunc func()

type CollectSend struct {
	TcpConnSuccess  chan net.IP
	HttpAuthSuccess chan net.IP
	HttpSendSuccess chan struct{}
}

type CollectSendInter interface {
	ScanIPHttpSend()
	ScanIPHttpAuth()
	LoopScanIpTcp()
	GetHttpAuthSuccess() chan<- net.IP
	GetHttpSendSuccess() <-chan struct{}
}

var _ CollectSendInter = &CollectSend{}

func NewCollectSend() *CollectSend {
	var TcpConnSuccess chan net.IP = make(chan net.IP, runtime.NumCPU()*10)
	var HttpAuthSuccess chan net.IP = make(chan net.IP, runtime.NumCPU()*10)
	var HttpSendSuccess chan struct{} = make(chan struct{}, 0)
	return &CollectSend{
		TcpConnSuccess:  TcpConnSuccess,
		HttpAuthSuccess: HttpAuthSuccess,
		HttpSendSuccess: HttpSendSuccess,
	}
}

func (cs *CollectSend) taskFuncWrapper(addr net.IP, tcpchan chan net.IP, wg *sync.WaitGroup) taskFunc {
	return func() {
		defer wg.Done()
		if err := tools.TcpTry(addr, conf.ServerPort, 3); err != nil {
			log.Glog.Debug(fmt.Sprintf("ip:%s tcp conn fail", addr.To4().String()))
		} else {
			tcpchan <- addr
			log.Glog.Info(fmt.Sprintf("!!!!!!ip:%s tcp conn sucess", addr.To4().String()))
		}

	}
}

func (cs *CollectSend) ScanIpTcp() {
	ipc := &tools.IpCidr{}
	ipc.GetAllIP()
	//fmt.Println(tools.ContainsIP(ipc.AllIP, "172.22.50.25"))
	//fmt.Println(ipc.AllIP)
	var wg sync.WaitGroup
	// 创建一个容量为10的goroutine池
	p, _ := ants.NewPool(maxPoolNum)
	defer p.Release() // xxx 使用完必须释放
	for _, ipt := range ipc.AllIP {
		wg.Add(1)
		err := p.Submit(cs.taskFuncWrapper(ipt, cs.TcpConnSuccess, &wg))
		if err != nil {
			log.Glog.Error(fmt.Sprintf("ants submit task Fail:%v", err))
		}
	}
	wg.Wait()
}

func (cs *CollectSend) LoopScanIpTcp() {
	for {
		cs.ScanIpTcp()
		time.Sleep(60 * time.Second)
	}
}

type HttpAuthResp struct {
	Msg  string `json:"msg"`
	Code uint64 `json:"code"`
}

func (cs *CollectSend) ScanIPHttpAuth() {
	for {
		select {
		case ip := <-cs.TcpConnSuccess:
			checkurl := tools.UrlJoin(ip, conf.ServerPort, conf.ServerAuthUrlSuffix)
			if checkurl == "" {
				return
			}
			httpauthResp := &HttpAuthResp{}
			err, resp := tools.HttpGetRes(checkurl, nil, resty.MethodGet, httpauthResp)
			if err != nil {
				log.Glog.Error(fmt.Sprintf("Access AuthUrl:%s Err:%v", checkurl, err))
			} else {
				log.Glog.Debug(fmt.Sprintf("Access AuthUrl:%s  Resp:%v", checkurl, resp))
				//fmt.Println(httpauthResp)
				if httpauthResp.Code == 200 && strings.Contains(strings.ToLower(httpauthResp.Msg), "collectservice") {
					log.Glog.Info(fmt.Sprintf("Access AuthUrl:%s  Sucess:%v", checkurl, httpauthResp))
					cs.HttpAuthSuccess <- ip
				} else {
					log.Glog.Info(fmt.Sprintf("Access AuthUrl:%s  Fail:%v", checkurl, httpauthResp))
				}
			}
		}
	}
}

type HttpSendResp struct {
	Msg  string `json:"msg"`
	Code uint64 `json:"code"`
}

func (cs *CollectSend) GetHttpAuthSuccess() chan<- net.IP {
	return cs.HttpAuthSuccess
}

func (cs *CollectSend) GetHttpSendSuccess() <-chan struct{} {
	return cs.HttpSendSuccess
}
func (cs *CollectSend) ScanIPHttpSend() {
	for {
		select {
		case ip := <-cs.HttpAuthSuccess:
			sendurl := tools.UrlJoin(ip, conf.ServerPort, conf.ServerSendUrlSuffix)
			if sendurl == "" {
				return
			}
			httpsendResp := &HttpSendResp{}

			lx := new(LinuxMetric)
			lx.Wlog = log.Wlog
			lx.RegMetrics()
			_ = conf.SetServerAddr(ip.To4().String()) // 写入配置 并 全局变量 outip

			var lxface collect.GetMetricInter = lx
			lxface.GetMetrics(conf.GlobalCtx)
			data := lxface.FormatData()

			err, resp := tools.HttpGetRes(sendurl, data, resty.MethodPost, httpsendResp)
			if err != nil {
				log.Glog.Error(fmt.Sprintf("Access SendUrl:%s Err:%v", sendurl, err))
			} else {
				log.Glog.Debug(fmt.Sprintf("Access SendUrl:%s  Resp:%v", sendurl, resp))
				//fmt.Println(httpauthResp)
				if httpsendResp.Code == 200 && strings.Contains(strings.ToLower(httpsendResp.Msg), "ok") {
					log.Glog.Info(fmt.Sprintf("Access SendUrl:%s  Sucess:%v", sendurl, httpsendResp))
					cs.HttpSendSuccess <- struct{}{}
				} else {
					log.Glog.Info(fmt.Sprintf("Access SendUrl:%s  Fail:%v", sendurl, httpsendResp))
				}
			}
		}
	}
}
