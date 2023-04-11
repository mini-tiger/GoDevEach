package collect

import (
	"collect_web/log"
	"collect_web/tools"
	"github.com/shirou/gopsutil/disk"
)

//Patrions:磁盘分区
//Usage：使用情况
//SerialNumber：序列号
//Label：标签
//IO：磁盘IO信息
type Disk struct {
	Physics bool        `json:"physics"`
	Disks   []*DiskPath `json:"disks"`
	//IO           []map[string]disk.IOCountersStat `json:"io"`
}
type DiskPath struct {
	Partitions   disk.PartitionStat `json:"partitions"`
	Usage        *disk.UsageStat    `json:"usage"`
	SerialNumber string             `json:"serialNumber"`
	Label        string             `json:"label"`
	//IO           map[string]disk.IOCountersStat `json:"io"`
}

func GetDisk() GetInfoInter {
	return &Disk{
		Physics: false,
	}
}
func (d *Disk) GetName() string {
	return "disk"
}

//获取磁盘使用情况
func (d *Disk) GetUsage(path string) *disk.UsageStat {
	usage, err := disk.Usage(path)
	if err == nil {
		return usage
	} else {
		return nil
	}
}

//得到磁盘序列号
func (d *Disk) GetSerialNumber(name string) string {
	//获取路径为name的磁盘的序列号
	sn := disk.GetDiskSerialNumber(name)
	return sn
}

//得到磁盘标签
func (d *Disk) GetLabel(name string) string {
	label := disk.GetLabel(name)
	return label
}

//得到磁盘IO信息
func (d *Disk) GetIO(name string) map[string]disk.IOCountersStat {
	IO, err := disk.IOCounters(name)
	if err == nil {
		return IO
	} else {
		return nil
	}
}

func (d *Disk) GetInfo(wraplog *log.Wraplog) (interface{}, ErrorCollect) {
	var errors tools.MapStr = make(map[string]interface{})

	//获取磁盘分区
	//如果all为false，只返回物理设备(如:硬盘、cd-rom驱动器、USB keys)，忽略其他所有设备(如:内存分区，如/dev/shm)
	DiskParti, err := disk.Partitions(d.Physics)
	if err != nil {
		errors.Set("DiskPartitions", err)
		return nil, ErrorCollect(errors)
	}

	//获取分区路径

	for _, n := range DiskParti {
		diskpath := &DiskPath{
			Partitions:   n,
			Usage:        d.GetUsage(n.Mountpoint),
			SerialNumber: d.GetSerialNumber(n.Device),
			Label:        d.GetLabel(n.Device),
		}
		d.Disks = append(d.Disks, diskpath)
	}
	if len(errors) > 0 {
		return d, ErrorCollect(errors)
	}
	return d, nil
}
