package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
)

func main() {
	infos, _ := cpu.Times(true)
	sum := float64(0)
	for _, info := range infos {
		sum = sum + info.User
	}
	fmt.Println(sum / float64(len(infos)))
}
