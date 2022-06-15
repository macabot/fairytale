package fairy

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/macabot/hypp"
	jsd "github.com/macabot/hypp/driver/js"
	"github.com/macabot/hypp/tag/html"
)

func triggerTaleEvent(key string) hypp.Action[*State] {
	return func(state *State, payload hypp.Payload) hypp.Dispatchable {
		event := payload.(hypp.Event)
		postMessageToTopFrame(Message[TaleEvent]{
			Type: MessageTaleEvent,
			Data: TaleEvent{Key: key, Event: event},
		})
		return state
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

func selectTale(state *State, payload hypp.Payload) hypp.Dispatchable {
	raw := payload.(json.RawMessage)
	var path []int
	if err := json.Unmarshal(raw, &path); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal selectTale data '%s': %w", string(raw), err))
	}
	if equalPaths(path, state.Current) {
		return state
	}
	newState := state.clone()
	newState.Current = path
	return newState
}

func operateControl(state *State, payload hypp.Payload) hypp.Dispatchable {
	raw := payload.(json.RawMessage)
	var data OperateControlData[json.RawMessage]
	if err := json.Unmarshal(raw, &data); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal operateControl data '%s': %w", string(raw), err))
	}
	tale := state.getTale(data.TalePath)
	control := tale.myControls[data.ControlIndex]
	tale.myState = control.UpdateFromMessage(tale.myState, data.EventData)
	return state.clone()
}

type MessageProps struct {
	Type         int
	Dispatchable hypp.Dispatchable
}

func onSelectTale(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: onMessage,
		Payload: MessageProps{
			Type:         MessageSelectTale,
			Dispatchable: dispatchable,
		},
	}
}

func onOperateControl(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: onMessage,
		Payload: MessageProps{
			Type:         MessageOperateControl,
			Dispatchable: dispatchable,
		},
	}
}

func runApp(state *State) {
	el := js.Global().Get("document").Call("querySelector", "html")
	if el.IsNull() {
		panic("Could not find <html> element.")
	}
	hypp.App(hypp.AppProps[*State]{
		Driver: jsd.Driver{},
		Init:   state,
		View: func(state *State) *hypp.VNode {
			var assets []*hypp.VNode
			currentTale := state.getTale(state.Current)
			if currentTale != nil {
				assets = state.Assets
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
		Subscriptions: func(state *State) []hypp.Subscription {
			return []hypp.Subscription{
				onSelectTale(hypp.Action[*State](selectTale)),
				onOperateControl(hypp.Action[*State](operateControl)),
			}
		},
	})

	select {} // run Go forever
}
