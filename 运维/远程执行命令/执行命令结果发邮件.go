package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"远程执行命令/funcs"
	"远程执行命令/modules"
)

func main() {

	MonitorHosts := map[string]string{"192.168.40.100": "W3b5Ev!c3"}
	var ResultHosts []*modules.HostMonitor = make([]*modules.HostMonitor, 0)

	for host, pass := range MonitorHosts {
		h := modules.NewSSHConn()
		h.Host = host
		h.Pass = pass

		if err := h.GetConn(); err != nil {
			log.Printf("GetConn Error:%v,host:%s\n", err, host)
		}

		//xxx 建立 链接 后 每个命令 单独session
		if err := h.GetSession(); err != nil {
			log.Printf("Error:%v,host:%s\n", err, host)
		}
		hm := new(modules.HostMonitor)
		hm.Host = host

		//xxx 磁盘使用率
		cmd := "df -h|awk 'NR!=1 {print $6 ,$5}'"
		combo, err := h.RunCmd(cmd)
		if err != nil {
			log.Printf("磁盘使用率执行出错:%s\n", err)
		}
		hm.DiskUsage = funcs.DiskFormat(string(combo))
		log.Printf("磁盘使用率命令输出:\n%s\n", hm.DiskUsage)

		h.CloseSession()

		//xxx 内存使用率
		if err := h.GetSession(); err != nil {
			log.Printf("Error:%v,host:%s\n", err, host)
		}

		cmd = "free -m|sed -n 2p|awk '{print $2,$3}'"
		combo, err = h.RunCmd(cmd)
		if err != nil {
			log.Printf("内存获取执行出错:%s\n", err)
		}
		memlist := bytes.Split(combo, []byte(" "))
		m := new(modules.MemStatus)
		if len(memlist) == 2 {
			m.Total, _ = strconv.Atoi(strings.TrimSpace(string(memlist[0])))
			m.Usage, _ = strconv.Atoi(strings.TrimSpace(string(memlist[1])))
		}
		m.UsageRate = fmt.Sprintf("%.2f", float64(m.Usage)/float64(m.Total)*100)
		//fmt.Println(m.UsageRate)
		log.Printf("内存获取命令输出:%+v\n", m)
		hm.Mem = m
		h.CloseSession()

		//xxx date
		if err := h.GetSession(); err != nil {
			log.Printf("Error:%v,host:%s\n", err, host)
		}
		cmd = "date '+%Y-%m-%d %H:%M:%S'"
		combo, err = h.RunCmd(cmd)
		if err != nil {
			log.Printf("date执行出错:%s\n", err)
		}

		log.Printf("date命令输出:%+v\n", string(combo))
		hm.CurrTime = string(combo)
		h.CloseSession()
		ResultHosts = append(ResultHosts, hm)
	}

	// xxx mail tpl
	funcs.MailHtml(ResultHosts)

	// xxx 发送邮件忽略

}
