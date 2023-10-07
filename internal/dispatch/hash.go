package dispatch

import (
	"net/url"

	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/js"
	"github.com/macabot/hypp/window"
)

func HashChangeSubscription[S hypp.State]() hypp.Subscription {
	return hypp.Subscription{
		Subscriber: func(dispatch hypp.Dispatch, _ hypp.Payload) hypp.Unsubscribe {
			listener := func(_ window.Event) {
				location := js.Global().Get("location").Call("toString").String()
				u, err := url.Parse(location)
				if err != nil {
					panic(err)
				}
				dispatch(updateCurrentFromLocationAction[S](u), nil)
			}
			id := window.AddEventListener("hashchange", listener)
			return func() {
				window.RemoveEventListener("hashchange", id)
			}
		},
	}
}

func updateCurrentFromLocationAction[S hypp.State](u *url.URL) hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], _ hypp.Payload) hypp.Dispatchable {
		newState := s.Clone()
		newState.UpdateCurrentFromURL(u)
		postWindowMessageToIFrame(windowMessage[[]int]{
			Type: windowMessageSelectTale,
			Data: newState.Current(),
		})
		return newState
	}
}
