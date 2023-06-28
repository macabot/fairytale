package control

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/macabot/fairytale"
	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
	"golang.org/x/exp/constraints"
)

// SelectOption represents a possible option for a Select.
type SelectOption[T any] struct {
	Label    string
	Value    T
	disabled bool
}

// Render the SelectOption as an <option> HTML element.
func (s SelectOption[T]) Render(selected bool) *hypp.VNode {
	b, err := json.Marshal(s.Value)
	if err != nil {
		panic(fmt.Errorf("fairytale: cannot JSON marshal SelectOption value of type %T", s.Value))
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

var _ fairytale.Control[hypp.EmptyState] = &Select[hypp.EmptyState, struct{}]{}

// Select is a Control that lets you update the state by selecting one
// of the available options.
type Select[S hypp.State, T comparable] struct {
	label   string
	update  func(S, T) hypp.Dispatchable
	value   func(S) T
	options []SelectOption[T]
}

// NewSelect creates a new Select.
func NewSelect[S hypp.State, T comparable](
	label string,
	update func(S, T) hypp.Dispatchable,
	value func(S) T,
	options []SelectOption[T],
) *Select[S, T] {
	return &Select[S, T]{
		label:   label,
		update:  update,
		value:   value,
		options: options,
	}
}

// Render renders the Select as a <select> HTML element.
func (s Select[S, T]) Render(
	state S,
	talePath []int,
	controlIndex int,
) *hypp.VNode {
	value := s.value(state)
	selectedIndex := -1
	for i, option := range s.options {
		if option.Value == value {
			selectedIndex = i
			break
		}
	}
	options := make([]*hypp.VNode, len(s.options))
	for i, option := range s.options {
		options[i] = option.Render(i == selectedIndex)
	}
	if selectedIndex < 0 {
		options = append(options, SelectOption[T]{
			Label:    "[Option not found]",
			disabled: true,
		}.Render(true))
	}
	return html.Label(
		nil,
		hypp.Text(s.label),
		html.Select(
			hypp.HProps{
				"onchange": dispatch.ChangeControlAction[S](talePath, controlIndex, func(event hypp.Event) json.RawMessage {
					return []byte(event.Target().Value())
				}),
			},
			options...,
		),
	)
}
func (s Select[S, T]) UpdateFromEvent(state S, event hypp.Event) hypp.Dispatchable {
	value := event.Target().Value()
	var t T
	if err := json.Unmarshal([]byte(value), &t); err != nil {
		panic(fmt.Errorf("fairytale: SelectControl cannot JSON unmarshal event data '%s' to type %T: %w", value, t, err))
	}
	return s.update(state, t)
}
func (s Select[S, T]) UpdateFromMessage(
	state S,
	data json.RawMessage,
) hypp.Dispatchable {
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		panic(fmt.Errorf("fairytale: SelectControl cannot JSON unmarshal message data '%s' to type %T: %w", data, t, err))
	}
	return s.update(state, t)
}

var _ fairytale.Control[hypp.EmptyState] = &Checkbox[hypp.EmptyState]{}

// Checkbox is a Control that lets you update the state by toggling a
// checkbox.
type Checkbox[S hypp.State] struct {
	label   string
	update  func(state S, checked bool) hypp.Dispatchable
	checked func(S) bool
}

// NewCheckbox creates a new Checkbox.
func NewCheckbox[S hypp.State](
	label string,
	update func(S, bool) hypp.Dispatchable,
	checked func(S) bool,
) *Checkbox[S] {
	return &Checkbox[S]{
		label:   label,
		update:  update,
		checked: checked,
	}
}

// Render renders the Checkbox as a <input type="checkbox"> HTML element.
func (c Checkbox[S]) Render(
	state S,
	path []int,
	controlIndex int,
) *hypp.VNode {
	return html.Label(
		nil,
		hypp.Text(c.label),
		html.Input(
			hypp.HProps{
				"type":    "checkbox",
				"checked": c.checked(state),
				"onchange": dispatch.ChangeControlAction[S](path, controlIndex, func(event hypp.Event) bool {
					return event.EscapeToValue().Get("target").Get("checked").Bool()
				}),
			},
		),
	)
}
func (c Checkbox[S]) UpdateFromEvent(state S, event hypp.Event) hypp.Dispatchable {
	checked := event.EscapeToValue().Get("target").Get("checked").Bool()
	return c.update(state, checked)
}
func (c Checkbox[S]) UpdateFromMessage(
	state S,
	data json.RawMessage,
) hypp.Dispatchable {
	var checked bool
	if err := json.Unmarshal(data, &checked); err != nil {
		panic(fmt.Errorf("fairytale: Checkbox cannot JSON unmarshal data '%s' to type %T: %w", data, checked, err))
	}
	return c.update(state, checked)
}

var _ fairytale.Control[hypp.EmptyState] = &NumberInput[hypp.EmptyState, float64]{}

type Number interface {
	constraints.Integer | constraints.Float
}

type NumberInput[S hypp.State, N Number] struct {
	label  string
	update func(state S, value N) hypp.Dispatchable
	value  func(S) N
	min    *N
	max    *N
}

func NewNumberInput[S hypp.State, N Number](
	label string,
	update func(state S, value N) hypp.Dispatchable,
	value func(state S) N,
) *NumberInput[S, N] {
	return &NumberInput[S, N]{
		label:  label,
		update: update,
		value:  value,
	}
}

func (n *NumberInput[S, N]) WithMin(min N) *NumberInput[S, N] {
	n.min = &min
	return n
}

func (n *NumberInput[S, N]) WithMax(max N) *NumberInput[S, N] {
	n.max = &max
	return n
}

func (n NumberInput[S, N]) parseNumber(b []byte) N {
	var number N
	if err := json.Unmarshal(b, &number); err != nil {
		panic(fmt.Errorf("fairytale: NumberInput cannot parse '%s' as type %T: %w", b, number, err))
	}
	return number
}

func (n NumberInput[S, N]) Render(
	state S,
	path []int,
	controlIndex int,
) *hypp.VNode {
	inputProps := hypp.HProps{
		"type":  "number",
		"value": fmt.Sprint(n.value(state)),
		"onchange": dispatch.ChangeControlAction[S](path, controlIndex, func(event hypp.Event) N {
			return n.parseNumber([]byte(event.Target().Value()))
		}),
	}
	if n.min != nil {
		inputProps["min"] = fmt.Sprint(*n.min)
	}
	if n.max != nil {
		inputProps["max"] = fmt.Sprint(*n.max)
	}
	return html.Label(
		nil,
		hypp.Text(n.label),
		html.Input(inputProps),
	)
}

func (n NumberInput[S, N]) keepInRange(number N) N {
	if n.min != nil && number < *n.min {
		number = *n.min
	}
	if n.max != nil && number > *n.max {
		number = *n.max
	}
	return number
}

func (n NumberInput[S, N]) UpdateFromEvent(
	state S,
	event hypp.Event,
) hypp.Dispatchable {
	number := n.parseNumber([]byte(event.Target().Value()))
	number = n.keepInRange(number)
	return n.update(state, number)
}

func (n NumberInput[S, N]) UpdateFromMessage(
	state S,
	data json.RawMessage,
) hypp.Dispatchable {
	number := n.parseNumber(data)
	number = n.keepInRange(number)
	return n.update(state, number)
}

var _ fairytale.Control[hypp.EmptyState] = &TextInput[hypp.EmptyState]{}

type TextInput[S hypp.State] struct {
	label     string
	update    func(state S, text string) hypp.Dispatchable
	value     func(S) string
	minLength *int
	maxLength *int
}

func NewTextInput[S hypp.State](
	label string,
	update func(state S, text string) hypp.Dispatchable,
	value func(state S) string,
) *TextInput[S] {
	return &TextInput[S]{
		label:  label,
		update: update,
		value:  value,
	}
}

func (t *TextInput[S]) WithMinLength(minLength int) *TextInput[S] {
	t.minLength = &minLength
	return t
}

func (t *TextInput[S]) WithMaxLength(maxLength int) *TextInput[S] {
	t.maxLength = &maxLength
	return t
}

func (t TextInput[S]) Render(state S, path []int, controlIndex int) *hypp.VNode {
	inputProps := hypp.HProps{
		"type":  "text",
		"value": t.value(state),
		"oninput": dispatch.ChangeControlAction[S](path, controlIndex, func(event hypp.Event) string {
			return event.Target().Value()
		}),
	}
	if t.minLength != nil {
		inputProps["minlength"] = strconv.Itoa(*t.minLength)
	}
	if t.maxLength != nil {
		inputProps["maxlength"] = strconv.Itoa(*t.maxLength)
	}
	return html.Label(
		nil,
		hypp.Text(t.label),
		html.Input(inputProps),
	)
}

func (t TextInput[S]) UpdateFromEvent(state S, event hypp.Event) hypp.Dispatchable {
	text := event.Target().Value()
	return t.update(state, text)
}

func (t TextInput[S]) UpdateFromMessage(
	state S,
	data json.RawMessage,
) hypp.Dispatchable {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		panic(fmt.Errorf("fairytale: TextInput cannot parse '%s' as type %T: %w", data, text, err))
	}
	return t.update(state, text)
}

var _ fairytale.Control[hypp.EmptyState] = &Button[hypp.EmptyState]{}

// Button is a Control that lets you update the state by clicking a
// button.
type Button[S hypp.State] struct {
	label  string
	update func(state S) hypp.Dispatchable
}

// NewButton creates a new Button.
func NewButton[S hypp.State](
	label string,
	update func(S) hypp.Dispatchable,
) *Button[S] {
	return &Button[S]{
		label:  label,
		update: update,
	}
}

// Render renders the Button as a <button> HTML element.
func (c Button[S]) Render(
	state S,
	path []int,
	controlIndex int,
) *hypp.VNode {
	return html.Button(
		hypp.HProps{
			"type": "button",
			"onclick": dispatch.ChangeControlAction[S](
				path,
				controlIndex,
				func(_ hypp.Event) struct{} {
					return struct{}{}
				},
			),
		},
		hypp.Text(c.label),
	)
}

func (c Button[S]) UpdateFromEvent(state S, _ hypp.Event) hypp.Dispatchable {
	return c.update(state)
}

func (c Button[S]) UpdateFromMessage(state S, _ json.RawMessage) hypp.Dispatchable {
	return c.update(state)
}
