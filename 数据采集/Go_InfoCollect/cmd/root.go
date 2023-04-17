package cmd

/**
 * @Author: Tao Jun
 * @Since: 2023/4/14
 * @Desc: root.go
**/

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	//	Use:   "git",
	//	Short: "Git is a distributed version control system.",
	//	Long: `Git is a free and open source distributed version control system
	//designed to handle everything from small to very large projects
	//with speed and efficiency.`,
	Run: func(cmd *cobra.Command, args []string) { // xxx 不加参数 往下执行
		//Error(cmd, args, errors.New("unrecognized command"))
		return
	},
}

func Execute() {
	rootCmd.Execute()
}
