package dispatch

import (
	"encoding/json"
	"fmt"

	"github.com/macabot/fairytale/fairy/internal/state"
	"github.com/macabot/hypp"
)

func OnTaleEvent(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: OnMessage,
		Payload: state.MessageProps{
			Type:         state.MessageTaleEvent,
			Dispatchable: dispatchable,
		},
	}
}

func AppendTaleEvent(s *state.State, payload hypp.Payload) hypp.Dispatchable {
	raw := payload.(json.RawMessage)
	var event state.TaleEvent
	if err := json.Unmarshal(raw, &event); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal appendTaleEvent data '%s': %w", string(raw), err))
	}
	newState := s.Clone()
	newState.TaleEvents = append(newState.TaleEvents, event)
	return newState
}

func TriggerTaleEvent(key string) hypp.Action[*state.State] {
	return func(s *state.State, payload hypp.Payload) hypp.Dispatchable {
		event := payload.(hypp.Event)
		PostMessageToTopFrame(state.Message[state.TaleEvent]{
			Type: state.MessageTaleEvent,
			Data: state.TaleEvent{Key: key, Event: event},
		})
		return s
	}
}

func SelectTaleByPath(path []int) hypp.Action[*state.State] {
	return func(s *state.State, _ hypp.Payload) hypp.Dispatchable {
		return selectTaleByPath(s, path)
	}
}

func selectTaleByPath(s *state.State, path []int) *state.State {
	if equalPaths(s.Current, path) {
		return s
	}
	newState := s.Clone()
	newState.Current = path
	newState.TaleEvents = nil
	PostMessageToIFrame(state.Message[[]int]{
		Type: state.MessageSelectTale,
		Data: path,
	})
	return newState
}

func equalPaths(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		if x != b[i] {
			return false
		}
	}
	return true
}

func OnSelectTale(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: OnMessage,
		Payload: state.MessageProps{
			Type:         state.MessageSelectTale,
			Dispatchable: dispatchable,
		},
	}
}

func SelectTale(s *state.State, payload hypp.Payload) hypp.Dispatchable {
	raw := payload.(json.RawMessage)
	var path []int
	if err := json.Unmarshal(raw, &path); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal selectTale data '%s': %w", string(raw), err))
	}
	if equalPaths(path, s.Current) {
		return s
	}
	newState := s.Clone()
	newState.Current = path
	return newState
}
