package main

import (
	"fmt"

	"github.com/dati-mipt/dhsbpp/hierarchy"
	"github.com/dati-mipt/dhsbpp/packing"
	"github.com/dati-mipt/dhsbpp/tree"
	"github.com/dati-mipt/dhsbpp/vizualize"
)

func main() {
	var region = "au1"
	var newHierarchy = hierarchy.NewHierarchy("datasets/"+region+"-baas-data.csv",
		"datasets/"+region+"-baas-tpd.csv")

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

	packing.PreprocessPartitionTree(pRoot)
	var bins = make([]*packing.Bin, 0)
	bins = packing.HierarchicalFirstFitDecreasing(pRoot, bins)

	err = vizualize.MakeVisualizationPicture(bins, "1distribution.png", "./outputPics/")
	if err != nil {
		fmt.Println(err)
		return
	}

	nameToPartitionNode, _ := pRoot.MapNameToPartitionNode()
	var loadedBin = packing.FindBinForRebalancing(bins, newHierarchy.TasksPerDay, nameToPartitionNode)

	err = vizualize.MakeVisualizationPicture(bins, "2distribution.png", "./outputPics/")
	if err != nil {
		fmt.Println(err)
		return
	}

	if loadedBin != nil {
		var migrationSize int64
		bins, migrationSize = packing.DynamicalHierarchicalFirstFitDecreasing(loadedBin, bins)
		fmt.Println(len(bins))
		fmt.Println("Bin Index:", loadedBin.Index)
		fmt.Println("Migration Size:", migrationSize)

		err = vizualize.MakeVisualizationPicture(bins, "3distribution.png", "./outputPics/")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
