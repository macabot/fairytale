package dispatch

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func OnRefreshApp(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: onMessage,
		Payload: messageProps{
			Type:         messageRefreshApp,
			Dispatchable: dispatchable,
		},
	}
}

func RefreshApp[S hypp.State](s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
	return s.Clone()
}
