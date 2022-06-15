package fairy

import "github.com/macabot/hypp"

var _ Node = &Tale{}

// Tale is Node that let's you develop and document a component.
type Tale struct {
	myName     string
	myState    any
	myView     func(any) *hypp.VNode
	myControls []Control
}

// NewTale creates a new Tale.
func NewTale[S any](
	name string,
	state S,
	view func(S) *hypp.VNode,
	controls []Control,
) *Tale {
	return &Tale{
		myName:     name,
		myState:    state,
		myView:     func(state any) *hypp.VNode { return view(state.(S)) },
		myControls: controls,
	}
}

func (t Tale) name() string               { return t.myName }
func (t Tale) children() []Node           { return nil }
func (t *Tale) tale() *Tale               { return t }
func (t Tale) isOpen() bool               { return false }
func (t *Tale) setIsOpen(isOpen bool)     { /* noop */ }
func (t Tale) view(state any) *hypp.VNode { return t.myView(state) }
