package main

import (
	"fmt"
	"github.com/netdata/go.d.plugin/pkg/iprange"

	"net"
)

// https://pkg.go.dev/github.com/netdata/go.d.plugin/pkg/iprange#section-readme

var ips string = "192.168.43.111/16 192.168.44.111/24 192.168.43.111 192.0.2.0/255.255.255.0 192.0.3.0-192.0.3.10" // 有效格式

func main() {
	r, _ := iprange.ParseRange("192.168.43.1/24")
	fmt.Println(r.Contains(net.ParseIP("192.168.43.28")),
		r.Contains(net.ParseIP("192.168.40.111")),
		r.Contains(net.ParseIP("192.168.0.22")),
		r.Contains(net.ParseIP("192.0.3.5")))
	fmt.Println(r.Family(), r.String(), r.Size())

	fmt.Println("------------")

	rr, _ := iprange.ParseRanges(ips)
	fmt.Printf("%#v\n", len(rr))

	fmt.Println("------------")

	p := iprange.Pool(rr)
	fmt.Println(p.Contains(net.ParseIP("192.168.43.28")),
		p.Contains(net.ParseIP("192.168.43.111")),
		p.Contains(net.ParseIP("192.168.0.22")),
		p.Contains(net.ParseIP("192.0.3.5")))

}
