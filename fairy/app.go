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
		postMessageToTopFrame(Message[TaleEvent]{
			Type: MessageTaleEvent,
			Data: TaleEvent{Key: key, Event: event},
		})
		return s
	}
}

func replaceEventHandlers(vNode *hypp.VNode) *hypp.VNode {
	if vNode.Kind() != hypp.SSRNode {
		return vNode
	}
	props := vNode.Props()
	if props == nil {
		return vNode
	}
	for key := range props {
		if key[0] == 'o' && key[1] == 'n' {
			props[key] = triggerTaleEvent(key)
		}
	}
	children := make([]*hypp.VNode, len(vNode.Children()))
	for i, child := range vNode.Children() {
		children[i] = replaceEventHandlers(child)
	}
	return hypp.H(
		vNode.Tag(),
		props,
		children...,
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
	var data OperateControlData[json.RawMessage]
	if err := json.Unmarshal(raw, &data); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal operateControl data '%s': %w", string(raw), err))
	}
	tale := s.getTale(data.TalePath)
	control := tale.myControls[data.ControlIndex]
	tale.myState = control.UpdateFromMessage(tale.myState, data.EventData)
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
			Type:         MessageSelectTale,
			Dispatchable: dispatchable,
		},
	}
}

func onOperateControl(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: onMessage,
		Payload: messageProps{
			Type:         MessageOperateControl,
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
			)
			return html.Html(
				nil,
				html.Head(
					nil,
					headChildren...,
				),
				html.Body(
					nil,
					renderCurrentTale(currentTale),
				),
			)
		},
		Node: jsd.Node(el),
		Subscriptions: func(_ *state) []hypp.Subscription {
			return []hypp.Subscription{
				onSelectTale(hypp.Action[*state](selectTale)),
				onOperateControl(hypp.Action[*state](operateControl)),
			}
		},
	})

	select {} // run Go forever
}
