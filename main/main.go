package main

import (
	"fmt"
	"github.com/dati-mipt/dhsbpp/hierarchy"
	"github.com/dati-mipt/dhsbpp/packing"
	"github.com/dati-mipt/dhsbpp/tree"
	"github.com/dati-mipt/dhsbpp/vizualize"
)

func main() {
	newHierarchy := hierarchy.NewHierarchy("datasets/au1-baas-data.csv",
		"datasets/au1-baas-tpd.csv")

	var root, err = tree.NewTree(newHierarchy.ChildToParent)
	if err != nil {
		fmt.Println(err)
		return
	}

	var pRoot = tree.NewPartitionTree(root)

	err = pRoot.SetInitialSize(newHierarchy.TasksPerDay, packing.InitDays)
	if err != nil {
		fmt.Println(err)
		return
	}

	var bins = make([]*packing.Bin, 0)
	packing.PreprocessPartitionTree(pRoot)
	packing.HierarchicalFirstFitDecreasing(pRoot, &bins)

	var rootDotNodeBefore = vizualize.NewDotTree(bins)
	err = vizualize.DrawBinTreeDot(rootDotNodeBefore, "DotTreeBefore.dot")
	// dot -Tpng DotTreeBefore.dot -o DotTreeBefore.png
	if err != nil {
		fmt.Println(err)
		return
	}

	nameToPartitionNode, _ := pRoot.MapNameToPartitionNode()
	packing.DynamicalHierarchicalFirstFitDecreasing(&bins, newHierarchy.TasksPerDay, nameToPartitionNode)

	var rootDotNodeAfter = vizualize.NewDotTree(bins)
	err = vizualize.DrawBinTreeDot(rootDotNodeAfter, "DotTreeAfter.dot")
	// dot -Tpng DotTreeAfter.dot -o DotTreeAfter.png
}
