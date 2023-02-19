package dispatch

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func SelectRotation(s *fairytale.State, payload hypp.Payload) hypp.Dispatchable {
	event := payload.(hypp.Event)
	value := event.Target().Value()
	rotation := fairytale.MustRotationFromString(value)

	newState := s.Clone()
	newState.Settings.Rotation = rotation
	postMessageToIFrame(message[struct{}]{
		Type: messageRefreshApp,
		Data: struct{}{},
	})
	return newState
}
