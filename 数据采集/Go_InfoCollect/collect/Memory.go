package collect

import (
	"collect_web/log"
	"collect_web/tools"
	"github.com/shirou/gopsutil/mem"
)

//SwapMemory	  交换内存
//VirtualMemory	  虚拟内存
//VirtualMemoryEx 虚拟交换内存
type Memory struct {
	SwapMemory      *mem.SwapMemoryStat      `json:"swapMemory"`
	VirtualMemory   *mem.VirtualMemoryStat   `json:"virtualMemory"`
	VirtualMemoryEx *mem.VirtualMemoryExStat `json:"virtualMemoryEx"`
}

//
func GetMemory() GetInfoInter {
	return new(Memory)
}

func (m *Memory) GetName() string {
	return "memory"
}

func (m *Memory) GetInfo(wlog *log.Wraplog) (interface{}, ErrorCollect) {
	var errors tools.MapStr = make(map[string]interface{})
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		errors.Set("swapInfo", err)
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		errors.Set("VirtualMemory", err)
	}
	memExinfo, err := mem.VirtualMemoryEx()
	if err != nil {
		errors.Set("VirtualMemoryEx", err)
	}

	m.SwapMemory = swapInfo
	m.VirtualMemory = memInfo

	m.VirtualMemoryEx = memExinfo
	if len(errors) > 0 {
		return m, ErrorCollect(errors)
	}
	return m, nil
}

//
//func main(){
//	fmt.Println(GetVirtualMemExInfo())
//	fmt.Println(GetVirtualMemInfo())
//}
