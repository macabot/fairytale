package dispatch

import (
	"encoding/json"
	"fmt"

	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func OnTaleEvent(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: onMessage,
		Payload: messageProps{
			Type:         messageTaleEvent,
			Dispatchable: dispatchable,
		},
	}
}

func AppendTaleEvent[S hypp.State](s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
	raw := payload.(json.RawMessage)
	var event fairytale.TaleEvent
	if err := json.Unmarshal(raw, &event); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal appendTaleEvent data '%s': %w", string(raw), err))
	}
	newState := s.Clone()
	newState.SetTaleEvents(append(newState.TaleEvents(), event))
	return newState
}

func TriggerTaleEvent[S hypp.State](tale *fairytale.Tale[S], key string, value any) hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
		event := payload.(hypp.Event)
		postMessageToTopFrame(message[fairytale.TaleEvent]{
			Type: messageTaleEvent,
			Data: fairytale.TaleEvent{Key: key, Event: event},
		})

		if dispatchable, ok := value.(hypp.Dispatchable); ok {
			tale.Dispatch(dispatchable)
		} else {
			// TODO warn
		}

		return s.Clone()
	}
}

func SelectTaleByPath[S hypp.State](path []int) hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], _ hypp.Payload) hypp.Dispatchable {
		return selectTaleByPath(s, path)
	}
}

func selectTaleByPath[S hypp.State](s *fairytale.State[S], path []int) *fairytale.State[S] {
	if equalPaths(s.Current(), path) {
		return s
	}
	newState := s.Clone()
	newState.SetCurrent(path)
	newState.SetTaleEvents(nil)
	postMessageToIFrame(message[[]int]{
		Type: messageSelectTale,
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
		Subscriber: onMessage,
		Payload: messageProps{
			Type:         messageSelectTale,
			Dispatchable: dispatchable,
		},
	}
}

func SelectTale[S hypp.State](s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
	raw := payload.(json.RawMessage)
	var path []int
	if err := json.Unmarshal(raw, &path); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal selectTale data '%s': %w", string(raw), err))
	}
	if equalPaths(path, s.Current()) {
		return s
	}
	newState := s.Clone()
	newState.SetCurrent(path)
	return newState
}
