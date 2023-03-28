package dispatch

import (
	"net/url"
	"syscall/js"

	"github.com/macabot/fairytale"
	"github.com/macabot/fairytale/internal/driver"
	"github.com/macabot/hypp"
)

func OnHashChange[S hypp.State]() hypp.Subscription {
	return hypp.Subscription{
		Subscriber: func(dispatch hypp.Dispatch, _ hypp.Payload) hypp.Unsubscribe {
			listener := func(_ hypp.Event) {
				location := js.Global().Get("location").Call("toString").String()
				u, err := url.Parse(location)
				if err != nil {
					panic(err)
				}
				dispatch(updateCurrentFromLocation[S](u), nil)
			}
			id := driver.Window.AddEventListener("hashchange", listener)
			return func() {
				driver.Window.RemoveEventListener("hashchange", id)
			}
		},
	}
}

func updateCurrentFromLocation[S hypp.State](u *url.URL) hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], _ hypp.Payload) hypp.Dispatchable {
		newState := s.Clone()
		newState.UpdateCurrentFromURL(u)
		postMessageToIFrame(message[[]int]{
			Type: messageSelectTale,
			Data: newState.Current(),
		})
		return newState
	}
}
