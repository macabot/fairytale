package fairy

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/macabot/hypp"
	jsd "github.com/macabot/hypp/driver/js"
	"github.com/macabot/hypp/tag/html"
)

func triggerTaleEvent(key string) hypp.Action[*state] {
	return func(s *state, payload hypp.Payload) hypp.Dispatchable {
		event := payload.(hypp.Event)
		postMessageToTopFrame(message[taleEvent]{
			Type: messageTaleEvent,
			Data: taleEvent{Key: key, Event: event},
		})
		return s
	}
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
			props[key] = triggerTaleEvent(key)
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

func renderCurrentTale(tale *Tale) *hypp.VNode {
	var content *hypp.VNode
	if tale == nil {
		content = hypp.Text("Select a tale")
	} else {
		content = replaceEventHandlers(tale.view(tale.myState))
	}
	return content
}

func selectTale(s *state, payload hypp.Payload) hypp.Dispatchable {
	raw := payload.(json.RawMessage)
	var path []int
	if err := json.Unmarshal(raw, &path); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal selectTale data '%s': %w", string(raw), err))
	}
	if equalPaths(path, s.Current) {
		return s
	}
	newState := s.clone()
	newState.Current = path
	return newState
}

func operateControl(s *state, payload hypp.Payload) hypp.Dispatchable {
	raw := payload.(json.RawMessage)
	var data operateControlData[json.RawMessage]
	if err := json.Unmarshal(raw, &data); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal operateControl data '%s': %w", string(raw), err))
	}
	tale := s.getTale(data.TalePath)
	control := tale.myControls[data.ControlIndex]
	tale.myState = control.UpdateFromMessage(tale.myState, data.EventData)
	return s.clone()
}

func changeTota11y(s *state, payload hypp.Payload) hypp.Dispatchable {
	raw := payload.(json.RawMessage)
	var enabled bool
	if err := json.Unmarshal(raw, &enabled); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal changeTota11y data '%s': %w", string(raw), err))
	}
	newState := s.clone()
	newState.Settings.tota11y = enabled
	if !enabled {
		document := js.Global().Get("document")
		toolbar := document.Call("getElementById", "tota11y-toolbar")
		if !toolbar.IsNull() {
			toolbar.Get("parentElement").Call("removeChild", toolbar)
		}
		elements := document.Call("querySelectorAll", ".tota11y")
		for i := 0; i < elements.Length(); i++ {
			element := elements.Index(i)
			element.Get("parentElement").Call("removeChild", element)
		}
	}
	return newState
}

func refreshApp(s *state, payload hypp.Payload) hypp.Dispatchable {
	return s.clone()
}

type messageProps struct {
	Type         int
	Dispatchable hypp.Dispatchable
}

func onSelectTale(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: onMessage,
		Payload: messageProps{
			Type:         messageSelectTale,
			Dispatchable: dispatchable,
		},
	}
}

func onOperateControl(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: onMessage,
		Payload: messageProps{
			Type:         messageOperateControl,
			Dispatchable: dispatchable,
		},
	}
}

func onToggleTota11y(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: onMessage,
		Payload: messageProps{
			Type:         messageToggleTota11y,
			Dispatchable: dispatchable,
		},
	}
}

func onRefreshApp(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: onMessage,
		Payload: messageProps{
			Type:         messageRefreshApp,
			Dispatchable: dispatchable,
		},
	}
}

func runApp(s *state) {
	el := js.Global().Get("document").Call("querySelector", "html")
	if el.IsNull() {
		panic("Could not find <html> element.")
	}
	hypp.App(hypp.AppProps[*state]{
		Driver: jsd.Driver{},
		Init:   s,
		View: func(s *state) *hypp.VNode {
			var assets []*hypp.VNode
			currentTale := s.getTale(s.Current)
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
				html.Title(nil, hypp.Text(taleToTitle(currentTale))),
			)
			var tota11yScript *hypp.VNode
			if s.Settings.tota11y {
				tota11yScript = html.Script(hypp.HProps{"src": "https://cdnjs.cloudflare.com/ajax/libs/tota11y/0.1.6/tota11y.min.js"})
			}
			return html.Html(
				nil,
				html.Head(
					nil,
					headChildren...,
				),
				html.Body(
					nil,
					renderCurrentTale(currentTale),
					tota11yScript,
				),
			)
		},
		Node: jsd.Node(el),
		Subscriptions: func(_ *state) []hypp.Subscription {
			return []hypp.Subscription{
				onSelectTale(hypp.Action[*state](selectTale)),
				onOperateControl(hypp.Action[*state](operateControl)),
				onToggleTota11y(hypp.Action[*state](changeTota11y)),
				onRefreshApp(hypp.Action[*state](refreshApp)),
			}
		},
	})

	select {} // run Go forever
}
