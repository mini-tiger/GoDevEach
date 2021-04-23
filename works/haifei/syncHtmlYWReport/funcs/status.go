package funcs

import (
	tu "gitee.com/taojun319/tjtools/utils"
	"github.com/shirou/gopsutil/process"
	"haifei/syncHtmlYWReport/g"
	"os"
	"runtime"
	"sync"
	"time"
)

var startRunTime = time.Now().Unix()
var RunStatt = time.NewTicker(time.Duration(60) * time.Second)

//var gct = time.NewTicker(time.Duration(300) * time.Second)
//
//type RunStat struct {
//	mm                 runtime.MemStats
//	d, h, m, s, sinter *uint64
//	tm                 time.Time
//}

//var OnceRunStat = sync.Pool{
//	New: func() interface{} {
//		return &RunStat{}
//	},
//}
//
//var meminfo *process.MemoryInfoExStat

func mem() {
	//var ps *process.Process

	ps, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic(err)
	}
	//ps = p

	meminfo, _ := ps.MemoryInfoEx()
	_ = g.GetLog().Warn("内存监控 RSS占用:%d Kb, VSS占用:%d Kb \n", meminfo.RSS>>10, meminfo.VMS>>10)

}

var RunStatFree = sync.Pool{
	New: func() interface{} {
		return &runtime.MemStats{}
	},
}

func RunStatusWheel() {

	var tc = tu.TimeCountFree.Get().(*tu.TimeCount)
	//var tc = &tu.TimeCount{Sinter:uint64(time.Now().Unix() - startRunTime)}
	tc.Sinter = uint64(time.Now().Unix() - startRunTime)
	defer func() {
		tu.TimeCountFree.Put(tc)
	}()

	if g.GetConfig().Debug {
		var mm = RunStatFree.Get().(*runtime.MemStats)
		defer func() {
			RunStatFree.Put(mm)
		}()
		runtime.ReadMemStats(mm)
		g.GetLog().Debug("Sys(从系统获取过的内存,虚拟内存):%+v kb\n", mm.Sys/1024)
		g.GetLog().Debug("Alloc(golang语言框架堆空间分配的内存,go虚拟机分配的内存,同HeapAlloc):%+v kb\n", mm.Alloc/1024)
		g.GetLog().Debug("HeapReleased(回收到OS的内存):%+v kb\n", mm.HeapReleased/1024)
		g.GetLog().Debug("HeapInuse(堆上目前使用的内存):%+v kb\n", mm.HeapInuse/1024)
		g.GetLog().Debug("HeapIdle(堆上目前未使用的内存):%+v kb\n", mm.HeapIdle/1024)
		g.GetLog().Debug("HeapSys(系统获得的堆内存,包含缓存空间):%+v kb\n", mm.HeapSys/1024)
		g.GetLog().Debug("StackInuse(栈正在使用的内存):%+v kb\n", mm.StackInuse/1024)
		g.GetLog().Debug("StackSys(系统分配给运行栈的内存):%+v kb\n", mm.StackSys/1024)
		g.GetLog().Debug("GCSys(GC垃圾回收标记元信息使用的内存):%+v kb\n", mm.GCSys/1024)
		g.GetLog().Debug("NumGC:%+v \n", mm.NumGC)
		g.GetLog().Debug("OtherSys(golang架构占用的额外内存):%+v kb\n", mm.OtherSys/1024)
		g.GetLog().Debug("BySize:%+v kb\n", mm.BySize[60])
		g.GetLog().Debug("MSpanSys:%+v kb\n", mm.MSpanSys/1024)
		g.GetLog().Debug("MCacheInuse:%+v kb\n", mm.MCacheInuse/1024)
		g.GetLog().Debug("Mallocs:%+v \n", mm.Mallocs)
		g.GetLog().Debug("StackInuse:%+v kb\n", mm.StackInuse/1024)
		g.GetLog().Debug("StackSys:%+v kb\n", mm.StackSys/1024)
	}

	mem()

	//utils.GetTime(time.Now().Unix() - startRunTime)
	//d, h, m, s := utils.GetCurrentTime()
	//g.GetLog().Info("开始运行时间 %s,已经运行%3d天%2d小时%2d分钟%2d秒\n", time.Unix(startRunTime, 0).Format(g.TimeLayoutChi), *d, *h, *m, *s)
	tc.ComputeTime()

	g.GetLog().Info("开始运行时间 %s,已经运行%3d天%2d小时%2d分钟%2d秒,内存地址: %p\n", time.Unix(startRunTime, 0).Format(g.TimeLayoutChi),
		tc.Day, tc.H, tc.M, tc.S, tc)
}

func RunStatus() {
	//gct := time.NewTicker(time.Duration(g.GetConfig().Timeinter*6) * time.Second)

	//go func() {
	//	for {
	//		time.Sleep(time.Duration(g.GetConfig().Timeinter) * 6 * time.Second)
	//		g.GetLog().Info("手动GC开始\n")
	//		runtime.GC()
	//		g.GetLog().Info("手动GC结束\n")
	//	}
	//}()
	//
	//for {
	//	time.Sleep(time.Duration(60) * time.Second)
	//	go RunStatusWheel()
	//	//debug.FreeOSMemory()
	//
	//}
	for {
		select {
		//case <-gct.C:
		//	g.GetLog().Info("手动GC开始\n")
		//	runtime.GC()
		//	//debug.FreeOSMemory()
		//	g.GetLog().Info("手动GC结束\n")
		case <-RunStatt.C:
			RunStatusWheel()
		}
	}
}
