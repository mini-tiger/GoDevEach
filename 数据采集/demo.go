package main

import (
	"fmt"
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
	"log"
)

func main() {
	process, err := sysinfo.Self()
	if err != nil {
		panic(err)
	}
	fmt.Println(process.Info())
	fmt.Println(process.Parent())
	if handleCounter, ok := process.(types.OpenHandleCounter); ok {
		count, err := handleCounter.OpenHandleCount()
		if err != nil {
			panic(err)
		}
		log.Printf("%d open handles", count)
	}

}
