package funcs

import (
	"cejing/TaskMulti/g"
	"cejing/TaskMulti/modules"
	"fmt"
	nsema "gitee.com/taojun319/tjtools/control"
	"github.com/wxnacy/wgo/arrays"
	"net/http"
	"strings"
	"time"
)

//var ParsingSync *sync.WaitGroup = new(sync.WaitGroup)

//var MoveBakFileChan chan string = make(chan string, 0)
//var MoveFailFileChan chan string = make(chan string, 0)
var RunChan chan string = make(chan string, 0)

//var FusionFileLock = make(chan struct{}, 0)
//var Log *nxlog4go.Logger

var ModelSema *nsema.Semaphore = nsema.NewSemaphore(1)
var TaskSema *nsema.Semaphore = nsema.NewSemaphore(1)
var ModelChan chan map[string]interface{} = make(chan map[string]interface{}, 0)
var mathQxwjm = make([]string, 0)
var NomathQxwjm = make([]string, 0)

func CronWork() {
	var ReqMod = &modules.ReqModelArgs{
		BaseReqArgs: &modules.BaseReqArgs{Page: 1, PageSize: 6000},
	}

	g.HttpClient = &http.Client{Timeout: 60 * time.Second}
	g.ReqModelUrl = fmt.Sprintf("http://%s/CJQX/crud/model_config/read", g.GetConfig().UrlIP)
	// xxx model
	modResult := modules.GetHtmlDataHandler(ReqMod, "crud").(*modules.PostResult)
	if modResult.Status == 200 && modResult.Data.Total > 0 {
		for index, value := range modResult.Data.Datalist {
			ModelSema.Acquire()
			fmt.Printf("第%d个模型:%+v\n", index+1, value["model_name"])
			ModelChan <- value
		}
	}
	fmt.Println(len(mathQxwjm))
	fmt.Println(len(NomathQxwjm))
	fmt.Println(mathQxwjm)

	fmt.Println("=============")
	fmt.Println(NomathQxwjm)
}
func WorkGroup() {
	for {
		select {
		case model := <-ModelChan:
			_, ok1 := model["model_name"]
			_, ok2 := model["id"]
			_, ok3 := model["inparamslist"]
			if ok1 && ok2 && ok3 {
				go rundo(model)
			}else{
				ModelSema.Release()
			}

		}
	}
}

func CheckQxlist(model map[string]interface{}, qxlist []string) bool {
	modelin := model["inparamslist"].(string)
	modellist := strings.Split(modelin, ",")
	//fmt.Println(11,modellist)
	//fmt.Println(22,qxlist)
	for _, mv := range modellist {
		if arrays.ContainsString(qxlist, mv) > 0 {
			continue
		} else {
			return false
		}
	}
	return true

}

