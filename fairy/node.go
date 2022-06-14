package fairy

type Node interface {
	Name() string
	Children() []Node
	Tale() *Tale
	IsOpen() bool
	SetIsOpen(bool)
}

var _ Node = &Branch{}

type Branch struct {
	name     string
	children []Node
	tale     *Tale
	isOpen   bool
}

func (b Branch) Name() string           { return b.name }
func (b Branch) Children() []Node       { return b.children }
func (b Branch) Tale() *Tale            { return b.tale }
func (b Branch) IsOpen() bool           { return b.isOpen }
func (b *Branch) SetIsOpen(isOpen bool) { b.isOpen = isOpen }

func NewTree(children ...Node) Node {
	return &Branch{children: children, isOpen: true}
}

func NewBranch(name string, children ...Node) Node {
	return &Branch{name: name, children: children}
}
