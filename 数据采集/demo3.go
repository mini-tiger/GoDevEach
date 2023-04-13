package main

import (
	"fmt"
	"github.com/jaypipes/ghw"
)

func main() {
	pci, err := ghw.PCI()
	if err != nil {
		fmt.Printf("Error getting PCI info: %v", err)
	}
	fmt.Printf("host PCI devices:\n")
	fmt.Println("====================================================")

	for _, device := range pci.Devices {
		//getDetail(pci, device.Address)
		fmt.Printf("%+v\n", device)
		//vendor := device.Vendor
		//vendorName := vendor.Name
		//if len(vendor.Name) > 200 {
		//	vendorName = string([]byte(vendorName)[0:198]) + "..."
		//}
		//product := device.Product
		//productName := product.Name
		//if len(product.Name) > 400 {
		//	productName = string([]byte(productName)[0:397]) + "..."
		//}
		//fmt.Printf("%-12s\t%-20s\t%-40s\n", device.Address, vendorName, productName)
	}
}
func getDetail(pci *ghw.PCIInfo, addr string) {

	fmt.Printf("PCI device information for %s\n", addr)
	fmt.Println("====================================================")
	deviceInfo := pci.GetDevice(addr)
	if deviceInfo == nil {
		fmt.Printf("could not retrieve PCI device information for %s\n", addr)
		return
	}

	vendor := deviceInfo.Vendor
	fmt.Printf("Vendor: %s [%s]\n", vendor.Name, vendor.ID)
	product := deviceInfo.Product
	fmt.Printf("Product: %s [%s]\n", product.Name, product.ID)
	subsystem := deviceInfo.Subsystem
	subvendor := pci.Vendors[subsystem.VendorID]
	subvendorName := "UNKNOWN"
	if subvendor != nil {
		subvendorName = subvendor.Name
	}
	fmt.Printf("Subsystem: %s [%s] (Subvendor: %s)\n", subsystem.Name, subsystem.ID, subvendorName)
	class := deviceInfo.Class
	fmt.Printf("Class: %s [%s]\n", class.Name, class.ID)
	subclass := deviceInfo.Subclass
	fmt.Printf("Subclass: %s [%s]\n", subclass.Name, subclass.ID)
	progIface := deviceInfo.ProgrammingInterface
	fmt.Printf("Programming Interface: %s [%s]\n", progIface.Name, progIface.ID)
}
