package main

import (
	"fmt"
	"github.com/dati-mipt/dhsbpp/hierarchy"
	"github.com/dati-mipt/dhsbpp/packing"
	"github.com/dati-mipt/dhsbpp/partitionTree"
	"github.com/dati-mipt/dhsbpp/tree"
	"github.com/dati-mipt/dhsbpp/vizualize"
)

func main() {
	newHierarchy := hierarchy.NewHierarchy("datasets/au1-baas-data.csv",
		"datasets/au1-baas-tpd.csv")

	root, err := tree.NewTree(newHierarchy.ChildToParent)
	if err != nil {
		fmt.Println(err)
		return
	}

	var pRoot = partitionTree.NewPartitionTree(root)

	const initDays = 100
	err = pRoot.SetInitialSize(newHierarchy.TasksPerDay, initDays)
	if err != nil {
		fmt.Println(err)
		return
	}

	const maxCapacity = 250
	const allocationFactor = 60
	const reallocationDelta = 20
	bins := make([]*packing.Bin, 0)
	packing.PreprocessPartitionTree(pRoot, maxCapacity, allocationFactor)
	packing.HierarchicalFirstFitDecreasing(pRoot, &bins, maxCapacity, allocationFactor, reallocationDelta)

	rootBin := vizualize.NewTreeOfBins(bins)
	err = vizualize.DrawBinTreeDot(rootBin, "BinTree.dot")
	if err != nil {
		fmt.Println(err)
		return
	}

	nameToPartitionNode, _ := pRoot.MapNameToPartitionNode()
	packing.DynamicalHierarchicalFirstFitDecreasing(&bins, newHierarchy.TasksPerDay,
		nameToPartitionNode, initDays, maxCapacity, allocationFactor, reallocationDelta)
	rootBinAfter := vizualize.NewTreeOfBins(bins)
	err = vizualize.DrawBinTreeDot(rootBinAfter, "AfterBinTree.dot")
}
