package fairy

import (
	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
)

type Node struct {
	node     state.Node
	children []*Node
}

func NewTree(children ...*Node) *Node {
	return &Node{
		node:     state.NewTree(),
		children: children,
	}
}

func NewBranch(name string, children ...*Node) *Node {
	return &Node{
		node:     state.NewBranch(name),
		children: children,
	}
}

func NewTale[S any](name string, s S, view func(S) *hypp.VNode) *Node {
	return &Node{
		node: state.NewTale(name, s, view),
	}
}
