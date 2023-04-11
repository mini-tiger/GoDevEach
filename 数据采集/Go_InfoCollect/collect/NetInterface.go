package collect

import (
	"collect_web/log"
	"collect_web/tools"
	"fmt"
	"net"
	"strings"
)

//Index			序号
//MTU			最大传输单元
//Name			名称
//HardwareAddr	硬件地址
//Flags			标志
type NetInterface struct {
	Index           int        `json:"index"`
	MTU             int        `json:"mtu"`
	Name            string     `json:"name"`
	HardwareAddr    string     `json:"hardwareAddr"`
	Flags           net.Flags  `json:"flags"`
	Ips             []net.Addr `json:"ips"`
	MulticastAdders []net.Addr `json:"multicastAdders"`
}
type NetInterfaces struct {
	Net []*NetInterface `json:"net"`
}

func GetNetIfaces() GetInfoInter {
	return new(NetInterfaces)
}

func (ni *NetInterfaces) GetName() string {
	return "NetInterface"
}

//网络接口信息采集
func (ni *NetInterfaces) GetInfo(wraplog *log.Wraplog) (interface{}, ErrorCollect) {
	var errors tools.MapStr = make(map[string]interface{})
	interfaces, err := net.Interfaces()
	if err != nil {
		errors.Set("GetInterface", err)
		return nil, ErrorCollect(errors)
	}
	ni.Net = make([]*NetInterface, 0, len(interfaces))
	for _, n := range interfaces {

		var addr []net.Addr
		var maddr []net.Addr
		var stringErrs strings.Builder
		if addr, err = n.Addrs(); err != nil {
			stringErrs.WriteString(fmt.Sprintf("getaddr err:%s;", err))
		}
		if maddr, err = n.MulticastAddrs(); err != nil {
			stringErrs.WriteString(fmt.Sprintf("getMulticastAddrs err:%s;", err))
		}
		temp := &NetInterface{
			Index:           n.Index,
			MTU:             n.MTU,
			Name:            n.Name,
			HardwareAddr:    n.HardwareAddr.String(),
			Flags:           n.Flags,
			Ips:             addr,
			MulticastAdders: maddr,
		}
		if stringErrs.Len() > 0 {
			errors.Set(n.Name, stringErrs.String())
		}
		ni.Net = append(ni.Net, temp)
	}
	if len(errors) > 0 {
		return ni, ErrorCollect(errors)
	}
	return ni, nil
}

//func main() {
//	interfaces:=GetNetInterfaceInfo()
//	fmt.Println(interfaces)
//	for _, inter := range interfaces {
//		fmt.Println("--------------------")
//		fmt.Println("接口名称:",inter.Name)
//		fmt.Println("最大传送单元:",inter.MTU)
//		fmt.Println("接口标志:",inter.Flags)
//		fmt.Println("接口地址:",inter.HardwareAddr)
//	}
//}
