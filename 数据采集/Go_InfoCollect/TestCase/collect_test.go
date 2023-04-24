package TestCase

import (
	"collect_web/collect"
	"collect_web/log"
	"github.com/jaypipes/ghw"
	"testing"
)

//var tLog *tlog = new(tlog)

func TestCollect(t *testing.T) {
	t.Run("cpu", func(t *testing.T) {
		//fmt.Println(t.Name())
		c := collect.GetCpu()

		data, err := c.GetInfo(log.Tloginst)
		if err != nil {
			t.Fatalf("%s err:%v\n", t.Name(), err)
		}
		_, ok := data.(*collect.Cpu)
		if !ok {
			t.Fatalf("%s err:%v\n", t.Name(), "struct err")
		}

	})
	t.Run("mem", func(t *testing.T) {
		c := collect.GetMemory()

		data, err := c.GetInfo(log.Tloginst)
		if err != nil {
			t.Fatalf("%s err:%v\n", t.Name(), err)
		}
		_, ok := data.(*collect.Memory)
		if !ok {
			t.Fatalf("%s err:%v\n", t.Name(), "struct err")
		}
	})
	t.Run("pci", func(t *testing.T) {
		c := collect.GetPCI()

		data, err := c.GetInfo(log.Tloginst)
		if err != nil {
			t.Fatalf("%s err:%v\n", t.Name(), err)
		}
		_, ok := data.(*collect.PCI)
		if !ok {
			t.Fatalf("%s err:%v\n", t.Name(), "struct err")
		}
	})
	t.Run("gpu", func(t *testing.T) {
		c := collect.GetGPU()

		data, err := c.GetInfo(log.Tloginst)
		if err != nil {
			t.Fatalf("%s err:%v\n", t.Name(), err)
		}
		_, ok := data.(*ghw.GPUInfo)
		if !ok {
			t.Fatalf("%s err:%v,data:%v\n", t.Name(), "struct err", data)
		}
	})
	t.Run("network", func(t *testing.T) {
		c := collect.GetNetIfaces()

		data, err := c.GetInfo(log.Tloginst)
		if err != nil {
			t.Fatalf("%s err:%v\n", t.Name(), err)
		}
		_, ok := data.(*collect.NetInterfaces)
		if !ok {
			t.Fatalf("%s err:%v,data:%v\n", t.Name(), "struct err", data)
		}
	})

}
