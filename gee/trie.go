package gee

import "strings"

type node struct {
	pattern string
	part    string
	// children []*node
	// this is an optimization
	children map[string]*node
	isWild   bool
}

func NewNode() *node {
	return &node{children: make(map[string]*node), isWild: false}
}

func NewNode2(part string, isWild bool) *node {
	return &node{children: make(map[string]*node), isWild: isWild, part: part}
}

// func (n *node) findChild(part string) *node {
// 	for _, child := range n.children {
// 		if child.part == part || child.isWild {
// 			return child
// 		}
// 	}

// 	return nil
//}

// func (n *node) findChildren(part string) []*node {
// 	results := make([]*node, 0)
// 	for _, child := range n.children {
// 		if child.part == part || child.isWild {
// 			results = append(results, child)
// 		}
// 	}
// 	return results
// }

func (n *node) findChildren(part string) []*node {
	results := make([]*node, 0)
	for key, child := range n.children {
		if key == part || child.isWild {
			results = append(results, child)
		}
	}
	return results
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	child, ok := n.children[parts[height]]
	if !ok {
		child = NewNode2(parts[height], parts[height][0] == ':' || parts[height][0] == '*')
		child.part = parts[height]
		n.children[parts[height]] = child
	}

	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern != "" {
			return n
		} else {
			return nil
		}
	}

	children := n.findChildren(parts[height])
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *node) travel(list *[]*node) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
