package partitionTree

import (
	"errors"
	"github.com/dati-mipt/dhsbpp/tree"
)

type PartitionNode struct {
	Name     string
	Parent   *PartitionNode
	Children []*PartitionNode

	NodeSize    int64
	SubTreeSize int64
}

func NewPartitionTree(root *tree.Node) *PartitionNode {
	var pNode = copyTree(root, nil)

	return pNode
}

func copyTree(root *tree.Node, parent *PartitionNode) *PartitionNode {
	var pNode = &PartitionNode{}
	pNode.Name = root.Name
	pNode.Parent = parent

	pNode.Children = make([]*PartitionNode, 0)
	for _, child := range root.Children {
		pNode.Children = append(pNode.Children, copyTree(child, pNode))
	}

	return pNode
}

func (pNode *PartitionNode) isRoot() bool {
	return pNode.Parent == nil
}

// return map : name -> *PartitionNode
func (pNode *PartitionNode) MapNameToPartitionNode() (map[string]*PartitionNode, error) {
	if !pNode.isRoot() {
		return nil, errors.New("partition tree : need partition root")
	}

	var nameToPartNode = make(map[string]*PartitionNode)
	nameToPartNode = pNode.mapNameToPartitionNodeFunc(nameToPartNode)

	return nameToPartNode, nil
}

func (pNode *PartitionNode) mapNameToPartitionNodeFunc(
	nameToPartNode map[string]*PartitionNode) map[string]*PartitionNode {

	nameToPartNode[pNode.Name] = pNode

	for _, child := range pNode.Children {
		nameToPartNode = child.mapNameToPartitionNodeFunc(nameToPartNode)
	}

	return nameToPartNode
}

func (pNode *PartitionNode) AddToNodeSize(tasks int64) { // tasks < 0 allowed
	pNode.NodeSize += tasks

	for ptr := pNode; ptr != nil; ptr = ptr.Parent {
		ptr.SubTreeSize += tasks
	}
}

func (pNode *PartitionNode) RemoveChild(child *PartitionNode) {
	for ptr := pNode; ptr != nil; ptr = ptr.Parent {
		ptr.SubTreeSize -= child.SubTreeSize
	}

	for idx := range pNode.Children {
		if pNode.Children[idx] == child {
			pNode.Children = append(pNode.Children[:idx], pNode.Children[idx+1:]...)
			break
		}
	}
}

func (pNode *PartitionNode) AppendChild(child *PartitionNode) {
	for ptr := pNode; ptr != nil; ptr = ptr.Parent {
		ptr.SubTreeSize += child.SubTreeSize
	}

	pNode.Children = append(pNode.Children, child) // order??
}

func (pNode *PartitionNode) SetInitialSize(tasksPerDay []map[string]int64, initDays int) error {
	if !pNode.isRoot() {
		return errors.New("partition tree : need partition root")
	}
	if pNode.SubTreeSize > 0 {
		return errors.New("partition tree : tree has already initial size ") // or reset to zero?
	}
	var nameToPartNode, err = pNode.MapNameToPartitionNode()
	if err != nil {
		return err
	}

	for i := 0; i < initDays; i++ {
		for name, tasks := range tasksPerDay[i] {
			nameToPartNode[name].AddToNodeSize(tasks)
		}
	}

	return nil
}

func (pNode *PartitionNode) Separate() []*PartitionNode {
	separate := make([]*PartitionNode, 0, len(pNode.Children)+1)
	for _, ch := range pNode.Children {
		separate = append(separate, ch)
	}

	pNode.Children = nil
	pNode.SubTreeSize = pNode.NodeSize
	separate = append(separate, pNode)

	return separate
}

func (pNode *PartitionNode) Unite(children []*PartitionNode) {
	pNode.Children = children

	for _, ch := range pNode.Children {
		pNode.SubTreeSize += ch.SubTreeSize
	}
}
