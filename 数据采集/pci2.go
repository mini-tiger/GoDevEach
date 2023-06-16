package main

import (
	"fmt"
	"github.com/jaypipes/ghw"
	"os/exec"
	"strings"
)

type LsPciInfo struct {
	DeviceID string            `json:"device_id"`
	Info     map[string]string `json:"info"`
}
type DeviceEntry struct {
	Addr      string
	Vendor    interface{} // "vendor"
	Product   interface{} // product
	Class     interface{}
	SubDevice interface{} // subsystem

}

func main() {
	pciall, _ := ghw.PCI()
	fmt.Println(pciall)

	//pciall, _ = ghw.PCI()

	for _, device := range pciall.Devices {
		de := &DeviceEntry{
			Addr:      device.Address,
			SubDevice: device.Subclass.Name,

			Product: device.Product.Name,
			Vendor:  device.Vendor.Name,
			Class:   device.Subclass.Name,
		}
		fmt.Printf("%+v\n", de)
		bj, _ := device.MarshalJSON()
		fmt.Printf("%+v\n", string(bj))
		fmt.Println("==================================")
	}

	//if len(os.Args) < 2 {
	//	fmt.Println("Usage: lspci2json <device_id>")
	//	return
	//}
	//
	//deviceID := os.Args[1]
	//lspciData, err := getLspciInfo(deviceID)
	//if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	return
	//}
	//
	//jsonData, err := json.MarshalIndent(lspciData, "", "  ")
	//if err != nil {
	//	fmt.Printf("Error encoding JSON: %v\n", err)
	//	return
	//}
	//
	//fmt.Printf("%s\n", jsonData)
}

func getLspciInfo(deviceID string) (LspciInfo, error) {
	cmd := exec.Command("lspci", "-v", "-s", deviceID)
	output, err := cmd.Output()
	if err != nil {
		return LspciInfo{}, fmt.Errorf("error executing command: %v", err)
	}

	lspciData := parseLspciOutput(string(output))

	return lspciData, nil
}

func parseLspciOutput(output string) LspciInfo {
	lines := strings.Split(output, "\n")
	deviceID := strings.TrimSpace(strings.Split(lines[0], ":")[0])

	info := make(map[string]string)
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			info[key] = value
		}
	}

	return LspciInfo{DeviceID: deviceID, Info: info}
}
