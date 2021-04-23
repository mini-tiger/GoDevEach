package main

import (
	"gitee.com/taojun319/tjtools/file"
	_ "github.com/mattn/go-oci8"
	"haifei/syncHtmlYWReport/funcs"
	"haifei/syncHtmlYWReport/g"
	"haifei/syncHtmlYWReport/modules"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

//https://github.com/360EntSecGroup-Skylar/excelize
// 中文 https://xuri.me/excelize/zh-hans/
// https://godoc.org/github.com/360EntSecGroup-Skylar/excelize#File.GetSheetName

//const ConfigJson = "/home/go/GoDevEach/works/haifei/syncHtmlYWReport/syncYWReport.json"
const ConfigJson = "syncYWReport.json"

func main() {

	_, basedir, _, _ := runtime.Caller(0)

	os.Chdir(path.Dir(basedir))
	//fmt.Println(path.Dir(basedir))
	g.LoadConfig(ConfigJson)
	// 初始化 日志
	g.InitLog()

	//xxx 由于导入包的时候，funcs,modeuls包的变量已经加载，需要日志与配置文件初始化后 赋值日志与配置文件对象给变量
	funcs.LoadLogAndCfg()
	modules.LoadLogAndCfg()
	//Log = logDiy.InitLog1(c.Logfile, c.LogMaxDays, true, "INFO")
	//Log = g.InitLog(c.Logfile,c.LogMaxDays)
	//Log = logDiy.Logger()

	// 打印配置
	//log.Println("读取配置: %+v\n", g.GetConfig())

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:8777", nil))
	}()

	//go funcs.WaitMoveFile()

	dir, err := filepath.Abs(filepath.Dir(g.GetConfig().HtmlFileReg))
	if err != nil || !file.IsExist(dir) {
		log.Fatalln("HTMLFile Dir 不存在")
	}

	go funcs.RunStatus()    // xxx 定时打印运行状态,GC
	go funcs.WaitHtmlFile() // 消费者
	go funcs.FindHtml()     // 生产者

	select {}
}
