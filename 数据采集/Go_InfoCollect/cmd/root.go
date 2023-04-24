package cmd

/**
 * @Author: Tao Jun
 * @Since: 2023/4/14
 * @Desc: root.go
**/

import (
	"collect_web/service"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	//	Use:   "git",
	//	Short: "Git is a distributed version control system.",
	//	Long: `Git is a free and open source distributed version control system
	//designed to handle everything from small to very large projects
	//with speed and efficiency.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Parent() == nil {
			//_ = cmd.Help() // 显示帮助信息
			//os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) { // xxx 不加参数 往下执行
		//Error(cmd, args, errors.New("unrecognized command"))
		service.CollectFlag = true
		service.CollectSrv()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func ErrorPrint(str interface{}) {
	fmt.Printf("[error] %v\n", str)
}

func ErrorCmdPrint(str interface{}, recode int) {
	if recode > 0 {
		fmt.Printf("[error] %v\n", str)
		os.Exit(recode)
	} else {
		fmt.Printf("[success] %v\n", str)
		os.Exit(0)
	}
}
