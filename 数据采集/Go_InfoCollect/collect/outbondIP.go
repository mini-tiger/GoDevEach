package collect

import (
	"collect_web/conf"
	"collect_web/log"
	"collect_web/tools"
	"net"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/7
 * @Desc: outbondIP.go
**/

type OutboundIP struct {
	IP string `json:"IP"`
}

func GetOutboundIP() GetInfoInter {
	return new(OutboundIP)
}

func (i *OutboundIP) GetName() string {
	return "OutboundIP"
}

func (i *OutboundIP) GetInfo(wraplog log.WrapLogInter) (interface{}, ErrorCollect) {
	var errors tools.MapStr = make(map[string]interface{})
	//host, port, err := tools.Url2Host(conf.SendHttpServer)
	//if err != nil {
	//	errors.Set("OutboundIP", err)
	//	return nil, ErrorCollect(errors)
	//}
	//fmt.Printf("abcd:%v %v %v\n", host, port, err)
	//time.Sleep(10)
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(conf.GetServerAddr(), conf.ServerPort), conf.DefaultTimeOut)

	if err != nil {
		errors.Set("OutboundIP", err)
		return nil, ErrorCollect(errors)
	}
	defer conn.Close()

	//localAddr := conn.LocalAddr().(*net.UDPAddr)
	localAddr := conn.LocalAddr().(*net.TCPAddr)
	i.IP = localAddr.IP.String()
	//fmt.Println(localAddr.IP)
	//return localAddr.IP
	//fmt.Printf("%T,%s\n",localAddr.IP,localAddr.IP)
	if len(errors) > 0 {
		return nil, ErrorCollect(errors)
	}
	return i, nil
}
