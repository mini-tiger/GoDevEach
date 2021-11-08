package services

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/host"

	"github.com/shirou/gopsutil/cpu"
)

type CpuCollector struct {
	cpuDesc *prometheus.Desc
}

func NewCpuCollector() prometheus.Collector {
	host, _ := host.Info()
	hostname = host.Hostname
	return &CpuCollector{
		cpuDesc: prometheus.NewDesc(
			"cpu_usage",
			"cpu",
			[]string{"DYNAMIC_HOST_NAME"}, //动态标签名称 xxx 这里要和 collect 方法里面的值对应
			prometheus.Labels{"STATIC_LABEL1": "静态值可以放在这里", "static_HOST_NAME": hostname}), // todo 静态label
	}
}

//实现采集器Describe接口
func (n *CpuCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.cpuDesc

}

// Collect returns the current state of all metrics of the collector.
//实现采集器Collect接口,真正采集动作
func (n *CpuCollector) Collect(ch chan<- prometheus.Metric) {

	host, _ := host.Info()
	hostname = host.Hostname
	infos, _ := cpu.Times(true)
	sum := float64(0)
	for _, info := range infos {
		sum = sum + info.User
	}
	fmt.Println(sum / float64(len(infos)))
	ch <- prometheus.MustNewConstMetric(n.cpuDesc, prometheus.GaugeValue, sum/float64(len(infos)), hostname) // xxx 对应动态标签名称

}
