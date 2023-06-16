package main

import (
	"encoding/json"
	"fmt"

	"github.com/jaypipes/ghw"
)

func main() {
	block, err := ghw.Block()
	if err != nil {
		fmt.Printf("Error getting block storage info: %v", err)
	}

	fmt.Printf("%v\n", block)

	for _, disk := range block.Disks {
		bjson, _ := json.Marshal(disk)
		fmt.Printf(" %v\n", string(bjson))
		//for _, part := range disk.Partitions {
		//	fmt.Printf("  %+v\n", part)
		//}
	}
}
