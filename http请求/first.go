package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

func main() {
	client := resty.New()

	resp, err := client.R().Get("https://baidu.com")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response Info:")
	fmt.Println("Status Code:", resp.StatusCode())
	fmt.Println("Status:", resp.Status())
	fmt.Println("Proto:", resp.Proto())
	fmt.Println("Time:", resp.Time())
	fmt.Println("Received At:", resp.ReceivedAt())
	fmt.Println("Size:", resp.Size())
	fmt.Println("Headers:")
	for key, value := range resp.Header() {
		fmt.Println(key, "=", value)
	}
	fmt.Println("Cookies:")
	for i, cookie := range resp.Cookies() {
		fmt.Printf("cookie%d: name:%s value:%s\n", i, cookie.Name, cookie.Value)
	}

	/*
		DNSLookup：DNS 查询时间，如果提供的是一个域名而非 IP，就需要向 DNS 系统查询对应 IP 才能进行后续操作；
		ConnTime：获取一个连接的耗时，可能从连接池获取，也可能新建；
		TCPConnTime：TCP 连接耗时，从 DNS 查询结束到 TCP 连接建立；
		TLSHandshake：TLS 握手耗时；
		ServerTime：服务器处理耗时，计算从连接建立到客户端收到第一个字节的时间间隔；
		ResponseTime：响应耗时，从接收到第一个响应字节，到接收到完整响应之间的时间间隔；
		TotalTime：整个流程的耗时；
		IsConnReused：TCP 连接是否复用了；
		IsConnWasIdle：连接是否是从空闲的连接池获取的；
		ConnIdleTime：连接空闲时间；
		RequestAttempt：请求执行流程中的请求次数，包括重试次数；
		RemoteAddr：远程的服务地址，IP:PORT格式。

	*/
	ti := resp.Request.TraceInfo()
	fmt.Println("Request Trace Info:")
	fmt.Println("DNSLookup:", ti.DNSLookup)
	fmt.Println("ConnTime:", ti.ConnTime)
	fmt.Println("TCPConnTime:", ti.TCPConnTime)
	fmt.Println("TLSHandshake:", ti.TLSHandshake)
	fmt.Println("ServerTime:", ti.ServerTime)
	fmt.Println("ResponseTime:", ti.ResponseTime)
	fmt.Println("TotalTime:", ti.TotalTime)
	fmt.Println("IsConnReused:", ti.IsConnReused)
	fmt.Println("IsConnWasIdle:", ti.IsConnWasIdle)
	fmt.Println("ConnIdleTime:", ti.ConnIdleTime)
	fmt.Println("RequestAttempt:", ti.RequestAttempt)
	//fmt.Println("RemoteAddr:", ti.RemoteAddr.String())
}
