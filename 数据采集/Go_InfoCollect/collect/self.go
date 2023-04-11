package collect

import (
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/6
 * @Desc: self.go
**/

func GetSelfProcess() (err error, pro types.ProcessInfo) {
	// self
	process, err := sysinfo.Self()
	if err != nil {
		return
	}
	pro, err = process.Info()
	if err != nil {
		return
	}
	return
}
