package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/hypp"
)

func CurrentTale(tale *fairytale.Tale) *hypp.VNode {
	var content *hypp.VNode
	if tale == nil {
		content = hypp.Text("Select a tale")
	} else {
		content = replaceEventHandlers(tale.View())
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

func taleToTitle(t *fairytale.Tale) string {
	if t == nil {
		return "No tale has been selected"
	}
	return "The " + t.Name() + " tale"
}
