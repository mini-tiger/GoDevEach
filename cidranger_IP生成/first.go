package main

import (
	"fmt"
	"github.com/yl2chen/cidranger"
	"log"
	"net"
)

type basicRangerEntry struct {
	ipNet net.IPNet
	asn   string
}

// get function for network
func (b *basicRangerEntry) Network() net.IPNet {
	return b.ipNet
}

// create customRangerEntry object using net and asn
func NewBasicRangerEntry(ipNet net.IPNet) cidranger.RangerEntry {
	return &basicRangerEntry{
		ipNet: ipNet,
	}
}
func main() {

	// instantiate NewPCTrieRanger
	ranger := cidranger.NewPCTrieRanger()
	_, network1, _ := net.ParseCIDR("192.168.1.0/24")
	_, network2, _ := net.ParseCIDR("10.11.0.0/16")
	ranger.Insert(NewBasicRangerEntry(*network1))
	ranger.Insert(NewBasicRangerEntry(*network2))

	contains, err := ranger.Contains(net.ParseIP("10.11.1.0")) // returns true, nil
	fmt.Println(contains)
	contains, err = ranger.Contains(net.ParseIP("192.168.1.255")) // returns false, nil
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(contains)

	containingNetworks, err := ranger.ContainingNetworks(net.ParseIP("192.168.1.0"))
	fmt.Println(containingNetworks[0])

	//var AllIPv4 *net.IPNet
	entries, err := ranger.CoveredNetworks(*network2) // for IPv4
	fmt.Println(entries[0])
}
