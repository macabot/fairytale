package component

import (
	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
)

func CurrentTale(tale *state.Tale) *hypp.VNode {
	var content *hypp.VNode
	if tale == nil {
		content = hypp.Text("Select a tale")
	} else {
		content = replaceEventHandlers(tale.View(tale.State()))
	}
	return content
}

func replaceEventHandlers(vNode *hypp.VNode) *hypp.VNode {
	if vNode == nil {
		return vNode
	}
	if vNode.Kind() != hypp.SSRNode {
		return vNode
	}
	props := vNode.Props()
	if props == nil {
		props = hypp.HProps{}
	}
	for key := range props {
		if key[0] == 'o' && key[1] == 'n' {
			props[key] = dispatch.TriggerTaleEvent(key)
		}
	}
	children := vNode.Children()
	newChildren := make([]*hypp.VNode, len(children))
	for i := 0; i < len(children); i++ {
		newChildren[i] = replaceEventHandlers(children[i])
	}
	return hypp.H(
		vNode.Tag(),
		props,
		newChildren...,
	)
}
