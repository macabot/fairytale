package fairytale

import "github.com/gosimple/slug"

// Node in the navigation tree.
type Node interface {
	Name() string
	Slug() string
	Children() []Node
	Tale() *Tale
	IsOpen() bool
	SetIsOpen(bool)
}

var _ Node = &Bundle{}

// Bundle forms a bundle of Nodes.
type Bundle struct {
	name     string
	slug     string
	children []Node
	isOpen   bool
}

func (b Bundle) Name() string           { return b.name }
func (b Bundle) Slug() string           { return b.slug }
func (b Bundle) Children() []Node       { return b.children }
func (b Bundle) Tale() *Tale            { return nil }
func (b Bundle) IsOpen() bool           { return b.isOpen }
func (b *Bundle) SetIsOpen(isOpen bool) { b.isOpen = isOpen }

// NewBundle creates a new Branch.
func NewBundle(name string, children ...Node) *Bundle {
	return &Bundle{
		name:     name,
		slug:     slug.Make(name),
		children: children,
	}
}
