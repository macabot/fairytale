package dispatch

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/macabot/fairytale/internal/driver"
	"github.com/macabot/hypp"
)

type messageProps struct {
	Type         messageType
	Dispatchable hypp.Dispatchable
}

func onMessage(dispatch hypp.Dispatch, payload hypp.Payload) hypp.Unsubscribe {
	props := payload.(messageProps)
	listener := func(event hypp.Event) {
		data := event.EscapeToValue().Get("data").String()
		var m message[json.RawMessage]
		if err := json.Unmarshal([]byte(data), &m); err != nil {
			panic(fmt.Errorf("Cannot unmarshal message with data: %s", data))
		}
		if m.Type == props.Type {
			dispatch(props.Dispatchable, m.Data)
		}
	}
	id := driver.Window.AddEventListener("message", listener)
	return func() {
		driver.Window.RemoveEventListener("message", id)
	}
}

type messageType int

const (
	messageSelectTale messageType = iota + 1
	messageOperateControl
	messageTaleEvent
	messageRefreshApp
)

type message[T any] struct {
	Type messageType
	Data T
}

func postMessageToIFrame[T any](m message[T]) {
	origin := js.Global().Get("window").Get("location").Get("origin")
	iframeEl := js.Global().Get("document").Call("querySelector", "iframe")
	b, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("fairy: cannot JSON marshal message with type '%d': %w", m.Type, err))
	}
	iframeEl.Get("contentWindow").Call("postMessage", string(b), origin)
}

func postMessageToTopFrame[T any](m message[T]) {
	origin := js.Global().Get("window").Get("location").Get("origin")
	top := js.Global().Get("top")
	b, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("fairy: cannot JSON marshal message with type '%d': %w", m.Type, err))
	}
	top.Call("postMessage", string(b), origin)
}
