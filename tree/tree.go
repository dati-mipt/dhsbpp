package tree

import (
	"DHSBPP/tenantsData"
)

type Node struct {
	name string
	children []*Node
	nodeSize int
}

func  NewTree( tenantsToParent map[string]string ) ( *Node, map[string]*Node, bool ) {
	root, nodeName := buildTree( tenantsToParent )

	ok := validateTree( root, nodeName )

	if ok {
		return root, nodeName, true
	} else {
		return nil, nil, false
	}
}

func buildTree( tenantsToParent map[string]string ) ( *Node, map[string]*Node){
	nodeName := make(map[string]*Node)
	var root *Node

	for key, value := range tenantsToParent {
		childName, parentName := key, value

		parentNode, ok := nodeName[parentName]
		if !ok {
			parentNode = &Node{name: parentName}
			nodeName[parentName] = parentNode
		}
		childNode, ok := nodeName[childName]
		if !ok {
			childNode = &Node{name: childName}
			nodeName[childName] = childNode
		}

		if childNode == parentNode { // root case
			root = childNode
		} else {
			parentNode.children = append(parentNode.children, childNode)
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

	for _, child := range node.children {
		if isAcyclicDFS(child, visited) == false {
			return false
		}
	}

	return true
}

func SetInitialSize( tpd []tenantsData.TasksPerDay,
	                 nodeName map[string]*Node, n_init_days int ) {

	n_days := 0
	curr_day := tpd[0].Day
	for _, item := range tpd {
		if item.Day != curr_day {
			curr_day = item.Day
			n_days++
		}
		if n_days >= n_init_days {
			break
		}

		node := nodeName[item.Tenant]
		node.nodeSize += item.CreatedTasks
	}
}