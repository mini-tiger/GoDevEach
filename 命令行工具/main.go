package main

import (
	"fmt"
	"os"

	"clitool/cmd"
)

// https://o-my-chenjian.com/2017/09/20/Using-Cobra-With-Golang/
//
func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
