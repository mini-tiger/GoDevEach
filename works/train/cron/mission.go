package cron

import (
	"time"
	"godev/works/train/g"
	"strconv"
	"fmt"
)

const Lasttime = "lasttime"

var lastTime *int = new(int) // todo 最近火车进入时间戳,可以从redis读取后 写入
//var tt *time.Ticker
//var tr *time.Timer // 一次性
var cTrain chan map[string]string = make(chan map[string]string, 5)

func TrainCrond() {
	//for{
	go GetTrainPic()
	LoadLasttime()
	CreateNextTrainMission(0)
	//if *lastTime != 0 {
	//	if b, resultTrain := diffTrainTime(); b {
	//		g.Logger().Debug("第一次启动，lastime < ti 创建图片任务")
	//		cTrain <- resultTrain
	//		nt := *lastTime + g.Config().Train.TrainInterval - g.GetNow()
	//		g.Logger().Printf("下次train任务开始时间%s\n", g.GetDateStr(nt+g.GetNow()))
	//		go CreateNextTrainMission(nt) // 下次一次性
	//
	//	} else {
	//		g.Logger().Debug("第一次启动，lastime == ti 创建循环任务")
	//		go CreateLoopMission() // todo 创建循环10s 判断一次是否有新的
	//	}
	//
	//} else {
	//	for {
	//		results, err := getTrainNew()
	//		if err != nil {
	//			time.Sleep(time.Duration(5) * time.Second)
	//			continue
	//		}
	//		firstRun(results[0])
	//		break
	//	}
	//}

	//go WaitMission() // loop wait mission


}
func LoadLasttime() {
	if g.Redis.StringExists(Lasttime) {
		val, _ := g.Redis.StringGet(Lasttime)
		i, err := strconv.Atoi(val)
		if err != nil {
			i = 0
			lastTime = &i
			return
		}
		lastTime = &i
	} else {
		i := 0
		lastTime = &i
	}
	g.Logger().Debug("获取到lasttime %d,时间:%s", *lastTime, g.GetDateStr(*lastTime))
}

func GetTrainPic() {
	for {
		select {
		case result := <-cTrain:
			g.Logger().Debug("获取图片启动,参数: %+v", result)
			go GetPicMission(result)
		}
	}
}

func firstRun(result map[string]string) {
	nowUnix := int(time.Now().Unix())
	//nowUnix = 111
	nextTime := *lastTime + g.Config().Train.TrainInterval
	fmt.Println("======", nowUnix, *lastTime, nextTime)
	switch {
	case nowUnix <= nextTime: //todo 当前时间小于 库中最新时间戳 加 时间间隔，可能第一次启动在 列车进入期间
		g.Logger().Debug("第一次启动，创建图片获取任务")
		cTrain <- result                              // todo 创建协程 获取图片任务， 不影响下次 获取最新时间的任务
		go CreateNextTrainMission(nextTime - nowUnix) // 距离下次的间隔秒数
		g.Logger().Debug("NextTime %s", g.GetDateStr(nextTime))

	case nowUnix > nextTime:
		g.Logger().Debug("第一次启动，创建循环任务")
		go CreateLoopMission() // todo 创建循环10s 判断一次是否有新的
	}
}

func CreateLoopMission() { // 循环获取时间任务，由于一次性任务没有取到
	//tt = time.NewTicker(time.Duration(g.Config().LoopInterval) * time.Second)
	go LoopMission(g.Config().Train.LoopInterval)
}

func CreateNextTrainMission(interval int) { // 一次性，下次获取时间任务
	//tr = time.NewTimer(time.Duration(interval) * time.Second)
	time.Sleep(time.Duration(interval) * time.Second)
	go NextTrainMission()
}

