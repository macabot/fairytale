package dispatch

import (
	"github.com/macabot/fairytale/fairy/internal/state"
	"github.com/macabot/hypp"
)

func SelectIFrameSize(s *state.State, payload hypp.Payload) hypp.Dispatchable {
	event := payload.(hypp.Event)
	value := event.Target().Value()
	size := state.MustIFrameSizeFromString(value)

	newState := s.Clone()
	newState.Settings.IFrameSize = size
	PostMessageToIFrame(state.Message[struct{}]{
		Type: state.MessageRefreshApp,
		Data: struct{}{},
	})
	return newState
}
