package collect

import (
	"collect_web/log"
	"collect_web/tools"
	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/host"
)

//type InfoStat struct {
//	Hostname             主机名称
//	Uptime               开机时间
//	BootTime             boot时间
//	Procs                进程数目
//	OS                   操作系统 如freebsd, linux
//	Platform             如： ubuntu, linuxmint
//	PlatformFamily       如: debian, rhel
//	PlatformVersion      操作系统版本
//	KernelVersion        操作系统内核版本
//	KernelArch           内核架构
//	VirtualizationSystem 虚拟系统
//	VirtualizationRole   虚拟角色 guest or host
//	HostID               hostid  // ex: uuid
//}
const HostInfoStatStr = "HostInfoStat"

type HostInfoStat struct {
	SN        string
	hostStat  *host.InfoStat
	Product   *ghw.ProductInfo
	Bios      *ghw.BIOSInfo
	BaseBoard *ghw.BaseboardInfo
}

func GetHostInfo() GetInfoInter {
	return new(HostInfoStat)
}
func (h *HostInfoStat) GetName() string {
	return HostInfoStatStr

}
func (h *HostInfoStat) GetInfo(wlog *log.Wraplog) (interface{}, ErrorCollect) {
	var errors tools.MapStr = make(map[string]interface{})

	HostInfo, err := host.Info()
	if err != nil {
		errors.Set("hostInfoStat", err)
		//return nil, ErrorCollect(errors)
	}
	bios, err := ghw.BIOS(&ghw.WithOption{Alerter: wlog})
	if err != nil {
		errors.Set("bios", err)
		//return nil, ErrorCollect(errors)
	}
	product, err := ghw.Product(&ghw.WithOption{Alerter: wlog})
	if err != nil {
		errors.Set("product", err)
		//return nil, ErrorCollect(errors)
	}

	baseboard, err := ghw.Baseboard(&ghw.WithOption{Alerter: wlog})
	if err != nil {
		errors.Set("baseboard", err)
		//return nil,ErrorCollect(errors)
	}
	h.SN = product.SerialNumber
	h.BaseBoard = baseboard
	h.hostStat = HostInfo
	h.Bios = bios
	h.Product = product
	if len(errors) > 0 {
		return h, nil
	}
	//return HostInfo
	return h, ErrorCollect(errors)
}

//
//func GetProduct(l *log.Wraplog) GetInfoInter {
//	return &HostProduct{}
//}
//
//func (p *HostProduct) GetInfo() (errors ErrorCollect) {
//	errors = make(map[string]interface{})
//	//var errors ErrorCollect
//	defer func() {
//		if err := recover(); err != nil {
//			errors.Set("product", err)
//		}
//
//	}()
//
//	product, err := ghw.Product(&ghw.WithOption{Alerter: p})
//	if err != nil {
//		errors.Set("product", err)
//		return nil, ErrorCollect(errors)
//	}
//	return product, nil
//}
//func (p *HostBios) GetInfo() (interface{}, ErrorCollect) {
//	var errors tools.MapStr = make(map[string]interface{})
//	defer func() {
//		if err := recover(); err != nil {
//			errors.Set("bios", err)
//		}
//
//	}()
//	bios, err := ghw.BIOS()
//	if err != nil {
//		errors.Set("bios", err)
//		return nil, ErrorCollect(errors)
//	}
//	return bios, nil
//}
