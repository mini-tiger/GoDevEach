package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

var logfile="/home/go/GoDevEach/run.log"
var logDir=path.Dir(logfile)
var logBaseFile=path.Base(logfile)
var NextSplit int64
var s  *sync.Mutex=new(sync.Mutex)
var F *os.File

//xxx 同时输出的标准输出和文件中

func main() {
	NextSplitlogTime() //计算下次分割的时间

	F, _ = os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755) //读写追加模式
	for {
		//第10秒写一次
		Wlog(time.Now().Format("2006-01-02 15:04:05")+"\n")
		time.Sleep(10*time.Minute)
	}
}

func Wlog(a string) {
	s.Lock()
	defer s.Unlock()
	// xxx 到了分割的时间
	if time.Now().Unix() >= NextSplit{
		// 分割文件  重新记录下次分割的时间
		F.Sync()
		F.Close()
		os.Chdir(logDir)
		os.Rename(logBaseFile,fmt.Sprintf("%s_%s",logBaseFile,time.Now().Format("2006-01-02")))
		F, _ = os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755) //读写模式
		wlog(a)
		NextSplitlogTime()
	}else{

		wlog(a)

	}

}

func wlog(a string)  {
	//_, _ = os.Stdout.WriteString(a)
	log.Printf(a)

	_, _ = F.Write([]byte(a))

}

func NextSplitlogTime()  {
	timeStr:=time.Now().Format("2006-01-02")
	//t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 23:59:59", time.Local)
	t2, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)

	//fmt.Println(t.Unix() + 1)
	//fmt.Println(t2.AddDate(0, 0, 1).Unix())
	NextSplit=t2.AddDate(0, 0, 1).Unix()
	//NextSplit=t2.Add(time.Duration(1*time.Minute)).Unix()

}