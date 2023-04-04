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
		panic(fmt.Errorf("Cannot parse href '%s' as url", href))
	}
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		panic(fmt.Errorf("href has invalid scheme '%s'", u.Scheme))
	}

	u.Path = "/ws"
	u.Fragment = ""
	u.RawQuery = ""
	socket := window.Get("WebSocket").New(u.String())

	listener := func(event hypp.Value) {
		data := event.Get("data").String()
		var m model.SocketMessage
		if err := json.Unmarshal([]byte(data), &m); err != nil {
			panic(fmt.Errorf("Cannot unmarshal message with data: %s", data))
		}
		switch m.Type {
		case model.SocketMessageReload:
			reloadPage()
		default:
			// TODO print warning
		}
	}
	f := driver.JavaScript.FuncOf(func(_ hypp.Value, args []hypp.Value) any {
		listener(args[0])
		return nil
	})
	socket.Call("addEventListener", "message", f)
	return func() {
		socket.Call("removeEventListener", "message", f)
	}
}

func reloadPage() {
	driver.Window.EscapeToValue().Get("location").Call("reload")
}
