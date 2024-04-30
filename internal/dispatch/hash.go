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
				dispatch(updateCurrentFromLocation[S], u)
			}
			id := window.AddEventListener("hashchange", listener)
			return func() {
				window.RemoveEventListener("hashchange", id)
			}
		},
	}
}

func updateCurrentFromLocation[S hypp.State](s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
	u := payload.(*url.URL)
	newState := s.Clone()
	newState.UpdateCurrentFromURL(u)
	postWindowMessageToIFrame(windowMessage[[]int]{
		Type: windowMessageSelectTale,
		Data: newState.Current(),
	})
	return newState
}
