package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/jmartin127/advent-of-code-2021/helpers"
)

type Node struct {
	visited bool
	value   string
}

func (n *Node) isStart() bool {
	return n.value == "start"
}

func (n *Node) isEnd() bool {
	return n.value == "end"
}

func (n *Node) isSmall() bool {
	if n.isStart() || n.isEnd() {
		return true
	}

	runes := []rune(n.value)
	return unicode.IsLower(runes[0])
}

type Graph struct {
	edgesByFromNode map[*Node][]*Node
	allNodes        []*Node
}

type Path struct {
	nodes []*Node
}

func (p *Path) print() {
	for _, n := range p.nodes {
		fmt.Printf("%s,", n.value)
	}
	fmt.Println()
}

func (p *Path) canIncludeInPath(n *Node) bool {
	if !n.isSmall() {
		return true
	}

	for _, v := range p.nodes {
		if v.value == n.value {
			return false
		}
	}
	return true
}

func (g *Graph) print() {
	for k, v := range g.edgesByFromNode {
		for _, n := range v {
			fmt.Printf("%s --> %+v\n", k.value, n.value)
		}
	}
}

func (g *Graph) addToMap(from, to *Node) {
	if _, ok := g.edgesByFromNode[from]; ok {
		g.edgesByFromNode[from] = append(g.edgesByFromNode[from], to)
	} else {
		g.edgesByFromNode[from] = []*Node{to}
	}
}

func (g *Graph) findAllPossiblePaths(a *Node, currentPath *Path, paths []*Path) ([]*Path, *Path) {
	originalPath := currentPath.copy()
	fmt.Printf("Visited %s\n", a.value)
	currentPath.nodes = append(currentPath.nodes, a)
	if a.isEnd() {
		fmt.Println("Appending a path:")
		currentPath.print()
		paths = append(paths, currentPath.copy())
		fmt.Println()
		return paths, originalPath
	}

	a.visited = true
	destinations := g.edgesByFromNode[a]
	for _, dest := range destinations {
		if currentPath.canIncludeInPath(dest) {
			paths, currentPath = g.findAllPossiblePaths(dest, currentPath.copy(), paths)
		}
	}

	return paths, originalPath
}

func (p *Path) copy() *Path {
	copy := make([]*Node, 0)
	for _, v := range p.nodes {
		newNode := &Node{value: v.value}
		copy = append(copy, newNode)
	}
	return &Path{nodes: copy}
}

func main() {
	filepath := "input.txt"
	list := helpers.ReadFile(filepath)

	// read the input
	edgesByFromNode := make(map[*Node][]*Node, 0)
	g := Graph{edgesByFromNode: edgesByFromNode}
	nodes := make(map[string]*Node, 0)
	for _, line := range list {
		from, to := parseLine(line, nodes)
		if !from.isEnd() && !to.isStart() {
			g.addToMap(from, to)
		}
		if !to.isEnd() && !from.isStart() {
			g.addToMap(to, from)
		}
	}
	g.print()

	// determine all possible paths
	startNode := nodes["start"]
	currentPath := &Path{nodes: make([]*Node, 0)}
	allPaths := make([]*Path, 0)
	result, _ := g.findAllPossiblePaths(startNode, currentPath, allPaths)

	// list the paths
	for _, p := range result {
		p.print()
	}
	fmt.Printf("num of paths %d\n", len(result))
}

func parseLine(line string, nodes map[string]*Node) (*Node, *Node) {
	vals := strings.Split(line, "-") // start-A

	fromNode := addToNodeMap(nodes, vals[0])
	toNode := addToNodeMap(nodes, vals[1])

	return fromNode, toNode
}

func addToNodeMap(nodes map[string]*Node, val string) *Node {
	if n, ok := nodes[val]; ok {
		return n
	} else {
		n := &Node{value: val}
		nodes[val] = n
		return n
	}
}
