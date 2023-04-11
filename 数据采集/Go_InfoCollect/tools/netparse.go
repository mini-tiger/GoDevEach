package tools

import (
	"net"
	"net/url"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/7
 * @Desc: netparse.go
**/
func Url2Host(urlstr string) (host string, port string, err error) {
	// 放置 HTTP URL
	//urlString := "http://example.com:8080/path"

	// 解析 URL
	parsedUrl, err := url.Parse(urlstr)
	if err != nil {
		panic(err)
	}

	// 获取主机名和端口号
	host, port, err = net.SplitHostPort(parsedUrl.Host)
	if err != nil {
		return
	}
	return
	// 解析端口号
	//portInt, err := net.LookupPort("tcp", port)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// 输出 IP 地址和端口号
	//fmt.Printf("IP: %s\n", net.ParseIP(host))
	//fmt.Printf("Port: %d\n", portInt)
}
