package tree

import (
	"errors"
	"fmt"
)

type Node struct {
	Name     string
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
