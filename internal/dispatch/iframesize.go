package dispatch

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func SelectIFrameSize(s *fairytale.State, payload hypp.Payload) hypp.Dispatchable {
	event := payload.(hypp.Event)
	value := event.Target().Value()
	size := fairytale.MustIFrameSizeFromString(value)

	newState := s.Clone()
	newState.Settings.IFrameSize = size
	postMessageToIFrame(message[struct{}]{
		Type: messageRefreshApp,
		Data: struct{}{},
	})
	return newState
}
