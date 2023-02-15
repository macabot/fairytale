package component

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
	"golang.org/x/exp/constraints"
)

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

var _ state.Control = &SelectControl[struct{}, struct{}]{}

// SelectControl is a Control that lets you update the state by selecting one
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
func (s SelectControl[S, T]) Render(
	state any,
	talePath []int,
	controlIndex int,
) *hypp.VNode {
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
				"onchange": dispatch.OnChangeControl(talePath, controlIndex, func(event hypp.Event) json.RawMessage {
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
func (s SelectControl[S, T]) UpdateFromMessage(
	state any,
	data json.RawMessage,
) any {
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		panic(fmt.Errorf("fairy: SelectControl cannot JSON unmarshal message data '%s' to type %T: %w", data, t, err))
	}
	return s.update(state.(S), t)
}

var _ state.Control = &CheckboxControl[struct{}]{}

// CheckboxControl is a Control that lets you update the state by toggling a
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
func (c CheckboxControl[S]) Render(
	state any,
	path []int,
	controlIndex int,
) *hypp.VNode {
	return html.Label(
		nil,
		hypp.Text(c.label),
		html.Input(
			hypp.HProps{
				"type":    "checkbox",
				"checked": c.checked(state.(S)),
				"onchange": dispatch.OnChangeControl(path, controlIndex, func(event hypp.Event) bool {
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
func (c CheckboxControl[S]) UpdateFromMessage(
	state any,
	data json.RawMessage,
) any {
	var checked bool
	if err := json.Unmarshal(data, &checked); err != nil {
		panic(fmt.Errorf("fairy: CheckboxControl cannot JSON unmarshal data '%s' to type %T: %w", data, checked, err))
	}
	return c.update(state.(S), checked)
}

var _ state.Control = &NumberInputControl[struct{}, float64]{}

type Number interface {
	constraints.Integer | constraints.Float
}

type NumberInputControl[S any, N Number] struct {
	label  string
	update func(state S, value N) S
	value  func(S) N
	min    *N
	max    *N
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

func (n NumberInputControl[S, N]) parseNumber(b []byte) N {
	var number N
	if err := json.Unmarshal(b, &number); err != nil {
		panic(fmt.Errorf("fairy: NumberInputControl cannot parse '%s' as type %T: %w", b, number, err))
	}
	return number
}

func (n NumberInputControl[S, N]) Render(
	state any,
	path []int,
	controlIndex int,
) *hypp.VNode {
	inputProps := hypp.HProps{
		"type":  "number",
		"value": fmt.Sprint(n.value(state.(S))),
		"onchange": dispatch.OnChangeControl(path, controlIndex, func(event hypp.Event) N {
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

func (n NumberInputControl[S, N]) keepInRange(number N) N {
	if n.min != nil && number < *n.min {
		number = *n.min
	}
	if n.max != nil && number > *n.max {
		number = *n.max
	}
	return number
}

func (n NumberInputControl[S, N]) UpdateFromEvent(
	state any,
	event hypp.Event,
) any {
	number := n.parseNumber([]byte(event.Target().Value()))
	number = n.keepInRange(number)
	return n.update(state.(S), number)
}

func (n NumberInputControl[S, N]) UpdateFromMessage(
	state any,
	data json.RawMessage,
) any {
	number := n.parseNumber(data)
	number = n.keepInRange(number)
	return n.update(state.(S), number)
}

var _ state.Control = &TextInputControl[struct{}]{}

type TextInputControl[S any] struct {
	label     string
	update    func(state S, text string) S
	value     func(S) string
	minLength *int
	maxLength *int
}

func NewTextInputControl[S any](
	label string,
	update func(state S, text string) S,
	value func(state S) string,
) *TextInputControl[S] {
	return &TextInputControl[S]{
		label:  label,
		update: update,
		value:  value,
	}
}

func (t *TextInputControl[S]) WithMinLength(minLength int) *TextInputControl[S] {
	t.minLength = &minLength
	return t
}

func (t *TextInputControl[S]) WithMaxLength(maxLength int) *TextInputControl[S] {
	t.maxLength = &maxLength
	return t
}

func (t TextInputControl[S]) Render(state any, path []int, controlIndex int) *hypp.VNode {
	inputProps := hypp.HProps{
		"type":  "text",
		"value": t.value(state.(S)),
		"oninput": dispatch.OnChangeControl(path, controlIndex, func(event hypp.Event) string {
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

func (t TextInputControl[S]) UpdateFromEvent(state any, event hypp.Event) any {
	text := event.Target().Value()
	return t.update(state.(S), text)
}

func (t TextInputControl[S]) UpdateFromMessage(
	state any,
	data json.RawMessage,
) any {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		panic(fmt.Errorf("fairy: TextInputControl cannot parse '%s' as type %T: %w", data, text, err))
	}
	return t.update(state.(S), text)
}

var _ state.Control = &ButtonControl[struct{}]{}

// ButtonControl is a Control that lets you update the state by clicking a
// button.
type ButtonControl[S any] struct {
	label  string
	update func(state S) S
}

// NewButtonControl creates a new ButtonControl.
func NewButtonControl[S any](label string, update func(S) S) *ButtonControl[S] {
	return &ButtonControl[S]{
		label:  label,
		update: update,
	}
}

// Render renders the ButtonControl as a <button> HTML element.
func (c ButtonControl[S]) Render(
	state any,
	path []int,
	controlIndex int,
) *hypp.VNode {
	return html.Button(
		hypp.HProps{
			"type": "button",
			"onclick": dispatch.OnChangeControl(
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

func (c ButtonControl[S]) UpdateFromEvent(state any, _ hypp.Event) any {
	return c.update(state.(S))
}

func (c ButtonControl[S]) UpdateFromMessage(state any, _ json.RawMessage) any {
	return c.update(state.(S))
}
