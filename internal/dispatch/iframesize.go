package dispatch

import (
	"fmt"

	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func SelectIFrameSize[S hypp.State](s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
	event := payload.(hypp.Event)
	value := event.Target().Value()
	size := mustIFrameSizeFromString(value)

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

func mustIFrameSizeFromString(s string) fairytale.IFrameSize {
	size, err := iFrameSizeFromString(s)
	if err != nil {
		panic(err)
	}
	return size
}

func iFrameSizeFromString(s string) (fairytale.IFrameSize, error) {
	for _, size := range fairytale.IFrameSizes {
		if size.String() == s {
			return size, nil
		}
	}
	return [2]int{}, fmt.Errorf("cannot create iFrameSize from string '%s'", s)
}
