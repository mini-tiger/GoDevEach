package main

import (
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	// "github.com/shirou/gopsutil/mem"  // to use v2
)

func main() {
	v, _ := mem.VirtualMemory()
	infos, _ := disk.Partitions(true)
	for _, info := range infos {
		data, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(data))
	}
	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	// convert to JSON. String() is also implemented
	fmt.Println(v)

}
