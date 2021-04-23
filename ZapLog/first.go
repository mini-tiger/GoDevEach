package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {
	// zap.NewDevelopment 格式化输出
	logger, _ := zap.NewDevelopment()
	//defer logger.Sync()
	logger.Info("无法获取网址",
		zap.String("url", "http://www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
	// zap.NewProduction json序列化输出
	logger2, _ := zap.NewProduction()
	//defer logger.Sync()
	logger2.Info("无法获取网址",
		zap.String("url", "http://www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
//
	config:=zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
		},
		//OutputPaths:      []string{"stdout", "/home/go/GoDevEach/ZapLog/log.txt"},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	l,_:=config.Build()
	l.Info("无法获取网址",
		zap.String("url", "http://www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

		l.Info("无法获取网址",
			zap.String("url", "http://www.baidu.com"),
			zap.Int("attempt", 3),
			zap.Duration("backoff", time.Second),
		)


	//
	config1:=zap.NewDevelopmentConfig()
	config1.EncoderConfig.EncodeLevel=zapcore.CapitalColorLevelEncoder
	logger3,_:=config1.Build()
	logger3.Info("color")
	logger3.Debug("debug")

}