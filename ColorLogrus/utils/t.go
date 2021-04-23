package utils

import (
	"errors"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func TestWriteLog(Log *log.Logger) {
	for {
		Log.Info("succeeded")
		Log.Warn("not correct")
		Log.Error("something error")

		time.Sleep(time.Duration(10) * time.Second)
	}
}

var logger *log.Logger
var lock = new(sync.RWMutex)
//var Lc *LogrusConfig
var shortfile bool

var (
	SimpleTerminalFormat *log.TextFormatter = &log.TextFormatter{ForceColors: true, TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true,
		CallerPrettyfier:TJLogrusCaller}
	SimpleFileJsonFormat log.Formatter      = &log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05",CallerPrettyfier:TJLogrusCaller}
	SimpleFileTextFormat log.Formatter      = &log.TextFormatter{ForceColors: true, TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true,
		CallerPrettyfier:TJLogrusCaller}
)

type LogrusConfig struct {
	TerminalFormat   log.Formatter
	FileOutputFormat log.Formatter
	ShortFile        bool
	LogLevel         log.Level
	MaxRemainCnt     uint
	LoopType         string
	LogFileNameAbs   string
}



func TJLogrusCaller(r *runtime.Frame) (string, string) {
	if shortfile{
		return fmt.Sprintf("func:%s", filepath.Base(runtime.FuncForPC(r.Entry).Name())),
			fmt.Sprintf(" %s:%d", filepath.Base(r.File), +r.Line)
	}else{
		return fmt.Sprintf("func:%s", filepath.Base(runtime.FuncForPC(r.Entry).Name())),
			fmt.Sprintf(" %s:%d", r.File, +r.Line)
	}

}

//func (l *LogrusConfig)reflectFormat(f interface{}) log.Formatter{
//	fmt.Printf("%+v\n",reflect.ValueOf(f))
//	v:=reflect.ValueOf(f).Elem()
//	v.MethodByName("CallerPrettyfier").Set(l.TJLogrusCaller)
//
//	return SimpleFileJsonFormat
//
//}

func InitLog(Lc *LogrusConfig) (logNew *log.Logger, e error) {
	logNew = log.New()
	var loopTime time.Duration
	switch Lc.LoopType {
	case "daily":
		loopTime = time.Hour * 24
		break
	case "hour":
		loopTime = time.Hour
		break
	case "minute":
		loopTime = time.Minute
		break
	default:
		e = errors.New(fmt.Sprintf("Input LoopTime: daily,hour,minute"))
		return
	}

	//  xxx 日志按天 或者 按个数 保存
	writer, err := rotatelogs.New(
		Lc.LogFileNameAbs+"_"+".%Y%m%d%H%M"+".log",
		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(Lc.LogFileNameAbs),


		// WithRotationTime设置日志分割的时间,这里设置为一小时分割一次
		rotatelogs.WithRotationTime(loopTime),

		/* WithMaxAge和WithRotationCount二者只能设置一个,

		WithMaxAge设置文件清理前的最长保存时间,
		 WithRotationCount设置文件清理前最多保存的个数.
		//rotatelogs.WithMaxAge(time.Hour*24) */
		rotatelogs.WithRotationCount(Lc.MaxRemainCnt),
	)

	if err != nil {
		e = errors.New(fmt.Sprintf("config local file system for logger error: %v", err))
		return
	}

	log.SetLevel(Lc.LogLevel)

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, Lc.FileOutputFormat)

	//lfsHook.SetFormatter(&log.TextFormatter{DisableColors: true,TimestampFormat:"2006-01-02 15:04:05",FullTimestamp:true})

	//	启用caller
	logNew.SetReportCaller(true)
	logNew.AddHook(lfsHook)
	// todo logrus 自定义格式
	logNew.SetFormatter(Lc.TerminalFormat)

	//Log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	logger = logNew // 赋值全局变量
	shortfile = Lc.ShortFile //  赋值 shourtfile
	return logNew, nil
}

func Logger() *log.Logger {
	lock.RLock()
	defer lock.RUnlock()
	return logger
}

func DiyLog(Log *log.Logger) *log.Entry {
	pc, fi, ln, _ := runtime.Caller(1)
	diylog1 := Log.WithFields(log.Fields{"prefix": "TestPrefix", "func1": runtime.FuncForPC(pc).Name(), "file1": fi, "LineNo1": ln})
	return diylog1
}
