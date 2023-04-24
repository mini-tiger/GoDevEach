package cmd

import (
	"collect_web/collect"
	"collect_web/conf"
	"collect_web/log"
	"collect_web/service"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var cmdCollect = &cobra.Command{
	Use:   "collect",
	Short: "Collect LocalHost Info [debug]",
	Long:  "Collect LocalHost Info [debug]",
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		lx := new(service.LinuxMetric)
		lx.RegMetrics()
		lx.Wlog = log.Tloginst

		var lxface collect.GetMetricInter = lx
		lxface.GetMetrics(conf.GlobalCtx)
		data := lxface.FormatData()
		str, err := json.Marshal(data)
		if err != nil {
			//fmt.Printf("json marshal err:%v\n", err)
			ErrorCmdPrint(fmt.Sprintf("json marshal err:%v\n", err), 1)
		} else {
			fmt.Println(string(str))
			//ErrorCmdPrint()
		}
		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(cmdCollect)
}
