package main

import (
	"fmt"
	"github.com/dati-mipt/dhsbpp/hierarchy"
	"github.com/dati-mipt/dhsbpp/packing"
	"github.com/dati-mipt/dhsbpp/tree"
)

func main() {
	newHierarchy := hierarchy.NewHierarchy("datasets/sg_data.csv",
		"datasets/sg_tpd.csv")

	root, err := tree.NewTree(newHierarchy)

	if err != nil {
		fmt.Println(err)
		return
	}

	const maxCapacity = 250.0
	const allocationFactor = 0.6

	bins := make([]*packing.Bin, 0)
	packing.PreprocessTree(root, maxCapacity, allocationFactor)
	packing.HierarchicalFirstFitDecreasing(root, &bins, maxCapacity, allocationFactor)
}
