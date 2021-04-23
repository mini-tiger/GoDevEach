package main

import (
	"cejing/TaskMulti/funcs"
	"cejing/TaskMulti/g"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
)

//https://github.com/360EntSecGroup-Skylar/excelize
// 中文 https://xuri.me/excelize/zh-hans/
// https://godoc.org/github.com/360EntSecGroup-Skylar/excelize#File.GetSheetName

//const ConfigJson = "/home/go/GoDevEach/works/haifei/syncHtmlYWReport/monitorEmail.json"
const ConfigJson = "monitorCommcell.json"

var sigs chan os.Signal = make(chan os.Signal, 1)


func main() {
	//f,_:=os.Create("trace.out")
	//defer f.Close()
	//trace.Start(f)
	//defer trace.Stop()
	//_,file,_,_:=runtime.Caller(0)
	g.LoadConfig(filepath.Join(g.Basedir, ConfigJson))
	_ = os.Chdir(g.Basedir)
	// 初始化 日志
	g.InitLog()

	//xxx 由于导入包的时候，funcs,modeuls包的变量已经加载，需要日志与配置文件初始化后 赋值日志与配置文件对象给变量

	//funcs.LoadLogAndCfg()

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:8882", nil))
	}()

	//go funcs.RunStatus() // xxx 定时打印运行状态,GC
	go funcs.WorkGroup() // 消费者
	funcs.CronWork()     // 生产者


	//select {}
}
