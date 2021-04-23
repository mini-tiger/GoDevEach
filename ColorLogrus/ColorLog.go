package main

import (
	"ColorLogrus/utils"
	logDiy "gitee.com/taojun319/tjtools/LogrusDiy"
	log "github.com/sirupsen/logrus"
	"time"
)

// todo https://github.com/sirupsen/logrus/blob/master/example_custom_caller_test.go
// todo

var Log *log.Logger

func main() {

	// caller

	//var log log.Logger
	var logconfig *logDiy.LogrusConfig = &logDiy.LogrusConfig{logDiy.SimpleFileTextFormat,
		logDiy.SimpleFileJsonFormat,
		true, log.InfoLevel, 2, "minute", "/home/go/GoDevEach/ColorLogrus/run.log"}

	Log, err := logDiy.InitLog(logconfig)
	if err != nil {
		Log.Error(err)
	}
	//log.SetOutput(os.Stdout)

	//LogDiy := utils.DiyLog(Log)

	// 记录不同文件caller的处理
	go utils.TestWriteLog(Log)

	i := 1
	//  封装创建Entry，方便输出临时变量 或者 固定变量

	for {
		Log.Debug("This is debug")
		log.Info("succeeded1")

		log.WithFields(log.Fields{"event": "A", "index": i}).Info("this is withFields1")
		Log.WithFields(log.Fields{"event": "A", "index": i}).Error("this is withFields1")
		//LogDiy.Info("this is withFields2")

		time.Sleep(time.Duration(10) * time.Second)
		i++
	}

}
