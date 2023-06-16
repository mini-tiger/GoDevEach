package main

import (
	"fmt"
	"sort"

	"github.com/jaypipes/pcidb"
)

type ByCountProducts []*pcidb.Vendor

func (v ByCountProducts) Len() int {
	return len(v)
}

func (v ByCountProducts) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ByCountProducts) Less(i, j int) bool {
	return len(v[i].Products) > len(v[j].Products)
}

func main() {
	pci, err := pcidb.New()
	if err != nil {
		fmt.Printf("Error getting PCI info: %v", err)
	}

	vendors := make([]*pcidb.Vendor, len(pci.Vendors))
	x := 0
	for _, vendor := range pci.Vendors {
		vendors[x] = vendor
		x++
	}

	sort.Sort(ByCountProducts(vendors))

	fmt.Println(" vendors by product")
	fmt.Println("====================================================")
	for _, vendor := range vendors {
		fmt.Printf("%v ('%v') has %d products\n", vendor.Name, vendor.ID, len(vendor.Products))
		for _, p := range vendor.Products {
			fmt.Printf("\t %+v\n", p)
		}
	}
}
