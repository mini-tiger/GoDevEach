package cmd

import (
	"collect_web/conf"
	"collect_web/service"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var serverFlag string

var cmdUpload = &cobra.Command{
	Use:       "upload ",
	Short:     "upload subcommand data to the combine service",
	Long:      "upload data to the combine service",
	ValidArgs: []string{"server"},
	//Args:      cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println(args)
		//fmt.Printf("%+v\n", cmd.Flag("server"))
		//fmt.Printf("%+v\n", cmd.Flags())
		//fmt.Printf("%+v\n", cmd)

		if serverFlag != "" {
			//ErrorPrint("missing required flag server")
			//cmd.Help()
			//os.Exit(1)
			uploadArgsValid(serverFlag)
		}
		service.CollectFlag = true
		service.CollectSrv()
	},
}

func init() {
	// 添加 required-flag
	//cmdUpload.PersistentFlags().StringVarP(&serverFlag, "server", "s", "", "ex. 172.60.3.139:30980, required flag")
	cmdUpload.Flags().StringVarP(&serverFlag, "server", "s", "", "ex. 172.60.3.139:30980")
	rootCmd.AddCommand(cmdUpload)
}

func uploadArgsValid(serverAddr string) {
	uriList := strings.Split(serverAddr, ":")
	if len(uriList) != 2 {
		//log.Fatalf()
		//os.Exit(1)
		ErrorCmdPrint(fmt.Sprintf("flag server: %s format Fail", serverFlag), 1)
	}

	if !conf.SetServerAddr(uriList[0]) {
		//log.Fatalf("flag server: %s format Fail", fl.Server)
		//os.Exit(1)
		ErrorCmdPrint(fmt.Sprintf("flag server: %s format Fail", serverFlag), 1)
	}
	//conf.SetServerAddr(uriList[0])

	if !conf.SetServerPort(uriList[1]) {
		//log.Fatalf("flag server: %s format Fail", fl.Server)
		//os.Exit(1)
		ErrorCmdPrint(fmt.Sprintf("flag server: %s format Fail", serverFlag), 1)
	}
	//conf.SetServerPort(uriList[1])
}
