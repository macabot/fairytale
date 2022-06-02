package fairy

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/macabot/hypp"
	jsd "github.com/macabot/hypp/driver/js"
	"github.com/macabot/hypp/tag/html"
)

var window hypp.Window = jsd.Driver{}.Window()

type AppState struct {
	hypp.EmptyState
	Tree    Node
	Current []int
	Assets  []*hypp.VNode
}

func (s AppState) getTale(path []int) *Tale {
	node := s.Tree
	for _, i := range path {
		node = node.Children()[i]
	}
	return node.Tale()
}

func (s AppState) clone() *AppState {
	return &s
}

func renderCurrentTale(tale *Tale) *hypp.VNode {
	var content *hypp.VNode
	if tale == nil {
		content = hypp.Text("Select a tale")
	} else {
		content = tale.View(tale.state)
	}
	return content
}

func selectTale(state *AppState, payload hypp.Payload) hypp.Dispatchable {
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

func operateControl(state *AppState, payload hypp.Payload) hypp.Dispatchable {
	raw := payload.(json.RawMessage)
	var data OperateControlData[json.RawMessage]
	if err := json.Unmarshal(raw, &data); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal operateControl data '%s': %w", string(raw), err))
	}
	tale := state.getTale(data.TalePath)
	control := tale.controls[data.ControlIndex]
	tale.state = control.UpdateFromMessage(tale.state, data.EventData)
	return state.clone()
}

type MessageProps struct {
	Type         int
	Dispatchable hypp.Dispatchable
}

func onMessage(dispatch hypp.Dispatch, payload hypp.Payload) hypp.Unsubscribe {
	props := payload.(MessageProps)
	listener := func(event hypp.Event) {
		data := event.EscapeToValue().Get("data").String()
		var message Message[json.RawMessage]
		if err := json.Unmarshal([]byte(data), &message); err != nil {
			panic(fmt.Errorf("Cannot unmarshal message with data: %s", data))
		}
		if message.Type == props.Type {
			dispatch(props.Dispatchable, message.Data)
		}
	}
	id := window.AddEventListener("message", listener)
	return func() {
		window.RemoveEventListener("message", id)
	}
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

func RunApp(state *AppState) {
	el := js.Global().Get("document").Call("querySelector", "html")
	if el.IsNull() {
		panic("Could not find <html> element.")
	}
	hypp.App(hypp.AppProps[*AppState]{
		Driver: jsd.Driver{},
		Init:   state,
		View: func(state *AppState) *hypp.VNode {
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
		Subscriptions: func(state *AppState) []hypp.Subscription {
			return []hypp.Subscription{
				onSelectTale(hypp.Action[*AppState](selectTale)),
				onOperateControl(hypp.Action[*AppState](operateControl)),
			}
		},
	})

	select {} // run Go forever
}
