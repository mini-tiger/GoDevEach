package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
	"os/exec"
)

var goDaemon bool //后台 服务  -d=true

func main() {
	flag.BoolVarP(&goDaemon, "daemon", "d", false, "run app as a daemon with -d=true.")
	flag.Parse()
	if goDaemon {
		fmt.Println(os.Args[0])
		fmt.Println(os.Args[1:])

		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		if err := cmd.Start(); err != nil {
			fmt.Printf("start %s failed, error: %v\n", os.Args[0], err)
			os.Exit(1)
		}
		fmt.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid)
		os.Exit(0)
	}
}
