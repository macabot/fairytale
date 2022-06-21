package fairy

import (
	"encoding/json"
	"fmt"

	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
	"golang.org/x/exp/constraints"
)

func onChangeControl[T any](
	talePath []int,
	controlIndex int,
	getEventData func(hypp.Event) T,
) hypp.Action[*state] {
	return func(s *state, payload hypp.Payload) hypp.Dispatchable {
		newState := s.clone()
		event := payload.(hypp.Event)
		tale := newState.getTale(talePath)
		control := tale.myControls[controlIndex]
		eventData := getEventData(event)
		// TODO pass eventData instead of event to Update method?
		tale.myState = control.UpdateFromEvent(tale.myState, event)
		postMessageToIFrame(message[operateControlData[T]]{
			Type: messageOperateControl,
			Data: operateControlData[T]{
				TalePath:     talePath,
				ControlIndex: controlIndex,
				EventData:    eventData,
			},
		})
		return newState
	}
}

// Control manages the state of a Tale. Typically, a Control manages a single
// property of the state, however a Control can change the whole state.
type Control interface {
	Render(state any, talePath []int, controlIndex int) *hypp.VNode
	UpdateFromEvent(state any, event hypp.Event) any
	UpdateFromMessage(state any, data json.RawMessage) any
}

// SelectOption represents a possible option for a SelectControl.
type SelectOption[T any] struct {
	Label    string
	Value    T
	disabled bool
}

// Render the SelectOption as an <option> HTML element.
func (s SelectOption[T]) Render(selected bool) *hypp.VNode {
	b, err := json.Marshal(s.Value)
	if err != nil {
		panic(fmt.Errorf("fairy: cannot JSON marshal SelectOption value of type %T", s.Value))
	}
	return html.Option(
		hypp.HProps{
			"value":    string(b),
			"selected": selected,
			"disabled": s.disabled,
		},
		hypp.Text(s.Label),
	)
}

var _ Control = &SelectControl[struct{}, struct{}]{}

// SelectControl is a Control that let's you update the state by selecting one
// of the available options.
type SelectControl[S, T any] struct {
	label         string
	update        func(S, T) S
	selectedIndex func(S) int
	options       []SelectOption[T]
}

// NewSelectControl creates a new SelectControl.
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

// Render renders the SelectControl as a <select> HTML element.
func (s SelectControl[S, T]) Render(state any, talePath []int, controlIndex int) *hypp.VNode {
	selectedIndex := s.selectedIndex(state.(S))
	options := make([]*hypp.VNode, len(s.options))
	for i, option := range s.options {
		options[i] = option.Render(i == selectedIndex)
	}
	if selectedIndex < 0 || selectedIndex >= len(s.options) {
		options = append(options, SelectOption[T]{
			Label:    "[Selected index out of range]",
			disabled: true,
		}.Render(true))
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

// CheckboxControl is a Control that let's you update the stage by toggling a
// checkbox.
type CheckboxControl[S any] struct {
	label   string
	update  func(state S, checked bool) S
	checked func(S) bool
}

// NewCheckboxControl creates a new CheckboxControl.
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

// Render renders the CheckboxControl as a <input type="checkbox"> HTML element.
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

var _ Control = &NumberInputControl[struct{}, float64]{}

type Number interface {
	constraints.Integer | constraints.Float
}

type NumberInputControl[S any, N Number] struct {
	label  string
	update func(state S, value N) S
	value  func(S) N
	min    *N
	max    *N
	step   *N
}

func NewNumberInputControl[S any, N Number](
	label string,
	update func(state S, value N) S,
	value func(state S) N,
) *NumberInputControl[S, N] {
	return &NumberInputControl[S, N]{
		label:  label,
		update: update,
		value:  value,
	}
}

func (n *NumberInputControl[S, N]) WithMin(min N) *NumberInputControl[S, N] {
	n.min = &min
	return n
}

func (n *NumberInputControl[S, N]) WithMax(max N) *NumberInputControl[S, N] {
	n.max = &max
	return n
}

func (n *NumberInputControl[S, N]) WithStep(step N) *NumberInputControl[S, N] {
	n.step = &step
	return n
}

func (n NumberInputControl[S, N]) parseNumber(b []byte) N {
	var number N
	if err := json.Unmarshal(b, &number); err != nil {
		panic(fmt.Errorf("fairy: NumberInputControl cannot parse '%s' as type %T: %w", b, number, err))
	}
	return number
}

func (n NumberInputControl[S, N]) Render(state any, path []int, controlIndex int) *hypp.VNode {
	inputProps := hypp.HProps{
		"type":  "number",
		"value": n.value(state.(S)),
		"onchange": onChangeControl(path, controlIndex, func(event hypp.Event) N {
			return n.parseNumber([]byte(event.Target().Value()))
		}),
	}
	if n.min != nil {
		inputProps["min"] = *n.min
	}
	if n.max != nil {
		inputProps["max"] = *n.max
	}
	if n.step != nil {
		inputProps["step"] = *n.step
	}
	return html.Label(
		nil,
		hypp.Text(n.label),
		html.Input(
			inputProps,
		),
	)
}

func (n NumberInputControl[S, N]) keepInRange(number N) N {
	if n.min != nil && number < *n.min {
		number = *n.min
	}
	if n.max != nil && number > *n.max {
		number = *n.max
	}
	return number
}

func (n NumberInputControl[S, N]) UpdateFromEvent(state any, event hypp.Event) any {
	number := n.parseNumber([]byte(event.Target().Value()))
	number = n.keepInRange(number)
	return n.update(state.(S), number)
}

func (n NumberInputControl[S, N]) UpdateFromMessage(state any, data json.RawMessage) any {
	number := n.parseNumber(data)
	number = n.keepInRange(number)
	return n.update(state.(S), number)
}
