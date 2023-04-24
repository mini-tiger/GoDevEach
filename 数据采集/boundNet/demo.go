package main

import (
	"fmt"
	"github.com/digitalocean/go-openvswitch/ovs"
	"os/exec"
)

func main() {
	// Create a connection to the OVS database
	client := ovs.New(ovs.Sudo())

	// Get all the interfaces from the OVS database
	bridges, err := client.VSwitch.ListBridges()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, bridge := range bridges {
		fmt.Println("Bridge:", bridge)
		b, err := exec.Command("ovs-vsctl", "list-ifaces", bridge).CombinedOutput()
		fmt.Println(string(b))
		fmt.Println(err)
	}
}
