package fairy

import (
	"fmt"
	"syscall/js"

	"github.com/macabot/fairytale/fairy/internal/dispatch"
	"github.com/macabot/fairytale/fairy/internal/state"
	"github.com/macabot/hypp"
	jsd "github.com/macabot/hypp/driver/js"
	"github.com/macabot/hypp/tag/html"
)

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

func renderCurrentTale(tale *state.Tale) *hypp.VNode {
	var content *hypp.VNode
	if tale == nil {
		content = hypp.Text("Select a tale")
	} else {
		content = replaceEventHandlers(tale.View(tale.State()))
	}
	return content
}

type messageProps struct {
	Type         int
	Dispatchable hypp.Dispatchable
}

func runApp(s *state.State) {
	el := js.Global().Get("document").Call("querySelector", "html")
	if el.IsNull() {
		panic("Could not find <html> element.")
	}
	hypp.App(hypp.AppProps[*state.State]{
		Driver: jsd.Driver{},
		Init:   s,
		View: func(s *state.State) *hypp.VNode {
			var assets []*hypp.VNode
			currentTale := s.GetTale(s.Current)
			if currentTale != nil {
				assets = s.Assets
			}
			headChildren := append(
				assets,
				html.Meta(hypp.HProps{"charset": "utf-8"}),
				html.Meta(hypp.HProps{
					"name":    "viewport",
					"content": "width=device-width, initial-scale=1.0",
				}),
				html.Title(nil, hypp.Text(state.TaleToTitle(currentTale))),
			)
			currentTaleNode := renderCurrentTale(currentTale)
			target := state.TaleInsideBody
			if currentTale != nil {
				target = currentTale.Settings().Target
			}
			var body *hypp.VNode
			switch target {
			case state.TaleInsideBody:
				body = html.Body(nil, currentTaleNode)
			case state.TaleAsBody:
				body = currentTaleNode
			default:
				panic(fmt.Errorf("invalid target %v", currentTale.Settings().Target))
			}
			return html.Html(
				nil,
				html.Head(
					nil,
					headChildren...,
				),
				body,
			)
		},
		Node: jsd.Node(el),
		Subscriptions: func(_ *state.State) []hypp.Subscription {
			return []hypp.Subscription{
				dispatch.OnSelectTale(hypp.Action[*state.State](dispatch.SelectTale)),
				dispatch.OnOperateControl(hypp.Action[*state.State](dispatch.OperateControl)),
				dispatch.OnRefreshApp(hypp.Action[*state.State](dispatch.RefreshApp)),
			}
		},
	})

	select {} // run Go forever
}
