package packing

import (
	"github.com/dati-mipt/dhsbpp/tree"
	"sort"
)

const (
	MaxCapacity       = 250
	AllocationFactor  = 60
	ReallocationDelta = 20
	InitDays          = 100

	Volume             = MaxCapacity * AllocationFactor / 100
	OverloadThreshold  = MaxCapacity * (AllocationFactor + ReallocationDelta) / 100
	UnderloadThreshold = MaxCapacity * (AllocationFactor - ReallocationDelta) / 100
)

type Bin struct {
	Index     int
	Size      int64
	PartNodes map[*tree.PartitionNode]bool
}

func NewBin(idx int) *Bin {
	var bin Bin

	bin.Index = idx
	bin.Size = 0
	bin.PartNodes = make(map[*tree.PartitionNode]bool)

	return &bin
}

func (bin *Bin) AddToBinSize(pNode *tree.PartitionNode, tasks int64) {
	pNode.AddToNodeSize(tasks)

	bin.Size += tasks
}

func (bin *Bin) MakeMapRootNodesOfBin() map[*tree.PartitionNode]bool {
	var tmpPartNodes = make(map[*tree.PartitionNode]bool)
	for node := range bin.PartNodes {
		tmpPartNodes[node] = true
	}

	for node := range tmpPartNodes {
		if _, ok := bin.PartNodes[node.Parent]; ok {
			delete(tmpPartNodes, node)
		}
	}

	return tmpPartNodes
}

func (bin *Bin) MakeSliceRootNodesOfBin() []*tree.PartitionNode {
	var sliceRootNodesOfBin = make([]*tree.PartitionNode, 0)
	for rootNode := range bin.MakeMapRootNodesOfBin() {
		sliceRootNodesOfBin = append(sliceRootNodesOfBin, rootNode)
	}

	return sliceRootNodesOfBin
}

func (bin *Bin) AddSubTree(pNode *tree.PartitionNode) {
	bin.Size += pNode.SubTreeSize
	bin.addNodes(pNode)
}

func (bin *Bin) addNodes(pNode *tree.PartitionNode) {
	bin.PartNodes[pNode] = true

	for _, child := range pNode.Children {
		bin.addNodes(child)
	}
}

func (bin *Bin) freeBin() {
	bin.PartNodes = make(map[*tree.PartitionNode]bool) // garbage collector?
	bin.Size = 0
}

func PreprocessPartitionTree(pRoot *tree.PartitionNode) {
	if pRoot.NodeSize > Volume {

		rootChunk := tree.PartitionNode{Name: pRoot.Name + "#", Parent: pRoot, Children: pRoot.Children,
			NodeSize: pRoot.NodeSize - Volume, SubTreeSize: pRoot.SubTreeSize - Volume}

		pRoot.NodeSize = Volume
		pRoot.Children = nil
		pRoot.Children = append(pRoot.Children, &rootChunk)
	}

	for _, child := range pRoot.Children {
		PreprocessPartitionTree(child)
	}
}

func HierarchicalFirstFitDecreasing(pNode *tree.PartitionNode, bins []*Bin) []*Bin {
	if pNode.SubTreeSize <= Volume {
		var bin = findBinForFit(bins, pNode)

		if bin != nil {
			bin.AddSubTree(pNode)
		} else {
			bin := NewBin(len(bins) + 1)
			bin.AddSubTree(pNode)
			bins = append(bins, bin)
		}
	} else {
		//	var separate, forUnite = pNode.SeparateMaxChild()
		var separate, forUnite = pNode.SeparateRoot()

		for _, node := range separate {
			bins = HierarchicalFirstFitDecreasing(node, bins)
		}

		//	pNode.UniteMaxChild(forUnite)
		pNode.UniteRoot(forUnite)
	}

	return bins
}

func findBinForFit(bins []*Bin, pNode *tree.PartitionNode) *Bin {
	for _, bin := range bins {
		if pNode.SubTreeSize <= Volume-bin.Size {
			return bin
		}
	}

	return nil
}

func FindBinForRebalancing(bins []*Bin, tasksPerDay []map[string]int64,
	nameToPartNode map[string]*tree.PartitionNode) *Bin {

	var initiallyUnderloadedBins = make(map[*Bin]bool)
	for _, bin := range bins {
		initiallyUnderloadedBins[bin] = bin.Size <= UnderloadThreshold
	}

	var loadedBin *Bin
	for i := InitDays; i < len(tasksPerDay) && loadedBin == nil; i++ {

		updateSizeInOneTimeInterval(bins, tasksPerDay[i], nameToPartNode, true)           //Add
		updateSizeInOneTimeInterval(bins, tasksPerDay[i-InitDays], nameToPartNode, false) //Sub

		loadedBin = findOverOrUnderloadedBin(bins, initiallyUnderloadedBins)
	}

	return loadedBin
}

func DynamicalHierarchicalFirstFitDecreasing(loadedBin *Bin, bins []*Bin) ([]*Bin, int64) {
	var oldSize = loadedBin.Size
	var untiedChildren = untieChildNodesOfBinFromOtherBins(loadedBin)
	var sliceRootNodesOfBin = loadedBin.MakeSliceRootNodesOfBin()
	loadedBin.freeBin()

	sort.Slice(sliceRootNodesOfBin, func(i, j int) bool {
		return sliceRootNodesOfBin[i].SubTreeSize > sliceRootNodesOfBin[j].SubTreeSize
	})

	for _, rootNode := range sliceRootNodesOfBin {
		PreprocessPartitionTree(rootNode) // think later
		bins = HierarchicalFirstFitDecreasing(rootNode, bins)
	}

	tieChildNodesToOtherBins(untiedChildren)

	var migrationSize = oldSize - loadedBin.Size

	return bins, migrationSize
}

func untieChildNodesOfBinFromOtherBins(bin *Bin) map[*tree.PartitionNode][]*tree.PartitionNode {
	var untiedChildren = make(map[*tree.PartitionNode][]*tree.PartitionNode)

	for pNode := range bin.PartNodes {
		for _, child := range pNode.Children {
			if ok := bin.PartNodes[child]; !ok {
				pNode.RemoveChild(child)
				untiedChildren[pNode] = append(untiedChildren[pNode], child)
			}
		}
	}

	return untiedChildren
}

func tieChildNodesToOtherBins(untiedChildren map[*tree.PartitionNode][]*tree.PartitionNode) {
	for pNode, children := range untiedChildren {
		for _, child := range children {
			pNode.AppendChild(child)
		}
	}
}

func updateSizeInOneTimeInterval(bins []*Bin, tasksPerOneDay map[string]int64,
	nameToPartNode map[string]*tree.PartitionNode, isPlus bool) {

	for name, tasks := range tasksPerOneDay { //Add
		var pNode = nameToPartNode[name]

		if !isPlus {
			tasks = -tasks
		}

		for _, bin := range bins {
			if ok := bin.PartNodes[pNode]; ok {
				bin.AddToBinSize(pNode, tasks)
				break
			}
		}
	}
}

func findOverOrUnderloadedBin(bins []*Bin, initiallyUnderloadedBins map[*Bin]bool) *Bin {
	for _, bin := range bins {
		if (bin.Size >= OverloadThreshold) ||
			(bin.Size <= UnderloadThreshold && !initiallyUnderloadedBins[bin]) {
			return bin
		}
	}

	return nil
}
