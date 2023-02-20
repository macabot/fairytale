package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func TreeView(s *fairytale.State) *hypp.VNode {
	children := s.Tree().Children()
	childNodes := make([]*hypp.VNode, len(children))
	for i, child := range children {
		childNodes[i] = Node(s, child, []int{i}, s.Current())
	}
	return html.Nav(
		hypp.HProps{"class": "tree-view"},
		childNodes...,
	)
}
