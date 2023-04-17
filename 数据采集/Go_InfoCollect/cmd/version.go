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

		//fmt.Fprint(os.Stdout, "version: %s", conf.Version)
		fmt.Printf("version : %s\n", conf.Version)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