func rundo(model map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			g.Logge.Errorf("任务失败,原因: %s", err)
			g.Logge.Warnf("任务结束")
		}
	}()

	defer func() {
		ModelSema.Release()
	}()

	g.ReqWellUrl = fmt.Sprintf("http://%s/CJQX/api/getJH", g.GetConfig().UrlIP)
	g.ReqWellQXWJMUrl = fmt.Sprintf("http://%s/CJQX/api/getwellflist", g.GetConfig().UrlIP)
	g.ReqQXWJMListUrl = fmt.Sprintf("http://%s/CJQX/api/getidxlist", g.GetConfig().UrlIP)
	g.ReqSaveTaskUrl = fmt.Sprintf("http://%s/CJQX/api/saveTask", g.GetConfig().UrlIP)

	// xxx well
	wellResult := modules.GetHtmlDataHandler(new(modules.ReqWellArgs), "api").(*modules.PostApiResult)
	//fmt.Printf("%+v\n",wellResult)
	if len(wellResult.Data) > 0 && wellResult.Status == 200 {
		for _, wv := range wellResult.Data {
			//fmt.Printf("wellname:%s\n",v)
			// xxx 井下曲线文件名
			wellQXWJM := modules.ReqWellQXWJMArgsFree.Get().(*modules.ReqWellQXWJMArgs)
			wellQXWJM.JH = wv

			qx := modules.GetHtmlDataHandler(wellQXWJM, "api").(*modules.PostApiResult)
			fmt.Printf("wellname:%s,qxwjm:%v\n", wv, qx.Data)

			// xxx 曲线文件名下特征名
			if qx.Status == 200 && len(qx.Data) > 0 {
				for _, qv := range qx.Data {
					TaskSema.Acquire()
					idx := modules.ReqQXWJMListArgsFree.Get().(*modules.ReqQXWJMListArgs)
					idx.QXWJM = qv
					qxlist := modules.GetHtmlDataHandler(idx, "api").(*modules.PostApiResult)
					fmt.Printf("qxwjm:%s,list:%v\n", qv, qxlist.Data)

					// xxx 判断checkedlist
					if CheckQxlist(model, qxlist.Data)  {    // && qv=="双庙1B组合一次中完.txt"
						g.Logge.Warnf("model:%s,well:%s,qxwjm:%s,math:%v", model["model_name"], wv, qv, "匹配")
						mathQxwjm = append(mathQxwjm, fmt.Sprintf("model:%s,well:%s,qxwjm:%s", model["model_name"], wv, qv))


						// xxx savetask,第一次拿 _id
						st:=modules.ReqReqSaveTask1ArgsFree.Get().(*modules.ReqSaveTask1Args)

						st.Tasktitle=fmt.Sprintf("%s_%s_%s",model["model_name"],wv,qv)
						st.Status="0"

						stResult := modules.GetHtmlDataHandler(st, "savetask").(*modules.PostSaveTaskResult)
						//g.Logge.Warnf("请求savetask:%+v", st)
						//g.Logge.Warnf("请求savetask,Result%+v", stResult)
						if stops,ok:=stResult.Data["ops"];ok {
							s := stops.([]interface{})[0]
							ss := s.(map[string]interface{})
							//fmt.Println(ss["_id"])

							st2:=&modules.ReqSaveTask2Args{}
							// 改变状态  启动任务
							st2.Model_config=model
							st2.Model_id=model["id"]
							st2.Wellname=wv
							st2.Qxwjm_list=qv
							st2.Qxwjmlist=[]string{qv}
							st2.Checkedlist=[]string{qv}
							st2.Yqtbm="02"
							st2.Share="0"
							st2.Taskdec=fmt.Sprintf("%s_%s_%s",model["model_name"],wv,qv)
							st2.Tasktitle=st2.Taskdec
							st2.Status="1"
							st2.Taskrole="admin"
							st2.Taskuser="admin"

							st2.Id = ss["_id"]

							stResult=modules.GetHtmlDataHandler(st2, "savetask").(*modules.PostSaveTaskResult)
							//g.Logge.Warnf("请求savetask2:%+v", st2)
							//g.Logge.Warnf("请求savetask2,Result%+v", stResult)
						}
						modules.ReqReqSaveTask1ArgsFree.Put(st)
						time.Sleep(5*time.Second)

					} else {
						g.Logge.Warnf("model:%s,well:%s,qxwjm:%s,math:%v  ", model["model_name"], wv, qv, "不匹配")
						NomathQxwjm = append(NomathQxwjm, fmt.Sprintf("model:%s,well:%s,qxwjm:%s", model["model_name"], wv, qv))
					}
					modules.ReqQXWJMListArgsFree.Put(idx)
					TaskSema.Release()
					//time.Sleep(1 * time.Second)
				}
			}

			modules.ReqWellQXWJMArgsFree.Put(wellQXWJM)
		}
	}

	//
	//	// xxx 循环请求 每个COMMCELL ReportTime时间,公共参数
	//	var rccs = modules.ReqCCArgsFree.Get().(*modules.ReqCCArgs)
	//	rccs.BaseReqArgs = &modules.BaseReqArgs{Page: 1, PageSize: 1}
	//	rccs.ReportTime = g.ReqReportTime
	//
	//	// xxx 插入历史记录,公共参数
	//	var rshcc = modules.ReqInsHisCCFree.Get().(*modules.ReqInsHisCC)
	//	rshcc.NormTime = g.ReqNormTime
	//
	//	for _, value := range ccres.Data.Datalist {
	//
	//		rccs.CommCell = (value["COMMCELL"]).(string)
	//
	//		//rccs.SaveStartTime() //xxx 生成时间范围
	//
	//		g.Logge.Debugf("请求 singe Commcell   : %+v", rccs.CommCell)
	//		g.Logge.Debugf("请求 singe Commcell 参数: %+v", string(rccs.JsonConvert()))
	//		singlecc := modules.GetHtmlDataHandler(rccs)
	//		g.Logge.Debugf("请求 singe Commcell 结果: %+v", singlecc)
	//
	//		//// xxx 插入历史记录
	//		//var rshcc = modules.RequestInsertHisCommcellFree.Get().(*modules.RequestInsertHisCommcell)
	//		//rshcc.NormTime = g.ReqNormTime
	//		rshcc.CommCell = rccs.CommCell
	//
	//		if singlecc.Data.Total > 0 && singlecc.Status == 200 {
	//			rshcc.ReportTime = singlecc.Data.Datalist[0]["REPORTTIME"]
	//			rshcc.EXIST = "已生成"
	//		} else {
	//			rshcc.EXIST = "未生成"
	//		}
	//
	//		singleHiscc := modules.GetHtmlDataHandler(rshcc)
	//
	//		g.Logge.Infof("插入历史CommCell 结果:%+v", *singleHiscc)
	//
	//		modules.PostResFree.Put(singlecc)
	//		modules.PostResFree.Put(singleHiscc)
	//
	//	}
	//	modules.ReqCCArgsFree.Put(rccs)
	//	modules.ReqInsHisCCFree.Put(rshcc)
	//}

	g.Logge.Warnf("任务结束")
}
