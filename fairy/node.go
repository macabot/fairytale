package fairy

import (
	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
)

type Node interface {
	node() state.Node
}

var _ Node = &Branch{}

type Branch struct {
	branch *state.Branch
}

func NewTree(children ...Node) *Branch {
	return NewBranch("", children...)
}

func NewBranch(name string, children ...Node) *Branch {
	c := make([]state.Node, len(children))
	for i, child := range children {
		c[i] = child.node()
	}
	return &Branch{
		branch: state.NewBranch(name, c...),
	}
}

func (b Branch) node() state.Node {
	return b.branch
}

type TaleTarget int

const (
	TaleInsideBody TaleTarget = iota
	TaleAsBody
)

type TaleSettings struct {
	Target TaleTarget
}

var _ Node = &Tale[struct{}]{}

type Tale[S any] struct {
	tale *state.Tale
}

func NewTale[S any](name string, s S, view func(S) *hypp.VNode) *Tale[S] {
	return &Tale[S]{
		tale: state.NewTale[S](name, s, view),
	}
}

func (t Tale[S]) node() state.Node {
	return t.tale
}

func (t *Tale[S]) WithControls(controls ...Control) *Tale[S] {
	c := make([]state.Control, len(controls))
	for i, control := range controls {
		c[i] = control.control()
	}
	t.tale = t.tale.WithControls(c...)
	return t
}
func (t *Tale[S]) WithSettings(settings TaleSettings) *Tale[S] {
	t.tale = t.tale.WithSettings(state.TaleSettings{
		Target: state.TaleTarget(settings.Target),
	})
	return t
}
