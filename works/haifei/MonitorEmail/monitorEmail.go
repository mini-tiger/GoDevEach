package main

import (
	"haifei/MonitorEmail/funcs"
	"haifei/MonitorEmail/g"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

//https://github.com/360EntSecGroup-Skylar/excelize
// 中文 https://xuri.me/excelize/zh-hans/
// https://godoc.org/github.com/360EntSecGroup-Skylar/excelize#File.GetSheetName

//const ConfigJson = "/home/go/GoDevEach/works/haifei/syncHtmlYWReport/monitorEmail.json"
const ConfigJson = "monitorEmail.json"

var sigs chan os.Signal = make(chan os.Signal, 1)

//var reload chan bool = make(chan bool, 1)

func main() {
	//_,file,_,_:=runtime.Caller(0)
	g.LoadConfig(filepath.Join(g.Basedir, ConfigJson))
	_ = os.Chdir(g.Basedir)
	// 初始化 日志
	g.InitLog()

	//xxx 由于导入包的时候，funcs,modeuls包的变量已经加载，需要日志与配置文件初始化后 赋值日志与配置文件对象给变量

	funcs.LoadLogAndCfg()

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:8888", nil))
	}()

	go funcs.RunStatus() // xxx 定时打印运行状态,GC
	go funcs.WorkEmail() // 消费者
	funcs.CronWork()     // 生产者

	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGTERM) // 第二个参数  接收的信号类型， 第三个参数  信号的动作

	go func() {
		for {
			select {
			case sig := <-sigs:
				_ = g.GetLog().Warn("接收到信号,重载配置文件:%s\n", sig)
				g.LoadConfig(filepath.Join(g.Basedir, ConfigJson))
				g.C.Stop() //xxx 重启 任务计划
				funcs.CronWork()
			}
		}

	}()

	select {}
}
