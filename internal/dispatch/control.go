package dispatch

import (
	"encoding/json"
	"fmt"

	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

type operateControlData[T any] struct {
	TalePath     []int
	ControlIndex int
	EventData    T
}

func ChangeControlAction[S hypp.State, T any](
	talePath []int,
	controlIndex int,
	getEventData func(hypp.Event) T,
) hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
		newState := s.Clone()
		event := payload.(hypp.Event)
		tale := newState.GetTale(talePath)
		control := tale.Controls()[controlIndex]
		eventData := getEventData(event)
		// TODO pass eventData instead of event to Update method?
		tale.Dispatch(control.UpdateFromEvent(tale.State(), event), payload)
		postWindowMessageToIFrame(windowMessage[operateControlData[T]]{
			Type: windowMessageOperateControl,
			Data: operateControlData[T]{
				TalePath:     talePath,
				ControlIndex: controlIndex,
				EventData:    eventData,
			},
		})
		return newState
	}
}

func OperateControlSubscription[S hypp.State]() hypp.Subscription {
	return hypp.Subscription{
		Subscriber: subscribeToWindowMessage,
		Payload: windowMessageProps{
			Type:         windowMessageOperateControl,
			Dispatchable: operateControlAction[S](),
		},
	}
}

func operateControlAction[S hypp.State]() hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
		raw := payload.(json.RawMessage)
		var data operateControlData[json.RawMessage]
		if err := json.Unmarshal(raw, &data); err != nil {
			panic(fmt.Errorf("fairy: cannot unmarshal operateControl data '%s': %w", string(raw), err))
		}
		tale := s.GetTale(data.TalePath)
		control := tale.Controls()[data.ControlIndex]
		tale.Dispatch(control.UpdateFromMessage(tale.State(), data.EventData), payload)
		return s.Clone()
	}
}
