package state

import (
	"encoding/json"

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

// Tale is Node that let's you develop and document a component.
type Tale struct {
	myName     string
	mySlug     string
	myState    any
	myView     func(any) *hypp.VNode
	myControls []Control
	mySettings TaleSettings
}

// NewTale creates a new Tale.
func NewTale[S any](
	name string,
	state S,
	view func(S) *hypp.VNode,
) *Tale {
	return &Tale{
		myName:  name,
		mySlug:  slug.Make(name),
		myState: state,
		myView:  func(state any) *hypp.VNode { return view(state.(S)) },
	}
}

func (t Tale) Name() string               { return t.myName }
func (t Tale) Slug() string               { return t.mySlug }
func (t Tale) Children() []Node           { return nil }
func (t *Tale) Tale() *Tale               { return t }
func (t Tale) IsOpen() bool               { return false }
func (t *Tale) SetIsOpen(isOpen bool)     { /* noop */ }
func (t Tale) View(state any) *hypp.VNode { return t.myView(state) }
func (t Tale) Controls() []Control {
	return t.myControls
}
func (t Tale) State() any {
	return t
}
func (t *Tale) SetState(s any) {
	t.myState = s
}
func (t *Tale) WithControls(controls ...Control) *Tale {
	t.myControls = controls
	return t
}
func (t *Tale) WithSettings(settings TaleSettings) *Tale {
	t.mySettings = settings
	return t
}
func (t Tale) Settings() TaleSettings {
	return t.mySettings
}

func TaleToTitle(t *Tale) string {
	if t == nil {
		return "No tale has been selected"
	}
	return "The " + t.Name() + " tale"
}