func LoopMission(interval int) {
	for {

		if b, resultTrain := diffTrainTime(); b {
			tn := g.GetNow() // 排除 执行sql使用的时间
			g.Logger().Printf("循环结束\n")
			//fmt.Println("go 启动获取图片")
			cTrain <- resultTrain
			nt := *lastTime + g.Config().Train.TrainInterval - tn
			//nt := g.Config().Train.TrainInterval
			g.Logger().Printf("下次train任务开始时间%s\n", g.GetDateStr(nt+tn))
			go CreateNextTrainMission(nt)
			return
		} else {
			g.Logger().Printf("下次循环任务开始时间%s\n", time.Now().Add(time.Duration(interval)*time.Second))
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func NextTrainMission() {

	if b, resultTrain := diffTrainTime(); b { //一次性 获取到新的时间数据
		now := g.GetNow()
		cTrain <- resultTrain
		nt := *lastTime + g.Config().Train.TrainInterval - now
		g.Logger().Printf("下次train任务开始时间%s\n", g.GetDateStr(nt+now))
		go CreateNextTrainMission(nt) // 下次一次性
	} else { // 没有获取到新的数据，进入循环获取
		g.Logger().Printf("创建循环任务开始时间%s\n", time.Now().Add(time.Duration(g.Config().Train.LoopInterval)))
		go CreateLoopMission()
	}
}

func diffTrainTime() (b bool, resultTrain map[string]string) { // 是否应该进行获取图片，库中时间
	resultTrainSlice, err := getTrainNew()
	if err != nil {
		g.Logger().Error("getTrainNew err:%s", err)
		return
	}

	if ti, b := g.Getint(resultTrainSlice[0]["UNIXSTAMP"]); b {
		g.Logger().Debug("获取到库中最新时间%d ,%s", ti, g.GetDateStr(ti))
		switch {
		//case *lastTime == 0:
		//	lastTime = &ti
		//	g.Redis.StringSet(Lasttime, *lastTime)
		//	firstRun(resultTrainSlice[0])
		case *lastTime == ti:
			return false, resultTrain
		case *lastTime < ti:
			lastTime = &ti
			g.Redis.StringSet(Lasttime, *lastTime)
			return true, resultTrainSlice[0]
		default: // 库中数据除了第一次以外，不应该比 lasttime 小
			lastTime = &ti
			g.Redis.StringSet(Lasttime, *lastTime)
			return false, resultTrain
		}
	} else { // 获取到的时间格式如果不对
		g.Logger().Error("获取到的时间格式错误 %d 秒后再次运行", g.Config().Train.TrainInterval)
		CreateNextTrainMission(g.Config().Train.TrainInterval)
		return false, resultTrain
	}
	return
}

func getTrainNew() (results []map[string]string, err error) {

	sql := "select * from (select train_serial,train_id,station_id,pass_time,vehicle_number,index_id,to_char(pass_time,'YYYY/MM/DD HH24:mi:ss') as pass_time_char," +
		"(pass_time - TO_DATE('1970-01-01 08:00:00', 'YYYY-MM-DD HH24:mi:ss')) * 86400 as unixstamp from tf_op_train order by pass_time desc) " +
		"where rownum = 1"

	g.Logger().Debug("sql :%s", sql)

	results, err = g.Engine.QueryString(sql)

	if len(results) == 0 {
		results[0]["UNIXSTMAP"] = "0"
		return
	}
	//fmt.Println(results)
	/*
	[map[TRAIN_SERIAL:32F6ACC0B0A243FC974BC2BCECF2AA51 TRAIN_ID:74326.0 STATION_ID:V28F04F01
	PASS_TIME:2018-12-24T22:46:46+08:00 VEHICLE_NUMBER:212 INDEX_ID:1614 PASS_TIME_CHAR:2018/12/24 22:46:46 UNIXSTAMP:1545662806]]

	*/
	//
	//if err != nil {
	//	g.Logger().Error("Exec sql:%s,err:%s\n", sql, err)
	//	return
	//}
	return
	//fmt.Printf("%T,%+v\n",results,results)
}
