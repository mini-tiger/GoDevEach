package collect

//
//import (
//	"errors"
//	"fmt"
//	"net"
//	"sort"
//	"strings"
//	"syscall"
//	"unsafe"
//)
//
////IP		IP地址
////MAC		MAC地址
////GateWay	网关
//type Network struct {
//	NetWork []map[string]string `json:"network"`
//	//MAC     map[string]string `json:"mac"`
//	//GateWay []GateWay         `json:"gateway"`
//}
//
////路由信息
//type rtInfo struct {
//	Dst              net.IPNet
//	Gateway, PrefSrc net.IP
//	OutputIface      uint32
//	Priority         uint32
//}
//
//type routeSlice []*rtInfo
//
////路由器信息
//type router struct {
//	ifaces []net.Interface
//	addrs  []net.IP
//	v4     routeSlice
//}
//
////网关
//type GateWay struct {
//	InterfaceName string `json:"interface_name"`
//	GateWay       string `json:"gate_way"`
//	Ip            string `json:"ip"`
//}
//
////获取路由信息
//func getRouteInfo() (*router, error) {
//	rtr := &router{}
//	tab, err := syscall.NetlinkRIB(syscall.RTM_GETROUTE, syscall.AF_INET)
//	if err != nil {
//		return nil, err
//	}
//	msgs, err := syscall.ParseNetlinkMessage(tab)
//	if err != nil {
//		return nil, err
//	}
//	for _, m := range msgs {
//		switch m.Header.Type {
//		case syscall.NLMSG_DONE:
//			break
//		case syscall.RTM_NEWROUTE:
//			rtmsg := (*syscall.RtMsg)(unsafe.Pointer(&m.Data[0]))
//			attrs, err := syscall.ParseNetlinkRouteAttr(&m)
//			if err != nil {
//				return nil, err
//			}
//			routeInfo := rtInfo{}
//			rtr.v4 = append(rtr.v4, &routeInfo)
//			for _, attr := range attrs {
//				switch attr.Attr.Type {
//				case syscall.RTA_DST:
//					routeInfo.Dst.IP = net.IP(attr.Value)
//					routeInfo.Dst.Mask = net.CIDRMask(int(rtmsg.Dst_len), len(attr.Value)*8)
//				case syscall.RTA_GATEWAY:
//					routeInfo.Gateway = net.IPv4(attr.Value[0], attr.Value[1], attr.Value[2], attr.Value[3])
//				case syscall.RTA_OIF:
//					routeInfo.OutputIface = *(*uint32)(unsafe.Pointer(&attr.Value[0]))
//				case syscall.RTA_PRIORITY:
//					routeInfo.Priority = *(*uint32)(unsafe.Pointer(&attr.Value[0]))
//				case syscall.RTA_PREFSRC:
//					routeInfo.PrefSrc = net.IPv4(attr.Value[0], attr.Value[1], attr.Value[2], attr.Value[3])
//				}
//			}
//		}
//	}
//	sort.Slice(rtr.v4, func(i, j int) bool {
//		return rtr.v4[i].Priority < rtr.v4[j].Priority
//	})
//	ifaces, err := net.Interfaces()
//	if err != nil {
//		return nil, err
//	}
//	for i, iface := range ifaces {
//		if i != iface.Index-1 {
//			break
//		}
//		if iface.Flags&net.FlagUp == 0 {
//			continue
//		}
//		rtr.ifaces = append(rtr.ifaces, iface)
//		ifaceAddrs, err := iface.Addrs()
//		if err != nil {
//			return nil, err
//		}
//		var addrs net.IP
//		for _, addr := range ifaceAddrs {
//			if inet, ok := addr.(*net.IPNet); ok {
//				if v4 := inet.IP.To4(); v4 != nil {
//					if addrs == nil {
//						addrs = v4
//					}
//				}
//			}
//		}
//		rtr.addrs = append(rtr.addrs, addrs)
//	}
//	return rtr, nil
//}
//
////获取网关信息
//func getGateWay() []GateWay {
//	newRoute, err := getRouteInfo()
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	var gateWays []GateWay
//	for _, rt := range newRoute.v4 {
//		var gateway GateWay
//		if rt.Gateway != nil {
//			gateway.InterfaceName = newRoute.ifaces[rt.OutputIface-1].Name
//			gateway.GateWay = rt.Gateway.String()
//			gateway.Ip = newRoute.addrs[rt.OutputIface-1].String()
//			gateWays = append(gateWays, gateway)
//		}
//	}
//	return gateWays
//}
//
////获取IP地址
//func GetIpAddrs() (myIpMac []map[string]string, err error) {
//
//	//获取网络接口
//	ifaces, err := net.Interfaces()
//	if err != nil {
//		//fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
//		return nil, err
//	}
//
//	mpIpMac := make([]map[string]interface{}, 0, len(ifaces))
//
//	var stringsErr strings.Builder
//	//遍历网络接口
//	for index, netInterface := range ifaces {
//		addrs, err := netInterface.Addrs()
//		if err != nil {
//			//fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
//			stringsErr.WriteString(fmt.Sprintf("index:%d , %s\n", index, err.Error()))
//			continue
//		}
//
//		macAddr := netInterface.HardwareAddr.String()
//		if len(macAddr) == 0 {
//			stringsErr.WriteString(fmt.Sprintf("index:%d , %s\n", index, "mac len 0"))
//			continue
//		}
//
//		//获取网络接口的地址
//		for _, a := range addrs {
//			mpIpMac[i.Name] = a.String()
//			//fmt.Printf("%v : %s \n", i.Name, a.String())
//		}
//	}
//	if stringsErr.Len() > 0 {
//		return mpIp, errors.New(stringsErr.String())
//	} else {
//		return mpIp, nil
//	}
//
//}
//
////获取硬件地址
//func GetMacAddrs() (map[string]string, error) {
//	mpMac := make(map[string]string)
//	netInterfaces, err := net.Interfaces()
//	if err != nil {
//		return nil, err
//	}
//	for _, netInterface := range netInterfaces {
//		macAddr := netInterface.HardwareAddr.String()
//		if len(macAddr) == 0 {
//			continue
//		}
//		mpMac[netInterface.Name] = macAddr
//	}
//	return mpMac
//}
//
//func GetNetInfo() GetInfoInter {
//	return new(Network)
//}
//
////获取网络信息
//
//func (n *Network) GetInfo() ErrorCollect {
//	var errors ErrorCollect
//	ip, err := GetIpAddrs()
//	if err != nil {
//		errors["ip"] = err
//	} else {
//		n.IP = ip
//	}
//	mac := GetMacAddrs()
//	gateways := getGateWay()
//	netinfo := Network{
//		IP:      ip,
//		MAC:     mac,
//		GateWay: gateways,
//	}
//	return netinfo
//
//}
