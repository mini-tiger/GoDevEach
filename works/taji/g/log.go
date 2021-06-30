package g

import (
	nxlog "github.com/ccpaging/nxlog4go"
	logDiy "github.com/mini-tiger/tjtools/logDiyNew"
	"sync"
)

var (
	lock  = new(sync.RWMutex)
	logge *nxlog.Logger
)

func InitLog() *nxlog.Logger {
	// 初始化 日志

	logge = logDiy.InitLog1(cfg.Logfile, cfg.LogMaxDays, true, cfg.LogLevel, cfg.Stdout)
	return logge

}

func GetLog() *nxlog.Logger {
	lock.RLock()
	defer lock.RUnlock()
	return logge
}
