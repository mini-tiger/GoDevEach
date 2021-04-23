package g

import (
	"fmt"
	nxlog "github.com/ccpaging/nxlog4go"
	"io"
	"os"
	"strings"
	"time"
)

//func InitLog(level string) (err error) {
//	switch level {
//	case "info":
//		log.SetLevel(log.InfoLevel)
//	case "debug":
//		log.SetLevel(log.DebugLevel)
//	case "warn":
//		log.SetLevel(log.WarnLevel)
//	default:
//		log.Fatal("log conf only allow [info, debug, warn], please check your confguire")
//	}
//	InitLog1()
//	return
//}

var (
	logge *nxlog.Logger
)
var logger *Log1

type Log1 struct{}

func (l *Log1) Println(m ...interface{}) {

	logge.Info(m)
}
func (l *Log1) Printf(arg0 interface{}, args ...interface{}) {
	//fmt.Printf("%T,%v\n", m, m)
	//arg0 = arg0.(string)

	//fmt.Printf("%T,%v\n",arg0,arg0)
	arg0 = strings.Trim(arg0.(string), "\n") // todo 去掉换行
	logge.Info(arg0, args...)
}
func (l *Log1) Error(arg0 interface{}, args ...interface{}) {
	logge.Error(arg0, args...)
}

func (l *Log1) Debug(arg0 interface{}, args ...interface{}) {
	logge.Debug(arg0, args...)
}

func (l *Log1) Fatalln(m ...interface{}) {
	logge.Error(m)
	os.Exit(1)
}
func (l *Log1) Fatalf(arg0 interface{}, args ...interface{}) {
	logge.Error(arg0, args)
	os.Exit(1)
}

func WLog(str string) { // 在配置文件没有加载，日志方法没有生效前，写入日志
	f, err1 := os.OpenFile("run.log", os.O_CREATE|os.O_SYNC|os.O_WRONLY|os.O_APPEND, 0666)
	defer f.Close()
	if err1 != nil {
		os.Stdout.Write([]byte(fmt.Sprintf("wLog file err:%s\n", err1)))
	}
	t1 := time.Now()
	//fmt.Println(t1.Format("2006 01-02 15:04:05"))

	str = fmt.Sprintf(" [%s] [%s] (%s) \n", t1.Format("2006 01-02 15:04:05"), "ERROR", str)
	f.Write([]byte(str))

}

func InitLog1() *nxlog.Logger {
	fileName := Config().Logfile

	//logFile, err := os.Create(fileName)
	//if err != nil {
	//	log.Fatalln("open file error !")
	//}

	nxlog.FileFlushDefault = 5                                                                  // 修改默认写入硬盘时间
	nxlog.LogCallerDepth = 3                                                                    //runtime.caller(3)  日志触发上报的层级
	rfw := nxlog.NewRotateFileWriter(fileName).SetDaily(true).SetMaxBackup(Config().LogMaxDays) //log保存最大天数

	var ww io.Writer
	if Config().Daemon {
		ww = io.MultiWriter(rfw) //todo 输出到rfw定义
	} else {
		ww = io.MultiWriter(os.Stdout, rfw) //todo 同时输出到rfw 与 系统输出
	}

	// Get a new logger instance
	// todo FINEST 级别最低
	// todo %p prefix, %N 行号
	logge = nxlog.New(nxlog.FINEST).SetOutput(ww).SetPattern("%P [%Y %T] [%L] (%S LineNo:%N) %M\n")
	//Log.SetPrefix("11111")
	logge.SetLevel(1)

	logge.Info("read config file ,successfully") // 走到这里代表配置文件已经读取成功
	logge.Info("日志文件最多保存%d天", Config().LogMaxDays)
	logge.Info("logging on %s", fileName)
	logge.Info("进程已启动, 当前进程PID:%d", os.Getpid())
	return logge
	// Log some experimental messages
	//for j := 0; j < 15; j++ {
	//	for i := 0; i < 400 / (j+1); i++ {
	//		Log.Finest("Everything is created now (notice that I will not be printing to the file)")
	//		Log.Info("%d. The time is now: %s", j, time.Now().Format("15:04:05 MST 2006/01/02"))
	//		Log.Critical("Time to close out!")
	//
	//		time.Sleep(1*time.Second)
	//	}
	//}
	//rfw.Close()

	//logger = log.New(logFile, "[Debug]", log.LstdFlags)

}

func Logger() *Log1 {
	lock.RLock()
	defer lock.RUnlock()
	return logger
}
