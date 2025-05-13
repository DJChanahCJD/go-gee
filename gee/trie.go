package gee

import "strings"

// 路由前缀树构建

// 1. 路由前缀树的节点结构体
type node struct {
	pattern  string // 待匹配路由，例如 /p/:lang， 为空表示非路由路径
	part     string // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool // 是否精确匹配，part 含有 : 或 * 时为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, depth int) {
	// 最后一个节点，标记为待匹配路由
	if len(parts) == depth {
		n.pattern = pattern
		return
	}

	part := parts[depth]
	child := n.matchChild(part)

	// 如果没有匹配的子节点，则创建一个新的子节点
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	// 继续向下查找
	child.insert(pattern, parts, depth+1)
}

func (n *node) search(parts []string, depth int) *node {
	// 如果当前节点是待匹配路由，则返回当前节点
	if len(parts) == depth || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[depth]
	children := n.matchChildren(part)

	// 从所有匹配成功的节点中继续查找，找到第一个匹配成功的节点
	for _, child := range children {
		result := child.search(parts, depth+1)
		if result != nil {
			return result
		}
	}

	return nil
}