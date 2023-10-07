package dispatch

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/window"
)

func SelectOrientationAction[S hypp.State]() hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
		event := payload.(window.Event)
		value := event.Target().Value()
		orientation := fairytale.MustOrientationFromString(value)

		settings := s.Settings()
		settings.Orientation = orientation
		newState := s.Clone()
		newState.SetSettings(settings)
		postWindowMessageToIFrame(windowMessage[struct{}]{
			Type: windowMessageRefreshApp,
			Data: struct{}{},
		})
		return newState
	}
}
