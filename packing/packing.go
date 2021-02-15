package packing

import (
	"dhsbpp/tree"
	"sort"
)

type Bin struct {
	Tenants  []*tree.Node
	MAX float64
	AF  float64
	RD  float64

	FreeSpace float64
	V        float64
}

func  NewBin( MAX float64, AF float64) *Bin {
	bin := &Bin{}

	bin.MAX = MAX
	bin.AF = AF

	bin.V = bin.AF * bin.MAX
	bin.FreeSpace  = bin.V

	return bin
}

func PreprocessTree( root *tree.Node, MAX float64, AF float64 ) {
	if root.NodeSize > MAX*AF {

		rootChunk := &tree.Node{ Name: root.Name+"#", Children: root.Children,
								 NodeSize: root.NodeSize - MAX*AF, TreeSize: root.TreeSize - MAX*AF }

		root.NodeSize = MAX*AF
		root.Children = nil
		root.Children = append( root.Children, rootChunk )

	}

	for _, ch := range root.Children {
		PreprocessTree( ch, MAX, AF )
	}
}

func ( bin *Bin) AddSubTree( root *tree.Node ) {
	bin.FreeSpace -= root.TreeSize
	bin.AddNodes( root )
}

func ( bin *Bin) AddNodes( node *tree.Node ) {
	bin.Tenants = append(bin.Tenants, node)

	for _, ch := range node.Children {
		bin.AddNodes( ch )
	}
}


func HFFD( root *tree.Node, bins *[]*Bin,
	       MAX float64, AF float64) {
	V := AF * MAX
	if root.TreeSize <= V  {
		isFit := false
		for _, bin := range *bins {
			if root.TreeSize <= bin.FreeSpace {
				bin.AddSubTree(root)
				isFit = true
				break
			}
		}

		if !isFit {
			bin := NewBin( MAX, AF)
			bin.AddSubTree(root)
			*bins = append( *bins, bin)
		}
	} else {
		children := root.Children

		separate := root.Separate()
		sort.Slice(separate, func ( i, j int) bool {
			        return separate[i].TreeSize > separate[j].TreeSize} )

		for _, node := range separate {
			HFFD(node, bins, MAX, AF)
		}
		root.Unite( children )
	}
}
