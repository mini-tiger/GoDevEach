package log

import (
	"go.uber.org/zap"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/4
 * @Desc: wraplog.go
**/
type Wraplog struct {
	ZapLog *zap.Logger
}

var Wlog *Wraplog = new(Wraplog)

func InitWrapLog() {
	Wlog.ZapLog = GlogSkip1
}

//cl.ZapLog = log.GlogSkip1
//

// cron
func (c *Wraplog) Info(msg string, keysAndValues ...interface{}) {
	c.ZapLog.Sugar().Info(msg, keysAndValues)
}

// cron
func (c *Wraplog) Error(err error, msg string, keysAndValues ...interface{}) {
	c.ZapLog.Sugar().Error(err, msg, keysAndValues)
}

// ghw
func (c *Wraplog) Printf(format string, v ...interface{}) {
	//str := fmt.Sprintf(format, v...)
	format = strings.Trim(format, "\n")
	c.ZapLog.Sugar().Errorf(format, v...)
}
