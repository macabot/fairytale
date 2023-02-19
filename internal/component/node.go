package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func subPath(full, sub []int) bool {
	if len(sub) > len(full) {
		return false
	}
	for i, x := range sub {
		if x != full[i] {
			return false
		}
	}
	return true
}

func Node(s *fairytale.State, n fairytale.Node, path []int, current []int) *hypp.VNode {
	isSubPath := subPath(current, path)

	children := n.Children()
	if len(children) == 0 {
		return html.A(
			hypp.HProps{
				"href": s.ToURL(path).String(),
				"class": map[string]bool{
					"selected": isSubPath,
				},
			},
			hypp.Text(n.Name()),
		)
	}

	childNodes := make([]*hypp.VNode, len(children)+1)
	childNodes[0] = html.Summary(nil, hypp.Text(n.Name()))
	for i, child := range n.Children() {
		childPath := make([]int, len(path)+1)
		copy(childPath, path)
		childPath[len(childPath)-1] = i
		childNodes[i+1] = Node(s, child, childPath, current)
	}

	return html.Details(
		hypp.HProps{
			"open": isSubPath,
		},
		childNodes...,
	)
}
