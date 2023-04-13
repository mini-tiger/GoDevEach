package tools

import (
	"collect_web/log"
	"fmt"
	"net"
	"net/url"
	"time"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/13
 * @Desc: tcpconn.go
**/

func TcpTry(addr net.IP, port string, timeout int64) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(addr.To4().String(), port), time.Duration(timeout)*time.Second)

	if err != nil {
		//fmt.Println("连接失败：", err)
		return err
	}
	conn.Close()
	return nil
}

type IpCidr struct {
	Ips   []*net.IPNet
	AllIP []net.IP
}

func (i *IpCidr) getNicIps() {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		log.Glog.Error(fmt.Sprintf("fail to get net interfaces ipAddress: %v\n", err))
		return
	}

	for _, address := range interfaceAddr {
		ipNet, isVailIpNet := address.(*net.IPNet)
		// 检查ip地址判断是否回环地址
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {

				i.Ips = append(i.Ips, ipNet)
			}

		}
	}
	return
}
func (i *IpCidr) GetAllIP() {
	i.getNicIps()
	i.getAllIP()
}
func (i *IpCidr) getAllIP() {
	for _, ipc := range i.Ips {
		//ip := "192.168.1.0"
		//mask := "255.255.255.0"

		ipAddr := ipc.IP
		//ipMask := net.IPMask(net.ParseIP(mask).To4())
		ipMask := ipc.Mask

		ipNet := &net.IPNet{
			IP:   ipAddr.Mask(ipMask),
			Mask: ipMask,
		}

		startIP := ipNet.IP
		endIP := make(net.IP, len(startIP))
		copy(endIP, startIP)
		for ii := 0; ii < len(ipNet.Mask); ii++ {
			endIP[ii] |= ^ipNet.Mask[ii]
		}

		startInt := ipToInt(startIP.To4())
		endInt := ipToInt(endIP.To4())
		mask, _ := ipNet.Mask.Size()
		log.Glog.Info(fmt.Sprintf("ipcidr: %v/%v", startIP, mask))
		for ii := startInt + 1; ii < endInt; ii++ {
			//fmt.Println(intToIP(i))
			i.AllIP = append(i.AllIP, intToIP(ii))
		}

	}

}

func ipToInt(ip net.IP) uint32 {
	return (uint32(ip[0]) << 24) | (uint32(ip[1]) << 16) | (uint32(ip[2]) << 8) | uint32(ip[3])
}

func intToIP(i uint32) net.IP {
	return net.IPv4(byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
}

func UrlJoin(ip net.IP, port string, relativePath string) string {
	baseURL := fmt.Sprintf("http://%s:%s", ip.To4().String(), port)

	// 解析基本URL
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Glog.Error(fmt.Sprintf("Error parsing base URL:", err))
		return ""
	}

	// 解析相对路径
	relativeURL, err := url.Parse(relativePath)
	if err != nil {
		log.Glog.Error(fmt.Sprintf("Error parsing relative path:", err))
		return ""
	}

	// 拼接URL
	u = u.ResolveReference(relativeURL)
	return u.String()
}
