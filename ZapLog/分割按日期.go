package main

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func main() {
	path := "/home/go/GoDevEach/ZapLog/go.log"

	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	WithMaxAge 和 WithRotationCount二者只能设置一个
	 `WithMaxAge` 设置文件清理前的最长保存时间
	 `WithRotationCount` 设置文件清理前最多保存的个数
	*/
	// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。

	hook, _ := rotatelogs.New(
		path+".%Y-%m-%d%H%M",
		rotatelogs.WithLinkName(path),
		//rotatelogs.WithMaxAge(time.Duration(180)*time.Second), //xxx 实际会保留4个
		rotatelogs.WithRotationCount(3),
		//rotatelogs.WithRotationSize(10),// 字节
		rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
	)

	encoderConfig := zapcore.EncoderConfig{

		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		//EncodeLevel:    zapcore.LowercaseLevelEncoder,  // level 字体 小写编码器
		EncodeLevel:zapcore.CapitalColorLevelEncoder,
		//EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeTime:zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder, //

		//EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeName:   zapcore.FullNameEncoder,
		//ConsoleSeparator: " ", //只对 console 方式打印有用
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)



	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), //xxx 文本输出格式
		//zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook)), // 打印到控制台和文件
		atomicLevel,                                                                     // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.NewDevelopmentConfig()
	d:=zap.Development()
	development.EncoderConfig.EncodeLevel=zapcore.CapitalColorLevelEncoder



	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", "serviceName"))
	// 构造日志
	logger := zap.New(core, caller,d, filed)
	logger.WithOptions()

	logger.Info("log 初始化成功")
	logger.Info("无法获取网址",
		zap.String("url", "http://www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second))

	sugarLogger := logger.Sugar()
	sugarLogger.Debugf("Trying to hit GET request for %s", "www.baidu.com")
}
