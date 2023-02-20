package dispatch

import (
	"net/url"
	"syscall/js"

	"github.com/macabot/fairytale"
	"github.com/macabot/fairytale/internal/driver"
	"github.com/macabot/hypp"
)

func OnHashChange() hypp.Subscription {
	return hypp.Subscription{
		Subscriber: func(dispatch hypp.Dispatch, _ hypp.Payload) hypp.Unsubscribe {
			listener := func(_ hypp.Event) {
				location := js.Global().Get("location").Call("toString").String()
				u, err := url.Parse(location)
				if err != nil {
					panic(err)
				}
				dispatch(updateCurrentFromLocation(u), nil)
			}
			id := driver.Window.AddEventListener("hashchange", listener)
			return func() {
				driver.Window.RemoveEventListener("hashchange", id)
			}
		},
	}
}

func updateCurrentFromLocation(u *url.URL) hypp.Action[*fairytale.State] {
	return func(s *fairytale.State, _ hypp.Payload) hypp.Dispatchable {
		newState := s.Clone()
		newState.UpdateCurrentFromURL(u)
		postMessageToIFrame(message[[]int]{
			Type: messageSelectTale,
			Data: newState.Current(),
		})
		return newState
	}
}
