package simplerouter

import (
	"fmt"
	"strings"
)

type node struct {
	children     []*node
	component    string
	isNamedParam bool
	methods      map[string]Handle
}

func checkNamedParam(component string) (string, bool) {
	if component[0] == ':' {
		return component[1:], true
	}
	return component, false
}

// addNode will add a new path into the tree. If there are multiple nodes that are not already in
// the tree, then it will add in those multiple nodes
func (n *node) addNode(method string, path string, handle Handle) {
	var component string
	var isNamedParam bool
	components := strings.Split(path, "/")[1:]
	size := len(components)
	aNode, depth := n.traverseTree(components)
	remaining := size - depth
	httpMethod := strings.ToUpper(method)

	// We found a match that was already in the tree
	if remaining == 0 {
		aNode.methods[httpMethod] = handle
	}

	// There are still nodes to be inserted into the tree
	for i := 0; i < remaining; i++ {
		component, isNamedParam = checkNamedParam(components[depth+i])
		newNode := node{component: component, isNamedParam: isNamedParam, methods: make(map[string]Handle)}
		aNode.children = append(aNode.children, &newNode)

		// Add the handler to the last node
		if i == remaining-1 {
			newNode.methods[httpMethod] = handle
		} else { // Continue adding to children when not the last node
			aNode = &newNode
		}
	}
}

// traverseTree will take in a list of components to traverse the tree with. It will stop at the
// last matching component in the list and return that node along with that name.
func (n *node) traverseTree(components []string) (*node, int) {
	return n.traverse(components, 0)
}

// traverse is called recursively to find the last node in the tree that matches the components
// sequence
func (n *node) traverse(components []string, depth int) (*node, int) {
	component := components[0]
	// Traverse only when there are children left to traverse
	if len(n.children) > 0 {
		// Iterate through all the children
		for _, child := range n.children {
			// Continue to traverse the tree now that we've found a match
			if child.component == component {
				if len(components) == 1 {
					return child, depth + 1
				}
				return child.traverse(components[1:], depth+1)
			}
		}
	}
	return n, depth
}

// searchPath will recursively search the tree to see if the components are within the array
// TODO: Need to figure a way to return the path params
func (n *node) searchPath(components []string) (*node, Params) {
	component := components[0]

	if len(n.children) > 0 {
		for _, child := range n.children {
			if child.component == component || child.isNamedParam {
				// Found a match in the path
				if len(components) == 1 {
					params := make(Params, 0, 5)
					if child.isNamedParam {
						params = append(params, Param{child.component, component})
					}
					return child, params
				}

				res, params := child.searchPath(components[1:])

				if res != nil {
					if child.isNamedParam {
						params = append(params, Param{child.component, component})
					}
					return res, params
				}
			}
		}
	}
	return nil, nil
}

// search checks if the path is contained in the tree structure. If there is a match in the
// tree structure, it will return corresponding handler.
func (n *node) search(method string, path string) (handle Handle, params Params, found bool) {
	components := strings.Split(path, "/")[1:]
	node, params := n.searchPath(components)
	if node == nil {
		handle = nil
		found = false
		return
	}

	httpMethod := strings.ToUpper(method)
	handle = node.methods[httpMethod]
	found = true

	return
}

func printParam(p Params) {
	fmt.Println("Printing param")
	for key, value := range p {
		fmt.Println("Key: ", key, " Value: ", value)
	}
}

// String prints out the stringified version of the tree
func (n *node) String() string {
	str := n.component + "{"
	for _, child := range n.children {
		str += child.String()
	}
	str += "}"
	return str
}
