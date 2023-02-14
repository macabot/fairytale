package dispatch

import (
	"net/url"
	"syscall/js"

	"github.com/macabot/fairytale/fairy/internal/driver"
	"github.com/macabot/fairytale/fairy/internal/state"
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

func updateCurrentFromLocation(u *url.URL) hypp.Action[*state.State] {
	return func(s *state.State, _ hypp.Payload) hypp.Dispatchable {
		newState := s.Clone()
		newState.UpdateCurrentFromURL(u)
		PostMessageToIFrame(state.Message[[]int]{
			Type: state.MessageSelectTale,
			Data: newState.Current,
		})
		return newState
	}
}
