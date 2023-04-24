package tools

import (
	"github.com/jaypipes/ghw/pkg/unitutil"
	"math"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/4/19
 * @Desc: mathUtils.go
**/
func ToGB(tpb int64) int64 {
	//tpb := memory.TotalPhysicalBytes
	unit, _ := unitutil.AmountString(tpb)
	return int64(math.Ceil(float64(tpb) / float64(unit)))
}
