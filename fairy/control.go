package fairy

import (
	"github.com/macabot/fairytale/internal/render/component"
	"github.com/macabot/fairytale/internal/state"
	"golang.org/x/exp/constraints"
)

// Control manages the state of a Tale. Typically, a Control manages a single
// property of the state, however a Control can change the whole state.
type Control interface {
	control() state.Control
}

// SelectOption represents a possible option for a SelectControl.
type SelectOption[T any] struct {
	Label string
	Value T
}

var _ Control = &SelectControl[struct{}, struct{}]{}

// SelectControl is a Control that lets you update the state by selecting one
// of the available options.
type SelectControl[S, T any] struct {
	selectControl *component.SelectControl[S, T]
}

// NewSelectControl creates a new SelectControl.
func NewSelectControl[S, T any](
	label string,
	update func(S, T) S,
	selectedIndex func(S) int,
	options []SelectOption[T],
) *SelectControl[S, T] {
	o := make([]component.SelectOption[T], len(options))
	for i, option := range options {
		o[i] = component.SelectOption[T]{
			Label: option.Label,
			Value: option.Value,
		}
	}
	return &SelectControl[S, T]{
		selectControl: component.NewSelectControl(label, update, selectedIndex, o),
	}
}

func (c SelectControl[S, T]) control() state.Control {
	return c.selectControl
}

var _ Control = &CheckboxControl[struct{}]{}

// CheckboxControl is a Control that lets you update the state by toggling a
// checkbox.
type CheckboxControl[S any] struct {
	checkboxControl *component.CheckboxControl[S]
}

func NewCheckboxControl[S any](
	label string,
	update func(S, bool) S,
	checked func(S) bool,
) *CheckboxControl[S] {
	return &CheckboxControl[S]{
		checkboxControl: component.NewCheckboxControl(label, update, checked),
	}
}

func (c CheckboxControl[S]) control() state.Control {
	return c.checkboxControl
}

var _ Control = &NumberInputControl[struct{}, float64]{}

type NumberInputControl[S any, N Number] struct {
	numberInputControl *component.NumberInputControl[S, N]
}

type Number interface {
	constraints.Integer | constraints.Float
}

func NewNumberInputControl[S any, N Number](
	label string,
	update func(state S, value N) S,
	value func(state S) N,
) *NumberInputControl[S, N] {
	return &NumberInputControl[S, N]{
		numberInputControl: component.NewNumberInputControl(label, update, value),
	}
}

func (c NumberInputControl[S, N]) control() state.Control {
	return c.numberInputControl
}

func (c *NumberInputControl[S, N]) WithMin(min N) *NumberInputControl[S, N] {
	c.numberInputControl = c.numberInputControl.WithMin(min)
	return c
}

func (c *NumberInputControl[S, N]) WithMax(max N) *NumberInputControl[S, N] {
	c.numberInputControl = c.numberInputControl.WithMax(max)
	return c
}

var _ Control = &TextInputControl[struct{}]{}

type TextInputControl[S any] struct {
	textInputControl *component.TextInputControl[S]
}

func NewTextInputControl[S any](
	label string,
	update func(state S, text string) S,
	value func(state S) string,
) *TextInputControl[S] {
	return &TextInputControl[S]{
		textInputControl: component.NewTextInputControl[S](label, update, value),
	}
}

func (c TextInputControl[S]) control() state.Control {
	return c.textInputControl
}

func (c *TextInputControl[S]) WithMinLength(minLength int) *TextInputControl[S] {
	c.textInputControl = c.textInputControl.WithMinLength(minLength)
	return c
}

func (c *TextInputControl[S]) WithMaxLength(maxLength int) *TextInputControl[S] {
	c.textInputControl = c.textInputControl.WithMaxLength(maxLength)
	return c
}

var _ Control = &ButtonControl[struct{}]{}

// ButtonControl is a Control that lets you update the state by clicking a
// button.
type ButtonControl[S any] struct {
	buttonControl *component.ButtonControl[S]
}

// NewButtonControl creates a new ButtonControl.
func NewButtonControl[S any](label string, update func(S) S) *ButtonControl[S] {
	return &ButtonControl[S]{
		buttonControl: component.NewButtonControl[S](label, update),
	}
}

func (c ButtonControl[S]) control() state.Control {
	return c.buttonControl
}
