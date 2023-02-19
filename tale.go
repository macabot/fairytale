package fairytale

import (
	"encoding/json"
	"fmt"

	"github.com/gosimple/slug"
	"github.com/macabot/hypp"
)

type TaleTarget int

const (
	TaleInsideBody TaleTarget = iota
	TaleAsBody
)

type TaleSettings struct {
	Target TaleTarget
}

var _ Node = &Tale{}

// Control manages the state of a Tale. Typically, a Control manages a single
// property of the state, however a Control can change the whole state.
type Control interface {
	Render(state any, talePath []int, controlIndex int) *hypp.VNode
	UpdateFromEvent(state any, event hypp.Event) any
	UpdateFromMessage(state any, data json.RawMessage) any
}

// Tale is a Node that let's you develop and document a component.
type Tale struct {
	name          string
	slug          string
	state         any
	view          func(any) *hypp.VNode
	validateState func(any)
	controls      []Control
	settings      TaleSettings
}

// New creates a new Tale.
func New[S any](
	name string,
	state S,
	view func(S) *hypp.VNode,
) *Tale {
	return &Tale{
		name:  name,
		slug:  slug.Make(name),
		state: state,
		view:  func(state any) *hypp.VNode { return view(state.(S)) },
		validateState: func(state any) {
			if _, ok := state.(S); !ok {
				var s S
				panic(fmt.Errorf("fairytale: invalid state for Tale '%s'. Expected '%T', got '%T'.", name, s, state))
			}
		},
	}
}

func (t Tale) Name() string           { return t.name }
func (t Tale) Slug() string           { return t.slug }
func (t Tale) Children() []Node       { return nil }
func (t *Tale) Tale() *Tale           { return t }
func (t Tale) IsOpen() bool           { return false }
func (t *Tale) SetIsOpen(isOpen bool) { /* noop */ }
func (t Tale) View() *hypp.VNode      { return t.view(t.state) }
func (t Tale) Controls() []Control    { return t.controls }
func (t Tale) State() any             { return t.state }
func (t Tale) Settings() TaleSettings { return t.settings }

func (t *Tale) SetState(s any) {
	t.validateState(s)
	t.state = s
}

func (t *Tale) WithControls(controls ...Control) *Tale {
	t.controls = controls
	return t
}

func (t *Tale) WithSettings(settings TaleSettings) *Tale {
	t.settings = settings
	return t
}
