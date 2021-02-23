package tree

import (
	"errors"
	"fmt"
	"github.com/dati-mipt/dhsbpp/hierarchy"
)

type Node struct {
	Name     string
	Children []*Node
	NodeSize float64
	TreeSize float64
}

func NewTree(hierarchy *hierarchy.Hierarchy) (*Node, error) {
	root, nodeName := buildTree(hierarchy.TenantToParent)
	ok := validateTree(root, nodeName)
	if !ok {
		return nil, errors.New("tree : tree is not valid")
	}

	const initDays = 30

	SetInitialNodeSize(hierarchy.TasksPerDay, nodeName, initDays)
	FindTreeSize(root)

	return root, nil
}

func buildTree(tenantsToParent map[string]string) (*Node, map[string]*Node) {
	var nodeName = make(map[string]*Node)
	var root *Node
	var i = 0
	for key, value := range tenantsToParent {
		childName, parentName := key, value

		parentNode, ok := nodeName[parentName]
		if !ok {
			parentNode = &Node{Name: parentName}
			nodeName[parentName] = parentNode
		}
		childNode, ok := nodeName[childName]
		if !ok {
			childNode = &Node{Name: childName}
			nodeName[childName] = childNode
		}

		if childName == parentName { // root case
			root = childNode
		} else {
			parentNode.Children = append(parentNode.Children, childNode)
		}
		i++
	}

	return root, nodeName
}

func validateTree(root *Node, nodeName map[string]*Node) bool {
	visited := make(map[*Node]bool)
	for _, node := range nodeName {
		visited[node] = false
	}

	acyclic := isAcyclicDFS(root, visited)

	allVisited := true
	for _, v := range visited {
		if v == false {
			allVisited = false
		}
	}
	fmt.Println(acyclic, allVisited)

	return acyclic && allVisited
}

func isAcyclicDFS(node *Node, visited map[*Node]bool) bool {
	if visited[node] == true {
		fmt.Println(node.Name)
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

func SetInitialNodeSize(tasksPerDay []map[string]float64,
	nodeName map[string]*Node, initDays int) {

	for i, dayChange := range tasksPerDay {
		if i >= initDays {
			break
		}

		for name, tasks := range dayChange {
			var node = nodeName[name]
			node.NodeSize += tasks
		}
	}
}

func FindTreeSize(node *Node) float64 {
	node.TreeSize += node.NodeSize
	for _, ch := range node.Children {
		node.TreeSize += FindTreeSize(ch)
	}

	return node.TreeSize
}

func (node *Node) Separate() []*Node {
	separate := make([]*Node, 0, len(node.Children)+1)
	for _, ch := range node.Children {
		separate = append(separate, ch)
	}

	node.Children = nil
	node.TreeSize = node.NodeSize
	separate = append(separate, node)

	return separate
}

func (node *Node) Unite(children []*Node) {
	node.Children = children

	for _, ch := range node.Children {
		node.TreeSize += ch.TreeSize
	}
}
