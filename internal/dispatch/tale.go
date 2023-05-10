package dispatch

import (
	"encoding/json"
	"fmt"

	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func TaleEventSubscription[S hypp.State]() hypp.Subscription {
	return hypp.Subscription{
		Subscriber: subscribeToWindowMessage,
		Payload: windowMessageProps{
			Type:         windowMessageTaleEvent,
			Dispatchable: appendTaleEventAction[S](),
		},
	}
}

func appendTaleEventAction[S hypp.State]() hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
		raw := payload.(json.RawMessage)
		var event fairytale.TaleEvent[S]
		if err := json.Unmarshal(raw, &event); err != nil {
			panic(fmt.Errorf("fairy: cannot unmarshal appendTaleEvent data '%s': %w", string(raw), err))
		}

		newState := s.Clone()
		tale := newState.GetTale(event.Path)
		tale.SetState(event.State)
		tale.AppendEvent(event)
		return newState
	}
}

func TaleStateSubscription[S hypp.State](tale *fairytale.Tale[S], path []int) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: func(dispatch hypp.Dispatch, payload hypp.Payload) hypp.Unsubscribe {
			tale.SetStateSubscriber(func(taleState S) {
				dispatch(hypp.Action[*fairytale.State[S]](func(s *fairytale.State[S], _ hypp.Payload) hypp.Dispatchable {
					postWindowMessageToTopFrame(windowMessage[fairytale.TaleEvent[S]]{
						Type: windowMessageTaleEvent,
						Data: fairytale.TaleEvent[S]{
							Path:  path,
							Label: "???", // FIXME
							State: taleState,
						},
					})
					return s.Clone()
				}), nil)
			})
			return func() {
				tale.SetStateSubscriber(nil)
			}
		},
	}
}

func taleDispatchAction[S hypp.State](
	tale *fairytale.Tale[S],
	dispatchable hypp.Dispatchable,
) hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
		tale.Dispatch(dispatchable, payload)
		return s
	}
}

func TriggerTaleEvent[S hypp.State](
	path []int,
	tale *fairytale.Tale[S],
	vNode *hypp.VNode,
	key string,
	value any,
) any {
	dispatchable, ok := value.(hypp.Dispatchable)
	if !ok {
		return value
	}

	return taleDispatchAction(tale, dispatchable)
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

func SelectTaleSubscription[S hypp.State]() hypp.Subscription {
	return hypp.Subscription{
		Subscriber: subscribeToWindowMessage,
		Payload: windowMessageProps{
			Type:         windowMessageSelectTale,
			Dispatchable: selectTaleAction[S](),
		},
	}
}

func selectTaleAction[S hypp.State]() hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
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
}
