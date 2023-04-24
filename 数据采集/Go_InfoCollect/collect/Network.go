package collect

import (
	"collect_web/log"
	"collect_web/tools"
	"fmt"
	"github.com/digitalocean/go-openvswitch/ovs"
	"os/exec"
	"strings"
)

type Ovs struct {
	Bridge map[string]interface{} `json:"Bridge"`
}

type OvsBr struct {
	Ports  []string `json:"Ports"`
	Ifaces []string `json:"IFaces"`
}

func GetOvs() GetInfoInter {
	return new(Ovs)
}

func (i *Ovs) GetName() string {
	return "Ovs"
}

func (i *Ovs) GetInfo(wraplog log.WrapLogInter) (interface{}, ErrorCollect) {
	var errors tools.MapStr = make(map[string]interface{})
	// Create a connection to the OVS database
	client := ovs.New(ovs.Sudo())

	// Get all the interfaces from the OVS database
	bridges, err := client.VSwitch.ListBridges()
	if err != nil {
		//fmt.Println(err)
		errors.Set(i.GetName(), fmt.Sprintf("out:%v,err:%v", bridges, err))
		return i, ErrorCollect(errors)
	}

	//i.Ifaces = make(map[string]interface{}, len(bridges))
	//i.Ports = make(map[string]interface{}, len(bridges))
	i.Bridge = make(map[string]interface{}, len(bridges))
	for _, bridge := range bridges {
		br := new(OvsBr)
		b, err := exec.Command("ovs-vsctl", "list-ifaces", bridge).CombinedOutput()
		if err != nil {
			continue
		}

		br.Ifaces = strings.Split(strings.TrimSpace(string(b)), "\n")

		b, err = exec.Command("ovs-vsctl", "list-ports", bridge).CombinedOutput()
		if err != nil {
			continue
		}
		br.Ports = strings.Split(strings.TrimSpace(string(b)), "\n")
		i.Bridge[bridge] = br
	}
	return i, ErrorCollect(errors)
}
