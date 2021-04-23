package main

import (
	"gitee.com/taojun319/tjtools/file"
	_ "github.com/mattn/go-oci8"
	"haifei/syncHtml/funcs"
	"haifei/syncHtml/g"
	"haifei/syncHtml/modules"
	"log"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
)

//https://github.com/360EntSecGroup-Skylar/excelize
// 中文 https://xuri.me/excelize/zh-hans/
// https://godoc.org/github.com/360EntSecGroup-Skylar/excelize#File.GetSheetName

//var c Config

//var OnceHtmlFusion = sync.Pool{ // xxx 只能取出一次，没有数据则使用new 方法
//	New: func() interface{} {
//		return &HtmlFusion{}
//	},
//}

const ConfigJson = "/home/go/GoDevEach/works/haifei/syncHtml/synchtml.json"

func main() {

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
		log.Println(http.ListenAndServe("0.0.0.0:7777", nil))
	}()


	dir, err := filepath.Abs(filepath.Dir(g.GetConfig().HtmlfileReg))
	if err != nil || !file.IsExist(dir) {
		log.Fatalln("HTMLFile Dir 不存在")
	}

	go funcs.RunStatus()    // xxx 定时打印运行状态,GC
	go funcs.WaitHtmlFile() // 消费者
	go funcs.FindHtml()     // 生产者

	select {}
}
