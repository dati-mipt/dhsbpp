package tree

import (
	"dhsbpp/tenantsData"
)

type Node struct {
	Name     string
	Children []*Node
	NodeSize float64
	TreeSize float64
}


func  NewTree( TD *tenantsData.TenantsData ) ( *Node, bool ) {
	root, nodeName := buildTree( TD.TenantToParent )

	ok := validateTree( root, nodeName )
	if !ok {
		return nil, false
	}

	SetInitialNodeSize( TD.Tpd, nodeName, 30 )
	FindTreeSize( root )
	

	return root, true
}

func buildTree( tenantsToParent map[string]string ) ( *Node, map[string]*Node){
	nodeName := make(map[string]*Node)
	var root *Node

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

		if childNode == parentNode { // root case
			root = childNode
		} else {
			parentNode.Children = append(parentNode.Children, childNode)
		}
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

func SetInitialNodeSize( tpd []tenantsData.TasksPerDay,
	                 nodeName map[string]*Node, nInitDays int ) {

	nDays := 0
	currDay := tpd[0].Day
	for _, item := range tpd {
		if item.Day != currDay {
			currDay = item.Day
			nDays++
		}
		if nDays >= nInitDays {
			break
		}

		node := nodeName[item.Tenant]
		node.NodeSize += item.CreatedTasks
	}
}

func FindTreeSize( node *Node ) float64 {
	node.TreeSize += node.NodeSize
	for _, ch := range node.Children {
		node.TreeSize += FindTreeSize( ch )
	}

	return node.TreeSize
}

func ( node *Node ) Separate() []*Node {
	separate := make( []*Node, 0, len(node.Children) + 1)
	for _, ch := range node.Children {
		separate = append(separate, ch)
	}

	node.Children = nil
	node.TreeSize = node.NodeSize
	separate = append(separate, node)

	return separate
}

func ( node *Node ) Unite( children []*Node ) {
	node.Children = children

	for _, ch := range node.Children {
		node.TreeSize += ch.TreeSize
	}
}
