package packing

import (
	"github.com/dati-mipt/dhsbpp/partitionTree"
	"sort"
)

type Bin struct {
	Index             int
	Tenants           map[*partitionTree.PartitionNode]bool
	MaxCapacity       int64
	AllocationFactor  int64 // in %s (0.6 -> 60)
	ReallocationDelta int64

	Size   int64
	Volume int64
}

func NewBin(idx int, maxCapacity int64, allocationFactor int64, reallocationDelta int64) *Bin {
	var bin Bin

	bin.Index = idx
	bin.Tenants = make(map[*partitionTree.PartitionNode]bool)
	bin.MaxCapacity = maxCapacity
	bin.AllocationFactor = allocationFactor
	bin.ReallocationDelta = reallocationDelta

	bin.Volume = bin.MaxCapacity * allocationFactor / 100
	bin.Size = 0

	return &bin
}

func (bin *Bin) AddToBinSize(pNode *partitionTree.PartitionNode, tasks int64) {
	pNode.AddToNodeSize(tasks)

	bin.Size += tasks
}

func (bin *Bin) MakeRootNodesOfBin() map[*partitionTree.PartitionNode]bool {
	var tmpTenants = make(map[*partitionTree.PartitionNode]bool)
	for node := range bin.Tenants {
		tmpTenants[node] = true
	}

	for node := range tmpTenants {
		if _, ok := bin.Tenants[node.Parent]; ok {
			delete(tmpTenants, node)
		}
	}

	return tmpTenants
}

func PreprocessPartitionTree(pRoot *partitionTree.PartitionNode, maxCapacity int64, allocationFactor int64) {
	var Volume = maxCapacity * allocationFactor / 100
	if pRoot.NodeSize > Volume {

		rootChunk := partitionTree.PartitionNode{Name: pRoot.Name + "#", Parent: pRoot, Children: pRoot.Children,
			NodeSize: pRoot.NodeSize - Volume, SubTreeSize: pRoot.SubTreeSize - Volume}

		pRoot.NodeSize = Volume
		pRoot.Children = nil
		pRoot.Children = append(pRoot.Children, &rootChunk)
	}

	for _, ch := range pRoot.Children {
		PreprocessPartitionTree(ch, maxCapacity, allocationFactor)
	}
}

func (bin *Bin) AddSubTree(pNode *partitionTree.PartitionNode) {
	bin.Size += pNode.SubTreeSize
	bin.addNodes(pNode)
}

func (bin *Bin) addNodes(pNode *partitionTree.PartitionNode) {
	bin.Tenants[pNode] = true

	for _, child := range pNode.Children {
		bin.addNodes(child)
	}
}

func (bin *Bin) freeBin() {
	bin.Tenants = nil // garbage collector?
	bin.Tenants = make(map[*partitionTree.PartitionNode]bool)
	bin.Size = 0
}

func HierarchicalFirstFitDecreasing(pRoot *partitionTree.PartitionNode, bins *[]*Bin,
	maxCapacity int64, allocationFactor int64, reallocationDelta int64) {
	Volume := maxCapacity * allocationFactor / 100
	if pRoot.SubTreeSize <= Volume {
		isFit := false
		for _, bin := range *bins {
			if pRoot.SubTreeSize <= bin.Volume-bin.Size {
				bin.AddSubTree(pRoot)
				isFit = true
				break
			}
		}

		if !isFit {
			bin := NewBin(len(*bins)+1, maxCapacity, allocationFactor, reallocationDelta)
			bin.AddSubTree(pRoot)
			*bins = append(*bins, bin)
		}
	} else {
		children := pRoot.Children

		separate := pRoot.Separate()
		sort.Slice(separate, func(i, j int) bool {
			return separate[i].SubTreeSize > separate[j].SubTreeSize
		})

		for _, node := range separate {
			HierarchicalFirstFitDecreasing(node, bins, maxCapacity, allocationFactor, reallocationDelta)
		}
		pRoot.Unite(children)
	}
}

func DynamicalHierarchicalFirstFitDecreasing(bins *[]*Bin, tasksPerDay []map[string]int64,
	nameToPartNode map[string]*partitionTree.PartitionNode, initDays int,
	maxCapacity int64, allocationFactor int64, reallocationDelta int64) { // global constants?

	var loadedBin *Bin
	for i := initDays; i < len(tasksPerDay) && loadedBin == nil; i++ {
		updateSizeInOneTimeInterval(*bins, tasksPerDay[i], nameToPartNode, true)           //Add
		updateSizeInOneTimeInterval(*bins, tasksPerDay[i-initDays], nameToPartNode, false) //Sub

		loadedBin = findOverOrUnderloadedBin(*bins)
	}

	if loadedBin != nil {
		var untiedChildren = untieChildNodesOfBinFromOtherBins(loadedBin)
		var rootNodesOfBin = loadedBin.MakeRootNodesOfBin()
		loadedBin.freeBin()

		for rootNode := range rootNodesOfBin {
			HierarchicalFirstFitDecreasing(rootNode, bins, maxCapacity, allocationFactor, reallocationDelta)
		}

		tieChildNodesToOtherBins(untiedChildren)
	}
}

func untieChildNodesOfBinFromOtherBins(bin *Bin) map[*partitionTree.PartitionNode][]*partitionTree.PartitionNode {
	var untiedChildren = make(map[*partitionTree.PartitionNode][]*partitionTree.PartitionNode)

	for pNode := range bin.Tenants {
		for _, child := range pNode.Children {
			if ok := bin.Tenants[child]; !ok {
				pNode.RemoveChild(child)
				untiedChildren[pNode] = append(untiedChildren[pNode], child)
			}
		}
	}

	return untiedChildren
}

func tieChildNodesToOtherBins(untiedChildren map[*partitionTree.PartitionNode][]*partitionTree.PartitionNode) {
	for pNode, children := range untiedChildren {
		for _, child := range children {
			pNode.AppendChild(child)
		}
	}
}

func updateSizeInOneTimeInterval(bins []*Bin, tasksPerOneDay map[string]int64,
	nameToPartNode map[string]*partitionTree.PartitionNode, isPlus bool) {

	for name, tasks := range tasksPerOneDay { //Add
		var pNode = nameToPartNode[name]

		if !isPlus {
			tasks = -tasks
		}

		for _, bin := range bins {
			if ok := bin.Tenants[pNode]; ok {
				bin.AddToBinSize(pNode, tasks)
				break
			}
		}
	}
}

func findOverOrUnderloadedBin(bins []*Bin) *Bin {
	for _, bin := range bins {
		var overloadedThreshold = bin.MaxCapacity * (bin.AllocationFactor + bin.ReallocationDelta) / 100
		var underloadedThreshold = bin.MaxCapacity * (bin.AllocationFactor - bin.ReallocationDelta) / 100
		if (bin.Size >= overloadedThreshold) || (bin.Size <= underloadedThreshold) {
			return bin
		}
	}

	return nil
}
