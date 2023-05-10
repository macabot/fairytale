package dispatch

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func SelectRotationAction[S hypp.State]() hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
		event := payload.(hypp.Event)
		value := event.Target().Value()
		rotation := fairytale.MustRotationFromString(value)

		settings := s.Settings()
		settings.Rotation = rotation
		newState := s.Clone()
		newState.SetSettings(settings)
		postWindowMessageToIFrame(windowMessage[struct{}]{
			Type: windowMessageRefreshApp,
			Data: struct{}{},
		})
		return newState
	}
}
