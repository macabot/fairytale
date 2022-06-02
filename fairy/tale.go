package fairy

import "github.com/macabot/hypp"

var _ Node = &Tale{}

type Tale struct {
	name       string
	isSelected bool
	state      any
	view       func(any) *hypp.VNode
	controls   []Control
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

func (t Tale) Name() string                   { return t.name }
func (t Tale) Children() []Node               { return nil }
func (t *Tale) Tale() *Tale                   { return t }
func (t Tale) IsSelected() bool               { return t.isSelected }
func (t *Tale) SetIsSelected(isSelected bool) { t.isSelected = isSelected }
func (t Tale) IsOpen() bool                   { return false }
func (t *Tale) SetIsOpen(isOpen bool)         {}
func (t Tale) View(state any) *hypp.VNode     { return t.view(state) }
