package collect

import (
	"collect_web/log"
	"collect_web/tools"
	gopci "github.com/hertg/gopci/pkg/pci"
	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/host"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/7
 * @Desc: outbondIP.go
**/

type PCI struct {
	Devices []*DeviceEntry
}

type DeviceEntry struct {
	Addr      string
	Vendor    interface{} // "vendor"
	Product   interface{} // product
	Class     interface{}
	SubDevice interface{} // subsystem

}

func GetPCI() GetInfoInter {
	return new(PCI)
}

func (i *PCI) GetName() string {
	return "PCI"
}

func (i *PCI) GetInfo(wlog log.WrapLogInter) (interface{}, ErrorCollect) {
	var errors tools.MapStr = make(map[string]interface{})

	var platform string
	hostInfo, err := host.Info()
	if err != nil {
		wlog.Error(err, "get pci err,in hostinfo")
	} else {
		platform = hostInfo.Platform
	}

	if !strings.Contains(strings.ToLower(platform), "openeuler") {
		devices, err := gopci.Scan()
		if err != nil {
			errors.Set("pci", err)
		}
		for _, device := range devices {
			de := &DeviceEntry{
				Addr: device.Address.Hex(),
				//SubDevice: device.Product,
				Product: device.Product.Label,
				Vendor:  device.Vendor.Label,
				Class:   device.Class.Label,
				//Driver:    *device.Driver,
				SubDevice: device.Subdevice.Label,
			}

			//fmt.Printf("%+v\n", device)
			//fmt.Printf("%+v\n", device.Address.Hex())
			//fmt.Printf("subdevice: %+v\n", device.Subdevice)
			//fmt.Printf("config: %+v\n", device.Config)
			//fmt.Printf("subvendor: %+v\n", device.Subvendor)
			//if device.Driver != nil {
			//	fmt.Printf("Drvice: %+v\n", *device.Driver)
			//}
			//fmt.Printf("%+v\n", de)
			i.Devices = append(i.Devices, de)
		}

	} else {
		pciall, err := ghw.PCI(&ghw.WithOption{Alerter: wlog})
		if err != nil {
			errors.Set("pci", err)
		}
		for _, device := range pciall.Devices {
			de := &DeviceEntry{
				Addr:      device.Address,
				SubDevice: device.Subclass.Name,

				Product: device.Product.Name,
				Vendor:  device.Vendor.Name,
				Class:   device.Subclass.Name,
			}
			i.Devices = append(i.Devices, de)
		}

	}

	if len(errors) > 0 {
		return nil, ErrorCollect(errors)
	}
	return i, nil
}
