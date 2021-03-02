package vizualize

import (
	"fmt"
	"github.com/dati-mipt/dhsbpp/packing"
	"github.com/dati-mipt/dhsbpp/tree"
	"os"
	"sort"
)

type DotNode struct {
	BinIndex int

	Children []*DotNode
	Parent   *DotNode

	NumOfNodes int64
	pRoot      *tree.PartitionNode
	pNodes     map[*tree.PartitionNode]bool
}

func NewDotTree(bins []*packing.Bin) *DotNode {
	var unconnectedDotNodes = make(map[*DotNode]bool)
	for idx := range bins {
		var rootNodes = bins[idx].MakeRootNodesOfBin()

		for rootNode := range rootNodes {
			var pNodes = findPartitionNodesOfDotNode(bins[idx].PartNodes, rootNode)
			var dotNode = DotNode{BinIndex: bins[idx].Index, pRoot: rootNode, pNodes: pNodes}
			unconnectedDotNodes[&dotNode] = true
		}
	}

	var rootDotNode = connectBinNodes(unconnectedDotNodes)

	return rootDotNode
}

func connectBinNodes(unconnectedDotNodes map[*DotNode]bool) *DotNode {
	var rootDotNode *DotNode

	for dotNode := range unconnectedDotNodes {
		var parent = findParent(dotNode, unconnectedDotNodes)

		if parent != nil {
			dotNode.Parent = parent
			parent.Children = append(parent.Children, dotNode)
		} else {
			rootDotNode = dotNode
		}
	}

	return rootDotNode
}

func findParent(child *DotNode, dotNodes map[*DotNode]bool) *DotNode {
	for dotNode := range dotNodes {
		if ok := dotNode.pNodes[child.pRoot.Parent]; ok {
			return dotNode
		}
	}

	return nil
}

func findPartitionNodesOfDotNode(allPartNodes map[*tree.PartitionNode]bool,
	pRoot *tree.PartitionNode) map[*tree.PartitionNode]bool {

	var pNodesOfDotNode = make(map[*tree.PartitionNode]bool)

	findPartitionNodesOfDotNodeFunc(pRoot, allPartNodes, pNodesOfDotNode)

	return pNodesOfDotNode
}

func findPartitionNodesOfDotNodeFunc(pNode *tree.PartitionNode, allPartNodes map[*tree.PartitionNode]bool,
	pNodesOfDotNode map[*tree.PartitionNode]bool) {
	if ok, _ := allPartNodes[pNode]; !ok {
		return
	}

	pNodesOfDotNode[pNode] = true

	for _, children := range pNode.Children {
		findPartitionNodesOfDotNodeFunc(children, allPartNodes, pNodesOfDotNode)
	}
}

func DrawBinTreeDot(binTree *DotNode, dotFileName string) error {
	var file, err = os.Create(dotFileName)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file,
		"digraph G                                                          \n"+
			"{                                                                     \n"+
			"  node[shape = \"box\", color=\"black\",fontsize=14,      	           \n"+
			"  style=\"filled\"];                                                  \n"+
			"  edge[color=\"black\",style=\"bold\"];\n")

	if err != nil {
		return err
	}

	err = drawNodeDot(binTree, file)
	if err != nil {
		return err
	}

	if _, err = fmt.Fprintf(file, "}"); err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}

	return nil
}

var m = map[int]string{
	1: "cadetblue1", 2: "darkorchid1", 3: "gold1", 4: "olivedrab2", 5: "violetred1", 6: "antiquewhite1",
	7: "coral", 8: "gray40", 9: "plum1", 10: "aquamarine4",
}

func drawNodeDot(binNode *DotNode, file *os.File) error {
	var _, err = fmt.Fprintf(file, "  ptr%p[label=\""+
		"bin index = %v\n"+
		"number of nodes = %v\n"+
		"size = %v\","+
		"fillcolor=\"%v\"];\n", binNode, binNode.BinIndex, len(binNode.pNodes),
		calculateSizeOfBinNode(binNode), m[binNode.BinIndex])

	if err != nil {
		return err
	}

	sort.Slice(binNode.Children, func(i, j int) bool {
		return len(binNode.Children[i].pNodes) < len(binNode.Children[j].pNodes)
	})

	for _, child := range binNode.Children {
		if _, err = fmt.Fprintf(file, "  ptr%p->ptr%p;\n", binNode, child); err != nil {
			return err
		}
		err = drawNodeDot(child, file)
	}

	return nil
}

func calculateSizeOfBinNode(binNode *DotNode) int64 {
	var size int64 = 0
	for pNode := range binNode.pNodes {
		size += pNode.NodeSize
	}
	return size
}
