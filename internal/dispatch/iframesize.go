package dispatch

import (
	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
)

func SelectIFrameSize(s *state.State, payload hypp.Payload) hypp.Dispatchable {
	event := payload.(hypp.Event)
	value := event.Target().Value()
	size := state.MustIFrameSizeFromString(value)

	newState := s.Clone()
	newState.Settings.IFrameSize = size
	postMessageToIFrame(message[struct{}]{
		Type: messageRefreshApp,
		Data: struct{}{},
	})
	return newState
}
