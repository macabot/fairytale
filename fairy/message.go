package fairy

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/macabot/hypp"
	jsd "github.com/macabot/hypp/driver/js"
)

var window hypp.Window = jsd.Driver{}.Window()

type message[T any] struct {
	Type int
	Data T
}

const (
	messageSelectTale = iota + 1
	messageOperateControl
	messageTaleEvent
)

type operateControlData[T any] struct {
	TalePath     []int
	ControlIndex int
	EventData    T
}

type taleEvent struct {
	Key   string
	Event any
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
	id := window.AddEventListener("message", listener)
	return func() {
		window.RemoveEventListener("message", id)
	}
}
