package main

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"log"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"
)

var L *sync.WaitGroup = new(sync.WaitGroup)
var ip1 chan []string = make(chan []string, 0)
var ip2 chan string = make(chan string, 0)
var ip3 chan string = make(chan string, 0)

func errlog(err error) {
	_, _, fileno, _ := runtime.Caller(1)
	log.Printf("RowNo: %d ,ERROR:%s\n", fileno, err)
}

func GetIntranetIp() {
	//L.Add(1)
	//defer func() {
	//	L.Done()
	//}()
	netIP := make([]string, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		ip1 <- netIP
		errlog(err)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println("ip:", ipnet.IP.String())
				netIP = append(netIP, ipnet.IP.String())
			}

		}
	}
	ip1 <- netIP
}

func get_external() {
	//L.Add(1)
	//defer func() {
	//	L.Done()
	//}()
	rq := httplib.Get("http://myexternalip.com/raw").SetTimeout(time.Duration(5)*time.Second, time.Duration(5)*time.Second)
	resp, err := rq.Response()
	if err != nil {
		ip2 <- ""
		errlog(err)
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ip2 <- ""
		errlog(err)
		return
	}
	tmp := string(content)
	tmp = strings.TrimSpace(tmp)
	tmp = strings.Trim(tmp, "\n")

	//resp, err := http.Get("http://myexternalip.com/raw")
	//if err != nil {
	//	ip2 <- ""
	//	errlog(err)
	//}
	//defer resp.Body.Close()
	//content, _ := ioutil.ReadAll(resp.Body)
	////buf := new(bytes.Buffer)
	////buf.ReadFrom(resp.Body)
	////s := buf.String()
	//tmp := string(content)
	//tmp = strings.TrimSpace(tmp)
	//tmp = strings.Trim(tmp, "\n")
	ip2 <- tmp
	//fmt.Println(string(content))
}

func GetOutboundIP() {
	//L.Add(1)
	//defer func() {
	//	L.Done()
	//}()
	//conn, err := net.Dial("udp", "8.8.8.8:80")
	conn, err := net.DialTimeout("tcp", "www.baidu.com:80", time.Duration(15)*time.Second)
	if err != nil {
		ip3 <- ""

		errlog(err)
	}
	defer conn.Close()

	//localAddr := conn.LocalAddr().(*net.UDPAddr)
	localAddr := conn.LocalAddr().(*net.TCPAddr)
	//fmt.Println(localAddr.IP)
	//return localAddr.IP
	//fmt.Printf("%T,%s\n",localAddr.IP,localAddr.IP)
	tmp := fmt.Sprintf("%s", localAddr.IP)
	tmp = strings.TrimSpace(tmp)
	tmp = strings.Trim(tmp, "\n")
	ip3 <- tmp
}

func main() {

	fmt.Println("======================== 除回环以外 所有IP=================================")

	go GetIntranetIp() // 除回环以外 所有IP

	fmt.Println("======================== 公网出网IP=================================")
	go get_external()

	fmt.Println("======================== 首选的出站内网IP地址=================================")
	go GetOutboundIP()

	//L.Wait()

	fmt.Printf("%-16s%v\n", "除回环以外所有IP    :", <-ip1)
	fmt.Printf("%-16s%s\n", "公网出网IP          :", <-ip2) // ip2 会带有换行字符
	ip33 := <-ip3
	fmt.Printf("%-16s%s\n", "首选的出站内网IP地址:", ip33)

}
