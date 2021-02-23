package packing

import (
	"github.com/dati-mipt/dhsbpp/tree"
	"sort"
)

type Bin struct {
	Tenants           []*tree.Node
	MaxCapacity       float64
	AllocationFactor  float64
	ReallocationDelta float64

	FreeSpace float64
	Volume    float64
}

func NewBin(maxCapacity float64, allocationFactor float64) *Bin {
	var bin Bin

	bin.MaxCapacity = maxCapacity
	bin.AllocationFactor = allocationFactor

	bin.Volume = bin.AllocationFactor * bin.MaxCapacity
	bin.FreeSpace = bin.Volume

	return &bin
}

func PreprocessTree(root *tree.Node, maxCapacity float64, allocationFactor float64) {
	if root.NodeSize > maxCapacity*allocationFactor {

		rootChunk := tree.Node{Name: root.Name + "#", Children: root.Children,
			NodeSize: root.NodeSize - maxCapacity*allocationFactor, TreeSize: root.TreeSize - maxCapacity*allocationFactor}

		root.NodeSize = maxCapacity * allocationFactor
		root.Children = nil
		root.Children = append(root.Children, &rootChunk)

	}

	for _, ch := range root.Children {
		PreprocessTree(ch, maxCapacity, allocationFactor)
	}
}

func (bin *Bin) AddSubTree(root *tree.Node) {
	bin.FreeSpace -= root.TreeSize
	bin.AddNodes(root)
}

func (bin *Bin) AddNodes(node *tree.Node) {
	bin.Tenants = append(bin.Tenants, node)

	for _, ch := range node.Children {
		bin.AddNodes(ch)
	}
}

func HierarchicalFirstFitDecreasing(root *tree.Node, bins *[]*Bin,
	maxCapacity float64, allocationFactor float64) {
	V := allocationFactor * maxCapacity
	if root.TreeSize <= V {
		isFit := false
		for _, bin := range *bins {
			if root.TreeSize <= bin.FreeSpace {
				bin.AddSubTree(root)
				isFit = true
				break
			}
		}

		if !isFit {
			bin := NewBin(maxCapacity, allocationFactor)
			bin.AddSubTree(root)
			*bins = append(*bins, bin)
		}
	} else {
		children := root.Children

		separate := root.Separate()
		sort.Slice(separate, func(i, j int) bool {
			return separate[i].TreeSize > separate[j].TreeSize
		})

		for _, node := range separate {
			HierarchicalFirstFitDecreasing(node, bins, maxCapacity, allocationFactor)
		}
		root.Unite(children)
	}
}
