package main

import (
	"fmt"
	"math"
)

type DiskInfo struct {
	Device     string `json:"device"`
	MountPoint string `json:"mountpoint"`
	FsType     string `json:"fstype"`
	Total      uint64 `json:"total"`
}

var a uint64

func main() {
	a = 18446744073692774000
	//b:=2**64
	fmt.Println(math.Pow(2, 64))
	fmt.Println(a)

}
