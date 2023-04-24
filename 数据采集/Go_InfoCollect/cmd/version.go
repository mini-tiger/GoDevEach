package cmd

import (
	"collect_web/conf"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version subcommand show app version info.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Flags().NFlag(), len(args), args)
		if cmd.Flags().NFlag() > 0 {
			fmt.Println("错误：此命令不接受任何 flag")
			_ = cmd.Help()
			os.Exit(1)
		}

		fmt.Printf("version : %s\n", conf.Version)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
