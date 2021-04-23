package funcs

import (
	"fmt"
	"github.com/robfig/cron"
	"haifei/MonitorCommcell/g"
	"haifei/MonitorCommcell/modules"
	"net/http"
	"time"
)

//var ParsingSync *sync.WaitGroup = new(sync.WaitGroup)

//var MoveBakFileChan chan string = make(chan string, 0)
//var MoveFailFileChan chan string = make(chan string, 0)
var RunChan chan string = make(chan string, 0)

//var FusionFileLock = make(chan struct{}, 0)
//var Log *nxlog4go.Logger

// xxx inithtml 临时变量
//var dom *goquery.Document
//var tmpSele1 *goquery.Selection
//var tmpSele2 *goquery.Selection
//var DataMap = make(map[string]interface{},25)

//var c *g.Config

//var DataMap = make(map[string]interface{}, 25)

//func LoadLogAndCfg() {
//	Log = g.Logge
//	//c = g.GetConfig()
//}

//type HtmlFusionFiles struct {
//	hs          []string
//	ParsingSync *sync.WaitGroup
//}

func CronWork() {
	cronFunc := func() {

		// xxx 去掉 秒
		RunChan <- time.Now().Format("2006-01-02 15:04:00")
	}
	g.C = cron.New()

	//AddFunc
	//spec := "0 */1 * * * " //每小时 只在 整分钟0秒执行
	_ = g.C.AddFunc(g.GetConfig().Schedule, cronFunc)

	g.C.Start()
}

func WorkGroup() {
	//reqArgs := &modules.RequestArgs{
	//	BaseReqArgs: &modules.BaseReqArgs{Page: 1, PageSize: 6000},
	//	RunType:     g.Runtype,
	//	SolveType:   g.Solvetype,
	//}

	//sj := modules.HtmlReqData{
	var RequestCC = &modules.ReqCCArrArgs{
		BaseReqArgs: &modules.BaseReqArgs{Page: 1, PageSize: 6000},
	}

	g.ReqCCUrl = fmt.Sprintf("http://%s/crud/HF_UNIQXCOMMCELL/read", g.GetConfig().UrlIP)
	g.ReqCCArrUrl = fmt.Sprintf("http://%s/crud/HF_YWREPORTDETAIL/read", g.GetConfig().UrlIP)
	g.ReqInsHisUrl = fmt.Sprintf("http://%s/crud/HF_HISCOMMCELL/create", g.GetConfig().UrlIP)

	g.HttpClient = &http.Client{Timeout: 60 * time.Second}

	runDo := func() {
		defer func() {
			if err := recover(); err != nil {
				g.Logge.Errorf("任务失败,原因: %s", err)
				g.Logge.Warnf("任务结束")
			}
		}()

		g.Logge.Info("开始请求COMMCELL")

		//xxx 请求总COMMCELLS
		ccres := modules.GetHtmlDataHandler(RequestCC)
		defer modules.PostResFree.Put(ccres)

		g.Logge.Debug("COMMCELL Result:%v", ccres)

		if ccres.Data.Total > 0 && ccres.Status == 200 {

			// xxx 循环请求 每个COMMCELL ReportTime时间,公共参数
			var rccs = modules.ReqCCArgsFree.Get().(*modules.ReqCCArgs)
			rccs.BaseReqArgs = &modules.BaseReqArgs{Page: 1, PageSize: 1}
			rccs.ReportTime = g.ReqReportTime

			// xxx 插入历史记录,公共参数
			var rshcc = modules.ReqInsHisCCFree.Get().(*modules.ReqInsHisCC)
			rshcc.NormTime = g.ReqNormTime

			for _, value := range ccres.Data.Datalist {

				rccs.CommCell = (value["COMMCELL"]).(string)

				//rccs.SaveStartTime() //xxx 生成时间范围

				g.Logge.Debugf("请求 singe Commcell   : %+v", rccs.CommCell)
				g.Logge.Debugf("请求 singe Commcell 参数: %+v", string(rccs.JsonConvert()))
				singlecc := modules.GetHtmlDataHandler(rccs)
				g.Logge.Debugf("请求 singe Commcell 结果: %+v", singlecc)

				//// xxx 插入历史记录
				//var rshcc = modules.RequestInsertHisCommcellFree.Get().(*modules.RequestInsertHisCommcell)
				//rshcc.NormTime = g.ReqNormTime
				rshcc.CommCell = rccs.CommCell

				if singlecc.Data.Total > 0 && singlecc.Status == 200 {
					rshcc.ReportTime = singlecc.Data.Datalist[0]["REPORTTIME"]
					rshcc.EXIST = "已生成"
				} else {
					rshcc.EXIST = "未生成"
				}

				singleHiscc := modules.GetHtmlDataHandler(rshcc)

				g.Logge.Infof("插入历史CommCell 结果:%+v", *singleHiscc)

				modules.PostResFree.Put(singlecc)
				modules.PostResFree.Put(singleHiscc)

			}
			modules.ReqCCArgsFree.Put(rccs)
			modules.ReqInsHisCCFree.Put(rshcc)
		}

		g.Logge.Warnf("任务结束")
	}

	for {
		select {
		case currentTimeStr := <-RunChan:
			g.Logge.Warnf("任务开始 当前时间:%s", time.Now().Format(g.TimeLayout))

			Tm, _ := time.ParseInLocation(g.TimeLayout, currentTimeStr, time.Local)
			g.CurrWorkTime = Tm.Unix() * 1000

			modules.GenReqTime()
			go runDo()

		}
	}
}
