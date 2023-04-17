package main

import (
	"fmt"
	"io/ioutil"
	"syscall"
)

func main() {
	// PCI设备的设备路径
	devicePath := "/sys/class/pci_bus"

	// 获取所有PCI设备的列表
	devices, err := ioutil.ReadDir(devicePath)
	if err != nil {
		fmt.Println("无法获取PCI设备列表：", err)
		return
	}

	// 遍历每个PCI设备并输出其信息
	for _, device := range devices {
		fmt.Println("设备名称：", device.Name())
		fmt.Println("设备ID：", device.Sys().(*syscall.Stat_t).Dev)
		// 其他信息...
	}
}
