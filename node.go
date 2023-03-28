package fairytale

import (
	"github.com/gosimple/slug"
	"github.com/macabot/hypp"
)

// Node in the navigation tree.
type Node[S hypp.State] interface {
	Name() string
	Slug() string
	Children() []Node[S]
	Tale() *Tale[S]
	IsOpen() bool
	SetIsOpen(bool)
}

var _ Node[hypp.EmptyState] = &Bundle[hypp.EmptyState]{}

// Bundle forms a bundle of Nodes.
type Bundle[S hypp.State] struct {
	name     string
	slug     string
	children []Node[S]
	isOpen   bool
}

func (b Bundle[S]) Name() string           { return b.name }
func (b Bundle[S]) Slug() string           { return b.slug }
func (b Bundle[S]) Children() []Node[S]    { return b.children }
func (b Bundle[S]) Tale() *Tale[S]         { return nil }
func (b Bundle[S]) IsOpen() bool           { return b.isOpen }
func (b *Bundle[S]) SetIsOpen(isOpen bool) { b.isOpen = isOpen }

// NewBundle creates a new Branch.
func NewBundle[S hypp.State](name string, children ...Node[S]) *Bundle[S] {
	return &Bundle[S]{
		name:     name,
		slug:     slug.Make(name),
		children: children,
	}
}
