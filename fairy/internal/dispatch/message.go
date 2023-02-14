package dispatch

import (
	"encoding/json"
	"fmt"

	"github.com/macabot/fairytale/fairy/internal/driver"
	"github.com/macabot/fairytale/fairy/internal/state"
	"github.com/macabot/hypp"
)

func OnMessage(dispatch hypp.Dispatch, payload hypp.Payload) hypp.Unsubscribe {
	props := payload.(state.MessageProps)
	listener := func(event hypp.Event) {
		data := event.EscapeToValue().Get("data").String()
		var m state.Message[json.RawMessage]
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
