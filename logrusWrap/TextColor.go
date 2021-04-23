package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var logger1 *logrus.Entry

func init() {
	// 设置日志格式
	// https://godoc.org/github.com/sirupsen/logrus#TextFormatter

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableColors:    false,
		CallerPrettyfier: CallerFormat_text, //xxx 自定义 runtime.Caller的格式
		PadLevelText:     true,              //xxx 是否完整显示 LEVEL文本
	})

	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	//log.SetFormatter(&log.TextFormatter{FullTimestamp:true,DisableColors:false})
	log.SetReportCaller(true)
	log.Info("test info") //xxx 默认方式 输出日志

	logger1 = log.WithFields(log.Fields{"请求id": "123444", "user_ip": "127.0.0.1"})

}

func CallerFormat_text(r *runtime.Frame) (function string, file string) {
	//fmt.Printf("%+v\n",r)
	return " ", fmt.Sprintf("%s:%d", path.Base(r.File), r.Line)
}

func main() {
	logger1.Infof("hello, logrus....%s", "1")
	logger1.Info("hello, logrus1....")
	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Debug("A group of walrus emerges from the ocean")
	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 100,
	}).Fatal("The ice breaks!")
}
