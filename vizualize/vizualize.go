package vizualize

import (
	"fmt"
	"github.com/dati-mipt/dhsbpp/packing"
	"github.com/dati-mipt/dhsbpp/partitionTree"
	"os"
	"sort"
)

type BinNode struct {
	idx int

	Children []*BinNode
	Parent   *BinNode

	NumOfNodes int64
	pRoot      *partitionTree.PartitionNode
	pNodes     map[*partitionTree.PartitionNode]bool
}

func NewTreeOfBins(bins []*packing.Bin) *BinNode {
	var unconnectedBinNodes = make(map[*BinNode]bool)
	for idx := range bins {
		var rootNodes = bins[idx].MakeRootNodesOfBin()

		for rootNode := range rootNodes {
			var pNodes = findPartitionNodes(rootNode, bins[idx].Tenants)
			var binNode = BinNode{idx: bins[idx].Index, pRoot: rootNode, pNodes: pNodes}
			unconnectedBinNodes[&binNode] = true
		}
	}

	var rootBinNode = connectBinNodes(unconnectedBinNodes)
	//	var lilRootBinNode = reduceNumOfNodes(rootBinNode)

	return rootBinNode
}

//func reduceNumOfNodes(bigBinNode *BinNode) *BinNode {

//}

func connectBinNodes(unconnectedBinNodes map[*BinNode]bool) *BinNode {
	var rootBinNode *BinNode

	for binNode := range unconnectedBinNodes {
		var parent = findParent(binNode, unconnectedBinNodes)

		if parent != nil {
			binNode.Parent = parent
			parent.Children = append(parent.Children, binNode)
		} else {
			rootBinNode = binNode
		}
	}

	return rootBinNode
}

func findParent(child *BinNode, binNodes map[*BinNode]bool) *BinNode {
	for binNode := range binNodes {
		if ok := binNode.pNodes[child.pRoot.Parent]; ok {
			return binNode
		}
	}

	return nil
}

func findPartitionNodes(pRoot *partitionTree.PartitionNode,
	tenants map[*partitionTree.PartitionNode]bool) map[*partitionTree.PartitionNode]bool {

	var pNodes = make(map[*partitionTree.PartitionNode]bool)

	pNodes = findPartitionNodesFunc(pRoot, tenants, pNodes)

	return pNodes
}

func findPartitionNodesFunc(pNode *partitionTree.PartitionNode, tenants map[*partitionTree.PartitionNode]bool,
	pNodes map[*partitionTree.PartitionNode]bool) map[*partitionTree.PartitionNode]bool {
	if ok, _ := tenants[pNode]; !ok {
		return pNodes
	}

	pNodes[pNode] = true

	for _, children := range pNode.Children {
		pNodes = findPartitionNodesFunc(children, tenants, pNodes)
	}

	return pNodes
}

func DrawBinTreeDot(binTree *BinNode, dotFileName string) error {
	file, err := os.Create(dotFileName)

	if err != nil {
		return err
	}

	defer file.Close()

	fmt.Fprintf(file,
		"digraph G                                                          \n"+
			"{                                                                     \n"+
			"  node[shape = \"box\", color=\"black\",fontsize=14,      	           \n"+
			"  style=\"filled\"];                                                  \n"+
			"  edge[color=\"black\",style=\"bold\"];\n")

	drawNodeDot(binTree, file)

	fmt.Fprintf(file, "}")

	return nil
}

var m = map[int]string{
	1: "cadetblue1",
	2: "darkorchid1",
	3: "gold1",
	4: "olivedrab2",
	5: "violetred1",
	6: "black",
}

func drawNodeDot(binNode *BinNode, file *os.File) {
	_, err := fmt.Fprintf(file, "  ptr%p[label=\"bin index = %v\nnumber of nodes = %v\nsize = %v\",fillcolor=\"%v\"];\n", binNode, binNode.idx,
		len(binNode.pNodes), calculateSizeOfBinNode(binNode), m[binNode.idx])
	if err != nil {
		fmt.Println(err)
	}

	sort.Slice(binNode.Children, func(i, j int) bool {
		return len(binNode.Children[i].pNodes) < len(binNode.Children[j].pNodes)
	})
	for _, child := range binNode.Children {
		fmt.Fprintf(file, "  ptr%p->ptr%p;\n", binNode, child)
		drawNodeDot(child, file)
	}
}

func calculateSizeOfBinNode(binNode *BinNode) int64 {
	var size int64 = 0
	for pNode := range binNode.pNodes {
		size += pNode.NodeSize
	}
	return size
}
