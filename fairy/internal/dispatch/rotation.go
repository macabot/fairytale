package dispatch

import (
	"github.com/macabot/fairytale/fairy/internal/state"
	"github.com/macabot/hypp"
)

func SelectRotation(s *state.State, payload hypp.Payload) hypp.Dispatchable {
	event := payload.(hypp.Event)
	value := event.Target().Value()
	rotation := state.MustRotationFromString(value)

	newState := s.Clone()
	newState.Settings.Rotation = rotation
	PostMessageToIFrame(state.Message[struct{}]{
		Type: state.MessageRefreshApp,
		Data: struct{}{},
	})
	return newState
}
