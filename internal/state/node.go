package state

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
	myName     string
	mySlug     string
	myChildren []Node
	myTale     *Tale
	myIsOpen   bool
}

func (b Branch) Name() string           { return b.myName }
func (b Branch) Slug() string           { return b.mySlug }
func (b Branch) Children() []Node       { return b.myChildren }
func (b Branch) Tale() *Tale            { return b.myTale }
func (b Branch) IsOpen() bool           { return b.myIsOpen }
func (b *Branch) SetIsOpen(isOpen bool) { b.myIsOpen = isOpen }

// NewTree creates a new navigation tree. The returned Node is the root of the
// tree.
func NewTree(children ...Node) Node {
	return &Branch{
		myChildren: children,
		myIsOpen:   true,
	}
}

// NewBranch creates a new Branch.
func NewBranch(name string, children ...Node) Node {
	return &Branch{
		myName:     name,
		mySlug:     slug.Make(name),
		myChildren: children,
	}
}
