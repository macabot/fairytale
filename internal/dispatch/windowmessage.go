package dispatch

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/macabot/fairytale/internal/driver"
	"github.com/macabot/hypp"
)

type windowMessageProps struct {
	Type         windowMessageType
	Dispatchable hypp.Dispatchable
}

func subscribeToWindowMessage(dispatch hypp.Dispatch, payload hypp.Payload) hypp.Unsubscribe {
	props := payload.(windowMessageProps)
	listener := func(event hypp.Event) {
		data := event.EscapeToValue().Get("data").String()
		var m windowMessage[json.RawMessage]
		if err := json.Unmarshal([]byte(data), &m); err != nil {
			panic(fmt.Errorf("fairytale: cannot unmarshal message with data: %s", data))
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

type windowMessageType string

const (
	windowMessageSelectTale     windowMessageType = "select-tale"
	windowMessageOperateControl windowMessageType = "operate-control"
	windowMessageTaleEvent      windowMessageType = "tale-event"
	windowMessageRefreshApp     windowMessageType = "refesh-app"
)

type windowMessage[T any] struct {
	Type windowMessageType
	Data T
}

func postWindowMessageToIFrame[T any](m windowMessage[T]) {
	origin := js.Global().Get("window").Get("location").Get("origin")
	iframeEl := js.Global().Get("document").Call("querySelector", "iframe")
	b, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("fairytale: cannot JSON marshal message with type '%s': %w", m.Type, err))
	}
	iframeEl.Get("contentWindow").Call("postMessage", string(b), origin)
}

func postWindowMessageToTopFrame[T any](m windowMessage[T]) {
	origin := js.Global().Get("window").Get("location").Get("origin")
	top := js.Global().Get("top")
	b, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("fairytale: cannot JSON marshal message with type '%s': %w", m.Type, err))
	}
	top.Call("postMessage", string(b), origin)
}
