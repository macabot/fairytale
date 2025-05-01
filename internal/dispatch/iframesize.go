package dispatch

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/window"
)

func SelectIFrameSize[S hypp.State](s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
	event := payload.(window.Event)
	value := event.Target().Value()
	size := fairytale.MustIFrameSizeFromString(value)

	settings := s.Settings()
	settings.IFrameSize = size
	newState := s.Clone()
	newState.SetSettings(settings)
	postWindowMessageToIFrame(windowMessage[struct{}]{
		Type: windowMessageRefreshApp,
		Data: struct{}{},
	})
	return newState
}
