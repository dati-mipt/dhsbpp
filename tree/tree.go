package tree

import (
	"errors"
	"fmt"
	"sort"
)

type Node struct {
	Name     string // name string
	Parent   *Node
	Children []*Node
}

func NewTree(childToParent map[string]string) (*Node, error) {
	var root = buildTree(childToParent)
	var ok = ValidateTree(root)
	if !ok {
		return nil, errors.New("tree : tree is not valid")
	}

	return root, nil
}

func buildTree(tenantsToParent map[string]string) *Node {
	var nameToNote = make(map[string]*Node)
	var root *Node
	for key, value := range tenantsToParent {
		childName, parentName := key, value

		_, ok := nameToNote[parentName]
		if !ok {
			nameToNote[parentName] = &Node{Name: parentName}
		}
		_, ok = nameToNote[childName]
		if !ok {
			nameToNote[childName] = &Node{Name: childName}
		}

		if childName == parentName { // root case
			root = nameToNote[childName]
		} else {
			nameToNote[childName].Parent = nameToNote[parentName]

			nameToNote[parentName].Children = append(nameToNote[parentName].Children,
				nameToNote[childName])
		}
	}

	return root
}

func ValidateTree(root *Node) bool {
	visited := make(map[*Node]bool)
	var allNodes, err = root.AllNodes()
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, node := range allNodes {
		visited[node] = false
	}

	acyclic := isAcyclicDFS(root, visited)

	allVisited := true
	for _, v := range visited {
		if v == false {
			allVisited = false
		}
	}

	return acyclic && allVisited
}

func isAcyclicDFS(node *Node, visited map[*Node]bool) bool {
	if visited[node] == true {
		return false
	}

	visited[node] = true

	for _, child := range node.Children {
		if isAcyclicDFS(child, visited) == false {
			return false
		}
	}

	return true
}

// return slice of all nodes
func (node *Node) AllNodes() ([]*Node, error) { // []Node or []*Node?
	if !node.isRoot() {
		return nil, errors.New("tree : need root")
	}

	var allNodes []*Node
	allNodes = node.allNodesFunc(allNodes)

	return allNodes, nil
}

func (node *Node) allNodesFunc(allNodes []*Node) []*Node {
	allNodes = append(allNodes, node)

	for _, child := range node.Children {
		allNodes = child.allNodesFunc(allNodes)
	}

	return allNodes
}

func (node *Node) isRoot() bool {
	return node.Parent == nil
}

//---------------------------Partition Tree----------------------

type PartitionNode struct {
	Name     string
	Parent   *PartitionNode
	Children []*PartitionNode

	NodeSize    int64
	SubTreeSize int64
}

func NewPartitionTree(root *Node) *PartitionNode {
	var pNode = copyTree(root, nil)

	return pNode
}

func copyTree(root *Node, parent *PartitionNode) *PartitionNode {
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
	pNode.mapNameToPartitionNodeFunc(nameToPartNode)

	return nameToPartNode, nil
}

func (pNode *PartitionNode) mapNameToPartitionNodeFunc(nameToPartNode map[string]*PartitionNode) {
	nameToPartNode[pNode.Name] = pNode

	for _, child := range pNode.Children {
		child.mapNameToPartitionNodeFunc(nameToPartNode)
	}
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

	pNode.Children = append(pNode.Children, child)
}

func (pNode *PartitionNode) SetInitialSize(tasksPerDay []map[string]int64, initDays int) error {
	if !pNode.isRoot() {
		return errors.New("partition tree : need partition root")
	}
	if pNode.SubTreeSize > 0 { // or reset to zero?
		return errors.New("partition tree : tree has already initial size ")
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

	sort.Slice(separate, func(i, j int) bool {
		return separate[i].SubTreeSize > separate[j].SubTreeSize
	})

	return separate
}

func (pNode *PartitionNode) Unite(children []*PartitionNode) {
	pNode.Children = children

	for _, ch := range pNode.Children {
		pNode.SubTreeSize += ch.SubTreeSize
	}
}
