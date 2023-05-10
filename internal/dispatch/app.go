package dispatch

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func RefreshAppSubscription[S hypp.State]() hypp.Subscription {
	return hypp.Subscription{
		Subscriber: subscribeToWindowMessage,
		Payload: windowMessageProps{
			Type:         windowMessageRefreshApp,
			Dispatchable: refreshAppAction[S](),
		},
	}
}

func refreshAppAction[S hypp.State]() hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], _ hypp.Payload) hypp.Dispatchable {
		return s.Clone()
	}
}
