package log

import (
	"log"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/23
 * @Desc: tempLog.go
**/

type Tlog struct {
}

var Tloginst *Tlog = new(Tlog)

//cl.ZapLog = log.GlogSkip1
//

// cron
func (c *Tlog) Info(msg string, keysAndValues ...interface{}) {
	log.Printf(msg, keysAndValues...)
}

// cron
func (c *Tlog) Error(err error, msg string, keysAndValues ...interface{}) {
	log.Println(err)
	log.Printf(msg, keysAndValues...)
}

// ghw
func (c *Tlog) Printf(format string, v ...interface{}) {
	//str := fmt.Sprintf(format, v...)
	format = strings.Trim(format, "\n")
	log.Printf(format, v...)
}
