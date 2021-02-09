package main

import (
	"DHSBPP/tenantsData"
	"DHSBPP/tree"
	"fmt"
)

func main() {
	td := tenantsData.NewTenantsData("tenantsData/sg_data.csv",
		                             "tenantsData/sg_tpd.csv")

	root, nodeName, ok := tree.NewTree( td.TenantToParent )
	fmt.Println( root, ok)
	tree.SetInitialSize( td.Tpd, nodeName, 30 )
}
