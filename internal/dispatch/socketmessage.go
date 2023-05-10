package dispatch

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/macabot/fairytale/internal/driver"
	"github.com/macabot/fairytale/internal/model"
	"github.com/macabot/hypp"
)

func SocketMessageSubscription() hypp.Subscription {
	return hypp.Subscription{
		Subscriber: subscribeToSocketMessage,
	}
}

func subscribeToSocketMessage(dispatch hypp.Dispatch, _ hypp.Payload) hypp.Unsubscribe {
	window := driver.Window.EscapeToValue()
	href := window.Get("location").Get("href").String()
	u, err := url.Parse(href)
	if err != nil {
		panic(fmt.Errorf("fairytale: cannot parse href '%s' as url", href))
	}
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		panic(fmt.Errorf("fairytale: href has invalid scheme '%s'", u.Scheme))
	}

	u.Path = "/ws"
	u.Fragment = ""
	u.RawQuery = ""
	socket := window.Get("WebSocket").New(u.String())

	closeListener := driver.JavaScript.FuncOf(func(_ hypp.Value, _ []hypp.Value) any {
		fmt.Println("[WebSocket] Received close event.")
		return nil
	})

	errorListener := driver.JavaScript.FuncOf(func(_ hypp.Value, args []hypp.Value) any {
		event := args[0]
		fmt.Printf("[WebSocket] Received error event '%v'", event)
		return nil
	})

	messageListener := driver.JavaScript.FuncOf(func(_ hypp.Value, args []hypp.Value) any {
		event := args[0]
		data := event.Get("data").String()
		fmt.Printf("[WebSocket] Received message event with data '%s'.\n", data)
		var m model.SocketMessage
		if err := json.Unmarshal([]byte(data), &m); err != nil {
			panic(fmt.Errorf("fairytale: cannot unmarshal message with data: %s", data))
		}
		switch m.Type {
		case model.SocketMessageReload:
			reloadPage()
		default:
			// TODO print warning
		}
		return nil
	})

	openListener := driver.JavaScript.FuncOf(func(_ hypp.Value, _ []hypp.Value) any {
		fmt.Println("[WebSocket] Receive open event.")
		return nil
	})

	socket.Call("addEventListener", "close", closeListener)
	socket.Call("addEventListener", "error", errorListener)
	socket.Call("addEventListener", "message", messageListener)
	socket.Call("addEventListener", "open", openListener)

	return func() {
		socket.Call("removeEventListener", "close", closeListener)
		socket.Call("removeEventListener", "error", errorListener)
		socket.Call("removeEventListener", "message", messageListener)
		socket.Call("removeEventListener", "open", openListener)
	}
}

func reloadPage() {
	driver.Window.EscapeToValue().Get("location").Call("reload")
}
