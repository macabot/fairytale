package fairy

import (
	"encoding/json"
	"fmt"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func onChangeControl[T any](
	talePath []int,
	controlIndex int,
	getEventData func(hypp.Event) T,
) hypp.Action[*State] {
	return func(state *State, payload hypp.Payload) hypp.Dispatchable {
		newState := state.clone()
		event := payload.(hypp.Event)
		tale := newState.getTale(talePath)
		control := tale.controls[controlIndex]
		eventData := getEventData(event)
		// TODO pass eventData instead of event to Update method?
		tale.state = control.UpdateFromEvent(tale.state, event)
		postMessage(Message[OperateControlData[T]]{
			Type: MessageOperateControl,
			Data: OperateControlData[T]{
				TalePath:     talePath,
				ControlIndex: controlIndex,
				EventData:    eventData,
			},
		})
		return newState
	}
}

type Control interface {
	Label() string
	Render(state any, talePath []int, controlIndex int) *hypp.VNode
	UpdateFromEvent(state any, event hypp.Event) any
	UpdateFromMessage(state any, data json.RawMessage) any
}

type SelectOption[T any] struct {
	Label string
	Value T
}

func (s SelectOption[T]) Render(selected bool) *hypp.VNode {
	b, err := json.Marshal(s.Value)
	if err != nil {
		panic(fmt.Errorf("fairy: cannot JSON marshal SelectOption value of type %T", s.Value))
	}
	return html.Option(
		hypp.HProps{
			"value":    string(b),
			"selected": selected,
		},
		hypp.Text(s.Label),
	)
}

var _ Control = &SelectControl[struct{}, struct{}]{}

type SelectControl[S, T any] struct {
	label         string
	update        func(S, T) S
	selectedIndex func(S) int
	options       []SelectOption[T]
}

func NewSelectControl[S, T any](
	label string,
	update func(S, T) S,
	selectedIndex func(S) int,
	options []SelectOption[T],
) *SelectControl[S, T] {
	return &SelectControl[S, T]{
		label:         label,
		update:        update,
		selectedIndex: selectedIndex,
		options:       options,
	}
}

func (s SelectControl[S, T]) Label() string { return s.label }
func (s SelectControl[S, T]) Render(state any, talePath []int, controlIndex int) *hypp.VNode {
	selectedIndex := s.selectedIndex(state.(S))
	options := make([]*hypp.VNode, len(s.options))
	for i, option := range s.options {
		options[i] = option.Render(i == selectedIndex)
	}
	return html.Label(
		nil,
		hypp.Text(s.label),
		html.Select(
			hypp.HProps{
				"onchange": onChangeControl(talePath, controlIndex, func(event hypp.Event) json.RawMessage {
					return []byte(event.Target().Value())
				}),
			},
			options...,
		),
	)
}
func (s SelectControl[S, T]) UpdateFromEvent(state any, event hypp.Event) any {
	value := event.Target().Value()
	var t T
	if err := json.Unmarshal([]byte(value), &t); err != nil {
		panic(fmt.Errorf("fairy: SelectControl cannot JSON unmarshal event data '%s' to type %T: %w", value, t, err))
	}
	return s.update(state.(S), t)
}
func (s SelectControl[S, T]) UpdateFromMessage(state any, data json.RawMessage) any {
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		panic(fmt.Errorf("fairy: SelectControl cannot JSON unmarshal message data '%s' to type %T: %w", data, t, err))
	}
	return s.update(state.(S), t)
}

var _ Control = &CheckboxControl[struct{}]{}

type CheckboxControl[S any] struct {
	label   string
	update  func(state S, checked bool) S
	checked func(S) bool
}

func NewCheckboxControl[S any](
	label string,
	update func(S, bool) S,
	checked func(S) bool,
) *CheckboxControl[S] {
	return &CheckboxControl[S]{
		label:   label,
		update:  update,
		checked: checked,
	}
}

func (c CheckboxControl[S]) Label() string { return c.label }
func (c CheckboxControl[S]) Render(state any, path []int, controlIndex int) *hypp.VNode {
	return html.Label(
		nil,
		hypp.Text(c.label),
		html.Input(
			hypp.HProps{
				"type":    "checkbox",
				"checked": c.checked(state.(S)),
				"onchange": onChangeControl(path, controlIndex, func(event hypp.Event) bool {
					return event.EscapeToValue().Get("target").Get("checked").Bool()
				}),
			},
		),
	)
}
func (c CheckboxControl[S]) UpdateFromEvent(state any, event hypp.Event) any {
	checked := event.EscapeToValue().Get("target").Get("checked").Bool()
	return c.update(state.(S), checked)
}
func (c CheckboxControl[S]) UpdateFromMessage(state any, data json.RawMessage) any {
	var checked bool
	if err := json.Unmarshal(data, &checked); err != nil {
		panic(fmt.Errorf("fairy: CheckboxControl cannot JSON unmarshal data '%s' to type %T: %w", data, checked, err))
	}
	return c.update(state.(S), checked)
}
