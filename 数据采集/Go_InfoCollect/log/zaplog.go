package log

import (
	"collect_web/conf"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"strings"
	"time"
)

var GlogSkip1 *zap.Logger

var Glog *zap.Logger

const (
	logTmFmtWithMS = "2006-01-02 15:04:05.000"
)

func levelCheck(level string) zapcore.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return zap.DebugLevel
	case "INFO":
		return zap.InfoLevel
	default:
		return zap.DebugLevel
	}
}

func JsonString(key string, val interface{}) zap.Field {
	var value string
	b, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		value = fmt.Sprintf("%v", val)
	} else {
		value = string(b)
	}
	//fmt.Println(value)
	return zap.Field{Key: key, Type: zapcore.StringType, String: value}
}

func JsonMarshalIndent(val interface{}) (value string) {
	//var value string
	b, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		value = fmt.Sprintf("%v", val)
	} else {
		value = string(b)
	}
	return
}
func CheckRunVersion() zapcore.Level {
	if conf.RunMode == "dev" {
		return zapcore.DebugLevel
	} else {

		return zapcore.InfoLevel
	}
}

func InitLog(flag string, filename string) {
	level := CheckRunVersion()

	//fmt.Println(CurrentDir)
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevelAt(level)
	//logfilename

	hook := lumberjack.Logger{
		Filename:   path.Join(conf.CurrentDir, filename), // 日志文件路径
		MaxSize:    128,                                  // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 0,                                    // 日志文件最多保存多少个备份
		MaxAge:     1,                                    // 文件最多保存多少天
		LocalTime:  true,
		Compress:   false, // 是否压缩
	}

	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format(logTmFmtWithMS) + "]")
	}
	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	// 自定义文件：行号输出项
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		//enc.AppendString("[" + caller.traceId + "]")
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}

	encoderConfig := zapcore.EncoderConfig{

		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "linenum",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		//EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		//EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		//EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder, //

		EncodeTime:   customTimeEncoder,   // 自定义时间格式
		EncodeLevel:  customLevelEncoder,  // 小写编码器
		EncodeCaller: customCallerEncoder, // 全路径编码

		//EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		//EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeName: zapcore.FullNameEncoder,
	}
	//fmt.Println(int8(atomicLevel.Level()))
	encoder := zapcore.NewJSONEncoder(encoderConfig) // prod
	if int8(atomicLevel.Level()) == -1 {             // debug
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	core := zapcore.NewCore(
		encoder, // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()

	// 开启文件及行号
	development := zap.Development()
	//zap.AddCallerSkip(1)
	// 设置初始化字段
	filed := zap.Fields(zap.String("flag", flag)) //区别不同的agent
	// 构造日志
	ZapLogOpt := []zap.Option{caller, development, filed}
	Glog = zap.New(core, ZapLogOpt...)
	Glog.Sugar().Infof("当前运行目录: %s, 发布级别: %s,  日志级别:%s", conf.CurrentDir, conf.RunMode, level)

	//GlogSkip0 = Glog.WithOptions(zap.AddCallerSkip(0))
	GlogSkip1 = Glog.WithOptions(zap.AddCallerSkip(1))
	//logger.Info("无法获取网址",
	//	zap.String("url", "http://www.baidu.com"),
	//	zap.Int("attempt", 3),
	//	zap.Duration("backoff", time.Second))

	//sugarLogger := logger.Sugar()
	//sugarLogger.Debugf("Trying to hit GET request for %s", "www.baidu.com")
}
