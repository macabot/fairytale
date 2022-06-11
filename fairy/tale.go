package fairy

import "github.com/macabot/hypp"

var _ Node = &Tale{}

type Tale struct {
	name     string
	state    any
	view     func(any) *hypp.VNode
	controls []Control
	parent   Node
}

func NewTale[S any](
	name string,
	state S,
	view func(S) *hypp.VNode,
	controls []Control,
) *Tale {
	return &Tale{
		name:     name,
		state:    state,
		view:     func(state any) *hypp.VNode { return view(state.(S)) },
		controls: controls,
	}
}

func (t Tale) Name() string               { return t.name }
func (t Tale) Children() []Node           { return nil }
func (t *Tale) Tale() *Tale               { return t }
func (t Tale) IsOpen() bool               { return false }
func (t *Tale) SetIsOpen(isOpen bool)     { /* noop */ }
func (t *Tale) Path() []int               { return getNodePath(t) }
func (t Tale) Parent() Node               { return t.parent }
func (t *Tale) SetParent(parent Node)     { t.parent = parent }
func (t Tale) View(state any) *hypp.VNode { return t.view(state) }
func (t Tale) State() any                 { return t.state }
func (t *Tale) SetState(state any)        { t.state = state }
