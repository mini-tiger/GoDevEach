package cmd

import (
	"collect_web/collect"
	"collect_web/conf"
	"collect_web/service"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

type tlog struct {
}

var tLog *tlog = new(tlog)

//cl.ZapLog = log.GlogSkip1
//

// cron
func (c *tlog) Info(msg string, keysAndValues ...interface{}) {
	log.Printf(msg, keysAndValues...)
}

// cron
func (c *tlog) Error(err error, msg string, keysAndValues ...interface{}) {
	log.Println(err)
	log.Printf(msg, keysAndValues...)
}

// ghw
func (c *tlog) Printf(format string, v ...interface{}) {
	//str := fmt.Sprintf(format, v...)
	format = strings.Trim(format, "\n")
	log.Printf(format, v...)
}

var cmdCollect = &cobra.Command{
	Use:   "collect",
	Short: "Collect once",
	Long:  `Collect once`,
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		lx := new(service.LinuxMetric)
		lx.RegMetrics()
		lx.Wlog = tLog

		var lxface collect.GetMetricInter = lx
		lxface.GetMetrics(conf.GlobalCtx)
		data := lxface.FormatData()
		str, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("json marshal err:%v\n", err)
		} else {
			fmt.Println(string(str))
		}
		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(cmdCollect)
}
