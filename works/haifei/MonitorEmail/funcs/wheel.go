package funcs

import (
	"bytes"
	"fmt"
	"haifei/MonitorEmail/modules"
	"net/http"
	"path/filepath"

	"syscall"

	"github.com/ccpaging/nxlog4go"
	"github.com/robfig/cron"
	"haifei/MonitorEmail/g"
	"time"
)

//var ParsingSync *sync.WaitGroup = new(sync.WaitGroup)

//var MoveBakFileChan chan string = make(chan string, 0)
//var MoveFailFileChan chan string = make(chan string, 0)
var RunChan chan int64 = make(chan int64, 0)


//var FusionFileLock = make(chan struct{}, 0)
var Log *nxlog4go.Logger

// xxx inithtml 临时变量
//var dom *goquery.Document
//var tmpSele1 *goquery.Selection
//var tmpSele2 *goquery.Selection
//var DataMap = make(map[string]interface{},25)

//var c *g.Config

//var DataMap = make(map[string]interface{}, 25)

func LoadLogAndCfg() {
	Log = g.GetLog()
	//c = g.GetConfig()
}

//type HtmlFusionFiles struct {
//	hs          []string
//	ParsingSync *sync.WaitGroup
//}

func CronWork() {
	cronFunc := func() {

		RunChan <- time.Now().Unix()
	}
	g.C = cron.New()

	//AddFunc
	//spec := "0 */1 * * * " //每小时 只在 整分钟0秒执行
	_ = g.C.AddFunc(g.GetConfig().Schedule, cronFunc)

	g.C.Start()
}

func WorkEmail() {
	//reqArgs := &modules.RequestArgs{
	//	BaseReqArgs: &modules.BaseReqArgs{Page: 1, PageSize: 6000},
	//	RunType:     g.Runtype,
	//	SolveType:   g.Solvetype,
	//}

	//sj := modules.HtmlReqData{
	//	RequestArgs: &modules.RequestArgs{
	//		BaseReqArgs: &modules.BaseReqArgs{Page: 1, PageSize: 6000},
	//		RunType:     g.Runtype,
	//		SolveType:   g.Solvetype,
	//	},
	//	// todo 以下和请求参数 没有关系
	//	DownLoadUrl:        fmt.Sprintf("http://%s/crud/HF_YWDETAIL_BASIC_VIEW_NEW/exportData", g.GetConfig().GetIP),
	//	RequestCommCellUrl: fmt.Sprintf("http://%s/crud/HF_YWGROUPCOMMCELL/read", g.GetConfig().GetIP),
	//	RequestTotalUrl:    fmt.Sprintf("http://%s/crud/HF_YWMONITORTOTAL/read", g.GetConfig().GetIP),
	//
	//
	//}
	g.Fp = filepath.Join(g.GetConfig().TmpDownloadDir, "monitor.xlsx")
	g.HtmlBuffer = new(bytes.Buffer)
	g.HttpClient = &http.Client{Timeout: 60 * time.Second}

	g.DownLoadUrl = fmt.Sprintf("http://%s/crud/HF_YWDETAIL_BASIC_VIEW_NEW/exportData", g.GetConfig().GetIP)
	g.RequestCommCellUrl = fmt.Sprintf("http://%s/crud/HF_YWGROUPCOMMCELL/read", g.GetConfig().GetIP)
	g.RequestTotalUrl = fmt.Sprintf("http://%s/crud/HF_YWMONITORTOTAL/read", g.GetConfig().GetIP)

	sj := modules.EmailGenSend{
		HtmlGetData: &modules.RequestArgs{
			BaseReqArgs: &modules.BaseReqArgs{Page: 1, PageSize: 6000},
			RunType:     g.Runtype,
			SolveType:   g.Solvetype,
		},
		SendMailOk: false,
	}

	runDo := func() {
		defer func() {
			if err := recover(); err != nil {
				_ = g.GetLog().Error("运维监控邮件发送任务失败,原因: %s\n", err)
			}
		}()

		//sj.HtmlData = HtmlDataFree.Get().(*HtmlData)

		defer func() {
			g.HtmlBuffer.Reset()
			//HtmlDataFree.Put(sj.HtmlData)
		}()

		//sji.DownFile(jsonBytes) //下载明细数据EXCEL

		if sj.GetCommCells() && sj.GetTotalData() && sj.DownFile() { //请求commcell, totalData ,excel

			modules.GenHtml()
			for _, mailto := range g.GetConfig().MailTo {
				g.SendMailSync.Add(1)
				go func(mailto string) {

					modules.SendMail(mailto)
					g.SendMailSync.Done()
				}(mailto)
			}

		}

		//if sj.DownloadSuccess && sj.RequestDataSuccess {
		//	sj.GenHtml()
		//	sj.SendMail()
		//}
		g.SendMailSync.Wait()
		_ = syscall.Unlink(g.Fp)
		_ = g.GetLog().Warn("运维监控邮件发送任务结束\n")
	}

	for {
		select {
		case currentUnix := <-RunChan:
			_ = g.GetLog().Warn("运维监控邮件发送任务开始\n")

			//sj.RequestArgs.StartTime = map[string]string{"key": "gt", "value": strconv.FormatInt((currentUnix-int64(g.GetConfig().ValidHours*3600))*1000, 10)}
			//sj.RequestArgs.StartTime = map[string]string{"key": "gt", "value": strconv.FormatInt(0, 10)}
			sj.SaveStartTime(currentUnix)

			g.JsonBytes = sj.JsonConvert()

			g.GetLog().Debug("请求参数 %s\n", string(g.JsonBytes))
			//if err != nil {
			//	_ = g.GetLog().Error("转换JSON 失败:%v\n", err)
			//} else {
			go runDo()
			//}
		}
	}
}
