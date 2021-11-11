package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"strings"
)

// xxx https://www.cnblogs.com/sparkdev/p/10856077.html
// imageCmd represents the image command

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Print images information",
	Long:  "Print all images information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("image one is ubuntu 16.04")
		fmt.Println("image two is ubuntu 18.04")
		fmt.Println("image args are : " + strings.Join(args, " "))
	},
}

var echoTimes int
var cmdTimes = &cobra.Command{
	Use:   "times [string to echo]",
	Short: "Echo anything to the screen more times",
	Long: `echo things multiple times back to the user by providing
a count and a string.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < echoTimes; i++ {
			fmt.Println("Echo: " + strings.Join(args, " "))
		}
	},
}
var rootCmd = &cobra.Command{}

func init() {
	rootCmd.AddCommand(imageCmd)
	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")
	imageCmd.AddCommand(cmdTimes)
	rootCmd.Execute()
}
