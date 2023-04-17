package collect

import (
	"collect_web/log"
	"collect_web/tools"
	"github.com/jaypipes/ghw"
)

type Gpu struct {
	*ghw.GPUInfo
}

//
func GetGPU() GetInfoInter {
	return &Gpu{}
}

func (g *Gpu) GetName() string {
	return "GPU"
}

func (g *Gpu) GetInfo(wraplog log.WrapLogInter) (interface{}, ErrorCollect) {
	var errors tools.MapStr = make(map[string]interface{})

	gpu, err := ghw.GPU(&ghw.WithOption{Alerter: wraplog})
	if err != nil {
		errors.Set("gpu", err)
		return nil, ErrorCollect(errors)
	}

	//return HostInfo
	if len(errors) > 0 {
		return nil, ErrorCollect(errors)
	}

	return gpu, nil
}
