package main

import (
	"dhsbpp/packing"
	"dhsbpp/tenantsData"
	"dhsbpp/tree"
)

func main() {
	TD := tenantsData.NewTenantsData("tenantsData/sg_data.csv",
		                             "tenantsData/sg_tpd.csv")

	root, _ := tree.NewTree( &TD )

	MAX, AF := 250.0,  0.6
	bins := make([]*packing.Bin, 0)
	packing.PreprocessTree( root, MAX, AF )
	packing.HFFD( root, &bins, MAX, AF, 0.2)
}
