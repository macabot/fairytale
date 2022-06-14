package fairy

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/macabot/hypp"
	jsd "github.com/macabot/hypp/driver/js"
)

var window hypp.Window = jsd.Driver{}.Window()

type Message[T any] struct {
	Type int
	Data T
}

const (
	MessageSelectTale = iota + 1
	MessageOperateControl
	MessageTaleEvent
)

type OperateControlData[T any] struct {
	TalePath     []int
	ControlIndex int
	EventData    T
}

type TaleEvent struct {
	Key   string
	Event any
}

func postMessageToIFrame[T any](message Message[T]) {
	origin := js.Global().Get("window").Get("location").Get("origin")
	iframeEl := js.Global().Get("document").Call("querySelector", "iframe")
	b, err := json.Marshal(message)
	if err != nil {
		panic(fmt.Errorf("fairy: cannot JSON marshal message with type '%d': %w", message.Type, err))
	}
	iframeEl.Get("contentWindow").Call("postMessage", string(b), origin)
}

func postMessageToTopFrame[T any](message Message[T]) {
	origin := js.Global().Get("window").Get("location").Get("origin")
	top := js.Global().Get("top")
	b, err := json.Marshal(message)
	if err != nil {
		panic(fmt.Errorf("fairy: cannot JSON marshal message with type '%d': %w", message.Type, err))
	}
	top.Call("postMessage", string(b), origin)
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
