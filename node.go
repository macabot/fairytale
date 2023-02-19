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

var _ Node = &Branch{}

// Branch is a Node that has at least 1 child node.
type Branch struct {
	name     string
	slug     string
	children []Node
	tale     *Tale
	isOpen   bool
}

func (b Branch) Name() string           { return b.name }
func (b Branch) Slug() string           { return b.slug }
func (b Branch) Children() []Node       { return b.children }
func (b Branch) Tale() *Tale            { return b.tale }
func (b Branch) IsOpen() bool           { return b.isOpen }
func (b *Branch) SetIsOpen(isOpen bool) { b.isOpen = isOpen }

// NewTree creates a new navigation tree. The returned Node is the root of the
// tree.
func NewTree(children ...Node) *Branch {
	return &Branch{
		children: children,
		isOpen:   true,
	}
}

// NewBranch creates a new Branch.
func NewBranch(name string, children ...Node) *Branch {
	return &Branch{
		name:     name,
		slug:     slug.Make(name),
		children: children,
	}
}
