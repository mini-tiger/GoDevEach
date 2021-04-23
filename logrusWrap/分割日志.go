package main

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"

	"os"

	"path"
	"runtime"
	"time"

	//rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

var log1 *log.Logger

func init() {
	pathlog := "/home/go/GoDevEach/logrusWrap/go.log"

	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	WithMaxAge 和 WithRotationCount二者只能设置一个
	  `WithMaxAge` 设置文件清理前的最长保存时间
	  `WithRotationCount` 设置文件清理前最多保存的个数
	*/
	// eee 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。

	writer, _ := rotatelogs.New(
		pathlog+".%Y-%m-%d%H%M",
		rotatelogs.WithLinkName(pathlog),
		//rotatelogs.WithMaxAge(time.Duration(180)*time.Second), //xxx 实际会保留4个
		rotatelogs.WithRotationCount(3),
		//rotatelogs.WithRotationSize(10),// 字节
		rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
	)
	tformat := &log.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableColors:    false,
		CallerPrettyfier: CallerFormat, //xxx 自定义 runtime.Caller的格式
		PadLevelText:     true,         //xxx 是否完整显示 LEVEL文本
	}
	log.SetFormatter(tformat)

	//F, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755) //读写追加模式

	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	//log.SetFormatter(&log.TextFormatter{FullTimestamp:true,DisableColors:false})
	log.SetReportCaller(true)
	//log.SetOutput(writer)
	log.Info("test info") //xxx 默认方式 输出日志

	// fixme 如果要同时输出 使用 log.New()    ， 默认log 除os.stdout外 日志不能输出颜色
	log1 = log.New()
	log1.SetFormatter(tformat)

	//log1.SetOutput(os.Stdout)
	log1.SetOutput(io.MultiWriter(os.Stdout, writer))
	log1.SetReportCaller(true)
	log1.Info("111")

	//log.SetFormatter(&log.JSONFormatter{})
}

func CallerFormat(r *runtime.Frame) (function string, file string) {
	//fmt.Printf("%+v\n",r)
	return " ", fmt.Sprintf("%s:%d", path.Base(r.File), r.Line)
}

func main() {
	for {
		log.Infof("%s\n", "hello, world!")
		log1.Infof("%s\n", "hello, world!")
		time.Sleep(time.Duration(15) * time.Second)
	}
}
