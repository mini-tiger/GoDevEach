package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"taji/g"
	"taji/modules"
	"taji/redisconn"
)

/**
 * @Author: Tao Jun
 * @Description: taji
 * @File:  syncRealData
 * @Version: 1.0.0
 * @Date: 2021/6/29 下午4:27
 */

var sigs chan os.Signal = make(chan os.Signal, 1)

const ConfigJson = "cfg.json"

func main() {

	g.LoadConfig(filepath.Join(g.Basedir, ConfigJson))
	_ = os.Chdir(g.Basedir)
	// 初始化 日志
	g.InitLog()

	err := redisconn.Conn()
	if err != nil {
		g.GetLog().Fatalf("redis:%s FAIL exit\n", g.GetConfig().RedisAddr)
	}

	err = modules.NewMongoConn()
	if err != nil {
		g.GetLog().Fatalf("mongo:%s FAIL exit\n", g.GetConfig().MongoUri)
	}

	redisconn.LoopRecvPubSub()

	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGTERM) // 第二个参数  接收的信号类型， 第三个参数  信号的动作

	go func() {
		for {
			select {
			case sig := <-sigs:
				_ = g.GetLog().Warn("接收到信号,关闭程序:%s\n", sig)
				redisconn.RDB.Close()
				modules.Mgo.DisableConn()
			}
		}

	}()

	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()
	select {}
}
