package packing

import (
	"github.com/dati-mipt/dhsbpp/tree"
)

const (
	MaxCapacity       = 250
	AllocationFactor  = 60
	ReallocationDelta = 20
	InitDays          = 100
	Volume            = MaxCapacity * AllocationFactor / 100
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

func (bin *Bin) MakeRootNodesOfBin() map[*tree.PartitionNode]bool {
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

func HierarchicalFirstFitDecreasing(pNode *tree.PartitionNode, bins *[]*Bin) {
	if pNode.SubTreeSize <= Volume {
		var bin = findBinForFit(*bins, pNode)

		if bin != nil {
			bin.AddSubTree(pNode)
		} else {
			bin := NewBin(len(*bins) + 1)
			bin.AddSubTree(pNode)
			*bins = append(*bins, bin)
		}
	} else {
		var children = pNode.Children
		var separate = pNode.Separate()

		for _, node := range separate {
			HierarchicalFirstFitDecreasing(node, bins)
		}

		pNode.Unite(children)
	}
}

func findBinForFit(bins []*Bin, pNode *tree.PartitionNode) *Bin {
	for _, bin := range bins {
		if pNode.SubTreeSize <= Volume-bin.Size {
			return bin
		}
	}

	return nil
}

func DynamicalHierarchicalFirstFitDecreasing(bins *[]*Bin, tasksPerDay []map[string]int64,
	nameToPartNode map[string]*tree.PartitionNode) { // global constants?

	var loadedBin *Bin
	for i := InitDays; i < len(tasksPerDay) && loadedBin == nil; i++ {
		updateSizeInOneTimeInterval(*bins, tasksPerDay[i], nameToPartNode, true)           //Add
		updateSizeInOneTimeInterval(*bins, tasksPerDay[i-InitDays], nameToPartNode, false) //Sub

		loadedBin = findOverOrUnderloadedBin(*bins)
	}

	if loadedBin != nil {
		var untiedChildren = untieChildNodesOfBinFromOtherBins(loadedBin)
		var rootNodesOfBin = loadedBin.MakeRootNodesOfBin()
		loadedBin.freeBin()

		for rootNode := range rootNodesOfBin {
			HierarchicalFirstFitDecreasing(rootNode, bins)
		}

		tieChildNodesToOtherBins(untiedChildren)
	}
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

func findOverOrUnderloadedBin(bins []*Bin) *Bin {
	for _, bin := range bins {
		var overloadedThreshold int64 = MaxCapacity * (AllocationFactor + ReallocationDelta) / 100
		var underloadedThreshold int64 = MaxCapacity * (AllocationFactor - ReallocationDelta) / 100
		if (bin.Size >= overloadedThreshold) || (bin.Size <= underloadedThreshold) {
			return bin
		}
	}

	return nil
}
