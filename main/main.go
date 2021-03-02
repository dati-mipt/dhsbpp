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

	var firstDistribution = vizualize.NewDotTree(bins)
	err = vizualize.DrawBinTreeDot(firstDistribution, "FirstDistribution.dot")
	// dot -Tpng FirstDistribution.dot -o FirstDistribution.png
	if err != nil {
		fmt.Println(err)
		return
	}

	nameToPartitionNode, _ := pRoot.MapNameToPartitionNode()
	var loadedBin = packing.FindBinForRebalancing(&bins, newHierarchy.TasksPerDay, nameToPartitionNode)

	var secondDistribution = vizualize.NewDotTree(bins)
	err = vizualize.DrawBinTreeDot(secondDistribution, "SecondDistribution.dot")
	// dot -Tpng SecondDistribution.dot -o SecondDistribution.png

	if loadedBin != nil {
		var migrationSize = packing.DynamicalHierarchicalFirstFitDecreasing(loadedBin, &bins)
		fmt.Println("Bin Index:", loadedBin.Index)
		fmt.Println("Migration Size:", migrationSize)
		var thirdDistribution = vizualize.NewDotTree(bins)
		err = vizualize.DrawBinTreeDot(thirdDistribution, "ThirdDistribution.dot")
		// dot -Tpng ThirdDistribution.dot -o ThirdDistribution.png
	}

}
