package main

import (
	"fmt"
	"github.com/jaypipes/ghw/pkg/unitutil"
	"math"

	"github.com/jaypipes/ghw"
)

func main() {
	memory, err := ghw.Memory()
	if err != nil {
		fmt.Printf("Error getting memory info: %v", err)
	}

	fmt.Println(memory.TotalPhysicalBytes/1024/1024/1024, memory.String(), memory)
	tpb := memory.TotalPhysicalBytes
	unit, _ := unitutil.AmountString(tpb)
	tpb = int64(math.Ceil(float64(memory.TotalPhysicalBytes) / float64(unit)))
	fmt.Println(tpb)
}
