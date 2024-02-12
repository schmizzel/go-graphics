package printer

import "fmt"

// Any tree that implements this interface can be printed
type Node interface {
	GetChildren() []Node
	GetName() string
}

// PrintTree prints a tree to a string
func PrintTree(root Node) string {
	if root == nil {
		return "empty"
	}

	printer := printer{
		divider: []bool{false},
	}

	return fmt.Sprintf("└── %s", printer.processNode(root, 1))
}

// PrintNode is a default implementation of a tree node.
// It can be used to build any tree and print it.
type PrintNode struct {
	children []Node
	name     string
}

func NewTree(name string) *PrintNode {
	return &PrintNode{
		name: name,
	}
}

func (n *PrintNode) AddNode(name string) *PrintNode {
	node := &PrintNode{name: name}
	n.children = append(n.children, node)
	return node
}

func (n *PrintNode) GetChildren() []Node {
	return n.children
}

func (n *PrintNode) GetName() string {
	return n.name
}

type printer struct {
	divider []bool
}

func (p *printer) processNode(n Node, depth int) string {
	if len(p.divider) < depth {
		p.divider = append(p.divider, true)
	}

	children := n.GetChildren()
	name := n.GetName()

	if len(children) == 0 {
		return fmt.Sprintf(name)
	}

	out := name + "\n"

	for i, child := range children {
		for j := 0; j < depth; j++ {
			if j < len(p.divider) && p.divider[j] {
				out += "│"
			} else {
				out += " "
			}

			out += "    "
		}

		if i == len(children)-1 {
			out += "└── " + p.processNode(child, depth+1)
			p.divider[depth-1] = false
			continue
		}

		out += "├── " + p.processNode(child, depth+1) + "\n"
	}

	return out
}
