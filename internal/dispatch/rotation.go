package dispatch

import (
	"fmt"

	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func SelectRotation[S hypp.State](s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
	event := payload.(hypp.Event)
	value := event.Target().Value()
	rotation := mustRotationFromString(value)

	settings := s.Settings()
	settings.Rotation = rotation
	newState := s.Clone()
	newState.SetSettings(settings)
	postMessageToIFrame(message[struct{}]{
		Type: messageRefreshApp,
		Data: struct{}{},
	})
	return newState
}

func mustRotationFromString(s string) fairytale.Rotation {
	rotation, err := rotationFromString(s)
	if err != nil {
		panic(err)
	}
	return rotation
}

func rotationFromString(s string) (fairytale.Rotation, error) {
	for _, rotation := range fairytale.Rotations {
		if rotation.String() == s {
			return rotation, nil
		}
	}
	return -1, fmt.Errorf("cannot create rotation from string '%s'", s)
}
