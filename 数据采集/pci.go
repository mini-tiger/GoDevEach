package main

import (
	"fmt"
	"github.com/hertg/gopci/pkg/pci"
	"github.com/jaypipes/ghw"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/14
 * @Desc: main.go
**/

func main() {
	devices, _ := pci.Scan()
	for index, device := range devices {
		fmt.Printf("[%d] %+v\n", index, *device)
		fmt.Printf("[%d] %+v\n", index, device.Address.Hex())
		//fmt.Printf("%+v\n", device)
		if device.Driver != nil {
			fmt.Printf("[%d] Drvice: %+v\n", index, *device.Driver)
		}
		fmt.Printf("[%d] subdevice: %+v\n", index, device.Subdevice)
		fmt.Printf("[%d] config: %+v\n", index, device.Config)
		fmt.Printf("[%d] vendor: %+v\n", index, device.Vendor)
		fmt.Printf("[%d] subvendor: %+v\n", index, device.Subvendor)
		fmt.Printf("[%d] class: %+v\n", index, device.Class)
		fmt.Printf("[%d] product: %+v\n", index, device.Product)

	}
	fmt.Println("================================================================")
	pciall, err := ghw.PCI()
	if err != nil {
		panic(err)
	}
	for index, device := range pciall.Devices {
		fmt.Printf("[%d] device: %+v\n", index, device)
		fmt.Printf("[%d] device.device: %+v\n", index, device.Driver)
		fmt.Printf("[%d] address: %+v\n", index, device.Address)
		fmt.Printf("[%d] vendor: %+v\n", index, device.Vendor.Name)
		fmt.Printf("[%d] product: %+v\n", index, device.Product.Name)
		fmt.Printf("[%d] subsystem: %+v\n", index, device.Subsystem.Name)
		fmt.Printf("[%d] class: %+v\n", index, device.Class.Name)
		fmt.Printf("[%d] subclass: %+v\n", index, device.Subclass.Name)
	}
}
