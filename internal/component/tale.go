package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/hypp"
)

func CurrentTale[S hypp.State](path []int, tale *fairytale.Tale[S]) *hypp.VNode {
	var content *hypp.VNode
	if tale == nil {
		content = hypp.Text("Select a tale")
	} else {
		content = replaceEventHandlers(path, tale, tale.View())
	}
	return content
}

func replaceEventHandlers[S hypp.State](
	path []int,
	tale *fairytale.Tale[S],
	vNode *hypp.VNode,
) *hypp.VNode {
	if vNode == nil {
		return vNode
	}
	if vNode.Kind() != hypp.ElementNode {
		return vNode
	}
	props := vNode.Props()
	if props == nil {
		props = hypp.HProps{}
	}
	for key := range props {
		if key[0] == 'o' && key[1] == 'n' {
			props[key] = dispatch.TriggerTaleEvent(
				path,
				tale,
				vNode,
				key,
				props[key],
			)
		}
	}
	children := vNode.Children()
	newChildren := make([]*hypp.VNode, len(children))
	for i := 0; i < len(children); i++ {
		newChildren[i] = replaceEventHandlers(path, tale, children[i])
	}
	return hypp.H(
		vNode.Tag(),
		props,
		newChildren...,
	)
}

func taleToTitle[S hypp.State](t *fairytale.Tale[S]) string {
	if t == nil {
		return "No tale has been selected"
	}
	return "The " + t.Name() + " tale"
}
